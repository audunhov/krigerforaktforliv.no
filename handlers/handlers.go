package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"krigerforaktforliv.no/db"
)

type Handler struct {
	store *db.Store
	tmpl  *template.Template
}

func NewHandler(s *db.Store, t *template.Template) *Handler {
	return &Handler{
		store: s,
		tmpl:  t,
	}
}

type PageData struct {
    Count int
		Error string
}

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
    count, err := h.store.GetSignatureCount(r.Context())
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }

		errorMessage := ""
    if r.URL.Query().Get("error") == "duplicate" {
        errorMessage = "Denne eposten har allerede vært brukt"
    }
    if r.URL.Query().Get("error") == "invalid" {
        errorMessage = "Skjemaet er ugyldig"
    }

    data := PageData{
        Count: count,
				Error: errorMessage,
    }

    err = h.tmpl.ExecuteTemplate(w, "index.html", data)
    if err != nil {
        http.Error(w, "Template error", http.StatusInternalServerError)
    }
}


func (h *Handler) SignHandler(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Redirect(w, r, "/?error=invalid", http.StatusSeeOther)
        return
    }

    // Honeypot check
    if r.FormValue("zip_code") != "" {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    name := strings.TrimSpace(r.FormValue("name"))
		var email *string

		formEmail := strings.TrimSpace(r.FormValue("email"))
		if formEmail != "" {
			email = &formEmail
		}

    if name == "" {
        http.Redirect(w, r, "/?error=invalid", http.StatusSeeOther)
        return
    }

    err := h.store.AddSignature(r.Context(), name, email)
    if err != nil {
        http.Redirect(w, r, "/?error=duplicate", http.StatusSeeOther)
        return
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}
