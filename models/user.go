package models

type User struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Pwd   string `json:"password"`
}

func (User) TableName() string {
	return "user"
}
