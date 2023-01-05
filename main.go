package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/sultanfariz/simple-grpc/hello"
	"google.golang.org/grpc"
)

type Server struct {
	hello.UnimplementedHelloWorldServer
}

// kita meng-implementasikan method SayHello
func (s *Server) SayHello(ctx context.Context, in *hello.SayHelloRequest) (*hello.SayHelloResponse, error) {
	return &hello.SayHelloResponse{
		Message: "Selamat datang " + in.GetName(),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50001")

	if err != nil {
		log.Fatalf("failed listen: %v", err)
	}

	fmt.Println("Server running on port :50001")

	srv := Server{}

	grpcServer := grpc.NewServer()
	hello.RegisterHelloWorldServer(grpcServer, &srv)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}