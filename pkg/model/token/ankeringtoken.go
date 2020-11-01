package token

import (
	"github.com/dermicha/goutils/database"
	"gorm.io/gorm"
)

type AnkerToken struct {
	gorm.Model
	Token       string `gorm:"uniqueIndex"`
	UsedCounter int    `gorm:"used_counter"`
}

func IsValidToken(token string) bool {
	db := database.GetDb()
	var ats = []AnkerToken{}

	db.
		Limit(1).
		Where("token= ?", token).
		Where("used_state > 0").
		Find(&ats)

	if len(ats) == 1 {
		return true
	} else {
		return false
	}
}

func UseToken(token string) int {
	db := database.GetDb()
	var ats = []AnkerToken{}

	db.
		Limit(1).
		Where("token= ?", token).
		Where("used_state > 0").
		Find(&ats)

	if len(ats) == 1 {
		at := ats[0]
		at.UsedCounter = at.UsedCounter - 1
		database.GetDb().Save(at)
		return at.UsedCounter
	} else {
		return 0
	}
}
