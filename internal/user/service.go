package user

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	user, err := s.rep.CreateUser(ctx, CreateUserParams{Email: input.Email,
		PasswordHash: string(hashedByte),
		Username:     input.Username})

	if err != nil {
		return User{}, err
	}
	return user, nil
}

type SignInOutput struct {
	User      User
	Token     string
	ExpiresAt time.Time
}

type CustomClaims struct {
	UserID   int64  `json:"uid"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *userService) SignIn(ctx context.Context, input SignInRequest) (SignInOutput, error) {
	user, err := s.rep.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return SignInOutput{}, ErrInvalidCredentials
		}
		return SignInOutput{}, fmt.Errorf("service signin get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return SignInOutput{}, ErrInvalidCredentials
	}

	now := time.Now()
	exp := now.Add(15 * time.Minute)

	claims := CustomClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "devlog-auth-service",
			Subject:   strconv.FormatInt(user.ID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return SignInOutput{}, ErrJWTSecretNotSet
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))

	if err != nil {
		return SignInOutput{}, fmt.Errorf("service signin sign token: %w", err)
	}

	return SignInOutput{
		User:      user,
		Token:     signed,
		ExpiresAt: exp,
	}, nil
}
