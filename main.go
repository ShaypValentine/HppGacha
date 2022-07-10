package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

var DB *sql.DB

func main() {
	db, err := databaseConnection()
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	dataToRoll(db)

	http.Handle("/ressources/", http.StripPrefix("/ressources/", http.FileServer(http.Dir("./ressources"))))

	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		emptyEntries()
		dataToRoll(db)
		http.Redirect(w, r, "/", http.StatusAccepted)
	})
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/signin", signin)
	http.HandleFunc("/roll", func(w http.ResponseWriter, r *http.Request) {
		tpl, err := template.ParseFiles("src/rollCard.gohtml")
		if err != nil {
			log.Fatalln(err)
		}
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			rolledItem := getRandom()
			connectedUser := getConnectedUser(r)
			addToInventory(connectedUser, rolledItem)
			if rolledItem.Rarity == 2 {
				tpl, err = template.ParseFiles("src/rollRareCard.gohtml")
				if err != nil {
					log.Fatalln(err)
				}
			}
			err = tpl.Execute(w, rolledItem)
			if err != nil {
				log.Fatalln(err)
			}
		}

	})

	http.HandleFunc("/login", loginPageHandler)
	http.HandleFunc("/inscription", inscriptionPageHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl, err := template.ParseFiles("src/index.gohtml")
		if err != nil {
			log.Fatalln(err)
		}
		connectedUser := getConnectedUser(r)
		err = tpl.Execute(w, connectedUser)
		if err != nil {
			log.Fatalln(err)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func loginPageHandler(w http.ResponseWriter, request *http.Request) {
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
func inscriptionPageHandler(w http.ResponseWriter, request *http.Request) {
	tpl, err := template.ParseFiles("src/inscriptionForm.html")
	if err != nil {
		log.Fatalln(err)
	}
	err = tpl.Execute(w, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func getConnectedUser(r *http.Request) user {
	var connectedUser user
	sessionCookie, err := r.Cookie("session_token")
	if err == nil {
		token := sessionCookie.Value
		userSession, exists := sessions[token]
		if !exists {
			log.Println("user n'existe pas")
		}
		if userSession.isExpired() {
			delete(sessions, token)
		}
		connectedUser.Username = userSession.username
		connectedUser.Id = userSession.id
	}
	return connectedUser
}
