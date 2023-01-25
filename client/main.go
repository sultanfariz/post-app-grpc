package main

import (
	"context"
	"log"

	"github.com/sultanfariz/simple-grpc/interface/grpc/post"
	"google.golang.org/grpc"
)

const (
	port = ":50001"
)

func main() {
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	postClient := post.NewPostServiceClient(conn)

	// Contact the server and print out its response.
	ctx := context.Background()
	post, err := postClient.GetPostById(ctx, &post.GetPostByIdRequest{Id: 4})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %+v", post)
}
