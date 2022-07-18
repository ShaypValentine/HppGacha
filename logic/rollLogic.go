package logic

import (
	"math/rand"
	"time"

	models "hppGacha/src/models"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Entry struct {
	Card              models.Card
	AccumulatedWeight uint
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var accumulatedWeight uint
var entries []Entry

func addEntry(card models.Card) {
	var entry Entry
	entry.Card = card
	accumulatedWeight += card.Weight
	entry.AccumulatedWeight = accumulatedWeight
	entries = append(entries, entry)
}

// func getRandom() Entry {
// 	r := rand.Intn(1 * int(accumulatedWeight))
// 	for _, entry := range entries {
// 		if int(entry.AccumulatedWeight) >= r {
// 			return entry
// 		}
// 	}
// 	return Entry{}
// }

func EmptyEntries() {
	entries = nil
	accumulatedWeight = 0
}

func DatabaseConnection() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("hppgachaOrm.db"), &gorm.Config{})
	return db, err
}

func DataToRoll(db *gorm.DB) {
	var cards []models.Card
	_ = db.Find(&cards)
	for _, card := range cards {
		addEntry(card)
	}
}
