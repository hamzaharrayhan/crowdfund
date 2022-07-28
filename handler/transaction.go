package handler

import (
	"crowdfund/helper"
	"crowdfund/transaction"
	"crowdfund/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) GetCampaignTransactions(c *gin.Context) {
	var input transaction.GetCampaignTransactionsInput
	err := c.ShouldBindUri(&input)

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	if err != nil {
		response := helper.JSONResponse("Failed to retrieve transactions history", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	transactions, err := h.service.GetByCampaignID(input)
	if err != nil {
		response := helper.JSONResponse("Failed to retrieve transactions history", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.JSONResponse("Campaign transactions history", http.StatusOK, "success", transaction.FormatCampaignTransactions(transactions))
	c.JSON(http.StatusOK, response)
}