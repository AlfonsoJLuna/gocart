package admin

import (
	"html/template"
	"net/http"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"

	"gocart/config"
	"gocart/models"
)

type currenciesEditData struct {
    Currency		*models.Currency
	OriginalName	string
	Error			string
	Success			string
}

func currenciesEdit(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data currenciesEditData

		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		data.Currency, err = models.CurrencyReadByID(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		data.OriginalName = data.Currency.Name

		renderPage(w, tmpl, "currencies_edit", data)
	}
}

func currenciesEditPost(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data currenciesEditData

		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		data.Currency, err = models.CurrencyReadByID(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		data.OriginalName = data.Currency.Name

		data.Currency.Name      = r.FormValue("name")
		data.Currency.NameAlt   = r.FormValue("name_alt")
		data.Currency.IsEnabled	= r.FormValue("is_enabled") == "on"

		if err := models.CurrencyUpdate(db, data.Currency); err != nil {
			data.Error = friendlyError(err)
		} else {
			data.OriginalName = data.Currency.Name
			data.Success = "Currency saved successfully."
		}
		
		renderPage(w, tmpl, "currencies_edit", data)
	}
}
