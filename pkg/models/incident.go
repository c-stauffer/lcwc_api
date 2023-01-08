package models

import "time"

type Incident struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PubDate     time.Time `json:"pubDate"`
}
