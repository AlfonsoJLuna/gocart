package admin

import (
	"html/template"
	"net/http"
	"database/sql"

	"gocart/config"
	"gocart/services"
)

func currenciesList(cfg *config.Config, db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currencies, err := services.CurrencyListAll(db, 0, -1, false)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		renderPage(w, tmpl, "currencies_list", currencies)
	}
}
