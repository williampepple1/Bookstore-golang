package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// validateStruct validates a struct using the validator package
func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			tag := err.Tag()
			param := err.Param()
			
			var message string
			switch tag {
			case "required":
				message = fmt.Sprintf("%s is required", field)
			case "min":
				message = fmt.Sprintf("%s must be at least %s characters long", field, param)
			case "max":
				message = fmt.Sprintf("%s must be at most %s characters long", field, param)
			case "email":
				message = fmt.Sprintf("%s must be a valid email address", field)
			case "len":
				message = fmt.Sprintf("%s must be exactly %s characters long", field, param)
			case "uuid":
				message = fmt.Sprintf("%s must be a valid UUID", field)
			default:
				message = fmt.Sprintf("%s is invalid", field)
			}
			validationErrors = append(validationErrors, message)
		}
		return fmt.Errorf(strings.Join(validationErrors, "; "))
	}
	return nil
}

// validateStruct is a wrapper for the exported ValidateStruct function
func validateStruct(s interface{}) error {
	return ValidateStruct(s)
}
