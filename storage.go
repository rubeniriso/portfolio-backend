package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*CreateAccountRequest) (*Account, error)
	DeleteAccountById(int) error
	RestoreAccountById(int) error
	UpdateAccount(*Account) error
	GetAccountById(int) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account  (
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		password varchar(50),
		created_at timestamp default CURRENT_TIMESTAMP,
		deleted boolean
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(account *CreateAccountRequest) (*Account, error) {
	query := `
		INSERT INTO account(first_name, last_name, password)
		VALUES ($1, $2, $3)
		RETURNING id, first_name, last_name, created_at, deleted
		`
	rows, err := s.db.Query(query,
		account.FirstName,
		account.LastName,
		account.Password)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("error creating account")
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccountById(id int) error {
	query := `UPDATE account
		SET deleted = true
		WHERE id = $1
	`
	_, err := s.db.Query(query, id)
	return err
}

func (s *PostgresStore) RestoreAccountById(id int) error {
	query := `UPDATE account
		SET deleted = false
		WHERE id = $1
	`
	_, err := s.db.Query(query, id)
	return err
}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	query := `
		SELECT 
			id, 
			first_name,
			last_name,
			created_at,
			deleted 
		FROM account
		WHERE id = $1
	`
	rows, err := s.db.Query(query, id)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query(`
		SELECT 
			id, 
			first_name,
			last_name,
			created_at,
			deleted 
		FROM account 
		WHERE deleted = false`)
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
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
		&account.CreatedAt,
		&account.Deleted)
	return account, err
}
