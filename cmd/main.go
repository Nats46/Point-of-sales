package main

import (
	"fmt"
	"log"
	"os"

    "point-of-sales/config"
    "point-of-sales/db"
	"point-of-sales/internal/router"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("config error: %v", err)
    }

    database, err := db.Connect(cfg.DSN)
    if err != nil {
        log.Fatalf("db error: %v", err)
    }
    defer database.Close()

    if err := db.RunMigrations(database); err != nil {
        log.Fatalf("migration error: %v", err)
    }

	port := os.Getenv("PORT")
    if port == "" {
        port = "8000"
    }

	r := router.SetupRouter(database, os.Getenv("JWT_SECRET"))

	// run server
	r.Run(fmt.Sprintf(":%s", port))

}    
