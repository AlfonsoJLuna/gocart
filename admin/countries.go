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

type countriesEditData struct {
    Country    *models.Country
    Currencies []*models.Currency
    Error      string
    Success    string
}

func getCountries(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		countries, err := models.CountryListAll(db, 0, 0, false)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render(w, tmpl, "countries", countries)
	}
}

func getCountriesEdit(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id, err := uuid.Parse(r.PathValue("id"))
        if err != nil {
            http.Error(w, "invalid id", http.StatusBadRequest)
            return
        }

        country, err := models.CountryReadByID(db, id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }

        currencies, err := models.CurrencyListAll(db, 0, 0, false)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        var enabled_currencies []*models.Currency
        for _, cur := range currencies {
            if cur.IsEnabled {
                enabled_currencies = append(enabled_currencies, cur)
            }
        }

        render(w, tmpl, "countries-edit", countriesEditData{
            Country:    country,
            Currencies: enabled_currencies,
        })
    }
}

func postCountriesEdit(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id, err := uuid.Parse(r.PathValue("id"))
        if err != nil {
            http.Error(w, "invalid id", http.StatusBadRequest)
            return
        }

        country, err := models.CountryReadByID(db, id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }

        currencies, err := models.CurrencyListAll(db, 0, 0, false)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        var enabled_currencies []*models.Currency
        for _, cur := range currencies {
            if cur.IsEnabled {
                enabled_currencies = append(enabled_currencies, cur)
            }
        }

        country.Name            = r.FormValue("name")
        country.NameAlt         = r.FormValue("name_alt")
        country.CurrencyISOCode = r.FormValue("currency_iso_code")
        country.VATRate         = parseFloat(r.FormValue("vat_rate"))
        country.IsEnabled       = r.FormValue("is_enabled") == "on"

        if err := models.CountryUpdate(db, country); err != nil {
            render(w, tmpl, "countries-edit", countriesEditData{
                Country:    country,
                Currencies: enabled_currencies,
                Error:      friendlyError(err),
            })
            return
        }

        render(w, tmpl, "countries-edit", countriesEditData{
            Country:    country,
            Currencies: enabled_currencies,
            Success:    "Saved successfully.",
        })
    }
}

func postCountriesReset(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := seeds.SeedCountries(db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/countries", http.StatusSeeOther)
	}
}
