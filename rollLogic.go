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
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var accumulatedWeight int
var entries []BannerRoll

func addEntry(name string, avatar string, weight int) {
	accumulatedWeight += weight
	var UserHpp BannerRoll
	UserHpp.Name = name
	UserHpp.Avatar = avatar
	UserHpp.accumulatedWeight = accumulatedWeight
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
	db, err := sql.Open("sqlite3", "hppgacha.db")
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db, err
}

func dataToRoll(db *sql.DB) {
	results, err := db.Query("Select name,pathToPic,weight FROM rollable_users")
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var accumulatedWeight int
		var name string
		var avatar string
		err = results.Scan(&name, &avatar, &accumulatedWeight)
		if err != nil {
			fmt.Println(err)
		}
		addEntry(name, avatar, accumulatedWeight)
	}
}
