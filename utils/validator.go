package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidateStruct(v interface{}) (validationErrors []string, ok bool) {
	validate := validator.New()
	if err := validate.Struct(v); err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			switch e.Tag() {
			case "required":
				validationErrors = append(validationErrors, fmt.Sprintf("%s is required", e.Field()))
			case "email":
				validationErrors = append(validationErrors, fmt.Sprintf("%s is not a valid email", e.Field()))
			default:
				validationErrors = append(validationErrors, fmt.Sprintf("%s is invalid", e.Field()))
			}
		}
		if len(validationErrors) > 0 {
			return validationErrors, false
		}
	}
	ok = true
	return
}
