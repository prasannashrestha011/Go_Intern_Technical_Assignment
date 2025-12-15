package models

import "time"

type Order struct {
	ID         uint64
	UserID     uint64
	Amount     float64
	Status     string
	CreatedAt time.Time
}