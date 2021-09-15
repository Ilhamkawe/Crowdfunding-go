package main

import (
	"crowdfunding-TA/auth"
	"crowdfunding-TA/campaign"
	"crowdfunding-TA/handler"
	"crowdfunding-TA/helper"
	"crowdfunding-TA/payment"
	"crowdfunding-TA/transaction"
	"crowdfunding-TA/user"
	webHandler "crowdfunding-TA/web/handler"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
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

	// * Middleware
	// !=================================================================================
	authService := auth.NewService()
	// !=================================================================================

	// * user dependencies
	// !=================================================================================
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService, authService)
	// !=================================================================================

	// * campaign dependencies
	// !=================================================================================
	campaignRepository := campaign.NewRepository(db)
	campaignService := campaign.NewService(campaignRepository)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	// !=================================================================================

	// * Transaction dependencies
	// !=================================================================================
	transactionRepository := transaction.NewRepository(db)
	paymentService := payment.NewService(campaignRepository)
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	// !=================================================================================

	//  * WebHandler dependencies
	// !=================================================================================
	userWebHandler := webHandler.NewUserHandler(userService)
	// !=================================================================================

	// ? test

	// membuat Router
	router := gin.Default()
	router.Use(cors.Default())
	// router.LoadHTMLGlob("web/templates/**/*")
	router.HTMLRender = loadTemplates("./web/templates")
	// Static Route
	router.Static("css", "./web/assets/css")
	router.Static("js", "./web/assets/js")
	router.Static("vendors", "./web/assets/vendors")
	router.Static("web/images", "./web/assets/images")
	router.Static("images", "./images")
	// grouping API
	api := router.Group("/api/v1")
	// User Route
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/session", userHandler.Login)
	api.POST("/email_chekers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	// Campaign Route
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-image", authMiddleware(authService, userService), campaignHandler.CreateCampaignImage)

	// transaction Route
	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransaction)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransaction)
	api.POST("/transactions", authMiddleware(authService, userService), transactionHandler.CreateTransaction)

	// CMS Route
	router.GET("/users", userWebHandler.Index)
	router.GET("/users/new", userWebHandler.New)
	router.POST("/users/new", userWebHandler.Create)
	router.GET("/users/:id/edit", userWebHandler.Edit)
	router.POST("/users/:id/update", userWebHandler.Update)
	router.GET("/users/:id/avatar", userWebHandler.Avatar)
	router.POST("/users/:id/avatar", userWebHandler.CreateAvatar)
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

// ? gin multitemplate

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/**/*")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}

// ambil nilai header Authorization (token saja)
// validasi token
// ambil user_id
// ambil user dari db berdasarkan id lewat service
// set context isinya user
