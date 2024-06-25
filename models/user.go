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
