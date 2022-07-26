package models

import (
	"time"

	"gorm.io/gorm"
)

type ShadowPortal struct {
	gorm.Model
	UserID               uint
	BaseCardLeft         uint `gorm:"default:75"`
	RareCardLeft         uint `gorm:"default:15"`
	HasAccess            bool `gorm:"default:false"`
	AvailableShadowRolls uint `gorm:"default:2"`
	SacrificeAvailableAt time.Time
}
