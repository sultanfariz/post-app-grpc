package users

import (
	"time"
)

type User struct {
	Id               int    `gorm:"primaryKey"`
	Name            string
	Email            string `gorm:"not null"`
	Password         string
	Role 		   string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}