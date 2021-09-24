package handler

import (
	"crowdfunding-TA/campaign"
	"crowdfunding-TA/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	campaignService campaign.Service
	userService     user.Service
}

func NewCampaignHandler(campaignService campaign.Service, userService user.Service) *campaignHandler {
	return &campaignHandler{campaignService, userService}
}

func (h *campaignHandler) Index(c *gin.Context) {
	campaigns, err := h.campaignService.GetCampaigns(0)
	if err != nil {
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}

	c.HTML(http.StatusOK, "campaign_index.html", gin.H{
		"campaigns": campaigns,
	})
}

func (h *campaignHandler) New(c *gin.Context) {

	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}

	input := campaign.FormCampaignInput{}
	input.Users = users

	c.HTML(http.StatusOK, "campaign_new.html", input)

}

func (h *campaignHandler) Create(c *gin.Context) {
	var input campaign.FormCampaignInput
	err := c.ShouldBind(&input)

	if err != nil {
		// nanti redirect dan kirim pesan
		fmt.Println(err)
		return
	}

	user, err := h.userService.GetUserByID(input.UserID)
	if err != nil {
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}

	createCampaignInput := campaign.CreateCampaignInput{}
	createCampaignInput.Name = input.Name
	createCampaignInput.ShortDescription = input.ShortDescription
	createCampaignInput.Description = input.Description
	createCampaignInput.GoalAmount = input.GoalAmount
	createCampaignInput.Perks = input.Perks
	createCampaignInput.User = user

	_, err = h.campaignService.CreateCampaign(createCampaignInput)
	if err != nil {
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/campaigns")

}

func (h *campaignHandler) Image(c *gin.Context) {
	idParams := c.Param("id")
	id, _ := strconv.Atoi(idParams)

	c.HTML(http.StatusOK, "campaign_image.html", gin.H{
		"ID": id,
	})
}

func (h *campaignHandler) CreateImage(c *gin.Context) {

	file, err := c.FormFile("images")

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err,
		})
		return
	}

	idParams := c.Param("id")
	id, _ := strconv.Atoi(idParams)

	existingCampaign, err := h.campaignService.GetCampaignByID(campaign.GetCampaignDetailInput{ID: id})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err,
		})
		return
	}

	userID := existingCampaign.UserID

	path := fmt.Sprintf("images/campaign/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err,
		})
		return
	}

	createCampaignImageInput := campaign.CreateCampaignImageInput{}
	createCampaignImageInput.IsPrimary = true
	createCampaignImageInput.CampaignID = id

	userCampaign, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err,
		})
		return
	}
	createCampaignImageInput.User = userCampaign

	_, err = h.campaignService.SaveCampaignImage(createCampaignImageInput, path)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err,
		})
		return
	}

	c.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) Edit(c *gin.Context) {
	idParams := c.Param("id")
	id, _ := strconv.Atoi(idParams)

	campaign, err := h.campaignService.GetCampaignByID(campaign.GetCampaignDetailInput{ID: id})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err,
		})
		return
	}

	c.HTML(http.StatusOK, "campaign_edit.html", gin.H{
		"campaign": campaign,
	})

}

func (h *campaignHandler) Update(c *gin.Context) {
	idParams := c.Param("id")
	id, _ := strconv.Atoi(idParams)

	var input campaign.FormCampaignUpdate
	err := c.ShouldBind(&input)

	if err != nil {
		input.ID = id
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err,
		})
		return
	}

	existingCampaign, err := h.campaignService.GetCampaignByID(campaign.GetCampaignDetailInput{ID: id})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err,
		})
		return
	}

	userID := existingCampaign.UserID

	userCampaign, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err,
		})
		return
	}

	updateInput := campaign.CreateCampaignInput{}
	updateInput.Name = input.Name
	updateInput.Perks = input.Perks
	updateInput.Description = input.Description
	updateInput.ShortDescription = input.ShortDescription
	updateInput.User = userCampaign
	updateInput.GoalAmount = input.GoalAmount

	_, err = h.campaignService.UpdateCampaign(campaign.GetCampaignDetailInput{ID: id}, updateInput)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err,
		})
		return
	}

	c.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) Detail(c *gin.Context) {
	idParams := c.Param("id")
	id, _ := strconv.Atoi(idParams)

	campaign, err := h.campaignService.GetCampaignByID(campaign.GetCampaignDetailInput{ID: id})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err,
		})
		return
	}

	c.HTML(http.StatusFound, "campaign_detail.html", gin.H{
		"campaign": campaign,
	})

}
