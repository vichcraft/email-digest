package email

import "time"

type Email struct {
	Id           string
	AccountEmail string
	From         string
	Subject      string
	Snippet      string
	Date         time.Time
	IsUnread     bool
	Labels       []string
}
