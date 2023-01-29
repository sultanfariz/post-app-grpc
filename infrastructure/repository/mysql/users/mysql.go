package users

import (
	"context"
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

func (r *UsersRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	data := model.User{}
	if err := r.DBConnection.Where("email = ?", email).First(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *UsersRepository) GetById(ctx context.Context, id int) (*model.User, error) {
	data := model.User{}
	if err := r.DBConnection.Where("id = ?", id).First(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *UsersRepository) Insert(ctx context.Context, user *model.User) (*model.User, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := r.DBConnection.Create(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
