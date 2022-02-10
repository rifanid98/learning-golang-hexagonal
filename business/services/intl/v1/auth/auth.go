package auth

import (
	"errors"
	port "messaging/business/port/intl/v1/auth"
	userPort "messaging/business/port/intl/v1/user"
	userTokenPort "messaging/business/port/intl/v1/user_token"

	"messaging/business/model"

	"messaging/utils/auth"
	"messaging/utils/crypto"
	"messaging/utils/helper"
)

type (
	service struct {
		userTokenRepository userTokenPort.Repository
		userRepository      userPort.Repository
		jwt                 auth.JWT
		jwtUtil             auth.JWTUtil
		helper              helper.Helper
		crypto              crypto.Crypto
	}
)

func New(
	userTokenRepository userTokenPort.Repository,
	userRepository userPort.Repository,
	jwt auth.JWT,
	jwtUtil auth.JWTUtil,
	helper helper.Helper,
	crypto crypto.Crypto,
) port.Service {
	return &service{
		userTokenRepository,
		userRepository,
		jwt,
		jwtUtil,
		helper,
		crypto,
	}
}

var (
	ErrGenerateToken     error = errors.New("generate token failed")
	ErrInvalidToken      error = errors.New("invalid token")
	ErrInvalidCredential error = errors.New("invalid credential")
)

func (s *service) Verify(tokenString string) (interface{}, error) {
	token, tokenClaims, err := s.jwt.Parse(tokenString)
	if err != nil {
		return nil, err
	}
	_, err = s.userTokenRepository.GetByTokenID(tokenClaims.ID)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return token, nil
}

func (s *service) Login(authUser *port.AuthUser) error {
	if !s.bindUserCredential(authUser) {
		return ErrInvalidCredential
	}

	tokenID := s.helper.GenerateRandomNumber()
	if err := s.jwtUtil.GenerateToken(tokenID, authUser); err == nil {
		userToken := new(model.UserToken)
		userToken.TokenID = tokenID
		userToken.UserID = authUser.UserID

		return s.userTokenRepository.Insert(userToken)
	}

	return ErrGenerateToken
}

func (s *service) Logout(tokenID uint) error {
	return s.userTokenRepository.DeleteByTokenID(tokenID)
}

func (s *service) bindUserCredential(authUser *port.AuthUser) bool {
	existingUser, _ := s.userRepository.GetByEmail(authUser.Email)
	if existingUser != nil {
		authUser.UserID = existingUser.ID
		authUser.Role = existingUser.Role
		return s.crypto.UserVerifyPassword(authUser.Password, existingUser.Password)
	}

	return false
}
