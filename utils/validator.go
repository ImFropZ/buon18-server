package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"server/models"

	"github.com/go-playground/validator/v10"
)

func ValidateStruct(v interface{}) (validationErrors []string, ok bool) {
	validate := validator.New()

	validate.RegisterValidation("json", func(fl validator.FieldLevel) bool {
		if err := json.Unmarshal([]byte(fl.Field().String()), &map[string]interface{}{}); err != nil {
			return false
		}
		return true
	})

	validate.RegisterValidation("gender", func(fl validator.FieldLevel) bool {
		gender := fl.Field().String()
		for _, validGender := range models.VALID_GENDER_TYPES {
			if gender == validGender {
				return true
			}
		}
		return false
	})

	validate.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		// Define a regex pattern for phone numbers
		phoneRegex := `^\+?[0-9]{10,15}$` // Example pattern: allows international numbers starting with + and 10-15 digits
		re := regexp.MustCompile(phoneRegex)
		return re.MatchString(fl.Field().String())
	})

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
			case "gender":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be one of %s", e.Field(), models.VALID_GENDER_TYPES))
			case "phone":
				validationErrors = append(validationErrors, fmt.Sprintf("%s is not a valid phone number", e.Field()))
			case "json":
				validationErrors = append(validationErrors, fmt.Sprintf("%s is not a valid json string", e.Field()))
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
