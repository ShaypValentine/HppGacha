package main

import (
	admin "hppGacha/src/admin"
	logic "hppGacha/src/logic"
	models "hppGacha/src/models"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	db, err := logic.DatabaseConnection()
	if err != nil {
		log.Panic(err)
	}
	logic.DB = db
	admin.DB = db
	db.AutoMigrate(&models.Card{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.CardInInventory{})
	db.AutoMigrate(&models.ShadowPortal{})
	db.AutoMigrate(&models.Banner{})

	logic.DataToRoll(db)
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("src/ressources")})
	http.Handle("/ressources/", http.StripPrefix("/ressources", fileServer))

	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		logic.EmptyEntries()
		logic.DataToRoll(db)
		http.Redirect(w, r, "/", http.StatusFound)
	})
	http.HandleFunc("/recycle", logic.RecycleCard)
	http.HandleFunc("/sacrifice", logic.SacrificeCard)
	http.HandleFunc("/inventory", logic.ShowInventory)
	http.HandleFunc("/signup", logic.Signup)
	http.HandleFunc("/signin", logic.Signin)
	http.HandleFunc("/roll", logic.Roll)
	http.HandleFunc("/shadowRoll", logic.ShadowRoll)
	http.HandleFunc("/login", logic.LoginPageHandler)
	http.HandleFunc("/inscription", logic.InscriptionPageHandler)
	http.HandleFunc("/", logic.Index)
	http.HandleFunc("/adminw", admin.Index)
	http.HandleFunc("/admin/new_card", admin.NewCard)
	http.HandleFunc("/admin/edit_card", admin.EditCard)
	http.HandleFunc("/admin/show_users", admin.ShowUser)
	http.HandleFunc("/admin/show_cards", admin.ShowCards)
	http.HandleFunc("/admin/process_card", admin.ProcessCard)
	http.HandleFunc("/admin/process_card_edit", admin.ProcessCardEdit)
	http.HandleFunc("/shadow/portal", logic.ShadowIndex)
	http.HandleFunc("/disconnect", logic.Disconnect)
	http.HandleFunc("/banner_list", logic.BannerList)

	// Launch app on OS PORT var or 8008
	env := os.Getenv("LOCALENV")
	if env == "prod" {
		if err := http.ListenAndServeTLS(":8443", "persistent/server.crt", "persistent/server.key", nil); err != nil {
			log.Panic(err)
		}
	} else {
		if err := http.ListenAndServe(":80", nil); err != nil {
			log.Panic(err)
		}
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
	if err != nil {
		return nil, err
	}
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
