package admin

import (
	logic "hppGacha/logic"
	"log"
)

func addNewCard(cardName string, cardRarity int, cardWeight int, pathToFile string) {
	_, err := logic.DB.Exec("INSERT INTO rollable_users (name,pathToPic,weight,rarity) VALUES (?, ?, ?, ?)", cardName, pathToFile+".png", cardWeight, cardRarity)
	if err != nil {
		log.Fatal(err)
	}
}
