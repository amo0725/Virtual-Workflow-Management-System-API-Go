package common

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func (baseModel *BaseModel) SetCreatedAt() {
	baseModel.CreatedAt = time.Now()
}

func (baseModel *BaseModel) SetUpdatedAt() {
	baseModel.UpdatedAt = time.Now()
}
