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

	"github.com/Masterminds/sprig/v3"
	"gorm.io/gorm"
)

var DB *gorm.DB
var templateAdminPath = "src/views/admin/"

func Index(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.New("indexAdmin.html").Funcs(sprig.FuncMap()).ParseFiles(
		templateAdminPath+"indexAdmin.html",
		templateAdminPath+"adminNavBar.html"))

	err := tpl.Execute(w, nil)
	if err != nil {
		log.Panicln(err)
	}
}

func NewCard(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.New("newCard.html").Funcs(sprig.FuncMap()).ParseFiles(
		templateAdminPath+"newCard.html",
		templateAdminPath+"adminNavBar.html"))

	var banners []models.Banner
	DB.Find(&banners)
	err := tpl.Execute(w, banners)
	if err != nil {
		log.Panicln(err)
	}
}

func EditCard(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.New("cardEdit.html").Funcs(sprig.FuncMap()).ParseFiles(
		templateAdminPath+"cardEdit.html",
		templateAdminPath+"adminNavBar.html"))

	value := r.FormValue("id")

	var card models.Card
	err := DB.First(&card, value).Error
	if err != nil {
		log.Panic(err)
	}

	err = tpl.Execute(w, card)
	if err != nil {
		log.Panicln(err)
	}
}

func ShowUser(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.New("userTable.html").Funcs(sprig.FuncMap()).ParseFiles(
		templateAdminPath+"userTable.html",
		templateAdminPath+"adminNavBar.html"))

	var users []models.User
	err := DB.Find(&users).Error
	if err != nil {
		log.Panic(err)
	}
	err = tpl.Execute(w, users)
	if err != nil {
		log.Panicln(err)
	}
}

func ShowCards(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.New("cardList.html").Funcs(sprig.FuncMap()).ParseFiles(
		templateAdminPath+"cardList.html",
		templateAdminPath+"adminNavBar.html"))

	var cards []models.Card
	err := DB.Order("cardname asc").Find(&cards).Error
	if err != nil {
		log.Panic(err)
	}
	err = tpl.Execute(w, cards)
	if err != nil {
		log.Panicln(err)
	}

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

func ProcessCardEdit(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	cardName := r.PostFormValue("cardName")
	cardID := r.PostFormValue("cardID")
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

	var card models.Card
	err = DB.First(&card, cardID).Error
	if err != nil {
		log.Panic(err)
	}
	card.Cardname = cardName
	card.Rarity = uint(cardRarity)
	card.Weight = uint(cardWeight)
	card.IsShadowCard = cardShadow
	card.IsEventCard = uint(cardEvent)

	err = DB.Save(&card).Error
	if err != nil {
		log.Panic(err)
	}
	logic.EmptyEntries()
	logic.DataToRoll(DB)
	http.Redirect(w, r, "/admin/show_cards", http.StatusFound)
}
