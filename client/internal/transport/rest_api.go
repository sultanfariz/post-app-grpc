package transport

import (
	"context"
	"fmt"

	"github.com/sultanfariz/simple-grpc/client/internal/adapter"
	postPB "github.com/sultanfariz/simple-grpc/client/pkg/grpc/post"
	rabbitmq "github.com/wagslane/go-rabbitmq"
	"google.golang.org/grpc"
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
	// func (p PostHandler) SubscribePostByTopic(ctx context.Context, in *postPB.SubscribePostByTopicRequest) (postPB.PostClientService_SubscribePostByTopicClient, error) {
	// func (p PostHandler) SubscribePostByTopic(in *postPB.SubscribePostByTopicRequest) (*postPB.PostClientService_SubscribePostByTopicClient, error) {
	topic := in.GetTopic()
	fmt.Println("topic: ", topic)
	consumer, err := rabbitmq.NewConsumer(
		p.conn,
		// func(d rabbitmq.Delivery) (action rabbitmq.Action) {
		func(d rabbitmq.Delivery) rabbitmq.Action {
			// post := &Post{}
			// err := json.Unmarshal(d.Body, post)
			// if err != nil {
			// 	log.Println(err)
			// 	return rabbitmq.Reject
			// }

			return rabbitmq.Ack
		},
		// fmt.Sprintf("post:%s", topic),
		"grpc_queue",
		rabbitmq.WithConsumerOptionsRoutingKey("grpc_queue"),
		rabbitmq.WithConsumerOptionsExchangeName("events"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		return err
		// return nil, err
	}
	defer consumer.Close()

	// sigs := make(chan os.Signal, 1)
	// done := make(chan bool, 1)
	// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// <-sigs

	// go func() {
	// 	sig := <-sigs
	// 	fmt.Println()
	// 	fmt.Println(sig)
	// 	done <- true
	// }()

	// fmt.Println("awaiting signal")
	// <-done
	// fmt.Println("stopping consumer")

	return nil
	// return nil, nil

	// return &postPB.SubscribePostByTopicResponse{
	// 	Meta: &postPB.GenericResponse{
	// 		Status:  "success",
	// 		Message: "success",
	// 	},
	// 	// Post: &postPB.Post{
	// 	// 	Id:        post.Id,
	// 	// 	Title:     post.Title,
	// 	// 	Content:   post.Content,
	// 	// 	Topic:     post.Topic,
	// 	// 	CreatedAt: post.CreatedAt,
	// 	// 	UpdatedAt: post.UpdatedAt,
	// 	// },
	// }, nil
}
