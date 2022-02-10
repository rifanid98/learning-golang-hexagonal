package routes

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"learning-golang-hexagonal/utils/auth"
	"learning-golang-hexagonal/utils/config"

	m "learning-golang-hexagonal/api/intl/v1/routes/middleware"
	cryptoUtil "learning-golang-hexagonal/utils/crypto"
	helperUtil "learning-golang-hexagonal/utils/helper"

	adminUserCtrl "learning-golang-hexagonal/api/intl/v1/admin/user"
	authCtrl "learning-golang-hexagonal/api/intl/v1/auth"

	authService "learning-golang-hexagonal/business/services/intl/v1/auth"
	userService "learning-golang-hexagonal/business/services/intl/v1/user"

	userRepository "learning-golang-hexagonal/modules/repository/mongodb/user"
	userTokenRepository "learning-golang-hexagonal/modules/repository/mongodb/user_token"
)

func acl(permission string) echo.MiddlewareFunc {
	return m.ACL(permission)
}

func API(e *echo.Echo) {
	// Instance DB
	db := config.DB

	// instance dependency
	jwt := auth.NewJWT()
	jwtUtil := auth.NewUtilImpl(jwt)
	helper := helperUtil.New()
	crypto := cryptoUtil.NewCrypto()

	// instance repo
	userRepo := userRepository.New(db)
	userTokenRepo := userTokenRepository.New(db)

	// instance service
	userServ := userService.New(userRepo, userTokenRepo, crypto)
	authServ := authService.New(userTokenRepo, userRepo, jwt, jwtUtil, helper, crypto)

	// set jwt
	customMiddleware := m.NewAuth(authServ)
	JWTCustomConfig := middleware.JWTConfig{
		Skipper:        m.AuthAPISkipper,
		ParseTokenFunc: customMiddleware.CustomParse,
	}

	baseUrl := "/api/v1/extl"

	// auth route
	authHandler := authCtrl.New(authServ)
	authRoute := e.Group(fmt.Sprintf("%s/auth", baseUrl))
	authRoute.POST("/login", authHandler.Login)
	authRoute.POST("/logout", authHandler.Logout, middleware.JWTWithConfig(JWTCustomConfig))

	// admin routes
	// user
	adminUserHandler := adminUserCtrl.New(userServ)
	adminUserRoute := e.Group(fmt.Sprintf("%s/user", baseUrl), acl("admin"))
	adminUserRoute.Use(middleware.JWTWithConfig(JWTCustomConfig))
	adminUserRoute.POST("", adminUserHandler.Create)
	adminUserRoute.GET("/:id", adminUserHandler.Read)
	adminUserRoute.PUT("/:id", adminUserHandler.Update)
	adminUserRoute.DELETE("/:id", adminUserHandler.Delete)
	adminUserRoute.GET("", adminUserHandler.List)
}
