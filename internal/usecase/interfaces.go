package usecase

import (
	"context"

	"github.com/fahrulrzi/score-match-backend/internal/entity"
)

type AuthUseCase interface {
	Login(ctx context.Context, userLogin *entity.UserLoginRequest) (*entity.User, string, error)
	Register(ctx context.Context, userRegister *entity.UserRegisterRequest) (*entity.User, string, error)
}

type CustomerUseCase interface {
	Create(ctx context.Context, customer *entity.Customer) error
	// GetAllCustomers(ctx context.Context) ([]entity.Customer, error)
	GetFinalScore(ctx context.Context, customer *entity.CustomerRequest) (*entity.Customer, error)
	// DeleteCustomer(ctx context.Context, id uint) error
	// GetCustomerById(ctx context.Context, id uint) (*entity.Customer, error)
	GetScore(ctx context.Context, category, value string) int
	GetDBRScore(ctx context.Context, dbr float64) int
	GetInform(ctx context.Context, score int) (*entity.CustomerScoreResponse, error)
}
