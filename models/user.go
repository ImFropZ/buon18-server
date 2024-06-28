package models

import "time"

type User struct {
	Id      uint      `json:"id"`
	Email   string    `json:"email"`
	Name    string    `json:"name"`
	Pwd     string    `json:"password"`
	Role    string    `json:"role"`
	Deleted bool      `json:"deleted"`
	CId     uint      `json:"cid"`
	CTime   time.Time `json:"ctime"`
	MId     uint      `json:"mid"`
	MTime   time.Time `json:"mtime"`
}

type UserResponse struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (u *User) PrepareForCreate(cid uint, mid uint) (err error) {
	u.CId = cid
	u.CTime = time.Now()
	u.MId = mid
	u.MTime = time.Now()
	return
}

func (u *User) PrepareForUpdate(mid uint) (err error) {
	u.MId = mid
	u.MTime = time.Now()
	return
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		Id:    u.Id,
		Name:  u.Name,
		Email: u.Email,
		Role:  u.Role,
	}
}

func UsersToResponse(users []User) []UserResponse {
	usersResponse := make([]UserResponse, 0)
	for _, user := range users {
		usersResponse = append(usersResponse, user.ToResponse())
	}
	return usersResponse
}
