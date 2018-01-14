package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/schigh/go-health"
	"github.com/schigh/go-health/checkers"
	"github.com/schigh/go-health/handlers"
	"github.com/schigh/go-health/loggers"
	log "github.com/sirupsen/logrus"
	"goji.io"
	"goji.io/pat"
)

var (
	healthCheck *health.Health
	db          *sql.DB
)

func main() {
	if err := setupDatabase(); err != nil {
		log.Fatalf("database setup error: %v", err)
	}

	if err := setupHealthCheck(); err != nil {
		log.Fatalf("health check setup error: %v", err)
	}

	mux := goji.NewMux()
	healthHandler := handlers.NewJSONHandlerFunc(healthCheck, map[string]interface{}{})
	mux.HandleFunc(pat.Get("/healthcheck"), healthHandler)

	http.ListenAndServe("0.0.0.0:80", mux)
}

func setupDatabase() error {
	log.Info("setting up database")

	dsn := fmt.Sprintf("%s:%s@tcp(mysql:3306)/?parseTime=true&timeout=10s", os.Getenv("DB_USER"), os.Getenv("DB_PASS"))
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	db = conn

	return nil
}

func setupHealthCheck() error {
	log.Info("setting up health check")

	sqlCheck, err := checkers.NewSQL(&checkers.SQLConfig{
		DB: db,
	})
	if err != nil {
		return err
	}

	healthCheck = health.New()
	healthCheck.Logger = loggers.NewBasic()
	healthCheck.AddCheck(&health.Config{
		Name:     "sql-check",
		Checker:  sqlCheck,
		Interval: time.Duration(3) * time.Second,
		Fatal:    true,
	})

	return healthCheck.Start()
}
