package user_test

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"messaging/modules/repository/mongodb"
	"messaging/utils/crypto"
	"messaging/utils/crypto/mock"

	drivenPort "messaging/business/model"
	driverPort "messaging/business/port/intl/v1/user"
	driverPortImpl "messaging/business/services/intl/v1/user"
	drivenPortImplMock1 "messaging/modules/repository/mongodb/user/mock"
	drivenPortImplMock2 "messaging/modules/repository/mongodb/user_token/mock"
)

var cryptoUtil = crypto.NewCrypto()

var _ = Describe("User", func() {
	var userService driverPort.Service
	var userRepo *drivenPortImplMock1.RepositoryMock
	var userTokenRepo *drivenPortImplMock2.RepositoryMock
	var drivenModel *drivenPort.User
	var driverModel *driverPort.User
	var cryptoMock *mock.CryptoMock

	BeforeEach(func() {
		userRepo = new(drivenPortImplMock1.RepositoryMock)
		userTokenRepo = new(drivenPortImplMock2.RepositoryMock)
		cryptoMock = new(mock.CryptoMock)
		userService = driverPortImpl.New(userRepo, userTokenRepo, cryptoMock)
		drivenModel = drivenPort.NewDummyUser()
		driverModel = driverPort.NewDummyUser()
	})

	Describe("Create", func() {
		var generatedPassword string

		BeforeEach(func() {
			generatedPassword = cryptoUtil.UserGeneratePassword("admin")
			driverModel.Password = generatedPassword
			drivenModel.ID = ""
			drivenModel.Password = generatedPassword
		})

		It("should create user successfully", func() {
			cryptoMock.On("UserGeneratePassword", generatedPassword).Return(generatedPassword)
			userRepo.On("Insert", drivenModel).Return(drivenModel)

			err := userService.Create(driverModel)
			Expect(err).To(BeNil())
		})

		It("should failed create user - failed to insert data", func() {
			cryptoMock.On("UserGeneratePassword", generatedPassword).Return(generatedPassword)
			userRepo.On("Insert", drivenModel).Return(nil)

			err := userService.Create(driverModel)
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(errors.New(mongodb.UnableToInsert)))
		})
	})

	Describe("GetByID", func() {
		It("should get user by id", func() {
			userRepo.On("GetById", driverModel.ID).Return(drivenModel, nil)

			user, err := userService.GetByID(driverModel.ID)
			Expect(err).To(BeNil())
			Expect(user).ToNot(BeNil())
		})

		It("should get user by id - failed get user by id", func() {
			userRepo.On("GetById", driverModel.ID).Return(nil, drivenPortImplMock1.ErrIdNotFound)

			_, err := userService.GetByID(driverModel.ID)
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(errors.New(fmt.Sprintf(driverPortImpl.ErrFailedGetUserString, driverModel.ID))))
		})
	})

	Describe("Update", func() {
		var generatedPassword string

		BeforeEach(func() {
			generatedPassword = cryptoUtil.UserGeneratePassword("admin")
			driverModel.Password = generatedPassword
			drivenModel.Password = generatedPassword
		})

		It("should update user", func() {
			userRepo.On("GetById", driverModel.ID).Return(drivenModel, nil)
			cryptoMock.On("UserGeneratePassword", generatedPassword).Return(generatedPassword)
			userRepo.On("Update", drivenModel).Return(drivenModel)
			userTokenRepo.On("DeleteByUserID", drivenModel.ID).Return(drivenModel)

			err := userService.Update(driverModel)
			Expect(err).To(BeNil())
		})

		It("should failed update user - failed parse data from db", func() {
			userRepo.On("GetById", driverModel.ID).Return(nil, drivenPortImplMock1.ErrIdNotFound)

			err := userService.Update(driverModel)
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(drivenPortImplMock1.ErrIdNotFound))
		})

		It("should failed update user - id not found", func() {
			userRepo.On("GetById", driverModel.ID).Return(drivenModel, drivenPortImplMock1.ErrIdNotFound)

			err := userService.Update(driverModel)
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(drivenPortImplMock1.ErrIdNotFound))
		})

		It("should failed update user - failed update", func() {
			userRepo.On("GetById", driverModel.ID).Return(drivenModel, nil)
			cryptoMock.On("UserGeneratePassword", generatedPassword).Return(generatedPassword)
			userRepo.On("Update", drivenModel).Return(nil)
			userTokenRepo.On("DeleteByUserID", drivenModel.ID).Return(drivenModel)

			err := userService.Update(driverModel)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal(mongodb.UnableToUpdate))
		})
	})

	Describe("Delete", func() {
		It("should delete user successfully", func() {
			userRepo.On("Delete", driverModel.ID).Return(drivenModel)

			err := userService.Delete(driverModel.ID)
			Expect(err).To(BeNil())
		})

		It("should failed delete user", func() {
			userRepo.On("Delete", driverModel.ID).Return(nil)

			err := userService.Delete(driverModel.ID)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal(mongodb.UnableToDelete))
		})
	})

	Describe("List", func() {
		It("should get list of user", func() {
			var users []driverPort.User
			users = append(users, *driverModel)

			userRepo.On("List").Return(users, nil)

			list, err := userService.List()
			Expect(err).To(BeNil())
			Expect(list).ToNot(BeNil())
			Expect(len(list)).To(Equal(1))
		})

		It("should get empty list of user", func() {
			var users []driverPort.User
			userRepo.On("List").Return(users, nil)

			list, err := userService.List()
			Expect(err).To(BeNil())
			Expect(list).ToNot(BeNil())
			Expect(len(list)).To(Equal(1))
		})

		It("should failed get list of user", func() {
			userRepo.On("List").Return(nil, errors.New(mongodb.UnableToDecodeFindOneRes))

			list, err := userService.List()
			Expect(err).ToNot(BeNil())
			Expect(list).ToNot(BeNil())
			Expect(len(list)).To(Equal(0))
			Expect(err.Error()).To(Equal(mongodb.UnableToDecodeFindOneRes))
		})
	})
})
