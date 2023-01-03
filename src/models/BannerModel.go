package models

import "gorm.io/gorm"

type Banner struct {
	gorm.Model
	ID             uint
	Bannername     string
	BannerImage    string
	Cards []Card `gorm:"many2many:banner_cards;"`
	IsAvailable    bool `gorm:"default:false"`
}
