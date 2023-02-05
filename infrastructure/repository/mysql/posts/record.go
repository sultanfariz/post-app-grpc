package posts

import (
	"time"

	"github.com/sultanfariz/simple-grpc/infrastructure/repository/mysql/users"
)

type Post struct {
	Id        int        `gorm:"primaryKey"`
	Title     string     `gorm:"type:varchar(255);not null"`
	Content   string     `gorm:"type:text;not null"`
	Topic     string     `gorm:"type:varchar(255);not null"`
	UserId    int        `gorm:"column:user_id;not null"`
	User      users.User `gorm:"foreignKey:UserId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
