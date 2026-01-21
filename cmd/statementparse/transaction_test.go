package statementparse

import (
	"testing"
	"time"
)

func TestTransactnion_PostProcess(t *testing.T) {
	tests := []struct {
		name     string
		input    *Transaction
		expected *Transaction
	}{
		{
			name: "Currency set",
			input: &Transaction{
				Currency: "USD",
				LocalAmount: 123.456,
				Amount: 34567.565,
			},
			expected: &Transaction{
				Currency: "USD",
				LocalAmount: 123.456,
				Amount: 34567.565,
			},
		},
		{
			name: "Currency empty",
			input: &Transaction{
				Currency: "",
				LocalAmount: -123,
				Amount: 567.89,
			},
			expected: &Transaction{
				Currency: "HKD",
				LocalAmount: 567.89,
				Amount: 567.89,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.PostProcess()
			if tt.input.Currency != tt.expected.Currency {
				t.Errorf("PostProcess() Currency = %v; want %v", tt.input.Currency, tt.expected.Currency)
			}
		})
	}
}

func TestStatement_PostProcess(t *testing.T) {
	tests := []struct {
		name     string
		input    *Statement
		expected *Statement
	}{
		{
			name: "Date in December",
			input: &Statement{
				Date: time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC),
				Transactions: []*Transaction{
					{
						PostDate:        time.Date(2025, 10, 30, 0, 0, 0, 0, time.UTC),
						TransactionDate: time.Date(2025, 10, 29, 0, 0, 0, 0, time.UTC),
					}, {
						PostDate:        time.Date(2025, 12, 5, 0, 0, 0, 0, time.UTC),
						TransactionDate: time.Date(2025, 11, 4, 0, 0, 0, 0, time.UTC),
					}, {
						PostDate:        time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
						TransactionDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			expected: &Statement{
				Date: time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC),
				Transactions: []*Transaction{
					{
						PostDate:        time.Date(2025, 10, 30, 0, 0, 0, 0, time.UTC),
						TransactionDate: time.Date(2025, 10, 29, 0, 0, 0, 0, time.UTC),
					}, {
						PostDate:        time.Date(2025, 12, 5, 0, 0, 0, 0, time.UTC),
						TransactionDate: time.Date(2025, 11, 4, 0, 0, 0, 0, time.UTC),
					}, {
						PostDate:        time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
						TransactionDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
		},
		{
			name: "Date in January",
			input: &Statement{
				Date: time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
				Transactions: []*Transaction{
					{
						PostDate:        time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
						TransactionDate: time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC),
					}, {
						PostDate:        time.Date(2025, 12, 5, 0, 0, 0, 0, time.UTC),
						TransactionDate: time.Date(2025, 11, 4, 0, 0, 0, 0, time.UTC),
					}, {
						PostDate:        time.Date(2025, 12, 2, 0, 0, 0, 0, time.UTC),
						TransactionDate: time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			expected: &Statement{
				Date: time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
				Transactions: []*Transaction{
					{
						PostDate:        time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
						TransactionDate: time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC),
					}, {
						PostDate:        time.Date(2024, 12, 5, 0, 0, 0, 0, time.UTC),
						TransactionDate: time.Date(2024, 11, 4, 0, 0, 0, 0, time.UTC),
					}, {
						PostDate:        time.Date(2024, 12, 2, 0, 0, 0, 0, time.UTC),
						TransactionDate: time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.PostProcess()
			for i, tr := range tt.input.Transactions {
				expTr := tt.expected.Transactions[i]
				if !tr.PostDate.Equal(expTr.PostDate) {
					t.Errorf("Transaction %d PostDate = %v; want %v", i, tr.PostDate, expTr.PostDate)
				}
				if !tr.TransactionDate.Equal(expTr.TransactionDate) {
					t.Errorf("Transaction %d TransactionDate = %v; want %v", i, tr.TransactionDate, expTr.TransactionDate)
				}
			}
		})
	}
}
