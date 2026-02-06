package user

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

// We ll use rep struct in here

type userService struct {
	rep *userRepository
}

func NewUserService(rep *userRepository) *userService {
	return &userService{
		rep: rep,
	}
}

func (s *userService) SignUpUser(ctx context.Context, input SignUpInput) (User, error) {
	passwordByte := []byte(input.Password)

	hashedByte, err := bcrypt.GenerateFromPassword(passwordByte, 12)
	if err != nil {
		return User{}, err
	}
	user, err := s.rep.CreateUser(ctx, CreateUserInput{Email: input.Email, PasswordHash: string(hashedByte), Username: input.Username})
	if err != nil {
		return User{}, err
	}
	return user, nil
}
