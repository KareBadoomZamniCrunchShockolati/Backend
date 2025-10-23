package validation

import (
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/gin-gonic/gin/binding"
)


func RegisterGinValidator() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		return v.RegisterValidation("password_policy", PasswordValidationFunc)
	}
	return nil
}

// PasswordValidationFunc implements the go-playground validator signature.
func PasswordValidationFunc(fl validator.FieldLevel) bool {
	pwd, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	if pwd == "" {
		return false 
	}
	ok2, _ := passwordChecks(pwd)
	return ok2
}

// passwordChecks returns whether the password is valid and a slice of human-readable reasons for failures.
func passwordChecks(pwd string) (bool, []string) {
	var reasons []string
	if len(pwd) < 8 {
		reasons = append(reasons, "must be at least 8 characters long")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	specialChars := "!@#$%^&*()-_=+[]{}|;:,.<>?/"
	for _, r := range pwd {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case strings.ContainsRune(specialChars, r):
			hasSpecial = true
		}
	}
	if !hasUpper {
		reasons = append(reasons, "must contain at least one uppercase letter")
	}
	if !hasLower {
		reasons = append(reasons, "must contain at least one lowercase letter")
	}
	if !hasDigit {
		reasons = append(reasons, "must contain at least one digit")
	}
	if !hasSpecial {
		reasons = append(reasons, "must contain at least one special character")
	}

	return len(reasons) == 0, reasons
}

// FormatValidationError converts a binding/validator error into a map of field->list of problems.
func FormatValidationError(err error) map[string][]string {
	out := make(map[string][]string)
	if err == nil {
		return nil
	}

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			field := fe.Field()
			tag := fe.Tag()
			
			switch tag {
			case "required":
				out[field] = []string{field + " is required"}
			case "password_policy":
				if val, ok := fe.Value().(string); ok {
					_, reasons := passwordChecks(val)
					if len(reasons) == 0 {
						out[field] = []string{"invalid password"}
					} else {
						out[field] = reasons
					}
				} else {
					out[field] = []string{"invalid password"}
				}
			default:
				out[field] = append(out[field], fe.Error())
			}
		}
		return out
	}

	out["_error"] = []string{err.Error()}
	return out
}
