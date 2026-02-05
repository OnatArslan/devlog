package user

import "github.com/go-playground/validator/v10"

func RegisterValidations(v *validator.Validate) {
	v.RegisterStructValidation(RegisterStructLevelValidation, RegisterRequest{})
}

func RegisterStructLevelValidation(sl validator.StructLevel) {
	registerReq := sl.Current().Interface().(RegisterRequest)

	if registerReq.Password != registerReq.PasswordConfirm {
		sl.ReportError(registerReq.PasswordConfirm, "PasswordConfirm", "passwordConfirm", "eqfield", "Password")
	}
}
