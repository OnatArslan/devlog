package user

// We ll use rep struct in here
type UserRepository interface{}

type userService struct {
	rep *userRepository
}

func NewUserService(rep *userRepository) *userService {
	return &userService{
		rep: rep,
	}
}
