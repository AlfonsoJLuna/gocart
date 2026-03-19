package admin

import (
	"net/http"
	"html/template"
	
	"go.etcd.io/bbolt"
	
	"gocart/config"
)

func getDashboard(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		render(w, tmpl, "dashboard", nil)
    }
}
