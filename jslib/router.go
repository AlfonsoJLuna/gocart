package jslib

import (
	"net/http"
	"go.etcd.io/bbolt"
	"gocart/config"
)

func Route(cfg *config.Config, db *bbolt.DB) http.Handler {
	mux := http.NewServeMux()

	return mux
}
