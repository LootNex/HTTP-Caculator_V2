package main

import (
	"Calculator_V2/internal/agent"
	"Calculator_V2/internal/db"
	"Calculator_V2/internal/orkestrator"
	"Calculator_V2/pkg/config"
	"Calculator_V2/pkg/migrations"
	"log"
)

func main() {
	db, err := db.InitDB()
	if err != nil{
		log.Fatalf("cannot init sqlite %v", err)
	}

	defer db.Close()

	err = migrations.InitTables(db)
	if err != nil{
		log.Fatalf("cannot init table %v", err)
	}

	port := config.New()

	log.Println("server is running on port: "+port.Port)
	go orkestrator.OrkestratorRun(db)

	agent.AgentRun()


}

