package logic

import (
	"encoding/json"
	models "hppGacha/src/models"
	"html/template"
	"log"
	"net/http"
)

type recycleResponse struct {
	ErrorString string `json:"error_string" db:"error_string"`
	NewQuantity uint   `json:"new_quantity" db:"new_quantity"`
}

type recycleTarget struct {
	RecycleTargetId uint `json:"id,string" db:"id"`
}

func ShowInventory(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("src/views/inventory.html")
	if err != nil {
		log.Panic(err)
	}
	connectedUser, exists := getConnectedUser(w, r)
	if exists {
		err = tpl.Execute(w, connectedUser)
		if err != nil {
			log.Panicln(err)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func RecycleCard(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var connectedUser models.User
		var recycleTarget recycleTarget
		var response recycleResponse
		var targetedCard models.Card

		err := json.NewDecoder(r.Body).Decode(&recycleTarget)
		if err != nil {
			log.Panicln(err)
		}
		DB.First(&targetedCard, recycleTarget.RecycleTargetId)
		connectedUser, exists := getConnectedUser(w, r)
		if exists {
			response.NewQuantity, response.ErrorString = consumeDoublonForRoll(connectedUser, targetedCard)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

func consumeDoublonForRoll(user models.User, card models.Card) (uint, string) {
	var storedCard models.CardInInventory
	err := DB.Where("user_id=? AND card_id=?", user.ID, card.ID).First(&storedCard).Error
	if err != nil {
		log.Panic(err)
	}
	if storedCard.Quantity >= 4 {
		user.AvailableRolls = user.AvailableRolls + 1
		err := DB.Save(&user).Error
		if err != nil {
			log.Panic(err)
		}
		storedCard.Quantity = storedCard.Quantity - 3
		err = DB.Save(&storedCard).Error
		if err != nil {
			log.Panicln(err)
		}
		return (storedCard.Quantity), ""
	}
	return 0, "An error occured while trying to recycle the card"
}
