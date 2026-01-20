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
	Type         string         `json:"type"`
	Date         time.Time      `json:"date,format:date"`
	Transactions []*Transaction `json:"transactions"`
}

func NewStatement(statementType string, date time.Time, transactions []*Transaction) *Statement {
	return &Statement{
		Type:         statementType,
		Date:         date,
		Transactions: transactions,
	}
}

// PostProcess adjusts transaction dates for year boundary cases.
// It depends on the statement date. If the statement month is smaller than post/transaction month,
// the post/transaction year should be decremented by 1.
func (s *Statement) PostProcess() {
	month := s.Date.Month()
	if month == 12 {
		return
	}

	for _, t := range s.Transactions {
		if t.PostDate.Month() > month {
			t.PostDate = t.PostDate.AddDate(-1, 0, 0)
		}
		if t.TransactionDate.Month() > month {
			t.TransactionDate = t.TransactionDate.AddDate(-1, 0, 0)
		}
	}
}
