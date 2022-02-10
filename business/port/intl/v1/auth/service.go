package auth

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	// Product data
	AuthUser struct {
		UserID   string
		Email    string
		Role     string
		Password string
		Token    string
	}

	// Service is inbount port
	Service interface {
		// Verify token
		Verify(tokenString string) (interface{}, error)

		// Login user
		Login(authUser *AuthUser) error

		// Logout user
		Logout(tokenID uint) error
	}
)

func NewDummyAuthUser() *AuthUser {
	return &AuthUser{
		UserID:   primitive.NewObjectID().Hex(),
		Email:    "email",
		Role:     "admin",
		Password: "admin",
		Token:    "token",
	}
}
