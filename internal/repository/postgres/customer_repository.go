package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/fahrulrzi/score-match-backend/internal/entity"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{
		db: db,
	}
}

func (r *CustomerRepository) Create(ctx context.Context, customer *entity.Customer) error {
	query := `
	INSERT INTO customers (name, job, income, age, score, status, describe, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id
	`

	now := time.Now()
	customer.CreatedAt = now
	customer.UpdatedAt = now

	err := r.db.QueryRowContext(ctx, query, customer.Username, customer.Job, customer.Income, customer.Age, customer.Score, customer.Status, customer.Describe, customer.CreatedAt, customer.UpdatedAt).Scan(&customer.ID)
	return err
}

func (r *CustomerRepository) GetAllCustomers(ctx context.Context) ([]entity.Customer, error) {
	query := `SELECT id, name, job, income, age, score, status, describe, created_at, updated_at FROM customers`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := []entity.Customer{}
	for rows.Next() {
		customer := entity.Customer{}
		err := rows.Scan(&customer.ID, &customer.Username, &customer.Job, &customer.Income, &customer.Age, &customer.Score, &customer.Status, &customer.Describe, &customer.CreatedAt, &customer.UpdatedAt)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}

	return customers, nil
}

func (r *CustomerRepository) DeleteCustomer(ctx context.Context, id uint) error {
	query := "DELETE FROM customers WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *CustomerRepository) GetCustomerById(ctx context.Context, id uint) (*entity.Customer, error) {
	query := `SELECT id, name, job, income, age, score, status,	describe, created_at, updated_at FROM customers WHERE id = $1`

	customer := &entity.Customer{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&customer.ID, &customer.Username, &customer.Job, &customer.Income, &customer.Age, &customer.Score, &customer.Status, &customer.Describe, &customer.CreatedAt, &customer.UpdatedAt)
	return customer, err
}
