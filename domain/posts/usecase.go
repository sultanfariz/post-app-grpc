package posts

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostsUsecase struct {
	PostsRepository PostsRepositoryInterface
	ContextTimeout  time.Duration
}

func NewPostsUsecase(pr PostsRepositoryInterface, timeout time.Duration) *PostsUsecase {
	return &PostsUsecase{
		PostsRepository: pr,
		ContextTimeout: timeout,
	}
}

func (pu *PostsUsecase) GetAllPosts(ctx context.Context) ([]*Post, error) {
	ctx, cancel := context.WithTimeout(ctx, pu.ContextTimeout)
	defer cancel()

	posts, err := pu.PostsRepository.GetAll(ctx)
	if err != nil {
		return posts, err
	}

	return posts, nil
}

func (pu *PostsUsecase) GetPostById(ctx context.Context, id int) (*Post, error) {
	ctx, cancel := context.WithTimeout(ctx, pu.ContextTimeout)
	defer cancel()

	post, err := pu.PostsRepository.GetById(ctx, id)
	if err != nil {
		return post, err
	}

	return post, nil
}

func (pu *PostsUsecase) CreatePost(ctx context.Context, post *Post) (*Post, error) {
	ctx, cancel := context.WithTimeout(ctx, pu.ContextTimeout)
	defer cancel()

	validator := validator.New()
	if err := validator.Struct(post); err != nil {
		return post, status.Errorf(codes.InvalidArgument, err.Error())
	}

	post, err := pu.PostsRepository.Insert(ctx, post)
	if err != nil {
		return post, err
	}

	return post, nil
}
