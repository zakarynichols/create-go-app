package postgres

import (
	"database/sql"
	"fmt"
	"time"
)

type Config struct {
	Password string
	User     string
	Name     string
	Host     string
	SSLMode  string
}

type psqlService struct {
	db *sql.DB
}

// Open opens a postgres database and returns a service to interact with it.
func Open(config Config) (*psqlService, error) {
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s sslmode=%s",
		config.User,
		config.Password,
		config.Name,
		config.Host,
		config.SSLMode,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &psqlService{db}, nil
}

func (psql *psqlService) Close() error {
	return psql.db.Close()
}

func (psql *psqlService) QueryNow() (time.Time, error) {
	var currentTime time.Time

	err := psql.db.QueryRow("SELECT NOW()").Scan(&currentTime)
	if err != nil {
		return time.Time{}, err
	}

	return currentTime, err
}
