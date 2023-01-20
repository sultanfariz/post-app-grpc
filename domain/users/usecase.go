package users

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
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

	validator := validator.New()
	if err := validator.Struct(user); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	if user.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}

	// hash password
	hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashed)

	// check if email already exists
	data, err := uu.UsersRepository.GetByEmail(ctx, user.Email)
	if err != nil && err.Error() != "record not found" {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if data != nil {
		return nil, status.Errorf(codes.AlreadyExists, "email already exists")
	}

	// insert user to db
	user, err = uu.UsersRepository.Insert(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return user, nil
}

func (uu *UsersUsecase) Login(ctx context.Context, user *User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, uu.ContextTimeout)
	defer cancel()

	validator := validator.New()
	if err := validator.Struct(User{
		Email:    user.Email,
		Password: user.Password,
	}); err != nil {
		return "", status.Errorf(codes.InvalidArgument, err.Error())
	}

	// check if user exists
	data, err := uu.UsersRepository.GetByEmail(ctx, user.Email)
	if err != nil && err.Error() != "record not found" {
		return "", status.Errorf(codes.Internal, err.Error())
	}
	if data == nil {
		return "", status.Errorf(codes.NotFound, "user not found")
	}

	// check if password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(user.Password)); err != nil {
		return "", status.Errorf(codes.InvalidArgument, "invalid email or password")
	}

	// generate token
	token, err := uu.jwtConfig.GenerateToken(user.Id, user.Email)
	if err != nil {
		return "", status.Errorf(codes.Internal, err.Error())
	}

	return token, nil
}