package admin

import (
	"html/template"
    "net/http"

	"go.etcd.io/bbolt"
	
	"gocart/config"
	"gocart/models"
)

func countriesList(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		countries, err := models.CountryListAll(db, 0, 0, false)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		renderPage(w, tmpl, "countries_list", countries)
	}
}
