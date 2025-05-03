package auth

import "database/sql"

type App struct {
	DB *sql.DB
}

func NewApp(db *sql.DB) *App {
	return &App{
		DB: db,
	}
}
