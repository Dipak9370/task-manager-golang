package models

import "time"

type Task struct {
	ID          string `gorm:"primaryKey"`
	Title       string
	Description string
	Status      string
	UserID      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
