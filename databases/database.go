package databases

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Resource struct {
	MongoDB *mongo.Database
	Redis   *redis.Client
}

func (r *Resource) Close() {
	logrus.Warning("Closing all db connections")
}

func InitResource() (*Resource, error) {
	err := godotenv.Load(".env")
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// Initialize MongoDB
	mongoHost := os.Getenv("MONGO_HOST")
	mongoDBName := os.Getenv("MONGO_DB_NAME")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoHost))
	defer cancel()

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		Addr:     os.Getenv("REDIS_HOST"),
	})

	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &Resource{
		MongoDB: mongoClient.Database(mongoDBName),
		Redis:   redisClient,
	}, nil
}
