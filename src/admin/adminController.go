package admin

import (
	"hppGacha/src/logic"
	"hppGacha/src/models"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

var DB *gorm.DB

func Index(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("src/views/admin/indexAdmin.html")
	if err != nil {
		log.Panic(err)
	}
	tpl.Execute(w, nil)
}

func NewCard(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("src/views/admin/newCard.html")
	if err != nil {
		log.Panic(err)
	}

	tpl.Execute(w, nil)
}

func ShowUser(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("src/views/admin/userTable.html")
	if err != nil {
		log.Panic(err)
	}
	var users []models.User
	err = DB.Find(&users).Error
	if err != nil {
		log.Panic(err)
	}
	tpl.Execute(w, users)
}

func ProcessCard(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	cardName := r.PostFormValue("cardName")
	cardRarity, err := strconv.Atoi(r.PostFormValue("rarity"))
	if err != nil {
		log.Panic(err)
	}
	cardWeight, err := strconv.Atoi(r.PostFormValue("weight"))
	if err != nil {
		log.Panic(err)
	}
	cardShadow, err := strconv.ParseBool(r.PostFormValue("shadow"))
	if err != nil {
		log.Panic(err)
	}
	cardEvent, err := strconv.Atoi(r.PostFormValue("event"))
	if err != nil {
		log.Panic(err)
	}
	pathToFile := strings.ToLower(strings.ReplaceAll(cardName, " ", "_"))
	file, header, err := r.FormFile("cardImage")
	fileName := header.Filename
	ext := filepath.Ext(fileName)
	if err != nil {
		log.Panic(err)
		return
	}
	defer file.Close()
	tempFile, err := ioutil.TempFile("src/ressources/temp-pics", pathToFile+"-*"+ext)
	if err != nil {
		log.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Panic(err)
	}
	tempFile.Write(fileBytes)
	err = os.Rename(tempFile.Name(), "src/ressources/pics/"+pathToFile+ext)
	if err != nil {
		log.Panic(err)
	}
	newCard := models.Card{Cardname: cardName, Rarity: uint(cardRarity), Weight: uint(cardWeight), Avatar: pathToFile + ext, IsShadowCard: cardShadow, IsEventCard: uint(cardEvent)}
	err = DB.Create(&newCard).Error
	if err != nil {
		log.Panic(err)
	}
	logic.AddEntry(newCard)
	http.Redirect(w, r, "/admin/new_card", http.StatusFound)
}
