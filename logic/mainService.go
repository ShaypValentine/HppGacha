package logic

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func Roll(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("src/rollCard.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		connectedUser, exists := getConnectedUser(w, r)
		if canRoll(connectedUser) {
			rolledItem := getRandom()
			consumeRoll(connectedUser)
			addToInventory(connectedUser, rolledItem)
			if rolledItem.Rarity == 2 {
				tpl, err = template.ParseFiles("src/rollRareCard.gohtml")
				if err != nil {
					log.Fatalln(err)
				}
			}
			if exists {
				Refresh(w, r)
			}
			err = tpl.Execute(w, rolledItem)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

}

type IndexInfo struct {
	User      user
	Inventory Inventory
}

func Index(w http.ResponseWriter, r *http.Request) {
	var indexInfos IndexInfo
	tpl, err := template.ParseFiles("src/index.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	connectedUser, exists := getConnectedUser(w, r)
	indexInfos.User = connectedUser
	if exists {
		indexInfos.Inventory = getInventoryForUser(connectedUser)
	}

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
	tpl, err := template.ParseFiles("src/inscriptionForm.html")
	if err != nil {
		log.Fatalln(err)
	}
	err = tpl.Execute(w, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func getConnectedUser(w http.ResponseWriter, r *http.Request) (user, bool) {
	var connectedUser user
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
		connectedUser.Username = userSession.username
		connectedUser.Id = userSession.id
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
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = DB.Exec("INSERT into users (userName,password) values ($1, $2)", creds.Username, string(hashedPassword)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)

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

	result := DB.QueryRow("select id,password from users where username=$1", creds.Username)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	storedCreds := &Credentials{}
	err = result.Scan(&storedCreds.Id, &storedCreds.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/login?error=NoUser", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		http.Redirect(w, r, "/login?error=BadCreds", http.StatusFound)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(600 * time.Second)
	sessions[sessionToken] = session{
		username: creds.Username,
		id:       storedCreds.Id,
		expiry:   expiresAt,
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
	http.Redirect(w, r, "/", http.StatusFound)

}
