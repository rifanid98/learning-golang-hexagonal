package mock

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"messaging/business/model"
	"messaging/modules/repository/mongodb"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) Insert(userToken *model.UserToken) error {
	args := r.Called(userToken)
	if args.Get(0) == nil {
		return errors.New(mongodb.UnableToInsert)
	}
	userToken.ID = primitive.NewObjectID().Hex()
	return nil
}

func (r *RepositoryMock) GetByTokenID(tokenID uint) (*model.UserToken, error) {
	args := r.Called(tokenID)
	if args.Get(0) == nil {
		return nil, errors.New(mongodb.UnableToDecodeFindOneRes)
	}
	return args.Get(0).(*model.UserToken), nil
}

func (r *RepositoryMock) DeleteByUserID(userID string) error {
	args := r.Called(userID)
	if args.Get(0) == nil {
		return errors.New(mongodb.UnableToDelete)
	}
	return nil
}

func (r *RepositoryMock) DeleteByTokenID(tokenID uint) error {
	args := r.Called(tokenID)
	if args.Get(0) == nil {
		return errors.New(mongodb.UnableToDelete)
	}
	return nil
}
