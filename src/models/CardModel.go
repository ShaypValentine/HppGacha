package models

import "gorm.io/gorm"

type Card struct {
	gorm.Model
	ID               uint
	Cardname         string
	Rarity           uint
	Avatar           string
	Weight           uint
	CardsInInventory []CardInInventory
	IsShadowCard     bool `gorm:"default:false"`
	IsEventCard     uint `gorm:"default:0"`
}
