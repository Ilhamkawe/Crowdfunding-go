package handler

import (
	"crowdfunding-TA/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) Index(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}
	c.HTML(http.StatusOK, "user_index.html", gin.H{"users": users})
}

func (h *userHandler) New(c *gin.Context) {
	c.HTML(http.StatusOK, "user_new.html", nil)
}

func (h *userHandler) Create(c *gin.Context) {
	var input user.FormInputRegister

	err := c.ShouldBind(&input)
	if err != nil {
		c.HTML(http.StatusOK, "user_new.html", gin.H{
			"input": input,
			"error": true,
		})
		return
	}

	registerInput := user.RegisterInputUser{}
	registerInput.Name = input.Name
	registerInput.Email = input.Email
	registerInput.Occupation = input.Occupation
	registerInput.Password = input.Password

	_, err = h.userService.RegisterUser(registerInput)

	if err != nil {
		c.HTML(http.StatusOK, "error.html", nil)
	}

	c.Redirect(http.StatusFound, "/users")
}

func (h *userHandler) Edit(c *gin.Context) {
	idParams := c.Param("id")
	id, _ := strconv.Atoi(idParams)

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}

	c.HTML(http.StatusOK, "user_edit.html", gin.H{
		"users": user,
	})

}

func (h *userHandler) Update(c *gin.Context) {
	idParams := c.Param("id")
	id, _ := strconv.Atoi(idParams)

	var input user.FormUpdateRegister

	err := c.ShouldBind(&input)

	if err != nil {
		c.HTML(http.StatusOK, "user_edit.html", gin.H{
			"error": true,
		})
		return
	}

	input.ID = id

	_, err = h.userService.UpdateUser(input)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/users")

}

func (h *userHandler) Avatar(c *gin.Context) {

	idParams := c.Param("id")
	id, _ := strconv.Atoi(idParams)

	c.HTML(http.StatusOK, "user_avatar.html", gin.H{
		"id": id,
	})

}

func (h *userHandler) CreateAvatar(c *gin.Context) {
	idParams := c.Param("id")
	id, _ := strconv.Atoi(idParams)

	file, err := c.FormFile("avatar")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	userID := id

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/users")

}
