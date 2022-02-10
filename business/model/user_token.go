package model

import (
	"time"
)

type (
	UserToken struct {
		ID        string    `json:"_id" bson:"_id,omitempty"`
		TokenID   uint      `json:"token_id" bson:"token_id,omitempty"`
		UserID    string    `json:"user_id" bson:"user_id,omitempty"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

func NewDummyUserToken() *UserToken {
	return &UserToken{
		ID:        "id",
		TokenID:   0,
		UserID:    "userId",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
}
