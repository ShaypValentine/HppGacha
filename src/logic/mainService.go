package logic

import (
	models "hppGacha/src/models"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var DB *gorm.DB
var ErrorString = map[string]string{
	"BadCreds":  "Wrong password or username",
	"UserExist": "An user already exist with this username",
	"NoUser":    "No account exists with this username",
}

type IndexInfo struct {
	User          models.User
	Inventory     models.Inventory
	TopCollectors []TopCollector
	MaxCards      int64
}

type TopCollector struct {
	Username   string `json:"username" db:"username"`
	UniqueCard int    `json:"uniqueCard" db:"uniqueCard"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	var indexInfos IndexInfo
	tpl, err := template.ParseFiles("src/views/index.html")
	if err != nil {
		log.Panicln(err)
	}
	connectedUser, _ := getConnectedUser(w, r)
	indexInfos.User = connectedUser
	indexInfos = getTopCollectors(indexInfos)

	err = tpl.Execute(w, indexInfos)
	if err != nil {
		log.Panicln(err)
	}
}

func Roll(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("src/views/rollCard.html")
	if err != nil {
		log.Panicln(err)
	}
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		rolledItem := getRandom()
		connectedUser, exists := getConnectedUser(w, r)
		if exists {
			Refresh(w, r)
			if connectedUser.AvailableRolls > 0 {
				connectedUser.AvailableRolls--
				DB.Save(&connectedUser)
				addToInventory(connectedUser, rolledItem)
			}
		}

		err = tpl.Execute(w, rolledItem.Card)
		if err != nil {
			log.Panicln(err)
		}

	}

}

func LoginPageHandler(w http.ResponseWriter, request *http.Request) {
	errorGet := request.URL.Query().Get("error")
	errorText := ErrorString[errorGet]
	tpl, err := template.ParseFiles("src/views/loginForm.html")
	if err != nil {
		log.Panic(err)
	}
	err = tpl.Execute(w, errorText)
	if err != nil {
		log.Panicln(err)
	}
}
func InscriptionPageHandler(w http.ResponseWriter, request *http.Request) {
	errorGet := request.URL.Query().Get("error")
	errorText := ErrorString[errorGet]
	tpl, err := template.ParseFiles("src/views/inscriptionForm.html")
	if err != nil {
		log.Panicln(err)
	}
	err = tpl.Execute(w, errorText)
	if err != nil {
		log.Panicln(err)
	}
}

func getConnectedUser(w http.ResponseWriter, r *http.Request) (models.User, bool) {
	var connectedUser models.User
	sessionCookie, err := r.Cookie("session_token")
	if err == nil {
		token := sessionCookie.Value
		userSession, exists := sessions[token]
		if !exists {
			return connectedUser, exists
		}
		if userSession.isExpired() {
			delete(sessions, token)
		}
		DB.Preload("CardsInInventory.Card").First(&connectedUser, userSession.id)
		return connectedUser, exists
	}
	return connectedUser, false
}

func Signup(w http.ResponseWriter, r *http.Request) {
	creds := &Credentials{}
	err := r.ParseForm()
	if err != nil {
		log.Panic(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	creds.Username = strings.ToLower(r.PostFormValue("username"))
	creds.Password = r.PostFormValue("password")
	if creds.Username != "" && creds.Password != "" {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
		if err != nil {
			log.Panic(err)
		}
		var newUser = models.User{Username: creds.Username, Password: string(hashedPassword)}
		if result := DB.Create(&newUser); result.Error != nil {
			http.Redirect(w, r, "/inscription?error=UserExist", http.StatusFound)
		}
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	http.Redirect(w, r, "/inscription?error=BadCreds", http.StatusFound)
}

func Signin(w http.ResponseWriter, r *http.Request) {
	creds := &Credentials{}
	err := r.ParseForm()
	if err != nil {
		log.Panic(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	creds.Username = strings.ToLower(r.PostFormValue("username"))
	creds.Password = r.PostFormValue("password")
	var user models.User
	err = DB.Where("username = ?", creds.Username).First(&user).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			http.Redirect(w, r, "/login?error=NoUser", http.StatusFound)
		}
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		http.Redirect(w, r, "/login?error=BadCreds", http.StatusFound)
	}
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(600 * time.Second)
	sessions[sessionToken] = session{
		id:     int(user.ID),
		expiry: expiresAt,
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
