package main

import (
	"log"
	"os"

	"jungle-rpg/internal/auth"
	"jungle-rpg/internal/repository"
	"jungle-rpg/internal/server"
)

func main() {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "/data/jungle.db"
	}

	// Ensure directory exists
	if dir := dbPath[:len(dbPath)-len("jungle.db")]; dir != "" {
		os.MkdirAll(dir, 0755)
	}

	db, err := repository.Open(dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	auth.InitStore()

	srv := server.New(db)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
