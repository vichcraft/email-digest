package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/vichcraft/email-digest/internal/email"
)

// DiscordMessage represents the JSON structure for Discord webhook
type DiscordMessage struct {
	Content string  `json:"content,omitempty"`
	Embeds  []Embed `json:"embeds,omitempty"`
}

// Embed represents a Discord embed (rich message format)
type Embed struct {
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	Color       int     `json:"color,omitempty"`
	Timestamp   string  `json:"timestamp,omitempty"`
	Footer      *Footer `json:"footer,omitempty"`
}

// Footer represents the footer of a Discord embed
type Footer struct {
	Text string `json:"text"`
}

// SendDetailedDigest sends a detailed email digest organized by account
// Each account is sent to its own Discord webhook/channel
func SendDetailedDigest(emailsByAccount map[string]AccountEmails) error {
	if len(emailsByAccount) == 0 {
		return fmt.Errorf("no emails to send")
	}

	// Send each account to its specific webhook
	accountNum := 1
	for accountEmail, accountData := range emailsByAccount {
		fmt.Printf("  → Sending message for %s (%d emails) to dedicated channel\n", accountEmail, len(accountData.Emails))

		if err := SendAccountMessage(accountData.WebhookURL, accountEmail, accountData.Emails, accountNum); err != nil {
			// Log error but continue with other accounts
			fmt.Printf("Warning: failed to send message for %s: %v\n", accountEmail, err)
		} else {
			fmt.Printf("Sent\n")
		}

		accountNum++

		// Small delay to avoid rate limiting
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

// AccountEmails holds emails and webhook URL for a specific account
type AccountEmails struct {
	Emails     []email.Email
	WebhookURL string
}

// SendAccountMessage sends a detailed message for a single account
func SendAccountMessage(webhookURL, account string, emails []email.Email, accountNum int) error {
	description := buildAccountDescription(emails)

	// Ensure description is not empty and within limits
	if description == "" {
		description = "No email details available"
	}
	if len(description) > 4096 {
		description = description[:4090] + "..."
	}

	embed := Embed{
		Title:       fmt.Sprintf("Account %d: %s", accountNum, account),
		Description: fmt.Sprintf("**Number of emails: %d**\n\n%s", len(emails), description),
		Color:       getColorForAccount(accountNum),
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Footer: &Footer{
			Text: "Email Monitor",
		},
	}

	message := DiscordMessage{
		Embeds: []Embed{embed},
	}

	return sendWebhook(webhookURL, message)
}

// buildAccountDescription creates the detailed email list for an account
func buildAccountDescription(emails []email.Email) string {
	var builder strings.Builder

	for i, e := range emails {
		builder.WriteString(fmt.Sprintf("**%d.**\n", i+1))

		// Handle empty fields safely
		from := e.From
		if from == "" {
			from = "Unknown sender"
		}

		subject := e.Subject
		if subject == "" {
			subject = "(No subject)"
		}

		snippet := e.Snippet
		if snippet == "" {
			snippet = "(No preview available)"
		}

		builder.WriteString(fmt.Sprintf("**From:** %s\n", cleanEmailAddress(from)))
		builder.WriteString(fmt.Sprintf("**Subject:** %s\n", truncate(subject, 150)))
		builder.WriteString(fmt.Sprintf("**Summary:** %s\n", truncate(snippet, 200)))

		// Format date safely
		if !e.Date.IsZero() {
			builder.WriteString(fmt.Sprintf("**Date:** %s\n", e.Date.Format("Jan 02, 3:04 PM")))
		}

		builder.WriteString("\n")

		// Discord has a 4096 character limit per embed description
		if builder.Len() > 3500 {
			remaining := len(emails) - i - 1
			if remaining > 0 {
				builder.WriteString(fmt.Sprintf("*...and %d more emails (list truncated)*\n", remaining))
			}
			break
		}
	}

	result := builder.String()

	// Ensure we return something even if empty
	if result == "" {
		return "No email details available"
	}

	return result
}

// SendSimpleMessage sends a simple text message to Discord
func SendSimpleMessage(webhookURL string, content string) error {
	if content == "" {
		content = "Empty message"
	}

	message := DiscordMessage{
		Content: content,
	}
	return sendWebhook(webhookURL, message)
}

// sendWebhook sends the actual HTTP request to Discord
func sendWebhook(webhookURL string, message DiscordMessage) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send webhook: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Read response body for error details
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("discord returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// Helper functions

func truncate(s string, maxLen int) string {
	// Remove any null bytes or invalid characters
	s = strings.Map(func(r rune) rune {
		if r == 0 {
			return -1
		}
		return r
	}, s)

	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func cleanEmailAddress(email string) string {
	// Remove null bytes
	email = strings.Map(func(r rune) rune {
		if r == 0 {
			return -1
		}
		return r
	}, email)

	// Extract just the email if it has a name like "John Doe <john@example.com>"
	if strings.Contains(email, "<") && strings.Contains(email, ">") {
		start := strings.Index(email, "<")
		end := strings.Index(email, ">")
		if start < end && start >= 0 && end < len(email) {
			name := strings.TrimSpace(email[:start])
			addr := email[start : end+1]
			if name != "" {
				return name + " " + addr
			}
			return addr
		}
	}
	return email
}

func getColorForAccount(accountNum int) int {
	colors := []int{
		3447003,  // Blue
		15158332, // Red
		10181046, // Purple
		15844367, // Gold
		3066993,  // Green
	}

	return colors[(accountNum-1)%len(colors)]
}
