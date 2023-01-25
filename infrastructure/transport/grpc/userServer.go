package grpc

import (
	"context"

	userDomain "github.com/sultanfariz/simple-grpc/domain/users"
	user "github.com/sultanfariz/simple-grpc/interface/grpc/user"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
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
}

func (s *UserServerGrpc) Register(ctx context.Context, in *user.RegisterRequest) (*user.RegisterResponse, error) {
	data := userDomain.User{
		Name:     in.GetName(),
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}

	userData, err := s.userUsecase.Register(ctx, &data)
	if err != nil {
		return nil, err
	}

	return &user.RegisterResponse{
		Meta: &user.GenericResponse{
			Status: "success",
			Message: "User created successfully",
		},
		User: &user.User{
			Id:        int32(userData.Id),
			Name:      userData.Name,
			Email:     userData.Email,
			CreatedAt: timestamppb.New(userData.CreatedAt),
			UpdatedAt: timestamppb.New(userData.UpdatedAt),
		},
	}, nil
}

func (s *UserServerGrpc) Login(ctx context.Context, in *user.LoginRequest) (*user.LoginResponse, error) {
	data := userDomain.User{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}

	token, err := s.userUsecase.Login(ctx, &data)
	if err != nil {
		return nil, err
	}

	return &user.LoginResponse{
		Meta: &user.GenericResponse{
			Status: "success",
			Message: "Login success",
		},
		Token: token,
	}, nil
}