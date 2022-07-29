package transaction

import (
	"time"
)

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

type UserTransactionsFormatter struct {
	ID        int               `json:"id"`
	Amount    int               `json:"amount"`
	Status    string            `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	Campaign  CampaignFormatter `json:"campaign"`
}

type CampaignFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func FormatUserTransaction(transaction Transaction) UserTransactionsFormatter {
	formatted := UserTransactionsFormatter{}
	formatted.ID = transaction.ID
	formatted.Amount = transaction.Amount
	formatted.Status = transaction.Status
	formatted.CreatedAt = transaction.CreatedAt

	formatted.Campaign.Name = transaction.Campaign.Name
	formatted.Campaign.ImageURL = ""
	if len(transaction.Campaign.CampaignImages) > 0 {
		formatted.Campaign.ImageURL = transaction.Campaign.CampaignImages[0].FileName
	}

	return formatted
}

func FormatUserTransactions(transactions []Transaction) []UserTransactionsFormatter {
	formatted := []UserTransactionsFormatter{}

	if len(transactions) == 0 {
		return formatted
	}

	for _, trx := range transactions {
		formatted = append(formatted, FormatUserTransaction(trx))
	}

	return formatted
}
