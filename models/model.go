package models

import "time"

var (
	// Setting user types
	SettingUserTypUser = "user"
	SettingUserTypBot  = "bot"
)

type CommonModel struct {
	CId   uint
	CTime time.Time
	MId   uint
	MTime time.Time
}

func (cm *CommonModel) PrepareForCreate(cid uint, mid uint) (err error) {
	cm.CId = cid
	cm.CTime = time.Now()
	cm.MId = mid
	cm.MTime = time.Now()
	return
}

func (cm *CommonModel) PrepareForUpdate(mid uint) (err error) {
	cm.MId = mid
	cm.MTime = time.Now()
	return
}
