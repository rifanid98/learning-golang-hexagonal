package user

type (
	// User data
	User struct {
		ID       string
		Email    string
		Password string
		Role     string
	}

	// Service is inbount port
	Service interface {
		// Create insert new data
		Create(user *User) error

		// GetByID get data by ID
		GetByID(ID string) (User, error)

		// Update update new data
		Update(user *User) error

		// Delete delete data
		Delete(ID string) error

		// List list data
		List() ([]User, error)
	}
)

func NewDummyUser() *User {
	return &User{
		ID:       "id",
		Email:    "email",
		Password: "admin",
		Role:     "admin",
	}
}
