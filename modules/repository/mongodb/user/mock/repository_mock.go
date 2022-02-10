package mock

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"learning-golang-hexagonal/business/model"
	"learning-golang-hexagonal/modules/repository/mongodb"
)

var (
	ErrIdNotFound = errors.New("id not found")
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) Insert(user *model.User) error {
	args := r.Called(user)
	if args.Get(0) == nil {
		return errors.New(mongodb.UnableToInsert)
	}
	user.ID = primitive.NewObjectID().Hex()
	return nil
}

func (r *RepositoryMock) GetById(ID string) (*model.User, error) {
	args := r.Called(ID)
	if args.Get(1) != nil {
		return nil, errors.New(mongodb.ErrIdNotFound)
	}
	if args.Get(0) == nil {
		return nil, ErrIdNotFound
	}
	return args.Get(0).(*model.User), nil
}

func (r *RepositoryMock) GetByEmail(email string) (*model.User, error) {
	args := r.Called(email)
	if args.Get(0) == nil {
		return nil, errors.New(mongodb.UnableToDecodeFindOneRes)
	}
	return args.Get(0).(*model.User), nil
}

func (r *RepositoryMock) Update(user *model.User) error {
	args := r.Called(user)
	if args.Get(0) == nil {
		return errors.New(mongodb.UnableToUpdate)
	}
	return nil
}

func (r *RepositoryMock) Delete(ID string) error {
	args := r.Called(ID)
	if args.Get(0) == nil {
		return errors.New(mongodb.UnableToDelete)
	}
	return nil
}

func (r *RepositoryMock) List() ([]model.UserView, error) {
	args := r.Called()
	if args.Get(0) == nil {
		return nil, errors.New(mongodb.UnableToDecodeFindOneRes)
	}

	var userViews []model.UserView
	userView := model.NewDummyUserView()
	userViews = append(userViews, *userView)

	return userViews, nil
}
