package admin

import (
    "html/template"
    "net/http"

    "go.etcd.io/bbolt"

    "gocart/config"
)

func Route(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /",                                         dashboard(cfg, db, tmpl))

    mux.HandleFunc("GET  /currencies",                              currenciesList(cfg, db, tmpl))
    mux.HandleFunc("GET  /currencies/{id}",                         currenciesEdit(cfg, db, tmpl))
    mux.HandleFunc("POST /currencies/{id}",                         currenciesEditPost(cfg, db, tmpl))
    mux.HandleFunc("POST /currencies/reset",                        currenciesResetPost(cfg, db, tmpl))

    mux.HandleFunc("GET  /countries",                               countriesList(cfg, db, tmpl))
    mux.HandleFunc("GET  /countries/{id}",                          countriesEdit(cfg, db, tmpl))
    mux.HandleFunc("POST /countries/{id}",                          countriesEditPost(cfg, db, tmpl))
    mux.HandleFunc("POST /countries/reset",                         countriesResetPost(cfg, db, tmpl))

    mux.HandleFunc("GET  /countries/{id}/region/new",               regionsNew(cfg, db, tmpl))
    mux.HandleFunc("POST /countries/{id}/region/new",               regionsNewPost(cfg, db, tmpl))
    mux.HandleFunc("GET  /countries/{id}/region/{index}",           regionsEdit(cfg, db, tmpl))
    mux.HandleFunc("POST /countries/{id}/region/{index}",           regionsEditPost(cfg, db, tmpl))
    mux.HandleFunc("POST /countries/{id}/region/{index}/delete",    regionsDeletePost(cfg, db, tmpl))

    return mux
}
