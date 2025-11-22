package main

import (
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/dtsong/battleforge/backend/internal/httpapi"
	"github.com/dtsong/battleforge/backend/internal/observability"
)

func main() {
	logger := observability.NewLogger()

	// TODO: Initialize database connection
	// dbConnString := getDBConnString()
	// db, err := db.NewDatabase(dbConnString)
	// if err != nil {
	// 	logger.Fatalf("failed to initialize database: %v", err)
	// }
	// defer db.Close()

	addr := getAddr()
	logger.Infof("starting battleforge-api on %s", addr)

	router := httpapi.NewRouter(logger)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatalf("server failed: %v", err)
	}
}

func getAddr() string {
	if v := os.Getenv("BATTLEFORGE_API_ADDR"); v != "" {
		return v
	}
	// default dev address
	return ":8080"
}
