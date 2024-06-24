package models

import (
	"time"
)

type Account struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	Code           string    `json:"code"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Gender         string    `json:"gender"`
	Address        string    `json:"address"`
	Phone          string    `json:"phone"`
	SecondaryPhone string    `json:"secondary_phone" gorm:"column:secondary_phone"`
	Deleted        bool      `json:"deleted"`
	CID            uint      `json:"cid" gorm:"column:cid"`
	CTime          time.Time `json:"ctime" gorm:"column:ctime"`
	MID            uint      `json:"mid" gorm:"column:mid"`
	MTime          time.Time `json:"mtime" gorm:"column:mtime"`

	// -- Associations
	SocialMedias []SocialMedia `json:"social_medias" gorm:"foreignKey:account_id"`
}

func (Account) TableName() string {
	return "account"
}

func (a *Account) PrepareForCreate(cid uint, mid uint) (err error) {
	a.CID = cid
	a.CTime = time.Now()
	a.MID = mid
	a.MTime = time.Now()
	return
}

func (a *Account) PrepareForUpdate(mid uint) (err error) {
	a.MID = mid
	a.MTime = time.Now()
	return
}
