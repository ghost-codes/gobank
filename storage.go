package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=password sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{db}, nil
}

func (s *PostgresStore) init() error {
	return s.createAccountTable()
}
func (s *PostgresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS Account(
 			id serial primary key,
			first_name varchar(50),
			last_name varchar(50),
			number serial,
			balance serial,
			created_at timestamp
	) `

	_, err := s.db.Exec(query)

	return err

}

func (s *PostgresStore) CreateAccount(account *Account) error {
	query := `INSERT INTO Account (first_name,last_name,number,balance,created_at) VALUES ($1,$2,$3,$4,$5)`

	_, err := s.db.Exec(query, account.FirstName, account.LastName, account.Number, account.Balance, account.CreatedAt)

	return err
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query("DELETE from Account where id=$1", id)
	return err
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	query := `SELECT * FROM Account WHERE id=$1`

	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account with %d not found", id)
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	accounts := []*Account{}
	query := `SELECT * FROM Account`

	results, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	for results.Next() {
		account, err := scanIntoAccount(results)

		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)

	}

	return accounts, nil

}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt)

	if err != nil {
		return nil, err
	}
	return account, nil
}
