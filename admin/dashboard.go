package admin

import (
	"html/template"
	"net/http"
	"database/sql"
	
	"gocart/config"
)

func dashboard(cfg *config.Config, db *sql.DB, tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		renderPage(w, tmpl, "dashboard", nil)
    }
}
