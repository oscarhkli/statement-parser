package statementparse

import (
	"math"
	"testing"
	"time"
)

func TestParseEmptyString(t *testing.T) {
	res := Parse("")
	if len(res.Transactions) > 0 {
		t.Errorf("Parse(\"\") = %v; want []", res)
	}
}

// func TestParseSampleStatement(t *testing.T) {
// 	t.Skip("incomplete test - implement parsing logic")

// 	data, err := os.ReadFile("testdata/hsbc-vs-001.txt")
// 	if err != nil {
// 		t.Fatalf("Failed to read test data: %v", err)
// 	}
// 	sampleText := string(data)

// 	res := Parse(sampleText)
// 	wanted := Statement{
// 		Type: "HSBC Statement",
// 		Date: time.Date(2025, 10, 20, 0, 0, 0, 0, time.UTC),
// 		Transactions: []Transaction{
// 			{
// 				PostDate:        time.Date(2025, 9, 12, 0, 0, 0, 0, time.UTC),
// 				TransactionDate: time.Date(2025, 9, 10, 0, 0, 0, 0, time.UTC),
// 				Description:     "Momo Kingdom Ltd",
// 			},
// 		},
// 	}

// 	if res.Date != wanted.Date {
// 		t.Errorf("Parse(sampleText) = %v; want %v", res, wanted)
// 	}
// }

func compareTransactions(t *testing.T, got, want []Transaction) {
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d, want %d", len(got), len(want))
	}

	for i := range got {
		g := got[i]
		w := want[i]

		if !g.PostDate.Equal(w.PostDate) {
			t.Errorf("element %d: PostDate mismatch: got %v, want %v", i, g.PostDate, w.PostDate)
		}
		if !g.TransactionDate.Equal(w.TransactionDate) {
			t.Errorf("element %d: TransactionDate mismatch: got %v, want %v", i, g.TransactionDate, w.TransactionDate)
		}
		if g.Description != w.Description {
			t.Errorf("element %d: Description mismatch: got %q, want %q", i, g.Description, w.Description)
		}
		if g.Location != w.Location {
			t.Errorf("element %d: Location mismatch: got %q, want %q", i, g.Location, w.Location)
		}
		if g.Currency != w.Currency {
			t.Errorf("element %d: Currency mismatch: got %q, want %q", i, g.Currency, w.Currency)
		}
		const eps = 0.0001
		if math.Abs(float64(g.LocalAmount-w.LocalAmount)) > eps {
			t.Errorf("element %d: LocalAmount mismatch: got %v, want %v", i, g.LocalAmount, w.LocalAmount)
		}
		if math.Abs(float64(g.Amount-w.Amount)) > eps {
			t.Errorf("element %d: Amount mismatch: got %v, want %v", i, g.Amount, w.Amount)
		}
	}
}

func TestParseTransactionsEmptyString(t *testing.T) {
	res, err := ParseTransactions("", 2025)
	if err != nil {
		t.Errorf("ParseTransactions(\"%s\", 2025) = %v; want []", "", err)
	}
	if len(res) != 0 {
		t.Errorf("ParseTransactions(\"%s\", 2025) = %v; want []", "", res)
	}
}

func TestParseTransactionsOneLineTransaction(t *testing.T) {
	text := " 12SEP      10SEP       Momo Kingdom Ltd            Ealing                            GB      GBP                       8.99                               97.03"
	res, err := ParseTransactions(text, 2025)
	if err != nil {
		t.Errorf("ParseTransactions(\"%s\", 2025) = %v; want []", text, err)
	}

	wanted := []Transaction{
		{
			PostDate:        time.Date(2025, 9, 12, 0, 0, 0, 0, time.UTC),
			TransactionDate: time.Date(2025, 9, 10, 0, 0, 0, 0, time.UTC),
			Description:     "Momo Kingdom Ltd",
			Location:        "Ealing, GB",
			Currency:        "GBP",
			LocalAmount:     8.99,
			Amount:          97.03,
		},
	}
	compareTransactions(t, []Transaction{*res[0]}, wanted)
}

func TestParseTransactionsTwoLineTransaction2(t *testing.T) {
	text := ` 12SEP      10SEP       Momo Kingdom Ltd            Ealing                            GB      GBP                       8.99                               97.03
 20SEP     18SEP       WH Smith Ealing            Ealing                 GB     GBP              4.49                      48.85`

	res, err := ParseTransactions(text, 2025)
	if err != nil {
		t.Errorf("ParseTransactions(\"%s\", 2025) = %v; want []", text, err)
	}

	wanted := []Transaction{
		{
			PostDate:        time.Date(2025, 9, 12, 0, 0, 0, 0, time.UTC),
			TransactionDate: time.Date(2025, 9, 10, 0, 0, 0, 0, time.UTC),
			Description:     "Momo Kingdom Ltd",
			Location:        "Ealing, GB",
			Currency:        "GBP",
			LocalAmount:     8.99,
			Amount:          97.03,
		}, {
			PostDate:        time.Date(2025, 9, 20, 0, 0, 0, 0, time.UTC),
			TransactionDate: time.Date(2025, 9, 18, 0, 0, 0, 0, time.UTC),
			Description:     "WH Smith Ealing",
			Location:        "Ealing, GB",
			Currency:        "GBP",
			LocalAmount:     4.49,
			Amount:          48.85,
		},
	}
	compareTransactions(t, []Transaction{*res[0], *res[1]}, wanted)
}
