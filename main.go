package main

import (
	"fmt"
	"math/rand"
	"time"
)

type BannerRoll struct {
	name            string
	acculatedWeight int
}

var accumulatedWeight int
var entries []BannerRoll

func main() {
	addEntry("⭐ Shayp", 50)
	addEntry("⭐ Madao", 50)
	addEntry("⭐ 'Drive'", 40)
	addEntry("⭐⭐ Koenji", 13)
	addEntry("⭐⭐ Cirrus", 12)
	addEntry("⭐⭐⭐ 18🥖 ", 2)
	for i := 0; i < 10; i++ {
		u := getRandom()
		fmt.Printf("\nYou just found a wild %s\n", u)
	}
}

func addEntry(name string, weight int) {
	accumulatedWeight += weight
	var UserHpp BannerRoll
	UserHpp.name = name
	UserHpp.acculatedWeight = accumulatedWeight
	entries = append(entries, UserHpp)
}

func getRandom() string {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(1 * accumulatedWeight)
	for _, entry := range entries {
		if entry.acculatedWeight >= r {
			return entry.name
		}
	}
	return "error"
}
