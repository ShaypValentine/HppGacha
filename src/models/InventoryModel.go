package models

import "gorm.io/gorm"

type CardInInventory struct {
	gorm.Model
	UserID   uint
	User     User
	CardID   uint
	Card     Card
	Quantity uint
}

type Inventory struct {
	User             User
	CardsInInventory []CardInInventory
}

type Tabler interface {
	TableName() string
}

func (CardInInventory) TableName() string {
	return "cards_in_inventory"
}
