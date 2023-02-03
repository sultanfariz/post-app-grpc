package post

import (
	"context"
	"time"
)

type Post struct {
	Id        int
	Title     string `validate:"required,min=3,max=128"`
	Content   string `validate:"required,min=16"`
	UserId    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostsUsecaseInterface interface {
	GetAllPosts(ctx context.Context) ([]*Post, error)
	GetPostById(ctx context.Context, id int) (*Post, error)
}
