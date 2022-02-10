package mock

import (
	"errors"
	"github.com/stretchr/testify/mock"
	port "messaging/business/port/intl/v1/auth"
)

type JWTUtilMock struct {
	mock.Mock
}

func (j *JWTUtilMock) GenerateToken(tokenID uint, authUser *port.AuthUser) (err error) {
	args := j.Called(tokenID, authUser)
	if args.Get(0) == nil {
		return errors.New("token not generated")
	}
	return nil
}
