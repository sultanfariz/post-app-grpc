package users

import (
	"context"
	"errors"
	"time"

	"github.com/sultanfariz/simple-grpc/domain"
	"golang.org/x/crypto/bcrypt"
)

type UsersUsecase struct {
	UsersRepository domain.UsersRepository
	ContextTimeout  time.Duration
}

func NewUsersUsecase(ur domain.UsersRepository, timeout time.Duration) *UsersUsecase {
	return &UsersUsecase{
		UsersRepository: ur,
		ContextTimeout: timeout,
	}
}

func (uu *UsersUsecase) Register(ctx context.Context, user *domain.User) (*domain.User, error) {
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