package models

import (
	"encoding/json"
)

var SettingCustomerAllowFilterFieldsAndOps = []string{"fullname-like", "gender-in", "email-like", "phone-like"}
var SettingCustomerAllowSortFields = []string{"fullname", "gender", "email", "phone"}

type SettingCustomer struct {
	*CommonModel
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
