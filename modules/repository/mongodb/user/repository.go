package user

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

	User struct {
		ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
		Email     string             `json:"email"`
		Password  string             `json:"password,omitempty"`
		Role      string             `json:"role"`
		CreatedAt time.Time          `json:"created_at"`
		UpdatedAt time.Time          `json:"updated_at"`
	}

	UserView struct {
		ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
		Email     string             `json:"email"`
		Password  string             `json:"password,omitempty"`
		Role      string             `json:"role"`
		CreatedAt time.Time          `json:"created_at"`
		UpdatedAt time.Time          `json:"updated_at"`
	}
)

func New(database *mongo.Database) *Repository {
	return &Repository{
		Collection: database.Collection("user"),
	}
}

func (repo *Repository) Insert(user *model.User) error {
	res, err := repo.Collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Errorf("unable to insert: %v", err)
		return errors.New("unable to insert")
	}

	if res.InsertedID == "" {
		log.Errorf("no data was inserted")
		return errors.New("no data was inserted")
	}

	user.ID = res.InsertedID.(primitive.ObjectID).Hex()

	return nil
}

func (repo *Repository) GetById(ID string) (*model.User, error) {
	var user model.User

	_id, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		log.Errorf("unable to convert id to _id: %v", err)
		return nil, errors.New("unable to convert id to _id")
	}

	var result User

	filter := bson.M{"_id": _id}
	res := repo.Collection.FindOne(context.Background(), filter)
	if err := res.Decode(&result); err != nil {
		log.Errorf("unable to decode FindOne res: %v", err)
		return nil, errors.New("unable to decode FindOne res")
	}

	user.ID = result.ID.Hex()
	user.Email = result.Email
	user.Password = result.Password
	user.Role = result.Role
	user.CreatedAt = result.CreatedAt

	return &user, nil
}

func (repo *Repository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	var result User

	filter := bson.M{"email": email}
	res := repo.Collection.FindOne(context.Background(), filter)
	if err := res.Decode(&result); err != nil {
		log.Errorf("unable to decode FindOne res: %v", err)
		return nil, errors.New("unable to decode FindOne res")
	}

	user.ID = result.ID.Hex()
	user.Email = result.Email
	user.Password = result.Password
	user.Role = result.Role
	user.CreatedAt = result.CreatedAt

	return &user, nil
}

func (repo *Repository) Update(user *model.User) error {
	// 1. find setting
	_id, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		log.Errorf("unable to convert id to _id: %v", err)
		return errors.New("unable to convert id to _id")
	}

	filter := bson.M{"_id": _id}
	res := repo.Collection.FindOne(context.Background(), filter)
	if res.Err() != nil {
		log.Errorf("unable to decode FindOne res: %v", res.Err())
		return errors.New("unable to find user data")
	}

	// 2. update data
	payload := new(User)
	payload.ID = _id
	payload.Email = user.Email
	payload.Role = user.Role
	payload.Password = user.Password

	result, err := repo.Collection.UpdateOne(context.Background(), filter, bson.M{"$set": payload})
	if err != nil {
		log.Errorf("unable to update: %v", err)
		return errors.New("unable to update")
	}

	if result.ModifiedCount < 1 {
		log.Errorf("user data not updated")
		return errors.New("user data not updated")
	}

	return nil
}

func (repo *Repository) Delete(ID string) error {
	_id, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		log.Errorf("unable to convert id to _id: %v", err)
		return errors.New("unable to convert id to _id")
	}

	res, err := repo.Collection.DeleteOne(context.Background(), bson.M{"_id": _id})
	if err != nil {
		log.Errorf("unable to delete data: %v", err)
		return errors.New("unable to delete data")
	}

	if res.DeletedCount < 1 {
		return errors.New("user data not deleted")
	}

	return nil
}

func (repo *Repository) List() ([]model.UserView, error) {
	var result []UserView

	cursor, err := repo.Collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Errorf("unable to find users : %v", err)
		return nil, errors.New("unable to find users")
	}

	err = cursor.All(context.Background(), &result)
	if err != nil {
		log.Errorf("unable to read the cursor : %v", err)
		return nil, errors.New("unable to read the cursor")
	}

	var users []model.UserView

	for _, item := range result {
		var user model.UserView
		user.ID = item.ID.Hex()
		user.Email = item.Email
		user.Password = item.Password
		user.Role = item.Role
		user.CreatedAt = item.CreatedAt

		users = append(users, user)
	}

	return users, nil
}
