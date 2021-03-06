package middleware

import (
	"github.com/labstack/echo/v4"
	authPort "learning-golang-hexagonal/business/port/intl/v1/auth"
	"strings"
)

type Auth struct {
	authService authPort.Service
}

func NewAuth(authService authPort.Service) *Auth {
	return &Auth{
		authService,
	}
}

func (handler *Auth) CustomParse(tokenString string, c echo.Context) (interface{}, error) {
	return handler.authService.Verify(tokenString)
}

func AuthAPISkipper(c echo.Context) bool {
	paths := []string{
		"/api/v1/auth/login",
	}

	requestURI := c.Request().RequestURI
	for i := range paths {
		if strings.Contains(requestURI, paths[i]) {
			return true
		}
	}
	indexAlwaysAllowed := requestURI == "/" || requestURI == ""

	return indexAlwaysAllowed
}
