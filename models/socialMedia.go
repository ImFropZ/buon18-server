package models

import "time"

type SocialMedia struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	AccountID uint      `json:"account_id" gorm:"column:account_id"`
	Platform  string    `json:"platform"`
	URL       string    `json:"url" gorm:"column:url"`
	CID       uint      `json:"cid" gorm:"column:cid"`
	CTime     time.Time `json:"ctime" gorm:"column:ctime"`
	MID       uint      `json:"mid" gorm:"column:mid"`
	MTime     time.Time `json:"mtime" gorm:"column:mtime"`
}

func (SocialMedia) TableName() string {
	return "social_media"
}

func (s *SocialMedia) PrepareForCreate(cid uint, mid uint) (err error) {
	s.CID = cid
	s.CTime = time.Now()
	s.MID = mid
	s.MTime = time.Now()
	return
}

func (s *SocialMedia) PrepareForUpdate(mid uint) (err error) {
	s.MID = mid
	s.MTime = time.Now()
	return
}
