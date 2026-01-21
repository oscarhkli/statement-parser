package statementparse

import (
	"encoding/csv"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type Transaction struct {
	PostDate        time.Time `json:"postDate,format:date"`
	TransactionDate time.Time `json:"transactionDate,format:date"`
	Description     string    `json:"description"`
	Location        string    `json:"location"`
	Currency        string    `json:"currency"`
	LocalAmount     float64   `json:"localAmount"`
	Amount          float64   `json:"amount"`
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
// Also, it sets the currency to "HKD" and local amount to amount if currency is empty.
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
		if t.Currency == "" {
			t.Currency = "HKD"
			t.LocalAmount = t.Amount
		}
	}
}

func (s Statement) ToJSON() (string, error) {
	jsonData, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (s Statement) ToCSV() (string, error) {
	var sb strings.Builder

	cw := csv.NewWriter(&sb)

	if err := cw.Write([]string{
		"post_date",
		"transaction_date",
		"description",
		"location",
		"currency",
		"local_amount",
		"amount",
	}); err != nil {
		return "", err
	}

	for _, t := range s.Transactions {
		record := []string{
			formatDate(t.PostDate),
			formatDate(t.TransactionDate),
			t.Description,
			t.Location,
			t.Currency,
			formatFloat(t.LocalAmount),
			formatFloat(t.Amount),
		}

		if err := cw.Write(record); err != nil {
			return "", err
		}
	}

	cw.Flush()
	if err := cw.Error(); err != nil {
		return "", err
	}

	return sb.String(), nil
}

func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}
