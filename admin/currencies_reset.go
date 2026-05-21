package admin

import (
	"html/template"
	"net/http"
	"database/sql"

	"gocart/config"
	"gocart/seeds"
)

func currenciesResetPost(cfg *config.Config, db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := seeds.SeedCurrencies(db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/currencies", http.StatusSeeOther)
	}
}
