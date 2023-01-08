package models

import "time"

type Incident struct {
	Title string `json:"title"`
	//Description string    `json:"description"`
	Township     string    `json:"township"`
	Intersection string    `json:"intersection"`
	Units        []string  `json:"units"`
	PubDateUtc   time.Time `json:"pubDateUtc"`
}
