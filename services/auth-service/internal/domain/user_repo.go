package domain

import "context"

type UserRepo interface {
	CreateUser(ctx context.Context, email, password string) error
	GetUserByEmail(ctx context.Context, email string) (User, error)
}
