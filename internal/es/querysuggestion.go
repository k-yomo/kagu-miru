package es

import "time"

type QuerySuggestion struct {
	Query     string    `json:"query"`
	CreatedAt time.Time `json:"createdAt"`
}
