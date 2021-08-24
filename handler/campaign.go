package handler

import (
	"crowdfunding-TA/campaign"
	"crowdfunding-TA/helper"
	"crowdfunding-TA/user"
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

	campaigns, err := h.campaignService.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Terjadi Kesalahan", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Berhasil Mengambil Data", http.StatusOK, "Berhasil", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) GetCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Terjadi Kesalahan", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := h.campaignService.GetCampaignByID(input)

	if err != nil {
		response := helper.APIResponse("Terjadi Kesalahan", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil Mengambil Data", http.StatusOK, "Berhasil", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Gagal Buat Campaign", http.StatusUnprocessableEntity, "gagal", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser

	newCampaign, err := h.campaignService.CreateCampaign(input)

	if err != nil {
		response := helper.APIResponse("Gagal buat Campaign", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil Buat Campaign", http.StatusOK, "sukses", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)
}

// tangkap input kedalam struct
// ambil current user dengan jwt
// panggil service , parameternya input struct
// panggil repo untuk simpan data
