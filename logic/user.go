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

type UserInfo struct {
	Username string `json:"username" db:"username"`
	Id       int    `json:"id" db:"id"`
	Rolls    int    `json:"rolls" db:"rolls"`
}

type UsersInfos struct {
	UserInfo []UserInfo
}

var sessions = map[string]session{}

func addToInventory(connectedUser UserInfo, card BannerRoll) {
	cardAlreadyExist := cardExistInInventory(DB, connectedUser.Id, card.Name)
	if cardAlreadyExist {
		DB.Exec("UPDATE inventory SET quantity = quantity + 1 where user = ? and cardName = ?", connectedUser.Id, card.Name)
	} else {
		DB.Exec("INSERT INTO inventory (user, cardName,quantity) VALUES (?, ?,1)", connectedUser.Id, card.Name)
	}
}

func canRoll(user UserInfo) bool {
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
	expiresAt := time.Now().Add(600 * time.Second)

	sessions[newSessionToken] = session{
		username: userSession.username,
		id:       userSession.id,
		expiry:   expiresAt,
	}

	delete(sessions, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   newSessionToken,
		Expires: time.Now().Add(600 * time.Second),
	})
}

func consumeRoll(user UserInfo) {
	_, err := DB.Exec("UPDATE users SET availableRoll = availableRoll - 1 WHERE id = ?", user.Id)
	if err != nil {
		log.Fatal(err)
	}
}

func GetUsersInfos() (infos UsersInfos) {
	var userInfo UserInfo
	rows, err := DB.Query("SELECT id, username,availableRoll FROM users")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		rows.Scan(&userInfo.Id, &userInfo.Username, &userInfo.Rolls)
		infos.UserInfo = append(infos.UserInfo, userInfo)
	}
	return infos
}

func getRollsForUser(user UserInfo) int {
	var nbrRoll int
	err := DB.QueryRow("SELECT availableRoll FROM users WHERE id = ?", user.Id).Scan(&nbrRoll)
	if err != nil {
		log.Fatal(err)
	}

	return nbrRoll
}

func getTopCollectors(infos IndexInfo) IndexInfo {
	var topCollector TopCollector
	var maxCards int
	err := DB.QueryRow("SELECT count(name) FROM rollable_users").Scan(&maxCards)
	if err != nil {
		log.Fatal(err)
	}
	infos.MaxCards = maxCards
	rows, err := DB.Query("SELECT username, count(cardName) FROM users INNER JOIN inventory ON inventory.user = users.id GROUP BY username ORDER BY count(cardName) DESC LIMIT 10 ")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		rows.Scan(&topCollector.Username, &topCollector.UniqueCard)
		infos.TopCollectors = append(infos.TopCollectors, topCollector)
	}
return infos
}
