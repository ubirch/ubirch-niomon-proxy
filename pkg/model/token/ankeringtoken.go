package token

import (
	"github.com/dermicha/goutils/database"
	"gorm.io/gorm"
)

type AnkerToken struct {
	gorm.Model
	Token     string `gorm:"uniqueIndex"`
	UsedState bool   `gorm:"used_state"`
}

func IsValidToken(token string) bool {
	db := database.GetDb()
	var ats = []AnkerToken{}

	db.
		Limit(1).
		Where("token= ?", token).
		Where("used_state = 0").
		Find(&ats)

	if len(ats) == 1 {
		return true
	} else {
		return false
	}
}

func UseToken(token string) bool {
	db := database.GetDb()
	var ats = []AnkerToken{}

	db.
		Limit(1).
		Where("token= ?", token).
		Where("used_state = 0").
		Find(&ats)

	if len(ats) == 1 {
		at := ats[0]
		at.UsedState = true
		database.GetDb().Save(at)

		return true
	} else {
		return false
	}
}
