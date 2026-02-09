package validatorx

import (
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// New builds and configures the shared validator instance.
func New() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	// Register all custom validators in here
	v.RegisterValidation("strong-password", strongPassword)
	return v
}

func strongPassword(fl validator.FieldLevel) bool {
	s := fl.Field().String()

	// policy: min 8, en az 1 buyuk, 1 kucuk, 1 rakam, 1 ozel karakter, bosluk yok
	if len(s) < 8 || strings.ContainsAny(s, " \t\n\r") {
		return false
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, r := range s {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}
