package statementparse

import "testing"

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
			},
			expected: &Transaction{
				Currency: "USD",
			},
		},
		{
			name: "Currency empty",
			input: &Transaction{
				Currency: "",
			},
			expected: &Transaction{
				Currency: "HKD",
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
