package admin

import (
	"net/http"
	"html/template"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"

	"gocart/config"
	"gocart/models"
	"gocart/seeds"
)

func getCurrencies(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currencies, err := models.CurrencyListAll(db, 0, 0, false)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render(w, tmpl, "currencies", currencies)
	}
}

func getCurrenciesEdit(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		c, err := models.CurrencyReadByID(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		render(w, tmpl, "currencies-edit", editData{Data: c})
	}
}

func postCurrenciesEdit(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		c, err := models.CurrencyReadByID(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		c.Name      = r.FormValue("name")
		c.NameAlt   = r.FormValue("name_alt")
		c.IsEnabled = r.FormValue("is_enabled") == "on"

		if err := models.CurrencyUpdate(db, c); err != nil {
			render(w, tmpl, "currencies-edit", editData{Data: c, Error: friendlyError(err)})
			return
		}

		render(w, tmpl, "currencies-edit", editData{Data: c, Success: "Saved successfully."})
	}
}

func postCurrenciesReset(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := seeds.SeedCurrencies(db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/currencies", http.StatusSeeOther)
	}
}
