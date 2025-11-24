package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/dtsong/vgccorner/backend/internal/httpapi"
	"github.com/dtsong/vgccorner/backend/internal/observability"
)

func main() {
	logger := observability.NewLogger()

	// Initialize database connection
	dbConnString := getDBConnString()
	logger.Infof("connecting to database at %s", dbConnString)
	// TODO: Uncomment when db.NewDatabase is ready
	// db, err := db.NewDatabase(dbConnString)
	// if err != nil {
	// 	logger.Fatalf("failed to initialize database: %v", err)
	// }
	// defer db.Close()

	addr := getAddr()
	logger.Infof("starting vgccorner-api on %s", addr)

	router := httpapi.NewRouter(logger)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatalf("server failed: %v", err)
	}
}

func getAddr() string {
	if v := os.Getenv("SERVER_PORT"); v != "" {
		return ":" + v
	}
	// default dev address
	return ":8080"
}

func getDBConnString() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "vgccorner")
	password := getEnv("DB_PASSWORD", "vgccorner_dev_password")
	dbName := getEnv("DB_NAME", "vgccorner")
	sslMode := getEnv("DB_SSL_MODE", "disable")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbName, sslMode)
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
