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
			case "gte":
				validationErrors = append(validationErrors, fmt.Sprintf("%s length must be greater than or equal to %s", e.Field(), e.Param()))
			case "gt":
				validationErrors = append(validationErrors, fmt.Sprintf("%s length must be greater than %s", e.Field(), e.Param()))
			case "lte":
				validationErrors = append(validationErrors, fmt.Sprintf("%s length must be less than or equal to %s", e.Field(), e.Param()))
			case "lt":
				validationErrors = append(validationErrors, fmt.Sprintf("%s length must be less than %s", e.Field(), e.Param()))
			case "max":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be less than or equal to %s", e.Field(), e.Param()))
			case "min":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be greater than or equal to %s", e.Field(), e.Param()))
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
