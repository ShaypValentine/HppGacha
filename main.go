package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

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

	http.HandleFunc("/roll", func(w http.ResponseWriter, r *http.Request) {
		tpl, err := template.ParseFiles("src/rollCard.gohtml")
		if err != nil {
			log.Fatalln(err)
		}
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			rolledItem := getRandom()
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl, err := template.ParseFiles("src/index.gohtml")

		if err != nil {
			log.Fatalln(err)
		}
		err = tpl.Execute(w, getRandom())
		if err != nil {
			log.Fatalln(err)
		}
	})

	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
