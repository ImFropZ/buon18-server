package models

import (
	"time"
)

type Client struct {
	Id        uint      `json:"id" `
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Note      string    `json:"note"`
	CId       uint      `json:"cid"`
	CTime     time.Time `json:"ctime"`
	MId       uint      `json:"mid"`
	MTime     time.Time `json:"mtime"`

	// -- Associations
	SocialMedias []SocialMediaData `json:"social_medias"`
}

type ClientResponse struct {
	Id           uint                  `json:"id"`
	Code         string                `json:"code"`
	Name         string                `json:"name"`
	Address      string                `json:"address"`
	Phone        string                `json:"phone"`
	Latitude     float64               `json:"latitude"`
	Longitude    float64               `json:"longitude"`
	Note         string                `json:"note"`
	SocialMedias []SocialMediaResponse `json:"social_medias"`
}

func (a *Client) PrepareForCreate(cid uint, mid uint) (err error) {
	a.CId = cid
	a.CTime = time.Now()
	a.MId = mid
	a.MTime = time.Now()
	return
}

func (a *Client) PrepareForUpdate(mid uint) (err error) {
	a.MId = mid
	a.MTime = time.Now()
	return
}

func (a *Client) ToResponse() ClientResponse {
	return ClientResponse{
		Id:        a.Id,
		Code:      a.Code,
		Name:      a.Name,
		Address:   a.Address,
		Phone:     a.Phone,
		Latitude:  a.Latitude,
		Longitude: a.Longitude,
		Note:      a.Note,
	}
}
