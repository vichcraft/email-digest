package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/vichcraft/email-digest/internal/config"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Gmail accounts: %d\n", len(cfg.GmailAccounts))

	fmt.Println(cfg.GmailAccounts[2].Email)

	fmt.Println("Config loaded successfully")
}
