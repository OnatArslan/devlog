package user

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// CREATING SERVICE STRUCT --- --- --- --- --- ---
type userService struct {
	rep *userRepository
}

func NewUserService(rep *userRepository) *userService {
	return &userService{
		rep: rep,
	}
}

// CREATING METHODS --- --- --- --- --- ---

// SignUp Stuff --- --- ---
type SignUpInput struct {
	Email    string
	Username string
	Password string
}

func (s *userService) SignUp(ctx context.Context, input SignUpInput) (User, error) {
	passwordByte := []byte(input.Password)

	hashedByte, err := bcrypt.GenerateFromPassword(passwordByte, 12)
	if err != nil {
		return User{}, err
	}
	user, err := s.rep.CreateUser(ctx, CreateUserParams{Email: input.Email, PasswordHash: string(hashedByte), Username: input.Username})
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *userService) SignIn(ctx context.Context, email string) (User, error) {

	user, err := s.rep.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return User{}, ErrInvalidCredentials
		}
		return User{}, fmt.Errorf("service signin get user: %w", err)
	}
	return user, nil
}
