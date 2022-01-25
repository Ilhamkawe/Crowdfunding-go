package main

import (
	"crowdfunding-TA/auth"
	"crowdfunding-TA/campaign"
	"crowdfunding-TA/handler"
	"crowdfunding-TA/middleware"
	"crowdfunding-TA/payment"
	"crowdfunding-TA/transaction"
	"crowdfunding-TA/user"
	webHandler "crowdfunding-TA/web/handler"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
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

	// !=================================================================================

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
	campaignWebHandler := webHandler.NewCampaignHandler(campaignService, userService)
	transactionWebHandler := webHandler.NewTransactionHandler(transactionService)
	sessionWebHandler := webHandler.NewSessionHandler(userService)
	// !=================================================================================

	// * Task Scheduler
	// !=================================================================================
	// set jadwal berdasarkan zona waktu indonesia
	jakarta, _ := time.LoadLocation("Asia/Jakarta")
	scheduler := cron.New(cron.WithLocation(jakarta))

	defer scheduler.Stop()

	scheduler.AddFunc("0 0 * * *", campaignService.IsCollectAbleByDate)

	go scheduler.Start()

	// ? test

	// membuat Router
	router := gin.Default()
	router.Use(cors.Default())
	router.Use(cors.New(
		cors.Config{
			AllowOrigins: []string{"http://localhost:3000"},
			AllowMethods: []string{"POST", "GET", "PATCH", "DELETE", "HEAD"},
			AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		}))
	cookieStore := cookie.NewStore(auth.SECRET_KEY)
	router.Use(sessions.Sessions("userID", cookieStore))
	// router.LoadHTMLGlob("web/templates/**/*")
	router.HTMLRender = loadTemplates("./web/templates")
	// Static Route
	router.Static("css", "./web/assets/css")
	router.Static("js", "./web/assets/js")
	router.Static("vendors", "./web/assets/vendors")
	router.Static("web/images", "./web/assets/images")
	router.Static("images", "./images")
	router.Static("attachment", "./attachment")
	// grouping API
	api := router.Group("/api/v1")
	// User Route
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/session", userHandler.Login)
	api.POST("/email_chekers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", middleware.AuthMiddleware(authService, userService), userHandler.UploadAvatar)
	api.PUT("/users/update", middleware.AuthMiddleware(authService, userService), userHandler.UpdateUserInfo)
	api.GET("/users/fetch", middleware.AuthMiddleware(authService, userService), userHandler.FetchUser)
	api.PUT("/users/changepass", middleware.AuthMiddleware(authService, userService), userHandler.ChangePassword)
	// Campaign Route
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.POST("/campaigns/paginate", campaignHandler.PaginateCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.GET("/campaigns/L/:limit", campaignHandler.Limit)
	api.GET("/campaigns/R/:id", campaignHandler.GetRewards)
	api.GET("/campaign/:id/user", middleware.AuthMiddleware(authService, userService), campaignHandler.GetUserCampaignByID)
	api.POST("/campaigns", middleware.AuthMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/attachment", middleware.AuthMiddleware(authService, userService), campaignHandler.UpdateAttachment)
	api.PUT("/campaigns/:id", middleware.AuthMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-image", middleware.AuthMiddleware(authService, userService), campaignHandler.CreateCampaignImage)
	api.POST("/campaign-reward", middleware.AuthMiddleware(authService, userService), campaignHandler.CreateCampaignReward)
	api.POST("/campaign-reward/delete", middleware.AuthMiddleware(authService, userService), campaignHandler.DeleteReward)
	api.DELETE("/campaign-image", middleware.AuthMiddleware(authService, userService), campaignHandler.DeleteImage)
	api.POST("/campaign/search/paginate", campaignHandler.SearchCampaignPaginate)
	api.POST("/campaign/search", campaignHandler.SearchCampaign)
	api.POST("/campaign/activity", middleware.AuthMiddleware(authService, userService), campaignHandler.CreateActivity)
	api.PUT("/campaign/activity", middleware.AuthMiddleware(authService, userService), campaignHandler.UpdateActivity)
	api.DELETE("/campaign/activity", middleware.AuthMiddleware(authService, userService), campaignHandler.DeleteActivity)
	api.GET("/campaign/user/activity/:id/:campaign_id", middleware.AuthMiddleware(authService, userService), campaignHandler.FindActivityByUser)
	api.GET("/campaign/activity/find/:id", campaignHandler.FindActivity)
	api.GET("/campaign/activity/:id", campaignHandler.FindAllActivityByCampaignID)
	// Cattegory Route
	api.GET("/cattegory", campaignHandler.FindAllCattegory)
	api.DELETE("/cattegory/:id/delete", campaignHandler.DeleteCattegory)
	api.POST("/cattegory", campaignHandler.CreateCattegory)
	// transaction Route
	api.GET("/campaigns/:id/transactions", middleware.AuthMiddleware(authService, userService), transactionHandler.GetCampaignTransaction)
	api.POST("/campaigns/transactions/:id/reward", middleware.AuthMiddleware(authService, userService), transactionHandler.GetAllByReward)
	api.GET("/transactions", middleware.AuthMiddleware(authService, userService), transactionHandler.GetUserTransaction)
	api.POST("/transactions", middleware.AuthMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.GetNotification)
	api.POST("/collect", middleware.AuthMiddleware(authService, userService), transactionHandler.CollectAmount)
	api.GET("/collect/:id", middleware.AuthMiddleware(authService, userService), transactionHandler.FindCollectData)

	// CMS Route
	router.GET("/users", middleware.AdminMiddleware(), userWebHandler.Index)
	router.GET("/users/new", middleware.AdminMiddleware(), userWebHandler.New)
	router.POST("/users/new", middleware.AdminMiddleware(), userWebHandler.Create)
	router.GET("/users/:id/edit", middleware.AdminMiddleware(), userWebHandler.Edit)
	router.POST("/users/:id/update", middleware.AdminMiddleware(), userWebHandler.Update)
	router.GET("/users/:id/avatar", middleware.AdminMiddleware(), userWebHandler.Avatar)
	router.POST("/users/:id/avatar", middleware.AdminMiddleware(), userWebHandler.CreateAvatar)
	router.GET("/campaigns", middleware.AdminMiddleware(), campaignWebHandler.Index)
	router.GET("/campaign/new", middleware.AdminMiddleware(), campaignWebHandler.New)
	router.POST("/campaign/new", middleware.AdminMiddleware(), campaignWebHandler.Create)
	router.GET("/campaign/:id/images", middleware.AdminMiddleware(), campaignWebHandler.Image)
	router.POST("/campaign/:id/images", middleware.AdminMiddleware(), campaignWebHandler.CreateImage)
	router.GET("/campaign/:id/edit", middleware.AdminMiddleware(), campaignWebHandler.Edit)
	router.POST("/campaign/:id/edit", middleware.AdminMiddleware(), campaignWebHandler.Update)
	router.GET("/campaign/:id/show", middleware.AdminMiddleware(), campaignWebHandler.Detail)
	router.GET("/campaign/:id/:status", middleware.AdminMiddleware(), campaignWebHandler.ChangeStatus)
	router.GET("/cattegory", middleware.AdminMiddleware(), campaignWebHandler.CattegoryIndex)
	router.POST("/cattegory", middleware.AdminMiddleware(), campaignWebHandler.StoreCattegory)
	router.GET("/cattegory/:id/delete", middleware.AdminMiddleware(), campaignWebHandler.DeleteCattegory)
	router.GET("/transactions", middleware.AdminMiddleware(), transactionWebHandler.Index)
	router.GET("/collect", middleware.AdminMiddleware(), transactionWebHandler.CollectList)
	router.GET("/collect/:id", middleware.AdminMiddleware(), transactionWebHandler.Collect)
	router.GET("/collect/:id/:status", middleware.AdminMiddleware(), transactionWebHandler.ChangeCollectStatus)
	router.GET("/login", sessionWebHandler.Index)
	router.POST("/login", sessionWebHandler.Login)
	router.GET("/logout", sessionWebHandler.Logout)
	router.Run()

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

	// Generate our templates map from our layouts/
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include, "web/templates/partial/partial.html")
		// fmt.Println(files)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}

// ambil nilai header Authorization (token saja)
// validasi token
// ambil user_id
// ambil user dari db berdasarkan id lewat service
// set context isinya user

func TestScheduler() {

	fmt.Println("Berhasil Menggunakan Cron Job")

}
