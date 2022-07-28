package transaction

import "crowdfund/user"

type GetCampaignTransactionsInput struct {
	ID   int       `uri:"id" binding:"required"`
	User user.User `json:"user" binding:"required"`
}
