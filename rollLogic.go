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

func linesToEntries(csvLines [][]string) {

	for _, line := range csvLines {
		weight, err := strconv.Atoi(line[1])
		if err != nil {
			fmt.Println(err)
		}
		addEntry(line[0], line[2], weight)
	}
}

func fileToLines(filePath string) {
	csvBanner, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}

	defer csvBanner.Close()
	csvLines, err := csv.NewReader(csvBanner).ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	linesToEntries(csvLines)
}

func emptyEntries() {
	entries = nil
	accumulatedWeight = 0
}
