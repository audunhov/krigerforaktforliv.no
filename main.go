package main

import (
	"embed"
	"html/template"
	"krigerforaktforliv.no/db"
	"krigerforaktforliv.no/handlers"
	"log"
	"os"
	"io/fs"
	"net/http"
)

//go:embed pages/*.html
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

func main() {
	tmpl, err := template.ParseFS(templateFS, "pages/*.html")
	if err != nil {
		log.Fatal(err)
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "petition.db"
	}

	store, err := db.NewStore(dbPath)
	if err != nil {
		log.Fatal(err)
	}


	h := handlers.NewHandler(store, tmpl)

	http.HandleFunc("/", h.IndexHandler)
	http.HandleFunc("/sign", h.SignHandler)

	assetsFS, _ := fs.Sub(staticFS, "static")
	fileServer := http.FileServer(http.FS(assetsFS))
	http.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
