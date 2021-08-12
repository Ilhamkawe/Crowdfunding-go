package main

import (
	"crowdfunding-TA/handler"
	"crowdfunding-TA/user"
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
	// handler
	userHandler := handler.NewUserHandler(userService)

	// test handler
	userService.SaveAvatar(3, "images/avatar1.png")
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
