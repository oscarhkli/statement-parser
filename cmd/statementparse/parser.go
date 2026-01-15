package statementparse

import (
	"log/slog"
	"strconv"
	"strings"
	"time"
)

func Parse(text string) Statement {
	return Statement{}
}

func findPhraseEndIndex(text string, start int) int {
	n := len(text)
	end := start
	for end < n {
		for end < n && text[end] != ' ' {
			end++
		}
		if end == n {
			return n - 1
		}
		if text[end+1] == ' ' {
			return end - 1
		}
		end++
	}
	return n - 1
}

func parseDate(dateStr string) (time.Time, error) {
	// Capitalize month to uppercase first letter, lowercase rest
	normalized := strings.ToUpper(dateStr[:3]) + strings.ToLower(dateStr[3:])
	return time.Parse("02Jan2006", normalized)
}

func parseAmount(amountStr string) (float32, error) {
	amountStr = strings.ReplaceAll(amountStr, ",", "")
	amount, err := strconv.ParseFloat(strings.TrimSpace(amountStr), 32)
	return float32(amount), err
}

func ParseTransactions(text string, year int) ([]*Transaction, error) {
	var transactions []*Transaction
	if len(text) == 0 {
		return transactions, nil
	}

	lines := strings.Split(text, "\n")
	slog.Info("", "Total lines", len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		var phrases []string
		for len(line) > 0 {
			endIdx := findPhraseEndIndex(line, 0)
			phrases = append(phrases, line[:endIdx+1])
			line = strings.TrimSpace(line[endIdx+1:])
		}

		if len(phrases) == 1 {
			transactions[len(transactions)-1].Description += "; " + phrases[0]
			continue
		}

		// skip credit transactions
		if strings.HasSuffix(phrases[len(phrases)-1], "CR") {
			continue
		}

		t := &Transaction{}
		transactions = append(transactions, t)

		// 1st and 2nd phrases must be postDate and transactionDate
		// 3rd must be part of description
		// Last phrase must be amount
		yearStr := strconv.Itoa(year)
		postDate, err := parseDate(phrases[0] + yearStr)
		if err != nil {
			return nil, err
		}
		t.PostDate = postDate
		transactionDate, err := parseDate(phrases[1] + yearStr)
		if err != nil {
			return nil, err
		}
		t.TransactionDate = transactionDate
		t.Description = phrases[2]
		amount, err := parseAmount(phrases[len(phrases)-1])
		if err != nil {
			return nil, err
		}
		t.Amount = amount

		phrases = phrases[3 : len(phrases)-1]

		if len(phrases) == 0 {
			continue
		}

		localAmount, err := parseAmount(phrases[len(phrases)-1])
		if err == nil {
			t.LocalAmount = localAmount
			t.Currency = phrases[len(phrases)-2]
			phrases = phrases[:len(phrases)-2]
		}

		if len(phrases) == 0 {
			continue
		}

		t.Location = strings.Join(phrases, ", ")
	}

	slog.Info("", "Total transactions parsed", len(transactions))
	return transactions, nil
}

// TODO: postprocessing the year Dec-Jan boundary cases
// TODO: postprocessing defaurlt currency