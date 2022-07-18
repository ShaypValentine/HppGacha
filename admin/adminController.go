package admin

// import (
// 	"fmt"
// 	logic "hppGacha/logic"
// 	"html/template"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"os"
// 	"strconv"
// 	"strings"
// )

// func Index(w http.ResponseWriter, r *http.Request) {
// 	tpl, err := template.ParseFiles("src/admin/indexAdmin.html")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	tpl.Execute(w, nil)
// }

// func NewCard(w http.ResponseWriter, r *http.Request) {
// 	tpl, err := template.ParseFiles("src/admin/newCard.html")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	tpl.Execute(w, nil)
// }

// func ShowUser(w http.ResponseWriter, r *http.Request) {
// 	tpl, err := template.ParseFiles("src/admin/userTable.html")
// 	usersInfos := logic.GetUsersInfos()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	tpl.Execute(w, usersInfos)
// }

// func ProcessCard(w http.ResponseWriter, r *http.Request) {
// 	r.ParseMultipartForm(10 << 20)

// 	cardName := r.PostFormValue("cardName")
// 	cardRarity, err := strconv.Atoi(r.PostFormValue("rarity"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	cardWeight, err := strconv.Atoi(r.PostFormValue("weight"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	pathToFile := strings.ToLower(strings.ReplaceAll(cardName, " ", "_"))
// 	file, _, err := r.FormFile("cardImage")
// 	if err != nil {
// 		fmt.Println("Error Retrieving the File")
// 		fmt.Println(err)
// 		return
// 	}
// 	defer file.Close()
// 	tempFile, err := ioutil.TempFile("temp-pics", pathToFile+"-*.png")
// 	fmt.Println(tempFile.Name())
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer tempFile.Close()

// 	fileBytes, err := ioutil.ReadAll(file)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	tempFile.Write(fileBytes)
// 	err = os.Rename(tempFile.Name(), "ressources/pics/"+pathToFile+".png")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	addNewCard(cardName, cardRarity, cardWeight, pathToFile)
// 	http.Redirect(w, r, "/admin/new_card", http.StatusFound)
// }
