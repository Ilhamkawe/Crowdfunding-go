package handler

import (
	"crowdfunding-TA/campaign"
	"crowdfunding-TA/helper"
	"crowdfunding-TA/transaction"
	"crowdfunding-TA/user"
	"fmt"
	"net/http"
	"strconv"

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

func (h *transactionHandler) GetAllByReward(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)
	fmt.Println("hay1")
	var input transaction.GetCampaignTransactionInput

	err := c.ShouldBindUri(&input)
	fmt.Println("hay2")

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Data Tidak Valid", http.StatusUnprocessableEntity, "gagal", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	var inputPaginate campaign.PaginateCampaignInput

	err = c.ShouldBindJSON(&inputPaginate)
	fmt.Println(input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Terjadi Kesalahan", http.StatusUnprocessableEntity, "gagal", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	Cid, _ := strconv.Atoi(c.Query("campaign_id"))
	page, _ := strconv.Atoi(c.Query("page"))
	inputPaginate.ActivePage = page

	input.User = currentUser

	transactions, err := h.transactionService.FindAllByReward(input.ID, Cid, currentUser.ID, inputPaginate)

	if err != nil {
		response := helper.APIResponse("Terjadi Kesalahan Saat Mengambil Data Transaksi", http.StatusBadRequest, "Gagal", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Data Transaksi", http.StatusOK, "Berhasil", transactions)
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

func (h *transactionHandler) GetNotification(c *gin.Context) {
	var input transaction.TransactionNotificationInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		response := helper.APIResponse("Gagal Mengirim Notification", http.StatusBadRequest, "Gagal", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = h.transactionService.ProcessPayment(input)

	if err != nil {
		response := helper.APIResponse("Gagal Mengirim Notification", http.StatusBadRequest, "Gagal", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, input)

}

func (h *transactionHandler) CollectAmount(c *gin.Context) {
	var input transaction.CollectInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Gagal Mencairkan Dana", http.StatusUnprocessableEntity, "Gagal", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	collect, err := h.transactionService.CollectAmount(input)

	if err != nil {
		response := helper.APIResponse("Gagal Mencairkan Dana", http.StatusBadRequest, "Gagal", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil Buat Transaksi", http.StatusOK, "Sukses", collect)
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) FindCollectData(c *gin.Context) {
	var input transaction.GetCampaignTransactionInput
	err := c.ShouldBindUri(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Gagal mendapatkan data pencairan dana", http.StatusUnprocessableEntity, "Gagal", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	collect, err := h.transactionService.FindCollectData(input.ID)

	if err != nil {
		response := helper.APIResponse("Gagal Menambil Data", http.StatusBadRequest, "Gagal", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil Ambil Data", http.StatusOK, "Sukses", collect)
	c.JSON(http.StatusOK, response)

}
