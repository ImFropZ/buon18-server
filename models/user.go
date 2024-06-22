package models

import "time"

type User struct {
	ID    uint      `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
	Pwd   string    `json:"password"`
	Role  string    `json:"role"`
	CID   uint      `json:"cid" gorm:"column:cid"`
	CTime time.Time `json:"ctime" gorm:"column:ctime"`
	MID   uint      `json:"mid" gorm:"column:mid"`
	MTime time.Time `json:"mtime" gorm:"column:mtime"`
}

func (User) TableName() string {
	return "user"
}

func (u *User) PrepareForCreate(cid uint, mid uint) (err error) {
	u.CID = cid
	u.CTime = time.Now()
	u.MID = mid
	u.MTime = time.Now()
	return
}

func (u *User) PrepareForUpdate(mid uint) (err error) {
	u.MID = mid
	u.MTime = time.Now()
	return
}
