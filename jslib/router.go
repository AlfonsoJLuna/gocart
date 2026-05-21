package jslib

import (
	"net/http"
	"database/sql"

	"gocart/config"
)

func Route(cfg *config.Config, db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	return mux
}
