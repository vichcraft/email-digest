package config

import (
	"fmt"
	"os"
)

type Config struct {
	GeminiAPIKey      string
	DiscordWebhookURL string
	GmailAccounts     []GmailAccount
	CredentialsPath   string
}

type GmailAccount struct {
	Email     string
	TokenPath string
}

func Load() (*Config, error) {
	geminiKey := os.Getenv("GEMINI_API_KEY")
	discordUrl := os.Getenv("DISCORD_WEBHOOK_URL")

	if geminiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}

	if discordUrl == "" {
		return nil, fmt.Errorf("DISCORD_WEBHOOK_URL not set")
	}

	accounts := []GmailAccount{
		{
			Email:     os.Getenv("GMAIL_ACCOUNT_1"),
			TokenPath: "credentials/token_account1.json",
		},
		{
			Email:     os.Getenv("GMAIL_ACCOUNT_2"),
			TokenPath: "credentials/token_account2.json",
		},
		{
			Email:     os.Getenv("GMAIL_ACCOUNT_3"),
			TokenPath: "credentials/token_account3.json",
		},
	}

	return &Config{
		GeminiAPIKey:      geminiKey,
		DiscordWebhookURL: discordUrl,
		GmailAccounts:     accounts,
		CredentialsPath:   "credentials/credentials.json",
	}, nil

}
