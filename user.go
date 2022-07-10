package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Id       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type session struct {
	username string
	id       int
	expiry   time.Time
}

type user struct {
	Username string `json:"username" db:"username"`
	Id       int    `json:"id" db:"id"`
}

var sessions = map[string]session{}

func addToInventory(connectedUser user, card BannerRoll) {
	cardAlreadyExist := cardExistInInventory(DB, connectedUser.Id, card.Name)
	if cardAlreadyExist {
		DB.Exec("UPDATE inventory SET quantity = quantity + 1 where user = ? and cardName = ?", connectedUser.Id, card.Name)
	} else {
		DB.Exec("INSERT INTO inventory (user, cardName,quantity) VALUES (?, ?,1)", connectedUser.Id, card.Name)
	}
}

func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

func signup(w http.ResponseWriter, r *http.Request) {
	// Parse and decode the request body into a new `Credentials` instance
	creds := &Credentials{}
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	creds.Username = r.PostFormValue("username")
	creds.Password = r.PostFormValue("password")
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	if err != nil {
		log.Fatal(err)
	}
	// Next, insert the username, along with the hashed password into the database
	if _, err = DB.Exec("INSERT into users (userName,password) values ($1, $2)", creds.Username, string(hashedPassword)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)

}

func signin(w http.ResponseWriter, r *http.Request) {
	// Parse and decode the request body into a new `Credentials` instance

	creds := &Credentials{}
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	creds.Username = r.PostFormValue("username")
	creds.Password = r.PostFormValue("password")
	// Get the existing entry present in the database for the given username
	result := DB.QueryRow("select id,password from users where username=$1", creds.Username)
	if err != nil {
		// If there is an issue with the database, return a 500 error
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// We create another instance of `Credentials` to store the credentials we get from the database
	storedCreds := &Credentials{}
	// Store the obtained password in `storedCreds`
	err = result.Scan(&storedCreds.Id, &storedCreds.Password)
	if err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/login?error=NoUser", 302)
			return
		}
		// If the error is of any other type, send a 500 status
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		http.Redirect(w, r, "/login?error=BadCreds", 302)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)
	sessions[sessionToken] = session{
		username: creds.Username,
		id:       storedCreds.Id,
		expiry:   expiresAt,
	}
	// If we reach this point, that means the users password was correct, and that they are authorized
	// The default 200 status is sent

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
	http.Redirect(w, r, "/", http.StatusFound)

}

func cardExistInInventory(db *sql.DB, userId int, cardname string) bool {
	sqlStmt := `SELECT cardName FROM inventory WHERE user = ? AND cardName = ?`
	err := db.QueryRow(sqlStmt, userId, cardname).Scan(&cardname)
	if err != nil {
		if err != sql.ErrNoRows {
			// a real error happened! you should change your function return
			// to "(bool, error)" and return "false, err" here
			log.Print(err)
		}

		return false
	}

	return true
}
