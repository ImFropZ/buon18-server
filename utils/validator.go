package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"system.buon18.com/m/models"
)

func ValidateAtLeastOneField(sl validator.StructLevel) bool {
	value := sl.Current().Interface()
	v := reflect.ValueOf(value)
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if field.Name == "ID" {
			continue
		}
		if !v.Field(i).IsZero() {
			return true
		}
	}
	return false
}

func ValidateStruct(v interface{}, isPatch bool) (validationErrors []string, ok bool) {
	validate := validator.New()

	if isPatch {
		validate.RegisterStructValidation(func(sl validator.StructLevel) {
			ms := sl.Current().Interface()
			if !ValidateAtLeastOneField(sl) {
				sl.ReportError(ms, "", "", "atleastone", "")
			}
		}, v)
	}

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

	validate.RegisterValidation("sales_quotation_status", func(fl validator.FieldLevel) bool {
		status := fl.Field().String()
		for _, validStatus := range models.VALID_SALES_QUOTATION_STATUS {
			if status == validStatus {
				return true
			}
		}
		return false
	})

	validate.RegisterValidation("accounting_account_typ", func(fl validator.FieldLevel) bool {
		typ := fl.Field().String()
		for _, validTyp := range models.VALID_ACCOUNTING_ACCOUNT_TYPES {
			if typ == validTyp {
				return true
			}
		}
		return false
	})

	validate.RegisterValidation("accounting_journal_typ", func(fl validator.FieldLevel) bool {
		typ := fl.Field().String()
		for _, validTyp := range models.VALID_ACCOUNTING_JOURNAL_TYPES {
			if typ == validTyp {
				return true
			}
		}
		return false
	})

	validate.RegisterValidation("accounting_journal_entry_typ", func(fl validator.FieldLevel) bool {
		typ := fl.Field().String()
		for _, validTyp := range models.VALID_ACCOUNTING_JOURNAL_ENTRY_TYPES {
			if typ == validTyp {
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
				validationErrors = append(validationErrors, fmt.Sprintf("%s is required", jsonFieldName(e.Namespace())))
			case "email":
				validationErrors = append(validationErrors, fmt.Sprintf("%s is not a valid email", jsonFieldName(e.Namespace())))
			case "gte":
				validationErrors = append(validationErrors, fmt.Sprintf("%s length must be greater than or equal to %s", jsonFieldName(e.Namespace()), e.Param()))
			case "gt":
				validationErrors = append(validationErrors, fmt.Sprintf("%s length must be greater than %s", jsonFieldName(e.Namespace()), e.Param()))
			case "lte":
				validationErrors = append(validationErrors, fmt.Sprintf("%s length must be less than or equal to %s", jsonFieldName(e.Namespace()), e.Param()))
			case "lt":
				validationErrors = append(validationErrors, fmt.Sprintf("%s length must be less than %s", jsonFieldName(e.Namespace()), e.Param()))
			case "max":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be less than or equal to %s", jsonFieldName(e.Namespace()), e.Param()))
			case "min":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be greater than or equal to %s", jsonFieldName(e.Namespace()), e.Param()))
			case "atleastone":
				validationErrors = append(validationErrors, "Need at least one field to perform update beside 'id' field")
			case "gender":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be one of %s", jsonFieldName(e.Namespace()), models.VALID_GENDER_TYPES))
			case "sales_quotation_status":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be one of %s", jsonFieldName(e.Namespace()), models.VALID_SALES_QUOTATION_STATUS))
			case "accounting_account_typ":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be one of %s", jsonFieldName(e.Namespace()), models.VALID_ACCOUNTING_ACCOUNT_TYPES))
			case "accounting_journal_typ":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be one of %s", jsonFieldName(e.Namespace()), models.VALID_ACCOUNTING_JOURNAL_TYPES))
			case "accounting_journal_entry_typ":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be one of %s", jsonFieldName(e.Namespace()), models.VALID_ACCOUNTING_JOURNAL_ENTRY_TYPES))
			case "phone":
				validationErrors = append(validationErrors, fmt.Sprintf("%s is not a valid phone number", jsonFieldName(e.Namespace())))
			case "json":
				validationErrors = append(validationErrors, fmt.Sprintf("%s is not a valid json string", jsonFieldName(e.Namespace())))
			default:
				validationErrors = append(validationErrors, fmt.Sprintf("%s is invalid", jsonFieldName(e.Namespace())))
			}
		}
		if len(validationErrors) > 0 {
			return validationErrors, false
		}
	}
	ok = true
	return
}

func jsonFieldName(str string) string {
	return strings.Join(strings.Split(str, ".")[1:], ".")
}
