package models

type User struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Pwd   string `json:"password"`
	Role  string `json:"role"`
}

func (User) TableName() string {
	return "user"
}
