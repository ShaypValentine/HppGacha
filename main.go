package main

import (
	"log"
	"net/http"
	"text/template"
)

func main() {
	fileToLines("web/banner_content.csv")

	http.Handle("/ressources/", http.StripPrefix("/ressources/", http.FileServer(http.Dir("./ressources"))))

	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		emptyEntries()
		fileToLines("web/banner_content.csv")
		http.Redirect(w, r, "/", http.StatusAccepted)
	})

	http.HandleFunc("/roll", func(w http.ResponseWriter, r *http.Request) {
		tpl, err := template.ParseFiles("src/rollCard.gohtml")
		if err != nil {
			log.Fatalln(err)
		}
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			err = tpl.Execute(w, getRandom())
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
