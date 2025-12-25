package config

import (
	"fmt"
	"os"
)

type Config struct {
	GeminiAPIKey    string
	GmailAccounts   []GmailAccount
	CredentialsPath string
}

type GmailAccount struct {
	Email          string
	TokenPath      string
	DiscordWebhook string
}

func Load() (*Config, error) {
	geminiKey := os.Getenv("GEMINI_API_KEY")

	if geminiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}

	accounts := []GmailAccount{
		{
			Email:          os.Getenv("GMAIL_ACCOUNT_1"),
			TokenPath:      "credentials/token_account1.json",
			DiscordWebhook: os.Getenv("DISCORD_WEBHOOK_1"),
		},
		{
			Email:          os.Getenv("GMAIL_ACCOUNT_2"),
			TokenPath:      "credentials/token_account2.json",
			DiscordWebhook: os.Getenv("DISCORD_WEBHOOK_2"),
		},
		{
			Email:          os.Getenv("GMAIL_ACCOUNT_3"),
			TokenPath:      "credentials/token_account3.json",
			DiscordWebhook: os.Getenv("DISCORD_WEBHOOK_3"),
		},
	}

	// Validate that all webhooks are set
	for i, acc := range accounts {
		if acc.DiscordWebhook == "" {
			return nil, fmt.Errorf("DISCORD_WEBHOOK_%d not set", i+1)
		}
	}

	return &Config{
		GeminiAPIKey:    geminiKey,
		GmailAccounts:   accounts,
		CredentialsPath: "credentials/credentials.json",
	}, nil

}
