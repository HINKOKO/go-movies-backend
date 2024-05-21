package main

import (
	"backend/internal/repository"
	"backend/internal/repository/dbrepo"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const port = 8080

type application struct {
	DSN          string
	Domain       string
	DB           repository.DatabaseRepo
	auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
	APIKey       string
}

func main() {
	// Set application config
	var app application

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")

	// read from command line - or default
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=12345 dbname=movies sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecret", "signing secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "signing audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "Cookie domain")
	flag.StringVar(&app.Domain, "domain", "example.com", "Domain")
	flag.StringVar(&app.APIKey, "api-key", apiKey, "Api key")
	flag.Parse()

	// Connect to Database
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	// populate the application struct relevant field
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()

	// __Host-refresh_token -- never gets sets on Chrome
	// https://datatracker.ietf.org/doc/html/draft-west-cookie-prefixes-05 section 3.2
	// works when you have HTTPS enabled - so forget dev mode - + no domain cookie definition
	// So for dev mode -> go for a cookie name with classic cool name, haha
	app.auth = Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.JWTSecret,
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "refresh_token",
		CookieDomain:  app.CookieDomain,
	}

	log.Println("Starting application on Port: ", port)

	// start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
