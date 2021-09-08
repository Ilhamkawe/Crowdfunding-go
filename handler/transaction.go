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

func (h *transactionHandler) GetUserTransaction(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	transactions, err := h.transactionService.GetTransactionsByUserID(userID)

	if err != nil {
		response := helper.APIResponse("Terjadi Kesalahan Saat Mengambil Data Transaksi", http.StatusBadRequest, "Gagal", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Data Transaksi", http.StatusOK, "Berhasil", transaction.FormatUserTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	var input transaction.CreateTransactionInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Gagal Buat Transaksi", http.StatusUnprocessableEntity, "Gagal", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	newTransaction, err := h.transactionService.CreateTransaction(input)

	if err != nil {
		response := helper.APIResponse("Gagal Buat Transaksi", http.StatusBadRequest, "Gagal", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil Buat Transaksi", http.StatusOK, "Sukses", transaction.FormatPaymentTransaction(newTransaction))
	c.JSON(http.StatusOK, response)

}
