package models

import "time"

type SocialMedia struct {
	Id uint
}

type SocialMediaData struct {
	Id            uint
	SocialMediaId uint
	Platform      string
	URL           string
	CId           uint
	CTime         time.Time
	MId           uint
	MTime         time.Time
}

type SocialMediaResponse struct {
	Id       uint   `json:"id"`
	Platform string `json:"platform"`
	URL      string `json:"url"`
}

type CreateSocialMediaRequest struct {
	Platform string `json:"platform" binding:"required"`
	URL      string `json:"url" binding:"required"`
}

type UpdateSocialMediaRequest struct {
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

func SocialMediasToResponse(s []SocialMediaData) []SocialMediaResponse {
	res := make([]SocialMediaResponse, 0)
	for _, v := range s {
		res = append(res, v.ToResponse())
	}
	return res
}
