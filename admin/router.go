package admin

import (
	"net/http"
	"html/template"
	
	"go.etcd.io/bbolt"
	"gocart/config"
)

func Route(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.Handler {
	mux := http.NewServeMux()

	// Dashboard
	mux.HandleFunc("GET /", 						getDashboard(cfg, db, tmpl))

	// Settings

	// Environment
    //mux.HandleFunc("GET    /environment",       		getEnvironment(cfg, db, tmpl))

	// Currencies
    mux.HandleFunc("GET		/currencies",       		getCurrencies(cfg, db, tmpl))
    mux.HandleFunc("GET		/currencies/edit/{id}",  	getCurrenciesEdit(cfg, db, tmpl))
    mux.HandleFunc("POST	/currencies/edit/{id}",		postCurrenciesEdit(cfg, db, tmpl))
	mux.HandleFunc("POST	/currencies/reset",			postCurrenciesReset(cfg, db, tmpl))

	// Countries
    mux.HandleFunc("GET		/countries",       			getCountries(cfg, db, tmpl))
    mux.HandleFunc("GET		/countries/edit/{id}",  	getCountriesEdit(cfg, db, tmpl))
    mux.HandleFunc("POST	/countries/edit/{id}",		postCountriesEdit(cfg, db, tmpl))
	mux.HandleFunc("POST	/countries/reset",			postCountriesReset(cfg, db, tmpl))

	return mux
}
