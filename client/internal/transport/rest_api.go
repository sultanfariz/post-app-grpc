package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sultanfariz/simple-grpc/client/internal/adapter"
	postPB "github.com/sultanfariz/simple-grpc/client/pkg/grpc/post"
	rabbitmq "github.com/wagslane/go-rabbitmq"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PostHandler struct {
	PostClient adapter.PostClient
	conn       *rabbitmq.Conn
	postPB.UnimplementedPostClientServiceServer
}

func NewPostHandler(gserver *grpc.Server, postClient adapter.PostClient, conn *rabbitmq.Conn) {
	postHandler := &PostHandler{
		PostClient: postClient,
		conn:       conn,
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

func (p PostHandler) SubscribePostByTopic(in *postPB.SubscribePostByTopicRequest, stream postPB.PostClientService_SubscribePostByTopicServer) error {
	topic := in.GetTopic()

	// create consumer
	consumer, err := rabbitmq.NewConsumer(
		p.conn,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			var data map[string]interface{}
			err := json.Unmarshal(d.Body, &data)
			if err != nil {
				return rabbitmq.NackRequeue
			}

			// set type assertion for id, createdAt, and updatedAt
			var id float64
			var ok bool
			if x, found := data["Id"]; found {
				if id, ok = x.(float64); !ok {
					fmt.Println("id is not float64")
					return rabbitmq.NackRequeue
				}
			} else {
				fmt.Println("id is not found")
				return rabbitmq.NackRequeue
			}
			createdAt, err := time.Parse(time.RFC3339, data["CreatedAt"].(string))
			if err != nil {
				fmt.Println("error parsing time: ", err)
				return rabbitmq.NackRequeue
			}
			updatedAt, err := time.Parse(time.RFC3339, data["UpdatedAt"].(string))
			if err != nil {
				fmt.Println("error parsing time: ", err)
				return rabbitmq.NackRequeue
			}

			// set stream response
			post := &postPB.Post{
				Id:        int32(id),
				Title:     data["Title"].(string),
				Content:   data["Content"].(string),
				Topic:     data["Topic"].(string),
				CreatedAt: timestamppb.New(createdAt),
				UpdatedAt: timestamppb.New(updatedAt),
			}

			response := &postPB.SubscribePostByTopicResponse{
				Meta: &postPB.GenericResponse{
					Status:  "success",
					Message: "success",
				},
				Post: post,
			}
			if err := stream.Send(response); err != nil {
				fmt.Println("error sending stream: ", err)
				return rabbitmq.NackRequeue
			}
			return rabbitmq.Ack
		},
		"grpc:"+topic,
		rabbitmq.WithConsumerOptionsRoutingKey("grpc:"+topic),
		rabbitmq.WithConsumerOptionsExchangeName("events"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		return err
	}
	defer consumer.Close()

	// block main thread - wait for shutdown signal
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("stopping consumer")

	return nil
}
