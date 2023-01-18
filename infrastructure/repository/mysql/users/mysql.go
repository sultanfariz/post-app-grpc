package users

import (
	"context"
	"errors"
	"time"

	model "github.com/sultanfariz/simple-grpc/domain/users"
	"golang.org/x/crypto/bcrypt"

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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := r.DBConnection.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) Login(ctx context.Context, email string, password string) (*model.User, error) {
	data := model.User{}
	if err := r.DBConnection.Where("email = ?", email).First(&data).Error; err != nil {
		return nil, errors.New("email or password is wrong")
	}

	if bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(password)) != nil {
		return nil, errors.New("email or password is wrong")
	}

	return &data, nil
}