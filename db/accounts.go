package db

import (
	"database/sql"
	"fmt"

	accounts "kangym/types"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*accounts.CreateAccountRequest) (*accounts.Account, error)
	DeleteAccountById(int) error
	RestoreAccountById(int) error
	GetAccountById(int) (*accounts.Account, error)
	GetAccounts() ([]*accounts.Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=prueba123 sslmode=disable"
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
	return s.createTables()
}

func (s *PostgresStore) createTables() error {
	query := `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

		CREATE TABLE IF NOT EXISTS account (
			id uuid DEFAULT uuid_generate_v4() primary key,
			email varchar(150),
			password varchar(50),
			first_name varchar(50),
			last_name varchar(50),
			created_at timestamp default CURRENT_TIMESTAMP,
			deleted boolean DEFAULT false
		);

		CREATE TABLE IF NOT EXISTS kanji (
			id uuid DEFAULT uuid_generate_v4() primary key,
			keyword varchar(50),
			character varchar(5),
			story varchar(250)
		);

		CREATE TABLE IF NOT EXISTS userkanji (
			id uuid DEFAULT uuid_generate_v4() primary key,
			account_id uuid REFERENCES account (id) ON DELETE CASCADE,
			kanji_id uuid REFERENCES kanji (id) ON DELETE CASCADE,
			userstory varchar(250)
		);

		`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(account *accounts.CreateAccountRequest) (*accounts.Account, error) {
	query := `
		INSERT INTO account(email, first_name, last_name, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id, first_name, last_name, created_at, deleted
		
		`
	rows, err := s.db.Query(query,
		account.Email,
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

func (s *PostgresStore) GetAccountById(id int) (*accounts.Account, error) {
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

func (s *PostgresStore) GetAccounts() ([]*accounts.Account, error) {
	rows, err := s.db.Query(`
		SELECT 
			id, 
			email,
			first_name,
			last_name,
			created_at,
			deleted 
		FROM account 
		WHERE deleted = false`)
	if err != nil {
		return nil, err
	}

	accounts := []*accounts.Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*accounts.Account, error) {
	account := new(accounts.Account)
	err := rows.Scan(
		&account.ID,
		&account.Email,
		&account.FirstName,
		&account.LastName,
		&account.CreatedAt,
		&account.Deleted)
	return account, err
}
