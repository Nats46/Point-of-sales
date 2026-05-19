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

	port := os.Getenv("PORT")
    if port == "" {
        port = "8000"
    }

	r := router.SetupRouter()

	// run server
	r.Run(fmt.Sprintf(":%s", port))

}    
