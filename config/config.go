package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPath        string
	APIEnabled    bool
	APIPort       int
	AdminEnabled  bool
	AdminPort     int
	JSLibEnabled  bool
	JSLibPort     int
	StoreName     string
	StoreCurrency string
	StoreCountry  string
	JWTSecret     string
}

func ConfigLoad() (*Config, error) {
	log.Println("Loading environment...")

	// Load .env file if present (ignore error, OS env vars are fine too)
	godotenv.Load()

	e := &Config{}
	var err error

	// Load strings
	stringVars := map[string]*string{
		"GOCART_DB_PATH":        &e.DBPath,
		"GOCART_STORE_NAME":     &e.StoreName,
		"GOCART_STORE_CURRENCY": &e.StoreCurrency,
		"GOCART_STORE_COUNTRY":  &e.StoreCountry,
		"GOCART_SECRET_JWT":     &e.JWTSecret,
	}
	for key, dest := range stringVars {
		val := os.Getenv(key)
		if val == "" {
			return nil, fmt.Errorf("missing required env var: %s", key)
		}
		*dest = val
	}

	// Load bools
	boolVars := map[string]*bool{
		"GOCART_API_EN":   &e.APIEnabled,
		"GOCART_ADMIN_EN": &e.AdminEnabled,
		"GOCART_JSLIB_EN": &e.JSLibEnabled,
	}
	for key, dest := range boolVars {
		val := os.Getenv(key)
		if val == "" {
			return nil, fmt.Errorf("missing required env var: %s", key)
		}
		*dest, err = strconv.ParseBool(val)
		if err != nil {
			return nil, fmt.Errorf("invalid bool value for %s: %w", key, err)
		}
	}

	// Load ints
	intVars := map[string]*int{
		"GOCART_API_PORT":   &e.APIPort,
		"GOCART_ADMIN_PORT": &e.AdminPort,
		"GOCART_JSLIB_PORT": &e.JSLibPort,
	}
	for key, dest := range intVars {
		val := os.Getenv(key)
		if val == "" {
			return nil, fmt.Errorf("missing required env var: %s", key)
		}
		parsed, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("invalid int value for %s: %w", key, err)
		}
		*dest = parsed
	}

	log.Println("Environment loaded successfully.")
	return e, nil
}
