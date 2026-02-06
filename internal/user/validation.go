package user

import "github.com/go-playground/validator/v10"

// This function used in main.go for RegisterValidators for user in one line
// user.RegisterValidations(validate) like this
func RegisterValidations(v *validator.Validate) {
	v.RegisterStructValidation(SignUpStructLevelValidation, SignUpRequest{})
}

// This is for signUp handler
func SignUpStructLevelValidation(sl validator.StructLevel) {
	registerReq := sl.Current().Interface().(SignUpRequest)

	if registerReq.Password != registerReq.PasswordConfirm {
		sl.ReportError(registerReq.PasswordConfirm, "PasswordConfirm", "passwordConfirm", "eqfield", "Password")
	}
}
