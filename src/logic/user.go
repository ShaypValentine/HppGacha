package logic

import (
	models "hppGacha/src/models"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Credentials struct {
	Id       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type session struct {
	id     int
	expiry time.Time
}

var sessions = map[string]session{}

func addToInventory(connectedUser models.User, entry Entry) {

	var err error
	if cardExistInInventory(connectedUser.ID, entry.Card.ID) {
		err = DB.Model(&models.CardInInventory{}).Where(" user_id=? AND card_id=?", connectedUser.ID, entry.Card.ID).Update("quantity", gorm.Expr("quantity  + ?", 1)).Error
	} else {
		newCard := models.CardInInventory{UserID: connectedUser.ID, CardID: entry.Card.ID, Quantity: 1}
		err = DB.Create(&newCard).Error
	}
	if err != nil {
		log.Panic(err)
	}
}

func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

func cardExistInInventory(userID uint, cardID uint) bool {
	var isCardInInventory bool
	err := DB.Model(&models.CardInInventory{}).Select("count(*) > 0").Where("user_id = ? AND card_id = ?", userID, cardID).Find(&isCardInInventory).Error
	if err != nil {
		log.Panic(err)
	}
	return isCardInInventory
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
	expiresAt := time.Now().Add(240 * time.Hour)

	sessions[newSessionToken] = session{
		id:     userSession.id,
		expiry: expiresAt,
	}

	delete(sessions, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   newSessionToken,
		Expires: time.Now().Add(240 * time.Hour),
	})
}

func getTopCollectors(infos IndexInfo) IndexInfo {
	var topCollector TopCollector
	var maxCards int64
	DB.Model(&models.Card{}).Count(&maxCards)
	infos.MaxCards = maxCards
	rows, err := DB.Raw("SELECT username, count(card_id) FROM users INNER JOIN cards_in_inventory ON cards_in_inventory.user_id = users.id GROUP BY username ORDER BY count(card_id) DESC LIMIT 10 ").Rows()
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&topCollector.Username, &topCollector.UniqueCard)
		infos.TopCollectors = append(infos.TopCollectors, topCollector)
	}
	return infos
}
