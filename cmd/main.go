package main

import (
	"log"

	"github.com/LootNex/HTTP-Caculator_V2/internal/agent"
	server "github.com/LootNex/HTTP-Caculator_V2/internal/agent/grpcserver"
	"github.com/LootNex/HTTP-Caculator_V2/internal/db"
	"github.com/LootNex/HTTP-Caculator_V2/internal/orkestrator"
	"github.com/LootNex/HTTP-Caculator_V2/pkg/config"
	"github.com/LootNex/HTTP-Caculator_V2/pkg/migrations"
)

func main() {
	db, err := db.InitDB()
	if err != nil {
		log.Fatalf("cannot init sqlite %v", err)
	}

	defer db.Close()

	err = migrations.InitTables(db)
	if err != nil {
		log.Fatalf("cannot init table %v", err)
	}

	port := config.New()

	log.Println("server is running on port: " + port.Port)
	go orkestrator.OrkestratorRun(db)

	go agent.AgentRun()

	log.Println("Agent run")

	err = server.StartServer()
	if err != nil {
		log.Fatalf("cannot start grpc server %v", err)
	}

}
