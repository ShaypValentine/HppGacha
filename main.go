package main

import (
	"fmt"
	logic "hppGacha/logic"
	"log"
	"net/http"
	"os"
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
	// Launch app on OS PORT var or 80
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
