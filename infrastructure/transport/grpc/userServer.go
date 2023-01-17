package grpc

import (
	"context"

	userDomain "github.com/sultanfariz/simple-grpc/domain/users"
	user "github.com/sultanfariz/simple-grpc/interface/grpc/user"
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