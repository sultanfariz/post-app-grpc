package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	userDomain "github.com/sultanfariz/simple-grpc/domain/users"
	"github.com/sultanfariz/simple-grpc/infrastructure/repository/mysql"
	userRepository "github.com/sultanfariz/simple-grpc/infrastructure/repository/mysql/users"
	user "github.com/sultanfariz/simple-grpc/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type UserServerGrpc struct {
	userUsecase userDomain.UsersUsecase
	user.UnimplementedAuthServiceServer
}

func NewUserServerGrpc(gserver *grpc.Server, userUcase userDomain.UsersUsecase) {
	userServer := &UserServerGrpc{
		userUsecase: userUcase,
	}
	user.RegisterAuthServiceServer(gserver, userServer)
	reflection.Register(gserver)
}

func (s *UserServerGrpc) RegisterUser(ctx context.Context, in *user.RegisterUserInput) (*user.GenericResponse, error) {
	data := userDomain.User{
		Name:     in.GetName(),
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}

	_, err := s.userUsecase.Register(ctx, &data)
	if err != nil {
		return nil, err
	}

	return &user.GenericResponse{
		Status: "success",
		Message: "User berhasil didaftarkan",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50001")
	timeoutContext := time.Duration(10 * time.Second)

	if err != nil {
		log.Fatalf("failed listen: %v", err)
	}
	fmt.Println("Server running on port :50001")

	db := mysql.InitDB()
	userRepo := userRepository.NewUsersRepository(db)
	userUsecase := userDomain.NewUsersUsecase(userRepo, timeoutContext)

	
	grpcServer := grpc.NewServer()
	NewUserServerGrpc(grpcServer, *userUsecase)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}