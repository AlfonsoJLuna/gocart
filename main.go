package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gocart/admin"
	"gocart/api"
	"gocart/config"
	"gocart/jslib"
)

func serveHTTP(name string, port int, handler http.Handler) {
	addr := fmt.Sprintf(":%d", port)

	log.Printf("%s server listening on port %d...", name, port)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Failed to initialize %s server: %v", name, err)
	}
}

func main() {
	log.Println("Welcome to GoCart.")

	cfg, err := config.ConfigLoad()
	if err != nil {
		log.Fatalf("Failed to load environment: %v", err)
	}

	db, err := dbInit(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	tmpl, err := admin.InitTemplates()
	if err != nil{
		log.Fatalf("Failed to init admin templates: %v", err)
	}

	if cfg.APIEnabled {
		go serveHTTP("API", cfg.APIPort, api.Route(cfg, db))
	}

	if cfg.AdminEnabled {
		go serveHTTP("Admin", cfg.AdminPort, admin.Route(cfg, db, tmpl))
	}

	if cfg.JSLibEnabled {
		go serveHTTP("JSLib", cfg.JSLibPort, jslib.Route(cfg, db))
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("GoCart is shutting down...")
}
