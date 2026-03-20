package admin

import (
	"html/template"
	"net/http"

	"go.etcd.io/bbolt"

	"gocart/config"
	"gocart/models"
)

func currenciesList(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currencies, err := models.CurrencyListAll(db, 0, 0, false)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		renderPage(w, tmpl, "currencies_list", currencies)
	}
}
