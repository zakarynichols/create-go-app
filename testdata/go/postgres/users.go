package postgres

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type userService struct {
	psql *psqlService
}

func NewUserService(psql *psqlService) *userService {
	return &userService{psql}
}

func (userService *userService) QueryUsers() ([]User, error) {
	a := User{
		ID: 1,
	}
	b := User{
		ID: 2,
	}
	c := User{
		ID: 3,
	}
	return []User{a, b, c}, nil
}
