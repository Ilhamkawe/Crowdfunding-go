package handler

import (
	"crowdfunding-TA/campaign"
	"crowdfunding-TA/transaction"
	"crowdfunding-TA/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type dashboardHandler struct {
	campaignService    campaign.Service
	userService        user.Service
	transactionService transaction.Service
}

func NewDashboardHandler(campaignService campaign.Service, userService user.Service, transactionService transaction.Service) *dashboardHandler {
	return &dashboardHandler{campaignService, userService, transactionService}
}

func (h *dashboardHandler) Index(c *gin.Context) {
	campaigns, err := h.campaignService.GetAllCampaign()
	if err != nil {
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}
	transactions, err := h.transactionService.FindAll()
	if err != nil {
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}
	cattegory, err := h.campaignService.FindAllCattegory()
	if err != nil {
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}
	pending, err := h.campaignService.GetCampaignByStatus("Pending")
	if err != nil {
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}
	Sukses, err := h.campaignService.GetCampaignByStatus("Berjalan")
	if err != nil {
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}
	cair, err := h.campaignService.GetCampaignByStatus("Dicairkan")
	if err != nil {
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}
	c.HTML(http.StatusOK, "dashboard_index.html", gin.H{
		"campaigns":    campaigns,
		"transactions": transactions,
		"users":        users,
		"cattegory":    cattegory,
		"countpending": len(pending),
		"countsukses":  len(Sukses),
		"countcair":    len(cair),
		"countuser":    len(users),
	})
}
