package data

import "time"

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`           // force it to be string
	Runtime   Runtime     `json:"runtime,omitempty"` // chose to parse field as a string
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}
