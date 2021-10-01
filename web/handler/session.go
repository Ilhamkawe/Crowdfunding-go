package handler

import (
	"crowdfunding-TA/user"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type sessionHandler struct {
	userService user.Service
}

func NewSessionHandler(userService user.Service) *sessionHandler {
	return &sessionHandler{userService}
}

func (h *sessionHandler) Index(c *gin.Context) {
	session := sessions.Default(c)
	userIDSession := session.Get("userID")

	if userIDSession != nil {
		c.Redirect(http.StatusFound, "/users")
	}
	c.HTML(http.StatusOK, "session_index.html", gin.H{
		"title": "Masuk Sebagai Admin | UTY Crowdfunding Admin",
	})
}

func (h *sessionHandler) Login(c *gin.Context) {
	var input user.LoginInput

	err := c.ShouldBind(&input)
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	user, err := h.userService.Login(input)
	if err != nil || user.Role != "admin" {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	session := sessions.Default(c)
	session.Set("userID", user.ID)
	session.Set("username", user.Name)
	session.Save()

	c.Redirect(http.StatusFound, "/users")
}

func (h *sessionHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Redirect(http.StatusFound, "/login")
}
