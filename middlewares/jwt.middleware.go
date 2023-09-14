package middlewares

import (
	"context"
	"errors"
	"os"
	"time"

	"virtual_workflow_management_system_gin/models"
	"virtual_workflow_management_system_gin/responses"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Claims struct {
	Username string          `json:"username"`
	Role     models.UserRole `json:"role"`
	jwt.StandardClaims
}

type RedisClientInterface interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

func JWTAuthMiddleware(redisClient RedisClientInterface) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		accessToken := ctx.GetHeader("Authorization")
		if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
			accessToken = accessToken[7:]
		}
		if accessToken == "" {
			responses.Error(ctx, "unauthorized")
			ctx.Abort()
			return
		}
		claims, err := ParseJWTToken(accessToken)
		if err != nil {
			responses.Error(ctx, "unauthorized")
			ctx.Abort()
			return
		}
		if time.Until(time.Unix(claims.ExpiresAt, 0)) < 30*time.Second {
			responses.ErrorWithToken(ctx, "token expired")
		}

		redisToken, err := redisClient.Get(context.Background(), "access_token_"+claims.Username).Result()
		if err != nil || redisToken != accessToken {
			responses.Error(ctx, "unauthorized")
			ctx.Abort()
			return
		}

		user := models.JWTUser{
			Username: claims.Username,
			Role:     claims.Role,
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}

func GenerateJWTToken(user models.User, redisClient RedisClientInterface) (map[string]string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		logrus.Error(err)
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "virtual_workflow_management_system_gin",
			IssuedAt:  time.Now().Unix(),
			Subject:   "access token",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	refreshExpirationTime := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpirationTime.Unix(),
			Issuer:    "virtual_workflow_management_system_gin",
			IssuedAt:  time.Now().Unix(),
			Subject:   "refresh token",
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	err = redisClient.Set(context.Background(), "access_token_"+user.Username, accessTokenString, 24*time.Hour).Err()
	if err != nil {
		logrus.Error("Could not set access token in Redis: ", err)
		return nil, err
	}

	err = redisClient.Set(context.Background(), "refresh_token_"+user.Username, refreshTokenString, 7*24*time.Hour).Err()
	if err != nil {
		logrus.Error("Could not set refresh token in Redis: ", err)
		return nil, err
	}

	return map[string]string{
		"access_token":  accessTokenString,
		"refresh_token": refreshTokenString,
	}, nil
}

func ParseJWTToken(accessTokenString string) (*Claims, error) {
	claims := &Claims{}
	accessToken, err := jwt.ParseWithClaims(accessTokenString, claims, func(accessToken *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	if accessToken.Valid {
		return claims, nil
	}

	return nil, err
}

func RefreshJWTToken(refreshTokenString string, redisClient RedisClientInterface) (map[string]string, error) {
	claims, err := ParseJWTToken(refreshTokenString)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	if claims.Subject != "refresh token" {
		return nil, errors.New("invalid refresh token")
	}

	if time.Until(time.Unix(claims.ExpiresAt, 0)) < 30*time.Second {
		return nil, errors.New("refresh token expired")
	}

	redisRefreshToken, err := redisClient.Get(context.Background(), "refresh_token_"+claims.Username).Result()
	if err != nil || redisRefreshToken != refreshTokenString {
		return nil, errors.New("invalid or expired refresh token")
	}

	newToken, err := GenerateJWTToken(models.User{Username: claims.Username}, redisClient)
	if err != nil {
		return nil, err
	}

	err = redisClient.Del(context.Background(), "refresh_token_"+claims.Username).Err()
	if err != nil {
		logrus.Error("failed to delete refresh token in Redis: ", err)
		return nil, err
	}

	err = redisClient.Set(context.Background(), "refresh_token_"+claims.Username, newToken["refresh_token"], time.Until(time.Unix(claims.ExpiresAt, 0))).Err()
	if err != nil {
		logrus.Error("failed to set refresh token in Redis: ", err)
		return nil, err
	}

	return newToken, nil
}

func DeleteJWTToken(username string, redisClient RedisClientInterface) error {
	_, err := redisClient.Del(context.Background(), "access_token_"+username).Result()
	if err != nil {
		logrus.Error("failed to delete access token in Redis: ", err)
		return err
	}

	_, err = redisClient.Del(context.Background(), "refresh_token_"+username).Result()
	if err != nil {
		logrus.Error("failed to delete refresh token in Redis: ", err)
		return err
	}

	return nil
}
