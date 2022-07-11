package logic

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
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

func canRoll(user user) bool {
	var nbrRoll int
	err := DB.QueryRow("SELECT availableRoll FROM users WHERE id = ?", user.Id).Scan(&nbrRoll)
	if err != nil {
		// Guest mode
		nbrRoll = 1
	}

	return (nbrRoll > 0)
}

func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

func cardExistInInventory(db *sql.DB, userId int, cardname string) bool {
	sqlStmt := `SELECT cardName FROM inventory WHERE user = ? AND cardName = ?`
	err := db.QueryRow(sqlStmt, userId, cardname).Scan(&cardname)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return false
	}
	return true
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value
	userSession, exists := sessions[sessionToken]
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newSessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	sessions[newSessionToken] = session{
		username: userSession.username,
		id:       userSession.id,
		expiry:   expiresAt,
	}

	delete(sessions, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   newSessionToken,
		Expires: time.Now().Add(120 * time.Second),
	})
}

func consumeRoll(user user) {
	_, err := DB.Exec("UPDATE users SET availableRoll = availableRoll - 1 WHERE id = ?", user.Id)
	if err != nil {
		log.Fatal(err)
	}
}
