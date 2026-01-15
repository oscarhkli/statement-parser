package statementparse

import (
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
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		transactions = append(transactions, &Transaction{})
		last := transactions[len(transactions)-1]

		var phrases []string
		for len(line) > 0 {
			endIdx := findPhraseEndIndex(line, 0)
			phrases = append(phrases, line[:endIdx+1])
			line = strings.TrimSpace(line[endIdx+1:])
		}

		if len(phrases) == 1 {
			last.Description += "; " + phrases[0]
			continue
		}

		// skip credit transactions
		if strings.HasSuffix(phrases[len(phrases)-1], "CR") {
			continue
		}

		// 1st and 2nd phrases must be postDate and transactionDate
		// 3rd must be part of description
		// Last phrase must be amount
		yearStr := strconv.Itoa(year)
		postDate, err := parseDate(phrases[0] + yearStr)
		if err != nil {
			return nil, err
		}
		last.PostDate = postDate
		transactionDate, err := parseDate(phrases[1] + yearStr)
		if err != nil {
			return nil, err
		}
		last.TransactionDate = transactionDate
		last.Description = phrases[2]
		amount, err := parseAmount(phrases[len(phrases)-1])
		if err != nil {
			return nil, err
		}
		last.Amount = amount

		phrases = phrases[3 : len(phrases)-1]

		localAmount, err := parseAmount(phrases[len(phrases)-1])
		if err == nil {
			last.LocalAmount = localAmount
			last.Currency = phrases[len(phrases)-2]
			phrases = phrases[:len(phrases)-2]
		}

		if len(phrases) > 0 {
			last.Location = strings.Join(phrases, ", ")
		}
	}
	return transactions, nil
}

// TODO: postprocessing the year Dec-Jan boundary cases
