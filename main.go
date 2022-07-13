package main

import (
	"fmt"
	admin "hppGacha/admin"
	logic "hppGacha/logic"
	"log"
	"net/http"
)

func main() {
	db, err := logic.DatabaseConnection()
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	logic.DataToRoll(db)

	http.Handle("/ressources/", http.StripPrefix("/ressources/", http.FileServer(http.Dir("./ressources"))))

	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		logic.EmptyEntries()
		logic.DataToRoll(db)
		http.Redirect(w, r, "/", http.StatusAccepted)
	})
	http.HandleFunc("/signup", logic.Signup)
	http.HandleFunc("/signin", logic.Signin)
	http.HandleFunc("/roll", logic.Roll)
	http.HandleFunc("/login", logic.LoginPageHandler)
	http.HandleFunc("/inscription", logic.InscriptionPageHandler)
	http.HandleFunc("/", logic.Index)
	http.HandleFunc("/adminw", admin.Index)
	http.HandleFunc("/admin/new_card", admin.NewCard)
	http.HandleFunc("/admin/show_users", admin.ShowUser)
	http.HandleFunc("/admin/process_card", admin.ProcessCard)
	// Launch app on OS PORT var or 8008
	if err := http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/hppgacha.art/fullchain.pem", "/etc/letsencrypt/live/hppgacha.art/privkey.pem", nil); err != nil {
		log.Fatal(err)
	}
}
