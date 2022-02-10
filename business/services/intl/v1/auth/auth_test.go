package auth_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"learning-golang-hexagonal/modules/repository/mongodb"
	"learning-golang-hexagonal/utils/auth/mock"
	"learning-golang-hexagonal/utils/helper"

	drivenPort "learning-golang-hexagonal/business/model"
	driverPort "learning-golang-hexagonal/business/port/intl/v1/auth"
	driverPortImpl "learning-golang-hexagonal/business/services/intl/v1/auth"
	drivenPortImplMock1 "learning-golang-hexagonal/modules/repository/mongodb/user/mock"
	drivenPortImplMock2 "learning-golang-hexagonal/modules/repository/mongodb/user_token/mock"
	authUtil "learning-golang-hexagonal/utils/auth"
	cryptoUtilMock "learning-golang-hexagonal/utils/crypto/mock"
	helperUtilMock "learning-golang-hexagonal/utils/helper/mock"
)

var _ = Describe("Auth", func() {
	var authService driverPort.Service
	var userRepo *drivenPortImplMock1.RepositoryMock
	var userTokenRepo *drivenPortImplMock2.RepositoryMock
	var jwtMock *mock.JWTMock
	var jwtUtilMock *mock.JWTUtilMock
	var helperMock *helperUtilMock.HelperMock
	var cryptoMock *cryptoUtilMock.CryptoMock

	var token string
	var tokenClaims *authUtil.JWTClaims
	var tokenID uint
	var userToken *drivenPort.UserToken
	var authUser *driverPort.AuthUser

	BeforeEach(func() {
		token = "token"
		tokenClaims = authUtil.NewDummyJWTClaims()
		tokenID = helper.New().GenerateRandomNumber()
		userToken = drivenPort.NewDummyUserToken()
		authUser = driverPort.NewDummyAuthUser()

		userRepo = new(drivenPortImplMock1.RepositoryMock)
		userTokenRepo = new(drivenPortImplMock2.RepositoryMock)
		jwtMock = new(mock.JWTMock)
		jwtUtilMock = new(mock.JWTUtilMock)
		helperMock = new(helperUtilMock.HelperMock)
		cryptoMock = new(cryptoUtilMock.CryptoMock)
		authService = driverPortImpl.New(userTokenRepo, userRepo, jwtMock, jwtUtilMock, helperMock, cryptoMock)
	})

	Describe("Verify", func() {
		It("should verify jwt token", func() {
			jwtMock.On("Parse", token).Return(token, tokenClaims, token)
			userTokenRepo.On("GetByTokenID", tokenClaims.ID).Return(userToken, nil)

			verify, err := authService.Verify(token)
			Expect(err).To(BeNil())
			Expect(verify).To(Equal(token))
		})

		It("should failed verify jwt token - failed parse token", func() {
			jwtMock.On("Parse", token).Return(nil, tokenClaims, token)

			verify, err := authService.Verify(token)
			Expect(err).ToNot(BeNil())
			Expect(verify).To(BeNil())
		})
	})

	Describe("Login", Ordered, func() {
		var currentUser *drivenPort.User

		BeforeEach(func() {
			// preparing data
			currentUser = drivenPort.NewDummyUser()
			currentUser.ID = authUser.UserID
			currentUser.Password = "$2a$08$cwJLLr.LfnKjUdpW6C3kE.KnEhdcQiVXGXcMd3iAXJ9IgMgLDnOci" // admin

			// preparing data
			userToken.UserID = authUser.UserID
			userToken.TokenID = tokenID
			userToken.ID = ""
		})

		It("should signing in user", func() {
			cryptoMock.On("UserVerifyPassword", "admin", currentUser.Password).Return(true)
			userRepo.On("GetByEmail", currentUser.Email).Return(currentUser, nil)
			helperMock.On("GenerateRandomNumber").Return(tokenID)
			jwtUtilMock.On("GenerateToken", tokenID, authUser).Return(token)
			userTokenRepo.On("Insert", userToken).Return(userToken)

			err := authService.Login(authUser)
			Expect(err).To(BeNil())
		})

		It("should failed signing in user - failed generate token", func() {
			cryptoMock.On("UserVerifyPassword", "admin", currentUser.Password).Return(true)
			userRepo.On("GetByEmail", currentUser.Email).Return(currentUser, nil)
			helperMock.On("GenerateRandomNumber").Return(tokenID)
			jwtUtilMock.On("GenerateToken", tokenID, authUser).Return(nil)

			err := authService.Login(authUser)
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(driverPortImpl.ErrGenerateToken))
		})
	})

	Describe("Logout", func() {
		It("should logged out user successfully", func() {
			userTokenRepo.On("DeleteByTokenID", tokenID).Return(tokenID)

			err := authService.Logout(tokenID)
			Expect(err).To(BeNil())
		})

		It("should failed logged out user", func() {
			userTokenRepo.On("DeleteByTokenID", tokenID).Return(nil)

			err := authService.Logout(tokenID)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal(mongodb.UnableToDelete))
		})
	})
})
