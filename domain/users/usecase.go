package users

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UsersUsecase struct {
	UsersRepository UsersRepositoryInterface
	ContextTimeout  time.Duration
}

func NewUsersUsecase(ur UsersRepositoryInterface, timeout time.Duration) *UsersUsecase {
	return &UsersUsecase{
		UsersRepository: ur,
		ContextTimeout: timeout,
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