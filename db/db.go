package db

import (
    _ "embed"
    "fmt"
    "log"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

//go:embed migration/000001_init.up.sql
var initSchema string

func Connect(dsn string) (*sqlx.DB, error) {
    db, err := sqlx.Connect("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(10)

    log.Println("✅ Database connected successfully")
    return db, nil
}

func RunMigrations(db *sqlx.DB) error {
    _, err := db.Exec(initSchema)
    if err != nil {
        return fmt.Errorf("failed to execute migration script: %w", err)
    }

    log.Println("✅ Database migrations applied successfully")
    return nil
}