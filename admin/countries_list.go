package admin

import (
	"database/sql"
	"html/template"
	"net/http"

	"gocart/config"
	"gocart/models"
	"gocart/services"
)

type countryListItem struct {
	*models.Country
	CurrencyISOCode string
	RegionCount     int
}

func countriesList(cfg *config.Config, db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		countries, err := services.CountryListAll(db, 0, -1, false)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		items := make([]*countryListItem, 0, len(countries))
		for _, c := range countries {
			item := &countryListItem{Country: c}

			if c.CurrencyID != nil {
				currency, err := services.CurrencyReadByID(db, *c.CurrencyID)
				if err == nil {
					item.CurrencyISOCode = currency.ISOCode
				}
			}

			regions, err := services.RegionListByCountryID(db, c.ID)
			if err == nil {
				item.RegionCount = len(regions)
			}

			items = append(items, item)
		}

		renderPage(w, tmpl, "countries_list", items)
	}
}