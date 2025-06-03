package database

import (
	"database/sql"
	"fmt"

	"github.com/fahrulrzi/score-match-backend/configs"
)

func NewPostgresConnection(config *configs.DatabaseConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

const CreateTablesSQL = `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
    	password VARCHAR(255) NOT NULL,
		username VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		job VARCHAR(255) NOT NULL,
		income INT NOT NULL,
		age INT NOT NULL,
		score INT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);
`

func InitDatabase(db *sql.DB) error {
	_, err := db.Exec(CreateTablesSQL)
	return err
}
