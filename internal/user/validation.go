package user

import "github.com/go-playground/validator/v10"

// RegisterValidations attaches user-specific struct-level validation rules to the validator instance.
func RegisterValidations(v *validator.Validate) {
	// Register struct-level validation hook for signup cross-field checks.
	v.RegisterStructValidation(SignUpStructLevelValidation, SignUpRequest{})
}

// SignUpStructLevelValidation enforces cross-field signup rules that tags alone cannot express.
func SignUpStructLevelValidation(sl validator.StructLevel) {
	// Cast current payload to signup request for field comparison.
	registerReq := sl.Current().Interface().(SignUpRequest)

	// Reject requests where confirmation password does not match the main password.
	if registerReq.Password != registerReq.PasswordConfirm {
		sl.ReportError(registerReq.PasswordConfirm, "PasswordConfirm", "passwordConfirm", "eqfield", "Password")
	}
}
