package users

import (
	"context"
	"time"
)

type User struct {
	Id        int    
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time  
	UpdatedAt time.Time
}

type UsersRepositoryInterface interface {
	Register(ctx context.Context, in *User) (*User, error)
	Login(ctx context.Context, email string, password string) (*User, error)
}

type UsersUsecaseInterface interface {
	Register(ctx context.Context, in *User) (*User, error)
	Login(ctx context.Context, in *User) (*User, error)
}