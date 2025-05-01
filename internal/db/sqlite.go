package db

import "database/sql"

func InitDB() (*sql.DB, error) {
	
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil{
		return nil, err
	}

	return db, nil
}
