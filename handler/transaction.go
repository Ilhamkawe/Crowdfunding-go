package handler

import (
	"crowdfunding-TA/helper"
	"crowdfunding-TA/transaction"
	"crowdfunding-TA/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	transactionService transaction.Service
}

func NewTransactionHandler(transactionService transaction.Service) *transactionHandler {
	return &transactionHandler{transactionService}
}

func (h *transactionHandler) GetCampaignTransaction(c *gin.Context) {
	var input transaction.GetCampaignTransactionInput

	err := c.ShouldBindUri(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Data Tidak Valid", http.StatusUnprocessableEntity, "gagal", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	transactions, err := h.transactionService.GetTransactionsByCampaignID(input)

	if err != nil {
		response := helper.APIResponse("Terjadi Kesalahan Saat Mengambil Data Transaksi", http.StatusBadRequest, "gagal", err.Error())
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.APIResponse("Data Transaksi", http.StatusOK, "Berhasil", transaction.FormatTransactions(transactions))
	c.JSON(http.StatusOK, response)
}
