package user

import (
	port "messaging/business/port/intl/v1/user"
	"messaging/utils/validator"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	userService port.Service
}

func New(userService port.Service) *Controller {
	return &Controller{
		userService,
	}
}

func (controller *Controller) List(c echo.Context) error {
	users, err := controller.userService.List()
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	datas := make([]*ResponseUserView, 0, len(users))
	for i := range users {
		user := populateResponseUserView(&users[i])
		datas = append(datas, user)
	}

	response := ResponseUsersView{
		User: datas,
	}

	return c.JSON(http.StatusOK, response)
}

func (controller *Controller) Create(c echo.Context) error {
	request := new(RequestUser)
	if err := c.Bind(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := validator.GetValidator().Struct(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user := new(port.User)
	user.Email = request.Email
	user.Password = request.Password
	user.Role = request.Role

	if err := controller.userService.Create(user); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(http.StatusCreated, populateResponseUser(user))
}

func (controller *Controller) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	request := new(RequestUserUpdate)
	if err := c.Bind(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := validator.GetValidator().Struct(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user := new(port.User)
	user.ID = id
	user.Email = request.Email
	user.Password = request.Password
	user.Role = request.Role

	if err := controller.userService.Update(user); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(http.StatusOK, populateResponseUser(user))
}

func (controller *Controller) Read(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	user, err := controller.userService.GetByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	if user.ID == "" {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return c.JSON(http.StatusOK, populateResponseUserView(&user))
}

func (controller *Controller) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	if err := controller.userService.Delete(id); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(http.StatusOK, "user deleted")
}

func populateResponseUserView(user *port.User) *ResponseUserView {
	response := new(ResponseUserView)
	response.ID = user.ID
	response.Email = user.Email
	response.Role = user.Role

	return response
}

func populateResponseUser(user *port.User) *ResponseUser {
	response := new(ResponseUser)

	response.ID = user.ID
	response.Email = user.Email
	response.Role = user.Role

	return response
}
