package logic

import (
	"html/template"
	"log"
	"net/http"
)

func ShadowIndex(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles(
		templatePath+"shadowPortal.html",
		templatePath+"navbar.html",
		templatePath+"_parts/head.html",
		templatePath+"_parts/footer.html",
		templatePath+"_parts/js.html")
	if err != nil {
		log.Panic(err)
	}
	connectedUser, exists := getConnectedUser(w, r)
	if exists {
		err = tpl.Execute(w, connectedUser)
		if err != nil {
			log.Panicln(err)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
