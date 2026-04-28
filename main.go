package main

import (
	"embed"
	"html/template"
	"krigerforaktforliv.no/db"
	"krigerforaktforliv.no/handlers"
	"log"
	"net/http"
)

//go:embed pages/*.html
var templateFS embed.FS

func main() {
	// 1. Load embedded templates
	tmpl, err := template.ParseFS(templateFS, "pages/*.html")
	if err != nil {
		log.Fatal(err)
	}

	// 2. Setup DB
	store, err := db.NewStore("petition.db")
	if err != nil {
		log.Fatal(err)
	}

	// 3. Initialize Handlers with dependencies
	h := handlers.NewHandler(store, tmpl)

	// 4. Setup Routes
	http.HandleFunc("/", h.IndexHandler)
	http.HandleFunc("/sign", h.SignHandler)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
