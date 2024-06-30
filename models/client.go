package models

import (
	"time"
)

type Client struct {
	Id        uint
	Code      string
	Name      string
	Address   string
	Phone     string
	Latitude  float64
	Longitude float64
	Note      string
	CId       uint
	CTime     time.Time
	MId       uint
	MTime     time.Time

	// -- Associations
	SocialMedias []SocialMediaData
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
		Id:           a.Id,
		Code:         a.Code,
		Name:         a.Name,
		Address:      a.Address,
		Phone:        a.Phone,
		Latitude:     a.Latitude,
		Longitude:    a.Longitude,
		Note:         a.Note,
		SocialMedias: SocialMediasToResponse(a.SocialMedias),
	}
}

func ClientsToResponse(clients []Client) []ClientResponse {
	clientsResponse := make([]ClientResponse, 0)
	for _, client := range clients {
		clientsResponse = append(clientsResponse, client.ToResponse())
	}
	return clientsResponse
}
