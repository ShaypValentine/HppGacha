package logic

import (
	models "hppGacha/src/models"
	"log"
)

// import (
// 	"encoding/json"
// 	"html/template"
// 	"log"
// 	"net/http"
// )

// type Inventory struct {
// 	Cards []CardInventory
// 	User  UserInfo
// }

// type recycleTarget struct {
// 	Name     string      `json:"name" db:"name"`
// 	Quantity json.Number `json:"quantity" db:"quantity"`
// }

// type recycleResponse struct {
// 	ErrorString string `json:"error_string" db:"error_string"`
// 	NewQuantity int64  `json:"new_quantity" db:"new_quantity"`
// }

func getInventoryForUser(user models.User) (inventory models.Inventory) {
	var cardsInInventory []models.CardInInventory
	result := DB.Where("user_id = ?", user.ID).Find(&cardsInInventory)
	if result.Error != nil {
		log.Println(result.Error)
	}
	inventory.CardsInInventory = cardsInInventory

	return inventory
}

// func ShowInventory(w http.ResponseWriter, r *http.Request) {

// 	var indexInfos IndexInfo
// 	tpl, err := template.ParseFiles("src/inventory.html")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	connectedUser, exists := getConnectedUser(w, r)
// 	indexInfos.User = connectedUser
// 	if exists {
// 		indexInfos.Inventory = getInventoryForUser(connectedUser)
// 		tpl.Execute(w, indexInfos)
// 	} else {
// 		http.Redirect(w, r, "/", http.StatusFound)
// 	}
// }

// func RecycleCard(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "POST" {
// 		var connectedUser UserInfo
// 		var recycleTarget recycleTarget
// 		var response recycleResponse
// 		err := json.NewDecoder(r.Body).Decode(&recycleTarget)
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
// 		connectedUser, exists := getConnectedUser(w, r)
// 		if exists {
// 			response.NewQuantity, response.ErrorString = consumeDoublonForRoll(connectedUser, recycleTarget)
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(response)
// 		} else {
// 			http.Redirect(w, r, "/", http.StatusFound)
// 		}
// 	}
// }

// func consumeDoublonForRoll(user UserInfo, card recycleTarget) (int64, string) {
// 	realQuantity := getQuantityForCard(user, card)
// 	sentQuantity, err := card.Quantity.Int64()
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	if realQuantity == sentQuantity {
// 		_, err := DB.Exec("UPDATE inventory SET quantity = quantity - 3 where user = ? and cardName = ?", user.Id, card.Name)
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
// 		_, err = DB.Exec("UPDATE users SET availableRoll = availableRoll + 1 where id = ? ", user.Id)
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
// 		return (realQuantity - 3), ""
// 	}
// 	return 0, "An error occured while trying to recycle the card"
// }

// func getQuantityForCard(user UserInfo, card recycleTarget) int64 {
// 	var realQuantity int64
// 	err := DB.QueryRow("select quantity from inventory where user = ? and cardName = ?", user.Id, card.Name).Scan(&realQuantity)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return realQuantity
// }
