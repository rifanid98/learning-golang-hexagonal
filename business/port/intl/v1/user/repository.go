package user

import "messaging/business/model"

// Repository is outbound port
type Repository interface {
	// Insert insert new data
	Insert(user *model.User) error

	// Read get data by ID
	GetById(ID string) (*model.User, error)

	// GetByEmail find with email
	GetByEmail(email string) (*model.User, error)

	// Update update new data
	Update(user *model.User) error

	// Delete delete data
	Delete(ID string) error

	// List list data
	List() ([]model.UserView, error)
}
