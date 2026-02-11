package models

type User struct {
	ID       string `gorm:"primaryKey"`
	Email    string `gorm:"unique"`
	Password string
	Role     string
}
