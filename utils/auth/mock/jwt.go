package mock

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"messaging/utils/auth"
)

var (
	ErrFailedParseJWT  = errors.New("failed parse jwt")
	ErrFailedCreateJWT = errors.New("failed create jwt")
)

type JWTMock struct {
	mock.Mock
}

func (j *JWTMock) Create(tokenClaims auth.JWTClaims) (string, error) {
	args := j.Called(tokenClaims)
	if args.Get(0) == nil {
		return "", ErrFailedCreateJWT
	}
	return "jwt_token", nil
}

func (j *JWTMock) Parse(tokenString string) (interface{}, *auth.JWTClaims, error) {
	args := j.Called(tokenString)
	if args.Get(0) == nil {
		return "", nil, ErrFailedParseJWT
	}
	return args.Get(0), args.Get(1).(*auth.JWTClaims), nil
}
