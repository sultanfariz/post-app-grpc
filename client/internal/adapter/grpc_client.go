package adapter

import (
	"context"
	"errors"

	"github.com/sultanfariz/simple-grpc/interface/grpc/post"
)

type PostClient struct {
	postClient post.PostServiceClient
}

func NewPostClient(postClient post.PostServiceClient) *PostClient {
	return &PostClient{
		postClient: postClient,
	}
}

func (p PostClient) GetAllPosts(ctx context.Context) ([]*post.Post, error) {
	posts, err := p.postClient.GetAllPosts(ctx, &post.GetAllPostsRequest{})
	if err != nil {
		return nil, err
	}

	if posts.Meta.Status != "success" {
		return nil, errors.New(posts.Meta.Message)
	}

	return posts.Posts, nil
}

func (p PostClient) GetPostById(ctx context.Context, id int32) (*post.Post, error) {
	post, err := p.postClient.GetPostById(ctx, &post.GetPostByIdRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	if post.Meta.Status != "success" {
		return nil, errors.New(post.Meta.Message)
	}

	return post.Post, nil
}
