package db

import (
    "database/sql"
    "fmt"

    _ "modernc.org/sqlite"
)

func InitDB() (*sql.DB, error) {
    db, err := sql.Open("sqlite", "C:/Nikita/Calc/HTTP-Caculator_V2/internal/db/store.db")
    if err != nil {
        return nil, fmt.Errorf("failed to open SQLite: %v", err)
    }

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping SQLite: %v", err)
    }

    fmt.Println("SQLite connection established")
    return db, nil
}
