package query

import "time"

type Document struct {
	ID        string
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
