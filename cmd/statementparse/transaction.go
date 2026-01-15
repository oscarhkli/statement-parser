package statementparse

import "time"

type Transaction struct {
	PostDate        time.Time `json:"postDate,format:date"`
	TransactionDate time.Time `json:"transactionDate,format:date"`
	Description     string    `json:"description"`
	Location        string    `json:"location"`
	Currency        string    `json:"currency"`
	LocalAmount     float32   `json:"localAmount"`
	Amount          float32   `json:"amount"`
}

func NewTransaction() *Transaction {
	return &Transaction{}
}

func (t *Transaction) PostProcess() {
	if t.Currency == "" {
		t.Currency = "HKD"
	}
}

type Statement struct {
	Type         string        `json:"type"`
	Date         time.Time     `json:"date,format:date"`
	Transactions []Transaction `json:"transactions"`
}

func NewStatement() *Statement {
	return &Statement{}
}

func (s *Statement) PostProcess() {

}

// TODO: postprocessing the year Dec-Jan boundary cases
