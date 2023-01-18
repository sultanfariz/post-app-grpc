package users

import (
	"time"
)

type User struct {
	Id               int    `gorm:"primaryKey"`
	Name             string	`gorm:"not null"`
	Email            string `gorm:"not null"`
	Password         string `gorm:"not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}