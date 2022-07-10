package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type BannerRoll struct {
	Name              string
	Avatar            string
	accumulatedWeight int
	Rarity            int
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var accumulatedWeight int
var entries []BannerRoll

func addEntry(name string, avatar string, weight int, rarity int) {
	accumulatedWeight += weight
	var UserHpp BannerRoll
	UserHpp.Name = name
	UserHpp.Avatar = avatar
	UserHpp.accumulatedWeight = accumulatedWeight
	UserHpp.Rarity = rarity
	entries = append(entries, UserHpp)
}

func getRandom() BannerRoll {
	r := rand.Intn(1 * accumulatedWeight)
	for _, entry := range entries {
		if entry.accumulatedWeight >= r {
			return entry
		}
	}
	return BannerRoll{}
}

func emptyEntries() {
	entries = nil
	accumulatedWeight = 0
}

func databaseConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "hppgacha.db?cache=shared&mode=rwc")
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	DB = db
	return db, err
}

func dataToRoll(db *sql.DB) {
	results, err := db.Query("Select name,pathToPic,weight,rarity FROM rollable_users")
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var accumulatedWeight int
		var name string
		var avatar string
		var rarity int
		err = results.Scan(&name, &avatar, &accumulatedWeight, &rarity)
		if err != nil {
			fmt.Println(err)
		}
		addEntry(name, avatar, accumulatedWeight, rarity)
	}
}
