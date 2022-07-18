package logic

import (
	"database/sql"
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

//  func Roll(w http.ResponseWriter, r *http.Request) {
//  	tpl, err := template.ParseFiles("src/rollCard.html")
//  	if err != nil {
//  		log.Fatalln(err)
//  	}
//  	if r.Method == "GET" {
//  		w.Header().Set("Content-Type", "application/json")
//  		connectedUser, exists := getConnectedUser(w, r)
//  		if exists {
//  			// Refresh(w, r)
//  		}
//  		if canRoll(connectedUser) {
//  			rolledItem := getRandom()
//  			consumeRoll(connectedUser)
//  			addToInventory(connectedUser, rolledItem)
//  			if rolledItem.Rarity == 2 {
//  				tpl, err = template.ParseFiles("src/rollRareCard.html")
//  				if err != nil {
//  					log.Fatalln(err)
//  				}
//  			}

//  			err = tpl.Execute(w, rolledItem)
//  			if err != nil {
//  				log.Fatalln(err)
//  			}
//  		}
//  	}

//  }

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
	tpl, err := template.ParseFiles("src/index.html")
	if err != nil {
		log.Fatalln(err)
	}
	connectedUser, exists := getConnectedUser(w, r)
	indexInfos.User = connectedUser
	if exists {
		indexInfos.Inventory = getInventoryForUser(connectedUser)
	}

	indexInfos = getTopCollectors(indexInfos)

	err = tpl.Execute(w, indexInfos)
	if err != nil {
		log.Fatalln(err)
	}
}

func LoginPageHandler(w http.ResponseWriter, request *http.Request) {
	error := request.URL.Query().Get("error")
	errorText := ""
	if error != "" {
		errorText = "No account with this username was found."
		if error == "BadCreds" {
			errorText = "wrong password or username."
		}
	}
	tpl, err := template.ParseFiles("src/loginForm.html")
	if err != nil {
		log.Fatalln(err)
	}
	err = tpl.Execute(w, errorText)
	if err != nil {
		log.Fatalln(err)
	}
}
func InscriptionPageHandler(w http.ResponseWriter, request *http.Request) {
	error := request.URL.Query().Get("error")
	errorText := ""
	if error != "" {
		errorText = "No account with this username was found."
		if error == "BadCreds" {
			errorText = "Empty Username or password "
		}
	}
	tpl, err := template.ParseFiles("src/inscriptionForm.html")
	if err != nil {
		log.Fatalln(err)
	}
	err = tpl.Execute(w, errorText)
	if err != nil {
		log.Fatalln(err)
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
		DB.First(&connectedUser, userSession.id)
		return connectedUser, exists
	}
	return connectedUser, false
}

func Signup(w http.ResponseWriter, r *http.Request) {
	creds := &Credentials{}
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
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
			log.Fatal(err)
		}
		var newUser = models.User{Username: creds.Username, Password: string(hashedPassword)}
		if result := DB.Create(&newUser); result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.Redirect(w, r, "/inscription?error=BadCreds", http.StatusFound)
		}
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	http.Redirect(w, r, "/inscription?error=BadCreds", http.StatusFound)
	return
}

func Signin(w http.ResponseWriter, r *http.Request) {
	creds := &Credentials{}
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	creds.Username = strings.ToLower(r.PostFormValue("username"))
	creds.Password = r.PostFormValue("password")
	var user models.User
	if result := DB.Where("username = ?", creds.Username).First(&user); result.Error != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err != nil {
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/login?error=NoUser", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		http.Redirect(w, r, "/login?error=BadCreds", http.StatusFound)
		w.WriteHeader(http.StatusUnauthorized)
		return
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
