package users

import (
	"context"
	"errors"
	"time"

	"github.com/sultanfariz/simple-grpc/infrastructure/commons"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UsersUsecase struct {
	UsersRepository UsersRepositoryInterface
	ContextTimeout  time.Duration
	jwtConfig       *commons.ConfigJWT
}

func NewUsersUsecase(ur UsersRepositoryInterface, timeout time.Duration, jwtConfig *commons.ConfigJWT) *UsersUsecase {
	return &UsersUsecase{
		UsersRepository: ur,
		ContextTimeout: timeout,
		jwtConfig: jwtConfig,
	}
}

func (uu *UsersUsecase) Register(ctx context.Context, user *User) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, uu.ContextTimeout)
	defer cancel()

	if user.Name == "" {
		return user, errors.New("name is required")
	}
	if user.Email == "" || user.Password == "" {
		return user, errors.New("invalid email or password")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashed)

	user, err := uu.UsersRepository.Register(ctx, user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (uu *UsersUsecase) Login(ctx context.Context, user *User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, uu.ContextTimeout)
	defer cancel()

	if user.Email == "" || user.Password == "" {
		return "", status.Errorf(codes.InvalidArgument, "invalid email or password")
	}

	user, err := uu.UsersRepository.Login(ctx, user.Email, user.Password)
	if err != nil {
		return "", err
	}
	token, err := uu.jwtConfig.GenerateToken(user.Id, user.Email)
	if err != nil {
		return "", status.Errorf(codes.Internal, err.Error())
	}

	return token, nil
}