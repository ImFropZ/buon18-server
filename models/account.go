package models

import (
	"time"
)

type Account struct {
	Id             uint      `json:"id" `
	Code           string    `json:"code"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Gender         string    `json:"gender"`
	Address        string    `json:"address"`
	Phone          string    `json:"phone"`
	SecondaryPhone string    `json:"secondary_phone"`
	Deleted        bool      `json:"deleted"`
	CId            uint      `json:"cid"`
	CTime          time.Time `json:"ctime"`
	MId            uint      `json:"mid"`
	MTime          time.Time `json:"mtime"`

	// -- Associations
	SocialMedias []SocialMediaData `json:"social_medias"`
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
	var accountsResponse []AccountResponse
	for _, account := range accounts {
		accountsResponse = append(accountsResponse, account.ToResponse())
	}
	return accountsResponse
}
