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
// userService contains business rules for user registration and authentication flows.
type userService struct {
	rep *userRepository
}

// NewUserService wires the service with its repository dependency.
func NewUserService(rep *userRepository) *userService {
	// Return a service instance bound to the repository implementation.
	return &userService{
		rep: rep,
	}
}

// CREATING METHODS --- --- --- --- --- ---

// SignUp Stuff --- --- ---
// SignUpInput carries raw registration fields received from the handler layer.
type SignUpInput struct {
	Email    string
	Username string
	Password string
}

// SignUp hashes the password and creates a new user record.
func (s *userService) SignUp(ctx context.Context, input SignUpInput) (User, error) {
	// Convert the plain password to bytes for bcrypt processing.
	passwordByte := []byte(input.Password)

	// Hash the password before persisting any user record.
	hashedByte, err := bcrypt.GenerateFromPassword(passwordByte, 12)
	if err != nil {
		return User{}, err
	}
	// Persist the new user with the generated password hash.
	user, err := s.rep.CreateUser(ctx, CreateUserParams{Email: input.Email,
		PasswordHash: string(hashedByte),
		Username:     input.Username})

	if err != nil {
		return User{}, fmt.Errorf("signup service : %w", err)
	}
	return user, nil
}

// SignInOutput contains the authenticated user and generated access token metadata.
type SignInOutput struct {
	User      User
	Token     string
	ExpiresAt time.Time
}

// CustomClaims extends JWT registered claims with application-specific user identity fields.
type CustomClaims struct {
	UserID int64  `json:"uid"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type SignInInput struct {
	Email    string
	Password string
}

// SignIn validates credentials and returns a signed short-lived JWT access token.
func (s *userService) SignIn(ctx context.Context, input SignInInput) (SignInOutput, error) {
	// Fetch the active user by email for credential verification.
	user, err := s.rep.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return SignInOutput{}, ErrInvalidCredentials
		}
		return SignInOutput{}, fmt.Errorf("service signin get user: %w", err)
	}

	// Compare the stored password hash with the provided raw password.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return SignInOutput{}, ErrInvalidCredentials
	}

	// Define token issuance and expiration timestamps.
	now := time.Now()
	exp := now.Add(15 * time.Minute)

	// Build application and standard JWT claims for this session.
	claims := CustomClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "devlog-auth-service",
			Subject:   strconv.FormatInt(user.ID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	// Read signing secret from environment and fail fast when missing.
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return SignInOutput{}, ErrJWTSecretNotSet
	}

	// Create and sign the JWT token with HS256.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedTokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return SignInOutput{}, fmt.Errorf("service signin sign token: %w", err)
	}

	// Return authenticated user metadata together with token payload.
	return SignInOutput{
		User:      user,
		Token:     signedTokenString,
		ExpiresAt: exp,
	}, nil
}

func (s *userService) GetMe(ctx context.Context, email string) (User, error) {

	user, err := s.rep.GetByEmail(ctx, email)
	if err != nil {
		return User{}, ErrUserNotFound
	}

	return user, nil
}
