package logic

import (
	"html/template"
	"log"
	"net/http"
)

type Inventory struct {
	Cards []CardInventory
	User  user
}

type CardInventory struct {
	Name     string
	Avatar   string
	Quantity int
	Rarity   int
}

func getInventoryForUser(user user) (inventory Inventory) {
	var card CardInventory
	var cardName string
	rows, err := DB.Query("Select cardName, quantity from inventory where user = ?", user.Id)
	if err != nil {
		log.Println("salut", err)
	}
	for rows.Next() {
		rows.Scan(&cardName, &card.Quantity)
		err := DB.QueryRow("Select name , pathToPic, rarity FROM rollable_users where name = ?", cardName).Scan(&card.Name, &card.Avatar, &card.Rarity)
		if err != nil {
			log.Println("hello", err)
		}
		inventory.Cards = append(inventory.Cards, card)
	}
	inventory.User = user
	return inventory
}

func ShowInventory(w http.ResponseWriter, r *http.Request) {

	var indexInfos IndexInfo
	tpl, err := template.ParseFiles("src/inventory.html")
	if err != nil {
		log.Fatalln(err)
	}
	connectedUser, exists := getConnectedUser(w, r)
	indexInfos.User = connectedUser
	if exists {
		indexInfos.Inventory = getInventoryForUser(connectedUser)
		tpl.Execute(w, indexInfos)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
