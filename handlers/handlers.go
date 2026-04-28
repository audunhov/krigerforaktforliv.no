package handlers

import (
	"html/template"
	"net/http"
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
        errorMessage = "This email has already signed the petition."
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
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")

	err = h.store.AddSignature(r.Context(),	name, email)
	if err != nil {
		http.Redirect(w, r, "/?error=duplicate", http.StatusSeeOther)
    return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
