package models

import "time"

type SocialMedia struct {
	Id        uint      `json:"id"`
	AccountId uint      `json:"account_id"`
	Platform  string    `json:"platform"`
	URL       string    `json:"url"`
	CId       uint      `json:"cid"`
	CTime     time.Time `json:"ctime"`
	MId       uint      `json:"mid"`
	MTime     time.Time `json:"mtime"`
}

func (s *SocialMedia) PrepareForCreate(cid uint, mid uint) (err error) {
	s.CId = cid
	s.CTime = time.Now()
	s.MId = mid
	s.MTime = time.Now()
	return
}

func (s *SocialMedia) PrepareForUpdate(mid uint) (err error) {
	s.MId = mid
	s.MTime = time.Now()
	return
}
