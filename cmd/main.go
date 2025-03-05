package main

import (
	"Calculator_V2/internal/agent"
	"Calculator_V2/internal/orkestrator"
	"Calculator_V2/pkg/config"
	"log"
)

func main() {
	port := config.New()
	log.Println("server is running on port: "+port.Port)
	go orkestrator.OrkestratorRun()

	agent.AgentRun()


}
