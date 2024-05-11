package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username         string `gorm:"uniqueIndex"`
	Password         string
	Role             string `gorm:"default:base_user"`
	AvailableRolls   uint   `gorm:"default:4"`
	CardsInInventory []CardInInventory
	ShadowPortal     ShadowPortal
}
