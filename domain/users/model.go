package users

import (
	"context"
	"time"
)

type User struct {
	Id        int    
	Name      string 
	Email     string `validate:"required,email"`
	Password  string `validate:"required,min=8,max=32"`
	CreatedAt time.Time  
	UpdatedAt time.Time
}

type UsersRepositoryInterface interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetById(ctx context.Context, id int) (*User, error)
	Insert(ctx context.Context, user *User) (*User, error)
}

type UsersUsecaseInterface interface {
	Register(ctx context.Context, in *User) (*User, error)
	Login(ctx context.Context, in *User) (*User, error)
}