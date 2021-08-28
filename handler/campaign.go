package handler

import (
	"crowdfunding-TA/campaign"
	"crowdfunding-TA/helper"
	"crowdfunding-TA/user"
	"fmt"
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

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	var inputID campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&inputID)

	if err != nil {
		response := helper.APIResponse("Terjadi Kesalahan Saat Mengupdate", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaign.CreateCampaignInput

	currentUser := c.MustGet("currentUser").(user.User)
	inputData.User = currentUser

	err = c.ShouldBindJSON(&inputData)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Gagal Update Campaign", http.StatusUnprocessableEntity, "gagal", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	updatedCampaign, err := h.campaignService.UpdateCampaign(inputID, inputData)

	// fmt.Println("errornya :", err)

	if err != nil {

		response := helper.APIResponse("Terjadi Kesalahan Saat Mengupdate", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil Buat Campaign", http.StatusOK, "sukses", campaign.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaignImage(c *gin.Context) {
	var input campaign.CreateCampaignImageInput

	err := c.ShouldBind(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Terjadi Kesalahan Saat Update Campaign", http.StatusUnprocessableEntity, "Gagal", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	file, err := c.FormFile("file")
	if err != nil {
		errorMessage := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Terjadi Kesalahan Saat Mengunggah Gambar", http.StatusBadRequest, "Gagal", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	path := fmt.Sprintf("images/%d-%s", currentUser.ID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		errorMessage := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Terjadi Kesalahan Saat Mengunggah Gambar", http.StatusBadRequest, "Gagal", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.campaignService.SaveCampaignImage(input, path)

	if err != nil {
		// errorMessage := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Terjadi Kesalahan Saat Mengunggah Gambar", http.StatusBadRequest, "Gagal", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Berhasil Mengunggah Gambar", http.StatusBadRequest, "sukses", data)
	c.JSON(http.StatusBadRequest, response)
}

// tangkap input kedalam struct
// ambil current user dengan jwt
// panggil service , parameternya input struct
// panggil repo untuk simpan data
