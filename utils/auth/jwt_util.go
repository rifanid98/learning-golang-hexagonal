package auth

import port "learning-golang-hexagonal/business/port/intl/v1/auth"

type JWTUtil interface {
	GenerateToken(tokenID uint, authUser *port.AuthUser) (err error)
}

type UtilImpl struct {
	jwt JWT
}

func NewUtilImpl(jwt JWT) *UtilImpl {
	return &UtilImpl{jwt: jwt}
}

func (u *UtilImpl) GenerateToken(tokenID uint, authUser *port.AuthUser) (err error) {
	var tokenClaims JWTClaims
	tokenClaims.ID = tokenID
	tokenClaims.Email = authUser.Email
	tokenClaims.Role = authUser.Role

	authUser.Token, err = u.jwt.Create(tokenClaims)
	return
}
