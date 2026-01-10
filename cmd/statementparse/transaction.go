package statementparse

type Transaction struct {
	PostDate        string  `json:"postDate"`
	TransactionDate string  `json:"transactionDate"`
	Description     string  `json:"description"`
	Currency        string  `json:"currency"`
	LocalAmount     float32 `json:"localAmount"`
	Amount          float32 `json:"amount"`
}
