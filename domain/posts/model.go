package posts

import (
	"context"
	"time"
)

type Post struct {
	Id        int
	Title     string `validate:"required,min=3,max=256"`
	Content   string `validate:"required,min=16"`
	Topic     string `validate:"required,min=3,max=128"`
	UserId    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostsRepositoryInterface interface {
	GetAll(ctx context.Context) ([]*Post, error)
	GetById(ctx context.Context, id int) (*Post, error)
	Insert(ctx context.Context, in *Post) (*Post, error)
	Delete(ctx context.Context, id int) error
}

type PostsUsecaseInterface interface {
	GetAllPosts(ctx context.Context) ([]*Post, error)
	GetPostById(ctx context.Context, id int) (*Post, error)
	CreatePost(ctx context.Context, in *Post) (*Post, error)
	DeletePost(ctx context.Context, id int) error
}
