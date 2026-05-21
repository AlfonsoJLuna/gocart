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
	"gocart/db"
	"gocart/jslib"
	"gocart/seeds"
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

	database, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	seeded, err := db.IsSeeded(database)
	if err != nil {
		log.Fatalf("Failed to check seed status: %v", err)
	}
	if !seeded {
		log.Println("Database is empty. Seeding...")
		if err := seeds.SeedAll(database); err != nil {
			log.Fatalf("Failed to seed database: %v", err)
		}
		if err := db.MarkSeeded(database); err != nil {
			log.Fatalf("Failed to mark database as seeded: %v", err)
		}
		log.Println("Database seeded successfully.")
	} else {
		log.Println("Database was already seeded.")
	}

	tmpl, err := admin.InitTemplates()
	if err != nil{
		log.Fatalf("Failed to init admin templates: %v", err)
	}

	if cfg.APIEnabled {
		go serveHTTP("API", cfg.APIPort, api.Route(cfg, database))
	}

	if cfg.AdminEnabled {
		go serveHTTP("Admin", cfg.AdminPort, admin.Route(cfg, database, tmpl))
	}

	if cfg.JSLibEnabled {
		go serveHTTP("JSLib", cfg.JSLibPort, jslib.Route(cfg, database))
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("GoCart is shutting down...")
}
