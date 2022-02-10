package user

import (
	"errors"
	"fmt"
	port "messaging/business/port/intl/v1/user"
	portUserToken "messaging/business/port/intl/v1/user_token"

	"messaging/business/model"
	"messaging/utils/crypto"
)

type (
	service struct {
		userRepo      port.Repository
		userTokenRepo portUserToken.Repository
		crypto        crypto.Crypto
	}
)

var (
	ErrFailedGetUserString = "cannot get user with id %v"
)

func New(
	userRepo port.Repository,
	userTokenRepo portUserToken.Repository,
	crypto crypto.Crypto) port.Service {
	return &service{
		userRepo,
		userTokenRepo,
		crypto,
	}
}

func (s *service) Create(user *port.User) error {
	data := new(model.User)
	data.Email = user.Email
	data.Role = user.Role
	data.Password = s.crypto.UserGeneratePassword(user.Password)

	if err := s.userRepo.Insert(data); err != nil {
		return err
	}
	user.ID = data.ID

	return nil
}

func (s *service) GetByID(ID string) (port.User, error) {
	var data port.User

	user, err := s.userRepo.GetById(ID)
	if err != nil {
		return data, errors.New(fmt.Sprintf(ErrFailedGetUserString, ID))
	}
	data.ID = user.ID
	data.Email = user.Email
	data.Password = user.Password

	return data, err
}

func (s *service) Update(user *port.User) error {
	existingUser, err := s.userRepo.GetById(user.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("id not found")
	}

	var data model.User
	data.ID = user.ID
	data.Email = user.Email
	data.Role = user.Role
	data.Password = existingUser.Password
	data.CreatedAt = existingUser.CreatedAt

	revokeToken := false

	if user.Password != "" {
		data.Password = s.crypto.UserGeneratePassword(user.Password)
		revokeToken = true
	}

	err = s.userRepo.Update(&data)
	if err == nil && revokeToken {
		err := s.userTokenRepo.DeleteByUserID(user.ID)
		if err != nil {
			return err
		}
	}

	return err
}

func (s *service) Delete(ID string) error {
	return s.userRepo.Delete(ID)
}

func (s *service) List() ([]port.User, error) {
	datas := make([]port.User, 0)

	users, err := s.userRepo.List()
	if err != nil {
		return datas, err
	}

	var user port.User
	for i := range users {
		user.ID = users[i].ID
		user.Email = users[i].Email
		user.Role = users[i].Role

		datas = append(datas, user)
	}

	return datas, nil
}
