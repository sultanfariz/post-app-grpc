package posts

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sultanfariz/simple-grpc/domain/users"
	rabbitmq "github.com/wagslane/go-rabbitmq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostsUsecase struct {
	PostsRepository PostsRepositoryInterface
	UsersRepository users.UsersRepositoryInterface
	ContextTimeout  time.Duration
	Publisher       *rabbitmq.Publisher
}

func NewPostsUsecase(pr PostsRepositoryInterface, ur users.UsersRepositoryInterface, timeout time.Duration, publisher *rabbitmq.Publisher) *PostsUsecase {
	return &PostsUsecase{
		PostsRepository: pr,
		UsersRepository: ur,
		ContextTimeout:  timeout,
		Publisher:       publisher,
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

	userEmail, ok := ctx.Value("email").(string)
	if userEmail == "" || !ok {
		return post, status.Errorf(codes.Unauthenticated, "user not found")
	}

	user, err := pu.UsersRepository.GetByEmail(ctx, userEmail)
	if user == nil || err != nil {
		return post, status.Errorf(codes.Unauthenticated, "user not found")
	}

	post.UserId = user.Id

	validator := validator.New()
	if err := validator.Struct(post); err != nil {
		return post, status.Errorf(codes.InvalidArgument, err.Error())
	}

	result, err := pu.PostsRepository.Insert(ctx, post)
	if err != nil {
		return post, err
	}

	// convert post data to []byte
	data, err := json.Marshal(post)
	if err != nil {
		return post, err
	}

	// publish to rabbitmq
	err = pu.Publisher.Publish(
		data,
		[]string{"grpc:" + post.Topic},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange("events"),
	)
	if err != nil {
		log.Println(err)
	}

	return result, nil
}

func (pu *PostsUsecase) DeletePost(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, pu.ContextTimeout)
	defer cancel()

	post, err := pu.PostsRepository.GetById(ctx, id)
	if err != nil {
		return err
	}
	if post == nil {
		return status.Errorf(codes.NotFound, "post not found")
	}

	// check if user is the owner of the post
	userEmail, ok := ctx.Value("email").(string)
	if userEmail == "" || !ok {
		return status.Errorf(codes.Unauthenticated, "user not found")
	}

	user, err := pu.UsersRepository.GetByEmail(ctx, userEmail)
	if user == nil || err != nil {
		return status.Errorf(codes.Unauthenticated, "user not found")
	}

	if post.UserId != user.Id {
		return status.Errorf(codes.PermissionDenied, "you are not the owner of this post")
	}

	err = pu.PostsRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
