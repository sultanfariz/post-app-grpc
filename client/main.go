package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sultanfariz/simple-grpc/client/internal/adapter"
	"github.com/sultanfariz/simple-grpc/client/internal/transport"
	postPB "github.com/sultanfariz/simple-grpc/client/pkg/grpc/post"
	"github.com/sultanfariz/simple-grpc/interface/grpc/post"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	serverPort  = ":50001"
	clientPort  = ":50051"
	gatewayPort = ":8080"
)

func main() {
	// conn, err := grpc.Dial(serverPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.DialContext(context.Background(), serverPort, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// defer conn.Close()

	postServiceClient := post.NewPostServiceClient(conn)
	postClient := adapter.NewPostClient(postServiceClient)
	// postHandler := transport.NewPostHandler(*postClient)

	ctx := context.Background()
	server := grpc.NewServer()
	transport.NewPostHandler(server, *postClient)

	// Start gRPC server
	lis, err := net.Listen("tcp", clientPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	reflection.Register(server)
	go func() {
		fmt.Println("Serving gRPC client on", clientPort)
		server.Serve(lis)
	}()
	// log.Printf("Serving gRPC client on %s", clientPort)

	// Start HTTP server
	// postPB.RegisterPostClientServiceServer(server, *postHandler)
	gatewayConn, err := grpc.DialContext(ctx, clientPort, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// defer gatewayConn.Close()

	muxServer := runtime.NewServeMux()
	if err = postPB.RegisterPostClientServiceHandler(ctx, muxServer, gatewayConn); err != nil {
		log.Fatalf("did not register: %v", err)
	}
	httpServer := &http.Server{
		Addr:    gatewayPort,
		Handler: muxServer,
	}
	fmt.Println("Serving gRPC-Gateway on", gatewayPort)
	httpServer.ListenAndServe()

}
