package grpc

import (
	"context"

	postDomain "github.com/sultanfariz/simple-grpc/domain/posts"
	post "github.com/sultanfariz/simple-grpc/interface/grpc/post"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PostServerGrpc struct {
	postUsecase postDomain.PostsUsecase
	post.UnimplementedPostServiceServer
}

func NewPostServerGrpc(gserver *grpc.Server, postUcase postDomain.PostsUsecase) {
	postServer := &PostServerGrpc{
		postUsecase: postUcase,
	}
	post.RegisterPostServiceServer(gserver, postServer)
}

func (s *PostServerGrpc) GetAllPosts(ctx context.Context, in *post.GetAllPostsRequest) (*post.GetAllPostsResponse, error) {
	postData, err := s.postUsecase.GetAllPosts(ctx)
	if err != nil {
		return nil, err
	}

	var posts []*post.Post
	for _, data := range postData {
		posts = append(posts, &post.Post{
			Id:        int32(data.Id),
			Title:     data.Title,
			Content:   data.Content,
			CreatedAt: timestamppb.New(data.CreatedAt),
			UpdatedAt: timestamppb.New(data.UpdatedAt),
		})
	}

	return &post.GetAllPostsResponse{
		Meta: &post.GenericResponse{
			Status:  "success",
			Message: "Successfully get all posts",
		},
		Posts: posts,
	}, nil
}

func (s *PostServerGrpc) GetPostById(ctx context.Context, in *post.GetPostByIdRequest) (*post.GetPostByIdResponse, error) {
	postData, err := s.postUsecase.GetPostById(ctx, int(in.GetId()))
	if err != nil {
		return nil, err
	}

	return &post.GetPostByIdResponse{
		Meta: &post.GenericResponse{
			Status:  "success",
			Message: "Successfully get post by id",
		},
		Post: &post.Post{
			Id:        int32(postData.Id),
			Title:     postData.Title,
			Content:   postData.Content,
			CreatedAt: timestamppb.New(postData.CreatedAt),
			UpdatedAt: timestamppb.New(postData.UpdatedAt),
		},
	}, nil
}

func (s *PostServerGrpc) CreatePost(ctx context.Context, in *post.CreatePostRequest) (*post.CreatePostResponse, error) {
	data := postDomain.Post{
		Title:   in.GetTitle(),
		Content: in.GetContent(),
	}

	postData, err := s.postUsecase.CreatePost(ctx, &data)
	if err != nil {
		return nil, err
	}

	return &post.CreatePostResponse{
		Meta: &post.GenericResponse{
			Status:  "success",
			Message: "Post created successfully",
		},
		Post: &post.Post{
			Id:        int32(postData.Id),
			Title:     postData.Title,
			Content:   postData.Content,
			CreatedAt: timestamppb.New(postData.CreatedAt),
			UpdatedAt: timestamppb.New(postData.UpdatedAt),
		},
	}, nil
}

func (s *PostServerGrpc) DeletePost(ctx context.Context, in *post.DeletePostRequest) (*post.DeletePostResponse, error) {
	err := s.postUsecase.DeletePost(ctx, int(in.GetId()))
	if err != nil {
		return nil, err
	}

	return &post.DeletePostResponse{
		Meta: &post.GenericResponse{
			Status:  "success",
			Message: "Post deleted successfully",
		},
	}, nil
}
