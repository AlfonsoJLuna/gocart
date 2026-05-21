package admin
 
import (
	"database/sql"
	"html/template"
	"net/http"
 
	"github.com/google/uuid"

	"gocart/config"
	"gocart/models"
	"gocart/services"
)

type countriesEditData struct {
	Country      	*models.Country
	OriginalName	string
	CurrencyID		string
	Currencies   	[]*models.Currency
	Regions      	[]*models.Region
	Error        	string
	Success      	string
}

func loadCountriesEdit(w http.ResponseWriter, db *sql.DB, r *http.Request) (countriesEditData, error) {
	var data countriesEditData

	countryID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return data, err
	}

	data.Country, err = services.CountryReadByID(db, countryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return data, err
	}
	data.OriginalName = data.Country.Name

	// Flatten CurrencyID for easy comparison in the template.
	if data.Country.CurrencyID != nil {
		data.CurrencyID = data.Country.CurrencyID.String()
	}
	
	currencies, err := services.CurrencyListAll(db, 0, -1, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return data, err
	}
	for _, c := range currencies {
		if c.IsEnabled {
			data.Currencies = append(data.Currencies, c)
		}
	}

	data.Regions, err = services.RegionListByCountryID(db, countryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return data, err
	}

	return data, nil
}

func countriesEdit(cfg *config.Config, db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := loadCountriesEdit(w, db, r)
		if err != nil {
			return
		}

		renderPage(w, tmpl, "countries_edit", data)
	}
}

func countriesEditPost(cfg *config.Config, db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := loadCountriesEdit(w, db, r)
		if err != nil {
			return
		}

		currencyID, err := uuid.Parse(r.FormValue("currency_id"))
		if err != nil {
			http.Error(w, "invalid currency id", http.StatusBadRequest)
			return
		}

		data.Country.Name       = r.FormValue("name")
		data.Country.NameAlt    = r.FormValue("name_alt")
		data.Country.CurrencyID = &currencyID
		data.Country.VATRate    = parseFloat(r.FormValue("vat_rate"))
		data.Country.IsEU       = r.FormValue("is_eu") == "on"
		data.Country.IsEnabled  = r.FormValue("is_enabled") == "on"

		if err := services.CountryUpdate(db, data.Country); err != nil {
			data.Error = friendlyError(err)
		} else {
			data.OriginalName = data.Country.Name
			data.Success = "Country saved successfully."
		}

		renderPage(w, tmpl, "countries_edit", data)
	}
}
