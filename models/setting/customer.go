package setting

import (
	"encoding/json"

	"system.buon18.com/m/models"
)

type SettingCustomer struct {
	*models.CommonModel
	Id                    *int
	FullName              *string
	Gender                *string
	Email                 *string
	Phone                 *string
	AdditionalInformation *string
}

func (SettingCustomer) AllowFilterFieldsAndOps() []string {
	return []string{"full-name:like", "full-name:ilike", "gender:in", "email:like", "phone:like"}
}

func (SettingCustomer) AllowSorts() []string {
	return []string{"full-name", "gender", "email", "phone"}
}

type SettingCustomerResponse struct {
	Id                    *int    `json:"id,omitempty"`
	FullName              *string `json:"full_name,omitempty"`
	Gender                *string `json:"gender,omitempty"`
	Email                 *string `json:"email,omitempty"`
	Phone                 *string `json:"phone,omitempty"`
	AdditionalInformation *any    `json:"additional_information,omitempty"`
}

func SettingCustomerToResponse(settingCustomer SettingCustomer) SettingCustomerResponse {
	var additionalInformation any
	if err := json.Unmarshal([]byte(*settingCustomer.AdditionalInformation), &additionalInformation); err != nil {
		additionalInformation = []byte{}
	}
	return SettingCustomerResponse{
		Id:                    settingCustomer.Id,
		FullName:              settingCustomer.FullName,
		Gender:                settingCustomer.Gender,
		Email:                 settingCustomer.Email,
		Phone:                 settingCustomer.Phone,
		AdditionalInformation: &additionalInformation,
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
