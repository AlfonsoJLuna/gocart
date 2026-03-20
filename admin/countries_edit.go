package admin

import (
	"html/template"
	"net/http"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"

	"gocart/config"
	"gocart/models"
)

type countriesEditData struct {
	Country    		*models.Country
	OriginalName	string
	Currencies 		[]*models.Currency
	Error      		string
	Success    		string
}

func countriesEdit(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data countriesEditData

		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		data.Country, err = models.CountryReadByID(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		data.OriginalName = data.Country.Name

		currencies, err := models.CurrencyListAll(db, 0, 0, false)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, c := range currencies {
			if c.IsEnabled {
				data.Currencies = append(data.Currencies, c)
			}
		}

		renderPage(w, tmpl, "countries_edit", data)
	}
}

func countriesEditPost(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data countriesEditData

		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		data.Country, err = models.CountryReadByID(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		data.OriginalName = data.Country.Name

		currencies, err := models.CurrencyListAll(db, 0, 0, false)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, c := range currencies {
			if c.IsEnabled {
				data.Currencies = append(data.Currencies, c)
			}
		}

		data.Country.Name = r.FormValue("name")
		data.Country.NameAlt = r.FormValue("name_alt")
		data.Country.CurrencyISOCode = r.FormValue("currency_iso_code")
		data.Country.VATRate = parseFloat(r.FormValue("vat_rate"))
		data.Country.IsEnabled = r.FormValue("is_enabled") == "on"

		if err := models.CountryUpdate(db, data.Country); err != nil {
			data.Error = friendlyError(err)
		} else {
			data.OriginalName = data.Country.Name
			data.Success = "Country saved successfully."
		}

		renderPage(w, tmpl, "countries_edit", data)
	}
}
