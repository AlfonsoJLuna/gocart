package admin

import (
	"html/template"
    "net/http"
	"database/sql"
	
	"gocart/config"
	"gocart/seeds"
)

func countriesResetPost(cfg *config.Config, db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := seeds.SeedCountries(db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/countries", http.StatusSeeOther)
	}
}
