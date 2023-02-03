package transport

import (
	"context"

	"github.com/sultanfariz/simple-grpc/client/internal/adapter"
	postPB "github.com/sultanfariz/simple-grpc/client/pkg/grpc/post"
	"google.golang.org/grpc"
)

type PostHandler struct {
	PostClient adapter.PostClient
	postPB.UnimplementedPostClientServiceServer
}

func NewPostHandler(gserver *grpc.Server, postClient adapter.PostClient) {
	postHandler := &PostHandler{
		PostClient: postClient,
	}
	postPB.RegisterPostClientServiceServer(gserver, *postHandler)
}

func (p PostHandler) GetAllPosts(ctx context.Context, in *postPB.GetAllPostsRequest) (*postPB.GetAllPostsResponse, error) {
	posts, err := p.PostClient.GetAllPosts(ctx)
	if err != nil {
		return nil, err
	}

	postsRes := make([]*postPB.Post, len(posts))
	for i, post := range posts {
		postsRes[i] = &postPB.Post{
			Id:        post.Id,
			Title:     post.Title,
			Content:   post.Content,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		}
	}

	return &postPB.GetAllPostsResponse{
		Meta: &postPB.GenericResponse{
			Status:  "success",
			Message: "Successfully get all posts",
		},
		Posts: postsRes,
	}, nil
}

func (p PostHandler) GetPostById(ctx context.Context, in *postPB.GetPostByIdRequest) (*postPB.GetPostByIdResponse, error) {
	post, err := p.PostClient.GetPostById(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &postPB.GetPostByIdResponse{
		Meta: &postPB.GenericResponse{
			Status:  "success",
			Message: "success",
		},
		Post: &postPB.Post{
			Id:        post.Id,
			Title:     post.Title,
			Content:   post.Content,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		},
	}, nil
}
