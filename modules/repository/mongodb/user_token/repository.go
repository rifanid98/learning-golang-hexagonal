package usertoken

import (
	"context"
	"errors"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"learning-golang-hexagonal/business/model"
	"time"
)

type (
	Repository struct {
		*mongo.Collection
	}

	UserToken struct {
		ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
		TokenID   uint               `json:"token_id" bson:"token_id,omitempty"`
		UserID    primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
		CreatedAt time.Time          `json:"created_at" bson:"created_at"`
		UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	}
)

func New(database *mongo.Database) *Repository {
	return &Repository{
		Collection: database.Collection("user_token"),
	}
}

func (repo *Repository) Insert(userToken *model.UserToken) error {
	userId, err := primitive.ObjectIDFromHex(userToken.UserID)
	if err != nil {
		log.Errorf("unable to convert userId to ObjectId(user_id): %v", err)
		return errors.New("unable to convert userId to ObjectId(user_id)")
	}

	payload := new(UserToken)
	payload.TokenID = userToken.TokenID
	payload.UserID = userId

	_, err = repo.Collection.InsertOne(context.Background(), payload)
	if err != nil {
		log.Errorf("unable to insert: %v", err)
		return errors.New("unable to insert")
	}

	return nil
}

func (repo *Repository) GetByTokenID(tokenID uint) (*model.UserToken, error) {
	var userToken model.UserToken

	var result UserToken

	filter := bson.M{"token_id": tokenID}
	res := repo.Collection.FindOne(context.Background(), filter)
	if err := res.Decode(&result); err != nil {
		log.Errorf("unable to decode FindOne res: %v", err)
		return nil, errors.New("unable to decode FindOne res")
	}

	userToken.ID = result.ID.Hex()
	userToken.TokenID = result.TokenID
	userToken.UserID = result.UserID.Hex()
	userToken.CreatedAt = result.CreatedAt
	userToken.UpdatedAt = result.UpdatedAt

	return &userToken, nil
}

func (repo *Repository) DeleteByUserID(userID string) error {
	userId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Errorf("unable to convert id to user_id: %v", err)
		return errors.New("unable to convert id to user_id")
	}

	res, err := repo.Collection.DeleteOne(context.Background(), bson.M{"user_id": userId})
	if err != nil {
		log.Errorf("unable to delete data: %v", err)
		return errors.New("unable to delete data")
	}

	if res.DeletedCount < 1 {
		return errors.New("user token data not deleted")
	}

	return nil
}

func (repo *Repository) DeleteByTokenID(tokenID uint) error {
	res, err := repo.Collection.DeleteOne(context.Background(), bson.M{"token_id": tokenID})
	if err != nil {
		log.Errorf("unable to delete data: %v", err)
		return errors.New("unable to delete data")
	}

	if res.DeletedCount < 1 {
		return errors.New("user data not deleteds")
	}

	return nil
}
