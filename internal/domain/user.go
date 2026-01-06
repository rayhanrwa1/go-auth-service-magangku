package domain

import "time"

type User struct {
	ID            string
	Name          string
	Email         string
	Password      string
	UserableType  string
	TermsAccepted bool
	CreatedAt     time.Time
}
