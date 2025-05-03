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

	port := config.New()
	log.Println("server is running on port: "+port.Port)
	go orkestrator.OrkestratorRun(db)

	agent.AgentRun()


}

