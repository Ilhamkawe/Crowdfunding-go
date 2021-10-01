package middleware

import (
	"crowdfunding-TA/auth"
	"crowdfunding-TA/helper"
	"crowdfunding-TA/user"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// ? Admin Middleware
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userIDSession := session.Get("userID")

		if userIDSession == nil {
			c.Redirect(http.StatusFound, "/login")
		}
	}
}

// ? API Middleware
func AuthMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Bearer Token
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")

		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)

		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}
