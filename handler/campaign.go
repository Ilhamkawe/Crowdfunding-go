package handler

import (
	"crowdfunding-TA/campaign"
	"crowdfunding-TA/helper"
	"crowdfunding-TA/user"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	campaignService campaign.Service
}

func NewCampaignHandler(campaignService campaign.Service) *campaignHandler {
	return &campaignHandler{campaignService}
}

func (h *campaignHandler) GetRewards(c *gin.Context) {
	var input campaign.GetCampaignDetailInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Terjadi Kesalahan", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	rewards, err := h.campaignService.GetRewards(input)
	if err != nil {
		response := helper.APIResponse("Terjadi Kesalahan", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil Mengambil Data", http.StatusOK, "Berhasil", campaign.FormatCampaignReward(rewards))
	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) SearchCampaignPaginate(c *gin.Context) {
	var input campaign.SearchCampaignPaginate
	page, _ := strconv.Atoi(c.Query("page"))
	input.ActivePage = page
	err := c.ShouldBindJSON(&input)

	fmt.Println(input)

	if err != nil {
		response := helper.APIResponse("Terjadi Kesalahan Input", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		fmt.Print(err)
		return
	}

	campaigns, err := h.campaignService.SearchCampaignPaginate(input)

	if err != nil {
		response := helper.APIResponse("Terjadi Kesalahan Saat Mencari Data", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil Mengambil Data", http.StatusOK, "Berhasil", campaigns)
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) SearchCampaign(c *gin.Context) {
	var input campaign.SearchCampaignInput
	err := c.ShouldBindJSON(&input)

	fmt.Println(input)

	if err != nil {
		response := helper.APIResponse("Terjadi Kesalahan Input", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		fmt.Print(err)
		return
	}

	campaigns, err := h.campaignService.SearchCampaign(input)

	if err != nil {
		response := helper.APIResponse("Terjadi Kesalahan Saat Mencari Data", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil Mengambil Data", http.StatusOK, "Berhasil", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) Limit(c *gin.Context) {
	var input campaign.GetLimitDataInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Terjadi Kesalahan", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	campaigns, err := h.campaignService.Limit(input.Limit)

	if err != nil {
		response := helper.APIResponse("Terjadi Kesalahan", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil Mengambil Data", http.StatusOK, "Berhasil", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
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

	err := c.ShouldBind(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Gagal Buat Campaign", http.StatusUnprocessableEntity, "gagal", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser

	//! upload attachment
	file, err := c.FormFile("attachment")
	fmt.Println(file.Filename)
	if err != nil {
		errorMessage := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Terjadi Kesalahan Saat Mengunggah Lampiran", http.StatusBadRequest, "Gagal", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	path := fmt.Sprintf("attachment/%d-%s", currentUser.ID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		errorMessage := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Terjadi Kesalahan Saat Mengunggah Lampiran", http.StatusBadRequest, "Gagal", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	input.Path = path

	// Save data
	newCampaign, err := h.campaignService.CreateCampaign(input)

	if err != nil {
		response := helper.APIResponse("Gagal buat Campaign", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil Buat Campaign", http.StatusOK, "sukses", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UpdateAttachment(c *gin.Context) {
	var input campaign.UpdateAttachmentInput

	err := c.ShouldBind(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Gagal Update Lampiran", http.StatusUnprocessableEntity, "gagal", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser
	input.Path = " "

	if strings.ToUpper(input.Action) == "UPLOAD" {
		//! upload attachment
		file, err := c.FormFile("attachment")
		fmt.Println(file.Filename)
		if err != nil {
			errorMessage := gin.H{"is_uploaded": false}

			response := helper.APIResponse("Terjadi Kesalahan Saat Mengunggah Lampiran", http.StatusBadRequest, "Gagal", errorMessage)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		path := fmt.Sprintf("attachment/%d-%s", currentUser.ID, file.Filename)

		err = c.SaveUploadedFile(file, path)
		if err != nil {
			errorMessage := gin.H{"is_uploaded": false}

			response := helper.APIResponse("Terjadi Kesalahan Saat Mengunggah Lampiran", http.StatusBadRequest, "Gagal", errorMessage)
			c.JSON(http.StatusBadRequest, response)
			return
		}
		input.Path = path
	}

	updateAttachment, err := h.campaignService.UpdateAttachment(input)
	if err != nil {
		response := helper.APIResponse("Gagal Update Lampiran", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil Ubah Lampiran", http.StatusOK, "sukses", updateAttachment)
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) GetUserCampaignByID(c *gin.Context) {
	var inputID campaign.GetCampaignDetailInput
	err := c.ShouldBindUri(&inputID)

	if err != nil {
		response := helper.APIResponse("Terjadi Kesalahan Saat Mengupdate", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	fmt.Println(inputID.ID)

	var campaignUser campaign.GetUserCampaign
	currentUser := c.MustGet("currentUser").(user.User)

	campaignUser.User = currentUser

	campaignDetail, err := h.campaignService.GetUserCampaignByID(inputID, campaignUser)

	if err != nil {
		response := helper.APIResponse("Terjadi Kesalahan", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil Mengambil Data", http.StatusOK, "Berhasil", campaign.FormatCampaignDetail(campaignDetail))
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

	var inputData campaign.UpdateCampaignInput

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
	response := helper.APIResponse("Berhasil Mengunggah Gambar", http.StatusOK, "sukses", data)
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaignReward(c *gin.Context) {
	var input campaign.CreateCampaignRewardInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Gagal Menambahkan Reward", http.StatusUnprocessableEntity, "gagal", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser
	newCampaignReward, err := h.campaignService.SaveCampaignReward(input)

	if err != nil {
		response := helper.APIResponse("Gagal Menambahkan Reward", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil Menambahkan Reward", http.StatusOK, "sukses", newCampaignReward)
	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) DeleteReward(c *gin.Context) {
	var input campaign.DeleteCampaignRewardInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Gagal Menghapus Reward", http.StatusUnprocessableEntity, "gagal", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser
	isDeleted, err := h.campaignService.DeleteReward(input)
	if err != nil {
		response := helper.APIResponse("Gagal Menghapus Reward", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil Menghapus Reward", http.StatusOK, "sukses", isDeleted)
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) DeleteImage(c *gin.Context) {
	var input campaign.DeleteCampaignImageInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Gagal Mennghapus Reward", http.StatusUnprocessableEntity, "gagal", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	isDeleted, err := h.campaignService.DeleteImage(input)
	if err != nil {
		response := helper.APIResponse("Gagal Menghapus Reward", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil Menghapus Reward", http.StatusOK, "sukses", isDeleted)
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateActivity(c *gin.Context) {
	var input campaign.CreateCampaignActivityInput

	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Gagal Buat Activity", http.StatusUnprocessableEntity, "gagal", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser
	fmt.Println(input)
	//! upload attachment
	file, err := c.FormFile("file")
	if err != nil {
		errorMessage := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Terjadi Kesalahan Saat Mengunggah Gambar", http.StatusBadRequest, "Gagal", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	path := fmt.Sprintf("images/activity/%d-%s", currentUser.ID, file.Filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		errorMessage := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Terjadi Kesalahan Saat Mengunggah Gambar", http.StatusBadRequest, "Gagal", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	input.ImageUrl = path

	newActivity, err := h.campaignService.CreateActivity(input)
	if err != nil {
		response := helper.APIResponse("Gagal buat Activity", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil Buat Campaign", http.StatusOK, "sukses", campaign.FormatCampaignActivity(newActivity))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UpdateActivity(c *gin.Context) {
	var input campaign.UpdateCampaignActivityInput
	err := c.ShouldBind(&input)
	fmt.Println(input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Gagal Ubah Activity", http.StatusUnprocessableEntity, "gagal", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser

	file, err := c.FormFile("file")
	if file != nil {
		if err != nil {
			errorMessage := gin.H{"is_uploaded": false}

			response := helper.APIResponse("Terjadi Kesalahan Saat Mengunggah Gambar", http.StatusBadRequest, "Gagal", errorMessage)
			c.JSON(http.StatusBadRequest, response)
			return
		}
		path := fmt.Sprintf("images/activity/%d-%s", currentUser.ID, file.Filename)
		err = c.SaveUploadedFile(file, path)
		if err != nil {
			errorMessage := gin.H{"is_uploaded": false}

			response := helper.APIResponse("Terjadi Kesalahan Saat Mengunggah Gambar", http.StatusBadRequest, "Gagal", errorMessage)
			c.JSON(http.StatusBadRequest, response)
			return
		}
		input.ImageUrl = path
	} else {
		// ambil imageURL terakhir jika tidak mengubah file
		currentActivity, err := h.campaignService.FindActivityByUser(
			campaign.GetCampaignActivityInput{
				ID:         input.ID,
				CampaignID: input.CampaignID,
			},
			campaign.GetUserCampaign{
				User: currentUser,
			},
		)
		if err != nil {
			response := helper.APIResponse("Gagal Ubah Activity", http.StatusBadRequest, "error", err.Error())
			c.JSON(http.StatusBadRequest, response)
			return
		}
		input.ImageUrl = currentActivity.ImageUrl
	}

	updateActivity, err := h.campaignService.UpdateActivity(input)
	if err != nil {
		response := helper.APIResponse("Gagal Ubah Activity", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil Buat Campaign", http.StatusOK, "sukses", campaign.FormatCampaignActivity(updateActivity))
	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) DeleteActivity(c *gin.Context) {
	var input campaign.DeleteCampaignActivityInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Gagal Menghapus Activity", http.StatusUnprocessableEntity, "gagal", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser

	isDeleted, err := h.campaignService.DeleteActivity(input)
	if err != nil {
		response := helper.APIResponse("Gagal Menghapus Activity", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil Menghapus Reward", http.StatusOK, "sukses", isDeleted)
	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) FindActivity(c *gin.Context) {
	var input campaign.GetCampaignActivityInput
	err := c.ShouldBindUri(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Gagal Menemukan Activity", http.StatusUnprocessableEntity, "gagal", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	activity, err := h.campaignService.FindActivity(input.ID)
	if err != nil {
		response := helper.APIResponse("Gagal Menemukan Activity", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil Menmukan Reward", http.StatusOK, "sukses", campaign.FormatCampaignActivity(activity))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) FindActivityByUser(c *gin.Context) {
	var input campaign.GetCampaignActivityInput
	err := c.ShouldBindUri(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Gagal Menemukan Activity", http.StatusUnprocessableEntity, "gagal", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	var userCampaign campaign.GetUserCampaign

	currentUser := c.MustGet("currentUser").(user.User)
	userCampaign.User = currentUser

	activity, err := h.campaignService.FindActivityByUser(input, userCampaign)
	if err != nil {
		response := helper.APIResponse("Gagal Menemukan Activity", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil Menemukan Activity", http.StatusOK, "sukses", campaign.FormatCampaignActivity(activity))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) FindAllActivityByCampaignID(c *gin.Context) {
	var input campaign.GetCampaignDetailInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Gagal Menemukan Activity", http.StatusUnprocessableEntity, "gagal", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	campaignActivity, err := h.campaignService.FindAllActivityByCampaignID(input.ID)
	if err != nil {
		response := helper.APIResponse("Gagal Menemukan Activity", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Berhasil Menmukan Activity", http.StatusOK, "sukses", campaign.FormatCampaignActivities(campaignActivity))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) PaginateCampaigns(c *gin.Context) {
	var input campaign.PaginateCampaignInput
	page, _ := strconv.Atoi(c.Query("page"))
	input.ActivePage = page
	err := c.ShouldBindJSON(&input)
	fmt.Println(input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Terjadi Kesalahan", http.StatusUnprocessableEntity, "gagal", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	paginate, err := h.campaignService.Paginate(input)
	if err != nil {
		response := helper.APIResponse("Gagal Mengambil Data", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil Menmukan Activity", http.StatusOK, "sukses", paginate)
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCattegory(c *gin.Context) {
	var input campaign.CattegoryInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Terjadi Kesalahan", http.StatusUnprocessableEntity, "gagal", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	cattegory, err := h.campaignService.CreateCattegory(input)

	if err != nil {
		response := helper.APIResponse("Gagal Input Cattegory", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil input cattegory", http.StatusOK, "sukses", cattegory)
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) DeleteCattegory(c *gin.Context) {
	var input campaign.CattegoryIdInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Terjadi Kesalahan", http.StatusUnprocessableEntity, "gagal", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	delete, err := h.campaignService.DeleteCattegory(input.ID)
	if err != nil {
		response := helper.APIResponse("Gagal delete Cattegory", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil delete Activity", http.StatusOK, "sukses", delete)
	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) FindAllCattegory(c *gin.Context) {
	cattegory, err := h.campaignService.FindAllCattegory()

	if err != nil {
		response := helper.APIResponse("Gagal Mengambil Data Cattegory", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Berhasil ambil cattegory", http.StatusOK, "sukses", cattegory)
	c.JSON(http.StatusOK, response)

}

// tangkap input kedalam struct
// ambil current user dengan jwt
// panggil service , parameternya input struct
// panggil repo untuk simpan data
