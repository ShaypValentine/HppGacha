package logic

import (
	"encoding/json"
	models "hppGacha/src/models"
	"html/template"
	"log"
	"net/http"
)

type sacrificeTarget struct {
	SacrificeTargetId     uint `json:"id,string" db:"id"`
	SacrificeTargetRarity uint `json:"rarity,string" db:"id"`
}

func ShadowIndex(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles(
		templatePath+"shadowPortal.html",
		templatePath+"navbar.html",
		templatePath+"_parts/head.html",
		templatePath+"_parts/footer.html",
		templatePath+"_parts/js.html")
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

func SacrificeCard(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var connectedUser models.User
		var sacrificeTarget sacrificeTarget
		var response recycleResponse
		var targetedCard models.Card

		err := json.NewDecoder(r.Body).Decode(&sacrificeTarget)
		if err != nil {
			log.Panicln(err)
		}
		DB.First(&targetedCard, sacrificeTarget.SacrificeTargetId)
		connectedUser, exists := getConnectedUser(w, r)
		if exists {
			response.NewQuantity, response.ErrorString = consumeSacrifice(connectedUser, targetedCard)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

func consumeSacrifice(user models.User, card models.Card) (uint, string) {
	var storedCard models.CardInInventory
	var shadowPortal models.ShadowPortal
	err := DB.Where("user_id=? AND card_id=?", user.ID, card.ID).Preload("Card").First(&storedCard).Error
	if err != nil {
		log.Panic(err)
	}
	err = DB.Where("user_id=?", user.ID).Assign(models.ShadowPortal{UserID: user.ID}).FirstOrCreate(&shadowPortal).Error
	if err != nil {
		log.Panic(err)
	}

	if storedCard.Quantity >= 2 {
		if storedCard.Card.Rarity == 2 {
			shadowPortal.RareCardLeft -= 1
		} else {
			shadowPortal.BaseCardLeft -= 1
		}
		if shadowPortal.BaseCardLeft <= 0 && shadowPortal.RareCardLeft <= 0 {
			shadowPortal.HasAccess = true
		}
		err := DB.Save(&shadowPortal).Error
		if err != nil {
			log.Panic(err)
		}
		storedCard.Quantity = storedCard.Quantity - 1
		err = DB.Save(&storedCard).Error
		if err != nil {
			log.Panicln(err)
		}
		return (storedCard.Quantity), ""
	}
	return 0, "An error occured while trying to sacrifice the card"
}

func ShadowRoll(w http.ResponseWriter, r *http.Request) {
	var rolledCard RolledCard
	tpl, err := template.ParseFiles("src/views/rollCard.html")
	if err != nil {
		log.Panicln(err)
	}
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		rolledItem := getRandomShadow()
		connectedUser, exists := getConnectedUser(w, r)
		if exists {
			Refresh(w, r)
			if connectedUser.ShadowPortal.AvailableShadowRolls > 0 {
				connectedUser.ShadowPortal.AvailableShadowRolls--
				DB.Save(&connectedUser.ShadowPortal)
				addToInventory(connectedUser, rolledItem)
			}
		}
		rolledCard.User = connectedUser
		rolledCard.Card = rolledItem.Card
		err = tpl.Execute(w, rolledCard)
		if err != nil {
			log.Panicln(err)
		}

	}

}
