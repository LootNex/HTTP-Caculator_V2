package migrations

import "database/sql"

func InitTables(db *sql.DB) error {

	// _, err := db.Exec(`DROP TABLE users`)

	// if err != nil{
	// 	return err
	// }

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users(
		user_id TEXT PRIMARY KEY, 
		login TEXT,
		password TEXT)`)

	if err != nil{
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS expressions(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		expression TEXT,
		result TEXT,
		user_id TEXT REFERENCES users(user_id))`)

	if err != nil{
		return err
	}
	return nil
}