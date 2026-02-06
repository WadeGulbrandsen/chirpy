package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/WadeGulbrandsen/chirpy/internal/database"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

const (
	default_port     = 8080
	default_prefix   = "/app/"
	default_app_path = "."
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	tokenSecret    string
	polkaKey       string
	port           int
	appPath        string
	appPrefix      string
}

//go:embed sql/schema/*.sql
var dbSchema embed.FS

func getConfig() *apiConfig {
	// Get configuration
	godotenv.Load()
	cfg := apiConfig{}
	cfg.fileserverHits.Store(0)
	cfg.platform = os.Getenv("PLATFORM")

	cfg.tokenSecret = os.Getenv("JWT_SECRET")
	if cfg.tokenSecret == "" {
		log.Fatalf("JWT_SECRET is required")
	}

	cfg.polkaKey = os.Getenv("POLKA_KEY")
	if cfg.tokenSecret == "" {
		log.Fatalf("POLKA_KEY is required")
	}

	appPath, err := getAppPath(os.Getenv("APP_PATH"))
	if err != nil {
		log.Fatalf("Invalid APP_PATH: %v\n", err)
	}
	cfg.appPath = appPath

	appPrefix, err := getAppPrefix(os.Getenv("APP_PREFIX"))
	if err != nil {
		log.Fatalf("Invalid APP_PREFIX: %v\n", err)
	}
	cfg.appPrefix = appPrefix

	port, err := getPort(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Invalid PORT: %v\n", err)
	}
	cfg.port = port

	cfg.dbQueries = getDatabase(os.Getenv("DB_URL"))
	return &cfg
}

func getAppPath(app_path string) (string, error) {
	if app_path == "" {
		app_path = default_app_path
	}
	info, err := os.Stat(app_path)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s is not a directory", app_path)
	}
	return app_path, nil
}

func getAppPrefix(prefix string) (string, error) {
	if prefix == "" {
		prefix = default_prefix
	}
	if len(prefix) < 3 || !strings.HasPrefix(prefix, "/") || !strings.HasSuffix(prefix, "/") || strings.Contains(prefix, " ") {
		return "", fmt.Errorf("prefix must be start and end with / and cannot contain spaces")
	}
	return prefix, nil
}

func getPort(port_string string) (int, error) {
	if port_string == "" {
		return default_port, nil
	}
	port, err := strconv.Atoi(port_string)
	if err != nil {
		return 0, err
	}
	if port <= 0 || port > 65535 {
		return port, fmt.Errorf("PORT must be between 1 and 65,535")
	}
	return port, nil
}

func getDatabase(dbURL string) *database.Queries {
	// Get Database connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Could not open DB_URL connection: %v\n", err)
	}

	// Run Database migrations
	goose.SetBaseFS(dbSchema)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Could not set database dialect: %v\n", err)
	}
	if err := goose.Up(db, "sql/schema"); err != nil {
		log.Fatalf("Could not run migrations on database: %v\n", err)
	}

	return database.New(db)
}
