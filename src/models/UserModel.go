package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username         string `gorm:"uniqueIndex"`
	Password         string
	AvailableRolls   uint `gorm:"default:4"`
	CardsInInventory []CardInInventory
	ShadowPortal     ShadowPortal
}
