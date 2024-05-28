package main

import (
	"time"
)

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

func NewAccount(firstName, lastName, password string) *Account {
	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		Password:  password,
		CreatedAt: time.Now().UTC(),
	}
}
