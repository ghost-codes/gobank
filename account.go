package main

import (
	"math/rand"
	"time"
)

type TransferRequest struct {
	ToAccount   int `json:"toAccount"`
	FromAccount int `json:"fr1omAccount"`
	Amount      int `json:"amount"`
}

type CreateAccount struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		Number:    rand.Int63n(1000000000),
		CreatedAt: time.Now().UTC(),
	}
}
