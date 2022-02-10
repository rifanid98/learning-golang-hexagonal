package model

import (
	"time"
)

type (
	User struct {
		ID        string    `json:"_id" bson:"_id,omitempty"`
		Email     string    `json:"email"`
		Password  string    `json:"password,omitempty"`
		Role      string    `json:"role"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	UserView struct {
		ID        string    `json:"_id" bson:"_id,omitempty"`
		Email     string    `json:"email"`
		Password  string    `json:"password,omitempty"`
		Role      string    `json:"role"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

func NewDummyUser() *User {
	return &User{
		ID:        "id",
		Email:     "email",
		Password:  "admin",
		Role:      "admin",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
}

func NewDummyUserView() *UserView {
	return &UserView{
		ID:        "id",
		Email:     "email",
		Password:  "admin",
		Role:      "admin",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
}
