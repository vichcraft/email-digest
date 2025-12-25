package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/vichcraft/email-digest/internal/config"
	"github.com/vichcraft/email-digest/internal/email"
)

func main() {
	godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()

	account := cfg.GmailAccounts[2]
	fmt.Printf("Connecting to %s....", account.Email)

	service, err := email.NewClient(ctx, cfg.CredentialsPath, account.TokenPath)
	if err != nil {
		log.Fatalf("Failed to create Gmail client: %v", err)
	}

	fmt.Println("Succesfully connected to Gmail!")

	profile, err := service.Users.GetProfile("me").Do()
	if err != nil {
		log.Fatalf("Failed to get profile: %v", err)
	}

	fmt.Printf("Email: %s\n", profile.EmailAddress)
	fmt.Printf("Total messages: %d\n", profile.MessagesTotal)
}
