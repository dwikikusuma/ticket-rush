package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dwikikusuma/ticket-rush/services/auth-service/internal/domain"
	userDB "github.com/dwikikusuma/ticket-rush/services/auth-service/internal/infra/postgres"
)

type userRepo struct {
	q *userDB.Queries
}

func NewUserRepo(q *sql.DB) domain.UserRepo {
	return &userRepo{q: userDB.New(q)}
}

func (r *userRepo) CreateUser(ctx context.Context, email, password string) error {
	param := userDB.CreateUserParams{
		Email:    email,
		Password: password,
	}

	_, err := r.q.CreateUser(ctx, param)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	userRow, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	user := domain.User{
		ID:        userRow.ID.String(),
		Email:     userRow.Email,
		CreatedAt: userRow.CreatedAt,
	}
	user.SetPassword(userRow.Password)
	return user, nil
}
