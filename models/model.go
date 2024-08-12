package models

import "time"

var (
	SettingGenderTypMale   = "m"
	SettingGenderTypFemale = "f"
	SettingGenderTypOther  = "o"

	// Setting user types
	SettingUserTypUser = "user"
	SettingUserTypBot  = "bot"
)

var VALID_GENDER_TYPES = []string{SettingGenderTypMale, SettingGenderTypFemale, SettingGenderTypOther}

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
