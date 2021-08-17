package handler

import (
	"crowdfunding-TA/auth"
	"crowdfunding-TA/helper"
	"crowdfunding-TA/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	// tangkap input user
	// map input dari user ke struct RegisterUserInput
	// struct di atas kita passing sebagai parameter service

	var input user.RegisterInputUser

	err := c.ShouldBindJSON(&input)

	if err != nil {

		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Gagal buat akun", http.StatusUnprocessableEntity, "gagal", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse("Gagal buat akun", http.StatusBadRequest, "gagal", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	token, err := h.authService.GenerateToken(newUser.ID)

	if err != nil {
		response := helper.APIResponse("Gagal buat akun", http.StatusBadRequest, "gagal", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, token)

	response := helper.APIResponse("Berhasil daftar akun", http.StatusOK, "sukses", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	// user input email dan password
	// input dintangkap handler
	// mapping dari input user ke input struct
	// input struct passing service
	// di sservice mencari dengan bantuan repository user dengan email
	// mencocokan password

	var input user.LoginInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Gagal Login", http.StatusUnprocessableEntity, "gagal", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	loggedInUser, err := h.userService.Login(input)

	if err != nil {
		errorMessage := gin.H{"errors :": err.Error()}
		response := helper.APIResponse("Gagal Login", http.StatusUnprocessableEntity, "gagal", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.authService.GenerateToken(loggedInUser.ID)

	if err != nil {
		response := helper.APIResponse("Gagal Login", http.StatusBadRequest, "gagal", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loggedInUser, token)

	response := helper.APIResponse("Login Berhasil", http.StatusOK, "Berhasil", formatter)

	c.JSON(http.StatusOK, response)

}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	// ambil input email dari user
	// mapping email ke struct input
	// passing struct input ke service
	// cek email dengan repository
	// repo-db

	var input user.CheckEmailInput

	err := c.ShouldBind(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Gagal saat check email", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)

		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)

	if err != nil {
		errorMessage := gin.H{"errors": "Terjadi Kesalahan"}

		response := helper.APIResponse("Gagal saat check email", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)

		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	metaMessage := "Email Sudah Terdaftar"

	if isEmailAvailable {
		metaMessage = "Email bisa digunakan"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "Berhasil", data)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	// ambil input dari user
	// simpan gambar ke folger image/
	// di service panggil repo
	// jwt sementara
	// repo ambil data user yang ID = 1
	// repo update data user simpan lokasi file

	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Terjadi kesalahan saat mengunggah avatar", http.StatusBadRequest, "Gagal", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	// ini diisi jwt
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)
	err = c.SaveUploadedFile(file, path)

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Terjadi kesalahan saat mengunggah avatar", http.StatusBadRequest, "Gagal", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Terjadi kesalahan saat mengunggah avatar", http.StatusBadRequest, "Gagal", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Terjadi kesalahan saat mengunggah avatar", http.StatusOK, "Berhasil", data)
	c.JSON(http.StatusBadRequest, response)

}
