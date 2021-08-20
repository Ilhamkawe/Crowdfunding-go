package handler

import (
	"crowdfunding-TA/campaign"
	"crowdfunding-TA/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	campaignService campaign.Service
}

func NewCampaignHandler(campaignService campaign.Service) *campaignHandler {
	return &campaignHandler{campaignService}
}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.campaignService.FindCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Terjadi Kesalahan", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Berhasil Mengambil Data", http.StatusOK, "Berhasil", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}

// tagkap parameter
// handler ke service
//  service yang menentukan repository mana yg di panggil
// repository : FindAll, FindByUserID
// db
