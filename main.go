package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type BannerRoll struct {
	name              string
	accumulatedWeight int
}

var accumulatedWeight int
var entries []BannerRoll

func main() {
	csvBanner, err := os.Open("banner_content.csv")
	if err != nil {
		fmt.Println(err)
	}

	defer csvBanner.Close()
	csvLines, err := csv.NewReader(csvBanner).ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	linesToEntries(csvLines)

	for i := 0; i < 10; i++ {
		u := getRandom()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("\nYou just found a wild %s\n", u)
	}
}

func addEntry(name string, weight int) {
	accumulatedWeight += weight
	var UserHpp BannerRoll
	UserHpp.name = name
	UserHpp.accumulatedWeight = accumulatedWeight
	entries = append(entries, UserHpp)
}

func getRandom() string {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(1 * accumulatedWeight)
	for _, entry := range entries {
		if entry.accumulatedWeight >= r {
			return entry.name
		}
	}
	return "error"
}

func linesToEntries(csvLines [][]string) {

	for _, line := range csvLines {
		weight, err := strconv.Atoi(line[1])
		if err != nil {
			fmt.Println(err)
		}
		addEntry(line[0], weight)
	}
}
