package transaction

import "time"

type CampaignTransactionsFormatter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func FormatCampaignTransaction(transaction Transaction) CampaignTransactionsFormatter {
	formatter := CampaignTransactionsFormatter{}
	formatter.ID = transaction.ID
	formatter.Name = transaction.User.Name
	formatter.Amount = transaction.Amount
	formatter.CreatedAt = transaction.CreatedAt

	return formatter
}

func FormatCampaignTransactions(transactions []Transaction) []CampaignTransactionsFormatter {
	formatted := []CampaignTransactionsFormatter{}
	if len(transactions) == 0 {
		return formatted
	}

	for _, campaign := range transactions {
		formatted = append(formatted, FormatCampaignTransaction(campaign))
	}
	return formatted
}
