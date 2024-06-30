package models

import (
	"time"
)

type Account struct {
	Id             uint
	Code           string
	Name           string
	Email          string
	Gender         string
	Address        string
	Phone          string
	SecondaryPhone string
	Deleted        bool
	CId            uint
	CTime          time.Time
	MId            uint
	MTime          time.Time

	// -- Associations
	SocialMedias []SocialMediaData
}

type AccountResponse struct {
	Id             uint                  `json:"id"`
	Code           string                `json:"code"`
	Name           string                `json:"name"`
	Gender         string                `json:"gender"`
	Email          string                `json:"email"`
	Address        string                `json:"address"`
	Phone          string                `json:"phone"`
	SecondaryPhone string                `json:"secondary_phone"`
	SocialMedias   []SocialMediaResponse `json:"social_medias"`
}

func (a *Account) PrepareForCreate(cid uint, mid uint) (err error) {
	a.CId = cid
	a.CTime = time.Now()
	a.MId = mid
	a.MTime = time.Now()
	return
}

func (a *Account) PrepareForUpdate(mid uint) (err error) {
	a.MId = mid
	a.MTime = time.Now()
	return
}

func (a *Account) ToResponse() AccountResponse {
	acc := AccountResponse{
		Id:             a.Id,
		Code:           a.Code,
		Name:           a.Name,
		Gender:         a.Gender,
		Email:          a.Email,
		Address:        a.Address,
		Phone:          a.Phone,
		SecondaryPhone: a.SecondaryPhone,
		SocialMedias:   make([]SocialMediaResponse, 0),
	}

	for _, smd := range a.SocialMedias {
		socialMedia := smd.ToResponse()
		acc.SocialMedias = append(acc.SocialMedias, socialMedia)
	}

	return acc
}

func AccountsToResponse(accounts []Account) []AccountResponse {
	accountsResponse := make([]AccountResponse, 0)
	for _, account := range accounts {
		accountsResponse = append(accountsResponse, account.ToResponse())
	}
	return accountsResponse
}
