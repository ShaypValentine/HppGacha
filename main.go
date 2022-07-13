package main

import (
	"fmt"
	admin "hppGacha/admin"
	logic "hppGacha/logic"
	"log"
	"net/http"
"strings"
)

func main() {
	db, err := logic.DatabaseConnection()
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	logic.DataToRoll(db)
http.FileServer(neuteredFileSystem{http.Dir("./ressources")})
       // fs := http.FileServer(http.Dir("ressources"))
	http.Handle("/ressources/", http.StripPrefix("/ressources",fileServer))

	fileServer := 
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

type neuteredFileSystem struct {
    fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
    f, err := nfs.fs.Open(path)
    if err != nil {
        return nil, err
    }

    s, err := f.Stat()
    if s.IsDir() {
        index := filepath.Join(path, "index.html")
        if _, err := nfs.fs.Open(index); err != nil {
            closeErr := f.Close()
            if closeErr != nil {
                return nil, closeErr
            }

            return nil, err
        }
    }

    return f, nil
}
