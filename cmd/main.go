package main

import (
	"Calculator_V2/internal/agent"
	"Calculator_V2/internal/db"
	"Calculator_V2/internal/orkestrator"
	"Calculator_V2/pkg/config"
	"fmt"
	"log"
)

func main() {
	db, err := db.InitDB()
	if err != nil{
		fmt.Errorf("cannot init sqlite")
	}
	log.Println("sqlite connection is succesful")
	
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users(
		uuid INTEGER PRIMARY KEY AUTOINCREMENT, 
		login TEXT,
		password TEXT)`)
	if err != nil{
		fmt.Errorf("cannot init table %v", err)
	}
	port := config.New()
	log.Println("server is running on port: "+port.Port)
	go orkestrator.OrkestratorRun(db)

	agent.AgentRun()


}

