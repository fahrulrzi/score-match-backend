package usecase

import (
	"context"
	"fmt"
	"strings"

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

func FormatRupiah(amount float64) string {
	// Ubah ke int biar nggak ribet sama desimal
	amountInt := int64(amount)
	result := fmt.Sprintf("Rp. %s", formatWithThousandSeparator(amountInt))
	return result
}

func formatWithThousandSeparator(n int64) string {
	str := fmt.Sprintf("%d", n)
	var result []string
	for len(str) > 3 {
		result = append([]string{str[len(str)-3:]}, result...)
		str = str[:len(str)-3]
	}
	if str != "" {
		result = append([]string{str}, result...)
	}
	return strings.Join(result, ".")
}

// GetInform implements CustomerUseCase.
func (c *customerUseCase) GetInform(ctx context.Context, score int, dbr float64) (*entity.CustomerScoreResponse, error) {
	var response entity.CustomerScoreResponse
	switch {
	case score >= 80 && score <= 100:
		response.Score = score
		response.Status = "Disetujui"
		response.Describe = "SKOR Kredit Adalah Baik. Kredit yang diajukan diterima sepenuhnya karena memenuhi semua parameter dan persyaratan yang telah ditentukan."
	case score >= 60 && score < 80:
		response.Score = score
		response.Status = "Diproses"
		response.Describe = `Pengajuan kredit dipertimbangkan. Bisa ditolak karena skor dan DBR tidak sesuai, bisa diterima maksimal sebesar ` + FormatRupiah(dbr)
	default:
		response.Score = score
		response.Status = "Ditolak"
		response.Describe = "SKOR Kredit Adalah Buruk. Kredit yang diajukan ditolak karena tidak memenuhi semua parameter dan persyaratan yang telah ditentukan."
	}

	return &response, nil
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
func (c *customerUseCase) GetFinalScore(ctx context.Context, customer *entity.CustomerRequest) (*entity.Customer, error) {
	var final entity.Customer
	jobScore := c.GetScore(ctx, "pekerjaan", customer.Job)
	lengthOfWorkScore := c.GetScore(ctx, "lama_kerja", customer.LengthOfWork)
	purposeScore := c.GetScore(ctx, "tujuan", customer.Purpose)
	collateralScore := c.GetScore(ctx, "jaminan", customer.Collateral)
	maritalStatusScore := c.GetScore(ctx, "status_perkawinan", customer.MaritalStatus)
	DBR := (float64(customer.Installment) / float64(customer.Income)) * 100
	DBRScore := c.GetDBRScore(ctx, float64(DBR))

	maxRent := (40 - DBR) / 100 * float64(customer.Income)

	final.Username = customer.Username
	final.Job = customer.Job
	final.Income = customer.Income
	final.Age = customer.Age
	final.Score = (jobScore + lengthOfWorkScore + purposeScore + collateralScore + maritalStatusScore + DBRScore) / 6

	inform, err := c.GetInform(ctx, final.Score, maxRent)
	if err != nil {
		return nil, err
	}
	final.Status = inform.Status
	final.Describe = inform.Describe

	return &final, nil
}

// Create implements CustomerUseCase.
func (c *customerUseCase) Create(ctx context.Context, customer *entity.Customer) error {
	customers := &entity.Customer{
		Username: customer.Username,
		Job:      customer.Job,
		Income:   customer.Income,
		Age:      customer.Age,
		Score:    customer.Score,
		Status:   customer.Status,
		Describe: customer.Describe,
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

// GetAllCustomers implements CustomerUseCase.
func (c *customerUseCase) GetAllCustomers(ctx context.Context) ([]entity.Customer, error) {
	customers, err := c.customerRepo.GetAllCustomers(ctx)
	if err != nil {
		return nil, err
	}
	return customers, nil
}

// // GetCustomerById implements CustomerUseCase.
// func (c *customerUseCase) GetCustomerById(ctx context.Context, id uint) (*entity.Customer, error) {
// 	panic("unimplemented")
// }
