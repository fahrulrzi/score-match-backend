package usecase

import (
	"context"

	"github.com/fahrulrzi/score-match-backend/internal/entity"
	"github.com/fahrulrzi/score-match-backend/internal/repository"
)

type customerUseCase struct {
	customerRepo repository.CustomerRepository
}

func NewCustomerUseCase(customerRepo repository.CustomerRepository) CustomerUseCase {
	return &customerUseCase{
		customerRepo: customerRepo,
	}
}

// GetDBRScore implements CustomerUseCase.
func (c *customerUseCase) GetDBRScore(ctx context.Context, dbr float64) int {
	for _, r := range entity.DBRScoreRanges {
		if dbr >= r.Min && dbr <= r.Max {
			return r.Score
		}
	}
	return 0
}

// GetScore implements CustomerUseCase.
func (c *customerUseCase) GetScore(ctx context.Context, category string, value string) int {
	var scoreMap map[string]int

	switch category {
	case "pekerjaan":
		scoreMap = entity.JobScore
	case "lama_kerja":
		scoreMap = entity.LengthOfWorkScore
	case "tujuan":
		scoreMap = entity.Purpose
	case "jaminan":
		scoreMap = entity.Collateral
	case "status_perkawinan":
		scoreMap = entity.MaritalStatusScore
	default:
		return 0
	}

	return scoreMap[value]
}

// GetFinalScore implements CustomerUseCase.
func (c *customerUseCase) GetFinalScore(ctx context.Context, customer *entity.CustomerRequest) int {
	jobScore := c.GetScore(ctx, "pekerjaan", customer.Job)
	lengthOfWorkScore := c.GetScore(ctx, "lama_kerja", customer.LengthOfWork)
	purposeScore := c.GetScore(ctx, "tujuan", customer.Purpose)
	collateralScore := c.GetScore(ctx, "jaminan", customer.Collateral)
	maritalStatusScore := c.GetScore(ctx, "status_perkawinan", customer.MaritalStatus)
	DBR := customer.Installment / customer.Income * 100

	DBRScore := c.GetDBRScore(ctx, float64(DBR))

	finalScore := jobScore + lengthOfWorkScore + purposeScore + collateralScore + maritalStatusScore + DBRScore/6

	return finalScore
}

// Create implements CustomerUseCase.
func (c *customerUseCase) Create(ctx context.Context, customer *entity.Customer) error {
	customers := &entity.Customer{
		Username:   customer.Username,
		Job:    customer.Job,
		Income: customer.Income,
		Age:    customer.Age,
		Score:  customer.Score,
	}

	err := c.customerRepo.Create(ctx, customers)
	if err != nil {
		return err
	}

	return nil
}

// // DeleteCustomer implements CustomerUseCase.
// func (c *customerUseCase) DeleteCustomer(ctx context.Context, id uint) error {
// 	panic("unimplemented")
// }

// // GetAllCustomers implements CustomerUseCase.
// func (c *customerUseCase) GetAllCustomers(ctx context.Context) ([]entity.Customer, error) {
// 	panic("unimplemented")
// }

// // GetCustomerById implements CustomerUseCase.
// func (c *customerUseCase) GetCustomerById(ctx context.Context, id uint) (*entity.Customer, error) {
// 	panic("unimplemented")
// }
