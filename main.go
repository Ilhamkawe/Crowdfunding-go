package main

import (
	"crowdfunding-TA/auth"
	"crowdfunding-TA/handler"
	"crowdfunding-TA/helper"
	"crowdfunding-TA/user"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:@tcp(127.0.0.1:3306)/crowdfunding?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}
	// mendeklarasikan repository untuk menyimpan data
	userRepository := user.NewRepository(db)
	// mendeklarasikan service untuk mapping data yang diinputkan menjadi struct user
	userService := user.NewService(userRepository)
	// mendeklarasikan service untuk meng generate jwt token
	authService := auth.NewService()
	// handler
	userHandler := handler.NewUserHandler(userService, authService)

	// test
	token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMH0.vlUN2r5bK4l7_BMf7jtFgDkWcify5grGLPU2XX2SwDo")

	if err != nil {
		fmt.Println("Error")
		fmt.Println("Error")
		fmt.Println("Error")
	}

	if token.Valid {
		fmt.Println("Valid")
		fmt.Println("Valid")
		fmt.Println("Valid")
	} else {
		fmt.Println("Invalid")
		fmt.Println("Invalid")
		fmt.Println("Invalid")
	}

	userService.SaveAvatar(3, "images/avatar1.png")
	fmt.Println(authService.GenerateToken(1001))
	// membuat Router
	router := gin.Default()
	// grouping API
	api := router.Group("/api/v1")
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/session", userHandler.Login)
	api.POST("/email_chekers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	router.Run()

}

// ? Middleware
func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
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

// ambil nilai header Authorization (token saja)
// validasi token
// ambil user_id
// ambil user dari db berdasarkan id lewat service
// set context isinya user
