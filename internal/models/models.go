package models

import "time"

type User struct {
	ID        uint
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Notes     []Note
}

type Note struct {
	ID        uint
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint
}
