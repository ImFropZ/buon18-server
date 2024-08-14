package setting

import (
	"encoding/json"
	"server/models"
	"strings"

	"github.com/nullism/bqb"
)

var SettingCustomerAllowFilterFieldsAndOps = []string{"fullname:like", "gender:in", "email:like", "phone:like"}
var SettingCustomerAllowSortFields = []string{"fullname", "gender", "email", "phone"}

type SettingCustomer struct {
	*models.CommonModel
	Id                    int
	FullName              string
	Gender                string
	Email                 string
	Phone                 string
	AdditionalInformation string
}

type SettingCustomerResponse struct {
	Id                    int    `json:"id"`
	FullName              string `json:"full_name"`
	Gender                string `json:"gender"`
	Email                 string `json:"email"`
	Phone                 string `json:"phone"`
	AdditionalInformation any    `json:"additional_information"`
}

func SettingCustomerToResponse(settingCustomer SettingCustomer) SettingCustomerResponse {
	var additionalInformation any
	if err := json.Unmarshal([]byte(settingCustomer.AdditionalInformation), &additionalInformation); err != nil {
		additionalInformation = []byte{}
	}
	return SettingCustomerResponse{
		Id:                    settingCustomer.Id,
		FullName:              settingCustomer.FullName,
		Gender:                settingCustomer.Gender,
		Email:                 settingCustomer.Email,
		Phone:                 settingCustomer.Phone,
		AdditionalInformation: additionalInformation,
	}
}

type SettingCustomerCreateRequest struct {
	FullName              string `json:"full_name" validate:"required"`
	Gender                string `json:"gender" validate:"required,gender"`
	Email                 string `json:"email" validate:"required,email"`
	Phone                 string `json:"phone" validate:"required,phone"`
	AdditionalInformation string `json:"additional_information" validate:"required,json"`
}

type SettingCustomerUpdateRequest struct {
	FullName              *string `json:"full_name"`
	Gender                *string `json:"gender" validate:"omitempty,gender"`
	Email                 *string `json:"email" validate:"omitempty,email"`
	Phone                 *string `json:"phone" validate:"omitempty,phone"`
	AdditionalInformation *string `json:"additional_information" validate:"omitempty,json"`
}

func (request SettingCustomerUpdateRequest) MapUpdateFields(bqbQuery *bqb.Query, fieldname string, value interface{}) error {
	switch strings.ToLower(fieldname) {
	case "fullname":
		bqbQuery.Comma("fullname = ?", value)
	case "gender":
		bqbQuery.Comma("gender = ?", value)
	case "email":
		bqbQuery.Comma("email = ?", value)
	case "phone":
		bqbQuery.Comma("phone = ?", value)
	case "additional_information":
		bqbQuery.Comma("additional_information = ?", value)
	default:
		return models.ErrInvalidUpdateField
	}
	return nil
}
