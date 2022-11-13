package handler

import (
	"crowdfunding-TA/transaction"
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

func (h *transactionHandler) Index(c *gin.Context) {
	transactions, err := h.transactionService.FindAll()
	if err != nil {
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}

	c.HTML(http.StatusOK, "transaction_index.html", gin.H{
		"transactions": transactions,
	})
}

func (h *transactionHandler) CollectList(c *gin.Context) {
	transactions, err := h.transactionService.FindAllCollectData()
	if err != nil {
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}
	c.HTML(http.StatusOK, "transaction_collect.html", gin.H{
		"transactions": transactions,
	})
}

func (h *transactionHandler) Collect(c *gin.Context) {
	idParams := c.Param("id")
	id, _ := strconv.Atoi(idParams)

	transactions, err := h.transactionService.FindCollectDataByCID(id)
	if err != nil {
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}
	var Status bool
	if transactions.Status == "Sukses" {
		Status = true
	} else {
		Status = false

	}
	c.HTML(http.StatusOK, "transaction_collect_detil.html", gin.H{
		"transactions": transactions,
		"status":       Status,
	})

}

func (s *transactionHandler) ChangeCollectStatus(c *gin.Context) {
	idParams := c.Param("id")
	id, _ := strconv.Atoi(idParams)

	status := c.Param("status")

	_, err := s.transactionService.ChangeCollectStatus(status, id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err,
		})
		return
	}
	path := fmt.Sprintf("/collect/%d", id)
	c.Redirect(http.StatusFound, path)
}

func (s *transactionHandler) DownloadReport(c *gin.Context) {
	transactions, err := s.transactionService.FindAllCollectData()

	_, err = s.transactionService.GPdfPendingCollectData()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error":       err,
			"transaction": transactions,
		})
		return
	}

	c.Redirect(http.StatusFound, "/report/files/init.pdf")

}
