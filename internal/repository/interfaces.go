package repository

import (
	"context"

	"github.com/fahrulrzi/score-match-backend/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserById(ctx context.Context, id uint) (*entity.User, error)
}

type CustomerRepository interface {
	Create(ctx context.Context, customer *entity.Customer) error
	GetAllCustomers(ctx context.Context) ([]entity.Customer, error)
	DeleteCustomer(ctx context.Context, id uint) error
	GetCustomerById(ctx context.Context, id uint) (*entity.Customer, error)
}

type TokenRepository interface {
	Create(ctx context.Context, token *entity.Token) error
	FindByToken(ctx context.Context, token string, tokenType entity.TokenType) (*entity.Token, error)
	Delete(ctx context.Context, id uint) error
	FindByUserID(ctx context.Context, userID uint) ([]*entity.Token, error)
	Exists(ctx context.Context, token string, tokenType entity.TokenType) (bool, error)
}
