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

	c := &Config{}
	var err error

	// Load strings
	stringVars := map[string]*string{
		"GOCART_DB_PATH":        &c.DBPath,
		"GOCART_STORE_NAME":     &c.StoreName,
		"GOCART_STORE_CURRENCY": &c.StoreCurrency,
		"GOCART_STORE_COUNTRY":  &c.StoreCountry,
		"GOCART_SECRET_JWT":     &c.JWTSecret,
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
		"GOCART_API_EN":   &c.APIEnabled,
		"GOCART_ADMIN_EN": &c.AdminEnabled,
		"GOCART_JSLIB_EN": &c.JSLibEnabled,
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
		"GOCART_API_PORT":   &c.APIPort,
		"GOCART_ADMIN_PORT": &c.AdminPort,
		"GOCART_JSLIB_PORT": &c.JSLibPort,
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
	return c, nil
}
