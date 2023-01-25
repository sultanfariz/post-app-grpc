package main

import (
	"fmt"
	"log"
	"net"
	"time"

	postDomain "github.com/sultanfariz/simple-grpc/domain/posts"
	userDomain "github.com/sultanfariz/simple-grpc/domain/users"
	"github.com/sultanfariz/simple-grpc/infrastructure/commons"
	"github.com/sultanfariz/simple-grpc/infrastructure/repository/mysql"
	postsRepository "github.com/sultanfariz/simple-grpc/infrastructure/repository/mysql/posts"
	userRepository "github.com/sultanfariz/simple-grpc/infrastructure/repository/mysql/users"
	grpcServerController "github.com/sultanfariz/simple-grpc/infrastructure/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", ":50001")
	timeoutContext := time.Duration(10 * time.Second)

	if err != nil {
		log.Fatalf("failed listen: %v", err)
	}
	fmt.Println("Server running on port :50001")

	configJWT := commons.ConfigJWT{
		SecretJWT:       "thisIs45ecretKey",
		ExpiresDuration: 72,
	}

	db := mysql.InitDB()
	userRepo := userRepository.NewUsersRepository(db)
	userUsecase := userDomain.NewUsersUsecase(userRepo, timeoutContext, &configJWT)
	postRepo := postsRepository.NewPostsRepository(db)
	postUsecase := postDomain.NewPostsUsecase(postRepo, timeoutContext)
	
	serverOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpcServerController.JWTInterceptor),
	}
	grpcServer := grpc.NewServer(serverOpts...)

	grpcServerController.NewUserServerGrpc(grpcServer, *userUsecase)
	grpcServerController.NewPostServerGrpc(grpcServer, *postUsecase)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}