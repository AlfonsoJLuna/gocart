package admin

import (
	"html/template"
	"net/http"

	"go.etcd.io/bbolt"

	"gocart/config"
)

func Route(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.Handler {
	mux := http.NewServeMux()

	// Dashboard
	mux.HandleFunc("GET /", getDashboard(cfg, db, tmpl))

	// Settings

	// Environment
	//mux.HandleFunc("GET    /environment",       		getEnvironment(cfg, db, tmpl))

	// Currencies
	mux.HandleFunc("GET  /currencies", getCurrencies(cfg, db, tmpl))
	mux.HandleFunc("GET  /currencies/{id}", getCurrenciesEdit(cfg, db, tmpl))
	mux.HandleFunc("POST /currencies/{id}", postCurrenciesEdit(cfg, db, tmpl))
	mux.HandleFunc("POST /currencies/reset", postCurrenciesReset(cfg, db, tmpl))

	// Countries
	mux.HandleFunc("GET  /countries", getCountries(cfg, db, tmpl))
	mux.HandleFunc("GET  /countries/{id}", getCountriesEdit(cfg, db, tmpl))
	mux.HandleFunc("POST /countries/{id}", postCountriesEdit(cfg, db, tmpl))
	mux.HandleFunc("POST /countries/reset", postCountriesReset(cfg, db, tmpl))

	// Regions
	mux.HandleFunc("GET  /countries/{id}/region/new", getRegionNew(cfg, db, tmpl))
	mux.HandleFunc("POST /countries/{id}/region/new", postRegionNew(cfg, db, tmpl))
	mux.HandleFunc("GET  /countries/{id}/region/{index}", getRegionEdit(cfg, db, tmpl))
	mux.HandleFunc("POST /countries/{id}/region/{index}", postRegionEdit(cfg, db, tmpl))
	mux.HandleFunc("POST /countries/{id}/region/{index}/delete", postRegionDelete(cfg, db, tmpl))

	return mux
}
