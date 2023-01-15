package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/sultanfariz/simple-grpc/domain"
	"github.com/sultanfariz/simple-grpc/repository/mysql"
	userRepository "github.com/sultanfariz/simple-grpc/repository/mysql/users"
	userUsecase "github.com/sultanfariz/simple-grpc/usecases/users"
	user "github.com/sultanfariz/simple-grpc/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type UserServerGrpc struct {
	userUsecase userUsecase.UsersUsecase
	user.UnimplementedAuthServiceServer
}

func NewUserServerGrpc(gserver *grpc.Server, userUcase userUsecase.UsersUsecase) {
	userServer := &UserServerGrpc{
		userUsecase: userUcase,
	}
	user.RegisterAuthServiceServer(gserver, userServer)
	reflection.Register(gserver)
}

func (s *UserServerGrpc) RegisterUser(ctx context.Context, in *user.RegisterUserInput) (*user.GenericResponse, error) {
	data := domain.User{
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
	userUsecase := userUsecase.NewUsersUsecase(userRepo, timeoutContext)

	
	grpcServer := grpc.NewServer()
	NewUserServerGrpc(grpcServer, *userUsecase)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}