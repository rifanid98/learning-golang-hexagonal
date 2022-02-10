package mock

import "github.com/stretchr/testify/mock"

type CryptoMock struct {
	mock.Mock
}

func (c *CryptoMock) UserGeneratePassword(password string) string {
	args := c.Called(password)
	if args.Get(0) == nil {
		return ""
	}
	return args.Get(0).(string)
}

func (c *CryptoMock) UserVerifyPassword(userPwd, userStoredPwd string) bool {
	args := c.Called(userPwd, userStoredPwd)
	if args.Get(0) == nil {
		return false
	}
	return args.Get(0).(bool)
}
