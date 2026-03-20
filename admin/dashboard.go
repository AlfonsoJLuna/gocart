package admin

import (
	"html/template"
	"net/http"
	
	"go.etcd.io/bbolt"
	
	"gocart/config"
)

func dashboard(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		renderPage(w, tmpl, "dashboard", nil)
    }
}
