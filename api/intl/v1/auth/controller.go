package auth

import (
	"github.com/golang-jwt/jwt/v4"
	port "messaging/business/port/intl/v1/auth"
	"messaging/utils/auth"
	"net/http"

	"github.com/labstack/echo/v4"
	"messaging/utils/validator"
)

type Controller struct {
	authService port.Service
}

func New(authService port.Service) *Controller {
	return &Controller{
		authService,
	}
}

func (controller *Controller) Login(c echo.Context) error {
	request := new(RequestAuth)
	if err := c.Bind(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := validator.GetValidator().Struct(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	authUser := new(port.AuthUser)
	authUser.Email = request.Email
	authUser.Password = request.Password

	if err := controller.authService.Login(authUser); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(http.StatusOK, ResponseAuth{Token: authUser.Token})
}

func (controller *Controller) Logout(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	tokenClaims := user.Claims.(*auth.JWTClaims)

	if err := controller.authService.Logout(tokenClaims.ID); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(http.StatusOK, "user token deleted")
}
