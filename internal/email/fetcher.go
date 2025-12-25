package email

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/gmail/v1"
)

func FetchRecentUnread(ctx context.Context, service *gmail.Service, hoursBack int) ([]Email, error) {
	since := time.Now().Add(-time.Duration(hoursBack) * time.Hour)
	sinceUnix := since.Unix()

	query := fmt.Sprintf("is:unread after:%d", sinceUnix)

	listCall := service.Users.Messages.List("me").Q(query).MaxResults(100)
	response, err := listCall.Do()
	if err != nil {
		return nil, fmt.Errorf("Error reading messages: %v", err)
	}

	if len(response.Messages) == 0 {
		return []Email{}, nil
	}

	profile, err := service.Users.GetProfile("me").Do()
	if err != nil {
		return nil, fmt.Errorf("Failed to get profile: %v", err)
	}

	accountEmail := profile.EmailAddress

	var emails []Email

	for _, msg := range response.Messages {
		email, err := fetchMessageDetails(ctx, service, msg.Id, accountEmail)
		if err != nil {
			fmt.Printf("Failed to fetch message %s: %v\n", msg.Id, err)
			continue
		}
		emails = append(emails, email)
	}

	return emails, nil

}

func fetchMessageDetails(ctx context.Context, service *gmail.Service, messageId string, accountEmail string) (Email, error) {
	msg, err := service.Users.Messages.Get("me", messageId).Format("full").Do()
	if err != nil {
		return Email{}, err
	}

	email := Email{
		Id:           msg.Id,
		AccountEmail: accountEmail,
		IsUnread:     IsUnread(msg.LabelIds),
		Labels:       msg.LabelIds,
	}

	for _, header := range msg.Payload.Headers {
		switch header.Name {
		case "From":
			email.From = header.Value
		case "Subject":
			email.Subject = header.Value
		case "Date":
			parsedDate, err := parseEmailDate(header.Value)
			if err == nil {
				email.Date = parsedDate
			}
		}
	}

	email.Snippet = msg.Snippet

	return email, nil
}

func IsUnread(labels []string) bool {
	for _, label := range labels {
		if label == "UNREAD" {
			return true
		}
	}
	return false
}

func parseEmailDate(dateStr string) (time.Time, error) {
	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"Mon, 02 Jan 2006 15:04:05 -0700",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("Unable to parse date: %s", dateStr)
}
