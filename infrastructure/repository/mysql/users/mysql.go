package users

import (
	"context"
	"errors"
	"time"

	model "github.com/sultanfariz/simple-grpc/domain/users"

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

func (r *UsersRepository) Register(ctx context.Context, in *model.User) (*model.User, error) {
	data := model.User{}
	if err := r.DBConnection.Where("email = ?", in.Email).First(&data).Error; err == nil {
		return nil, errors.New("email already registered")
	}

	user := model.User{
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