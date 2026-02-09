package user

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// RegisterValidations attaches user-specific struct-level validation rules to the validator instance.
func RegisterValidations(v *validator.Validate) {
	// Register struct-level validation hook for signup cross-field checks.
	v.RegisterStructValidation(SignUpStructLevelValidation, SignUpRequest{})

	v.RegisterStructValidation(SignInStructLevelValidation, SignInRequest{})
}

// SignUpStructLevelValidation enforces cross-field signup rules that tags alone cannot express.
func SignUpStructLevelValidation(sl validator.StructLevel) {
	// Cast current payload to signup request for field comparison.
	registerReq, _ := sl.Current().Interface().(SignUpRequest)

	// Reject requests where confirmation password does not match the main password.
	if registerReq.Password != registerReq.PasswordConfirm {
		sl.ReportError(registerReq.PasswordConfirm, "PasswordConfirm", "passwordConfirm", "eqfield", "Password")
	}
}

// This is for example actually
func SignInStructLevelValidation(sl validator.StructLevel) {
	signInReq, _ := sl.Current().Interface().(SignInRequest)

	if strings.HasPrefix(signInReq.Email, "admin") {
		sl.ReportError(signInReq.Email, "Email", "email", "forbidden_str", "Email")
	}

}
