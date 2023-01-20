package posts

import (
	"context"
	"time"
)

type Post struct {
	Id        int
	Title     string `validate:"required,min=3,max=32"`
	Content   string `validate:"required,min=16"`
	CreatedAt time.Time  
	UpdatedAt time.Time
}

type PostsRepositoryInterface interface {
	GetAllPosts(ctx context.Context) ([]*Post, error)
	GetPostById(ctx context.Context, id int) (*Post, error)
	CreatePost(ctx context.Context, in *Post) (*Post, error)
}

type PostsUsecaseInterface interface {
	GetAllPosts(ctx context.Context) ([]*Post, error)
	GetPostById(ctx context.Context, id int) (*Post, error)
	CreatePost(ctx context.Context, in *Post) (*Post, error)
}