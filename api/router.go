package api

import (
	"net/http"
	"database/sql"

	"gocart/config"
	"gocart/middleware"
)

func Route(cfg *config.Config, db *sql.DB) http.Handler {
    mux := http.NewServeMux()

    // Currencies
    mux.HandleFunc("POST   /currencies",           middleware.WithLogging(CreateCurrency(db)))
    mux.HandleFunc("GET    /currencies/{id}",      middleware.WithLogging(ReadCurrencyByID(db)))
    mux.HandleFunc("GET    /currencies",           middleware.WithLogging(ReadCurrencies(db)))
    mux.HandleFunc("PUT    /currencies/{id}",      middleware.WithLogging(UpdateCurrencyByID(db)))
    mux.HandleFunc("DELETE /currencies/{id}",      middleware.WithLogging(DeleteCurrencyByID(db)))

    // Countries
    mux.HandleFunc("POST   /countries",            middleware.WithLogging(CreateCountry(db)))
    mux.HandleFunc("GET    /countries/{id}",       middleware.WithLogging(ReadCountryByID(db)))
    mux.HandleFunc("GET    /countries",            middleware.WithLogging(ReadCountries(db)))
    mux.HandleFunc("PUT    /countries/{id}",       middleware.WithLogging(UpdateCountryByID(db)))
    mux.HandleFunc("DELETE /countries/{id}",       middleware.WithLogging(DeleteCountryByID(db)))

    // Auth
    mux.HandleFunc("POST   /auth/signup",          middleware.WithLogging(AuthSignUp(db)))
    mux.HandleFunc("POST   /auth/signin",          middleware.WithLogging(AuthSignIn(db)))
    mux.HandleFunc("POST   /auth/verify",          middleware.WithLogging(AuthVerify(db)))
    mux.HandleFunc("POST   /auth/verify/resend",   middleware.WithLogging(AuthVerifyResend(db)))
    mux.HandleFunc("POST   /auth/password/forgot", middleware.WithLogging(AuthPasswordForgot(db)))
    mux.HandleFunc("POST   /auth/password/reset",  middleware.WithLogging(AuthPasswordReset(db)))

    return mux
}
