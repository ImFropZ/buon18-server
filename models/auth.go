package models

type TokenAndRefreshToken struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
