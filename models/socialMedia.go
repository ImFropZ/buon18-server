package models

import "time"

type SocialMedia struct {
	Id uint `json:"id"`
}

type SocialMediaData struct {
	Id            uint      `json:"id"`
	SocialMediaId uint      `json:"social_media_id"`
	Platform      string    `json:"platform"`
	URL           string    `json:"url"`
	CId           uint      `json:"cid"`
	CTime         time.Time `json:"ctime"`
	MId           uint      `json:"mid"`
	MTime         time.Time `json:"mtime"`
}

type SocialMediaResponse struct {
	Id       uint   `json:"id"`
	Platform string `json:"platform"`
	URL      string `json:"url"`
}

func (s *SocialMediaData) PrepareForCreate(cid uint, mid uint) (err error) {
	s.CId = cid
	s.CTime = time.Now()
	s.MId = mid
	s.MTime = time.Now()
	return
}

func (s *SocialMediaData) PrepareForUpdate(mid uint) (err error) {
	s.MId = mid
	s.MTime = time.Now()
	return
}

func (s *SocialMediaData) ToResponse() SocialMediaResponse {
	return SocialMediaResponse{
		Id:       s.Id,
		Platform: s.Platform,
		URL:      s.URL,
	}
}
