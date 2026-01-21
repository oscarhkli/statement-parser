package statementparse

import (
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Parse(text string) Statement {
	lines := strings.Split(text, "\n")

	statementDate, err := extractStatementDate(lines)
	if err != nil {
		slog.Error("Failed to parse statement date", "error", err)
	}
	transactionLines := preprocessTransactionText(lines)
	transactionsText := strings.Join(transactionLines, "\n")
	transactions, err := parseTransactions(transactionsText, statementDate.Year())
	if err != nil {
		slog.Error("Failed to parse transactions", "error", err)
	}

	statement := NewStatement(
		extractStatementType(lines),
		statementDate,
		transactions,
	)
	statement.PostProcess()
	return *statement
}

func extractStatementType(lines []string) string {
	for _, line := range lines {
		lineUpper := strings.ToUpper(line)
		if !strings.Contains(lineUpper, "STATEMENT") {
			continue
		}

		if strings.Contains(lineUpper, "VISA SIGNATURE") {
			return "HSBC Visa Signature"
		}
		if strings.Contains(lineUpper, "HSBC RED") {
			return "HSBC Red"
		}
	}
	return ""
}

// extractStatementDate parses the statement date.
// It tries to find a line containing "Statement Date:" and extract the date following it in the next line.
// Returns zero time if not found or parsing fails.
func extractStatementDate(lines []string) (time.Time, error) {
	for i, line := range lines {
		if !strings.Contains(strings.ToUpper(line), "STATEMENT DATE") {
			continue
		}
		if i == len(lines)-1 {
			break
		}

		dateLine := strings.TrimSpace(lines[i+1])
		re := regexp.MustCompile(`\b\d{1,2}\s+[A-Z]{3}\s+\d{4}\b`)
		match := re.FindString(dateLine)
		if match == "" {
			slog.Warn("Statement date pattern not found in line", "line", dateLine)
			return time.Time{}, nil
		}

		return time.Parse("02 Jan 2006", match)
	}

	slog.Warn("Statement date not found in text")
	return time.Time{}, nil
}

// TODO: Preprocess text to extract transaction section
func preprocessTransactionText(lines []string) []string {
	var results []string

	inSection := false
	inTransaction := false
	re := regexp.MustCompile(`^\s*\d{2}[A-Z]{3}\s+\d{2}[A-Z]{3}`)

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		trimmedLineUpper := strings.ToUpper(trimmedLine)

		if strings.Contains(trimmedLineUpper, "POST DATE") && strings.Contains(trimmedLineUpper, "TRANS DATE") {
			inSection = true
			inTransaction = false
			continue
		}

		if !inSection {
			continue
		}

		if trimmedLine == "" {
			if inTransaction {
				inTransaction = false
			}
			continue
		}

		if re.MatchString(trimmedLine) {
			inTransaction = true
			results = append(results, trimmedLine)
			continue
		}

		if inTransaction {
			results = append(results, trimmedLine)
		}
	}

	return results
}

// findPhraseEndIndex finds the end index of a phrase starting from 'start' index.
// A phrase is defined as a sequence of non-space characters possibly separated by single spaces.
// The phrase ends when two consecutive spaces are found or end of string is reached.
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

// parseDate parses date by capitalizing month to uppercase first letter, lowercase rest
func parseDate(dateStr string) (time.Time, error) {
	normalized := strings.ToUpper(dateStr[:3]) + strings.ToLower(dateStr[3:])
	return time.Parse("02Jan2006", normalized)
}

func parseAmount(amountStr string) (float64, error) {
	amountStr = strings.ReplaceAll(amountStr, ",", "")
	amount, err := strconv.ParseFloat(strings.TrimSpace(amountStr), 64)
	return amount, err
}

// parseTransactions parses transaction lines from the given text for the specified year.
func parseTransactions(text string, year int) ([]*Transaction, error) {
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

		t := NewTransaction()
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
