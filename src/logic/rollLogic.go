package logic

import (
	"log"
	"math/rand"
	"time"

	models "hppGacha/src/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Entry struct {
	Card              models.Card
	AccumulatedWeight uint
	BannerID          uint
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var BannersWeight = make(map[uint]uint)
var accumulatedWeight uint
var shadowWeight uint
var entries []Entry
var shadowEntries []Entry

func AddEntry(card models.Card) {
	var entry Entry
	var shadowEntry Entry
	if card.IsShadowCard {
		shadowEntry.Card = card
		shadowWeight += card.Weight
		shadowEntry.AccumulatedWeight = shadowWeight
		shadowEntries = append(shadowEntries, shadowEntry)
	} else {
		accumulatedWeight += card.Weight
		if len(card.Banners) > 0 {
			for _, banner := range card.Banners {
				entry.Card = card
				BannersWeight[banner.ID] += card.Weight
				entry.AccumulatedWeight = BannersWeight[banner.ID]
				entry.BannerID = banner.ID
				entries = append(entries, entry)
			}
		}
	}

}

func getRandom(bannerID uint) Entry {
	var entriesBanner []Entry
	for _, entry := range entries {
		if bannerID == entry.BannerID {
			entriesBanner = append(entriesBanner, entry)
		}
	}
	r := rand.Intn(1 * int(BannersWeight[bannerID]))
	for _, entry := range entriesBanner {
		if int(entry.AccumulatedWeight) >= r {
			return entry
		}
	}
	return Entry{}
}
func getRandomShadow() Entry {
	r := rand.Intn(1 * int(shadowWeight))
	for _, entry := range shadowEntries {
		if int(entry.AccumulatedWeight) >= r {
			return entry
		}
	}
	return Entry{}
}

func EmptyEntries() {
	entries = nil
	shadowEntries = nil
	accumulatedWeight = 0
	shadowWeight = 0
	BannersWeight = make(map[uint]uint)
}

func DatabaseConnection() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("persistent/hppgacha.db"), &gorm.Config{})
	return db, err
}

func DataToRoll(db *gorm.DB) {
	var cards []models.Card
	err := db.Model(&models.Card{}).Preload("Banners").Find(&cards).Error
	if err != nil {
		log.Panic(err)
	}
	for _, card := range cards {
		AddEntry(card)
	}
}
