package statementparse

import "time"

type Statement struct {
	Type         string        `json:"type"`
	Date         time.Time     `json:"date,format:date"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	PostDate        time.Time `json:"postDate,format:date"`
	TransactionDate time.Time `json:"transactionDate,format:date"`
	Description     string    `json:"description"`
	Location        string    `json:"location"`
	Currency        string    `json:"currency"`
	LocalAmount     float32   `json:"localAmount"`
	Amount          float32   `json:"amount"`
}
