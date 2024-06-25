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
	SocialMedias []SocialMedia `json:"social_medias"`
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
