package users

import (
	"context"
	"errors"
	"time"

	domain "github.com/sultanfariz/simple-grpc/domain"

	"gorm.io/gorm"
)

type UsersRepository struct {
	DBConnection *gorm.DB
}

func NewUsersRepository(db *gorm.DB) *UsersRepository {
	return &UsersRepository{
		DBConnection: db,
	}
}

func (r *UsersRepository) Register(ctx context.Context, in *domain.User) (*domain.User, error) {
	data := domain.User{}
	if err := r.DBConnection.Where("email = ?", in.Email).First(&data).Error; err == nil {
		return nil, errors.New("email already registered")
	}

	user := domain.User{
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
		Role:     in.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := r.DBConnection.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}