package main

import (
	"crowdfunding-TA/auth"
	"crowdfunding-TA/handler"
	"crowdfunding-TA/user"
	"fmt"
	"log"

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
	token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMH0.MwlOWX5KQmsDhcpAuINO9b8oXNVI4dvPJ1URiNZ7UIE")

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
	api.POST("/avatars", userHandler.UploadAvatar)

	router.Run()

}
