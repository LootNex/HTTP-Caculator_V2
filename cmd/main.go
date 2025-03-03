package main

import (
	"Calculator_V2/internal/agent"
	"Calculator_V2/internal/orkestrator"
)

func main() {

	go orkestrator.OrkestratorRun()

	agent.AgentRun()

}