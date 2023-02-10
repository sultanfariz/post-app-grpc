package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/sultanfariz/simple-grpc/client/internal/adapter"
	"github.com/sultanfariz/simple-grpc/client/internal/transport"
	"github.com/sultanfariz/simple-grpc/interface/grpc/post"
	"github.com/wagslane/go-rabbitmq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type gRPCServer struct {
	ctx        context.Context
	serverPort string
	clientPort string
}

func NewgRPCServer(ctx context.Context, serverPort, clientPort string) *gRPCServer {
	return &gRPCServer{
		ctx:        ctx,
		serverPort: serverPort,
		clientPort: clientPort,
	}
}

func (g *gRPCServer) Start() error {
	// start connection to gRPC server
	conn, err := grpc.DialContext(g.ctx, g.serverPort, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return errors.New("did not connect: %v" + err.Error())
	}

	// start connection to RabbitMQ
	amqpConn, err := rabbitmq.NewConn(
		"amqp://guest:guest@localhost",
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	server := grpc.NewServer()
	postServiceClient := post.NewPostServiceClient(conn)
	postClient := adapter.NewPostClient(postServiceClient)
	transport.NewPostHandler(server, *postClient, amqpConn)

	// Start gRPC server
	lis, err := net.Listen("tcp", g.clientPort)
	if err != nil {
		return errors.New("failed to listen: %v" + err.Error())
	}
	reflection.Register(server)
	fmt.Println("Serving gRPC client on", g.clientPort)
	return server.Serve(lis)
}
