package db

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
    db, err := sql.Open("sqlite3", "/db/store.db")
    if err != nil {
        log.Fatal(err)
        return nil, fmt.Errorf("failed to open SQLite connection: %v", err)
    }
	
    if err := db.Ping(); err != nil {
        log.Fatal(err)
        return nil, fmt.Errorf("failed to ping SQLite: %v", err)
    }

    fmt.Println("SQLite connection established")
    return db, nil
}
