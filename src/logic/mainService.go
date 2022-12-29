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
var templatePath = "src/views/"

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

type RolledCard struct {
	Card models.Card
	User models.User
}

func Index(w http.ResponseWriter, r *http.Request) {
	var indexInfos IndexInfo
	tpl, err := template.ParseFiles(
		templatePath+"index.html",
		templatePath+"navbar.html",
		templatePath+"_parts/head.html",
		templatePath+"_parts/footer.html",
		templatePath+"_parts/js.html")
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
	var rolledCard RolledCard
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
				rolledCard.User = connectedUser
				rolledCard.Card = rolledItem.Card
				err = tpl.Execute(w, rolledCard)
				if err != nil {
					log.Panicln(err)
				}
			}
		} else {
			rolledCard.User = connectedUser
			rolledCard.Card = rolledItem.Card
			err = tpl.Execute(w, rolledCard)
			if err != nil {
				log.Panicln(err)
			}
		}
	}

}

func LoginPageHandler(w http.ResponseWriter, request *http.Request) {
	errorGet := request.URL.Query().Get("error")
	errorText := ErrorString[errorGet]
	tpl, err := template.ParseFiles(templatePath+"loginForm.html",
		templatePath+"navbar.html",
		templatePath+"_parts/head.html",
		templatePath+"_parts/footer.html",
		templatePath+"_parts/js.html")
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
	tpl, err := template.ParseFiles(templatePath+"inscriptionForm.html",
		templatePath+"navbar.html",
		templatePath+"_parts/head.html",
		templatePath+"_parts/footer.html",
		templatePath+"_parts/js.html")
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
		DB.Preload("ShadowPortal").Preload("CardsInInventory.Card").First(&connectedUser, userSession.id)
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
	expiresAt := time.Now().Add(240 * time.Hour)
	sessions[sessionToken] = session{
		id:     int(user.ID),
		expiry: expiresAt,
	}
	var shadowPortal models.ShadowPortal
	err = DB.Where("user_id=?", user.ID).Assign(models.ShadowPortal{UserID: user.ID}).FirstOrCreate(&shadowPortal).Error
	if err != nil {
		log.Panic(err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func Disconnect(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
	sessionToken := c.Value
	// remove the users session from the session map
	delete(sessions, sessionToken)
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
