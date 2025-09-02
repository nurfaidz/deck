package routes

import (
	"deck/controllers"
	"deck/database"
	"deck/middlewares"
	"deck/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	router.Static("/uploads", "./uploads")

	router.SetTrustedProxies(nil)

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	// Initialize services
	transactionService := services.NewTransactionService(database.DB)
	notificationService := services.NewNotificationService(database.DB)
	productService := services.NewProductService(database.DB)

	// Initialize controllers
	transactionController := controllers.NewTransactionController(database.DB, transactionService, notificationService)
	notificationController := controllers.NewNotificationController(notificationService)
	productController := controllers.NewProductController(productService)

	apiRouter := router.Group("/api/")

	apiRouter.POST("login", controllers.Login)

	// router user
	apiRouter.GET("users", middlewares.AuthMiddleware(), controllers.GetUsers)
	apiRouter.POST("users", middlewares.AuthMiddleware(), controllers.CreateUser)
	apiRouter.GET("users/:id", middlewares.AuthMiddleware(), controllers.GetUserById)
	apiRouter.PUT("users/:id", middlewares.AuthMiddleware(), controllers.UpdateUser)
	apiRouter.DELETE("users/:id", middlewares.AuthMiddleware(), controllers.DeleteUser)

	// router product
	apiRouter.GET("products", productController.GetProducts)
	apiRouter.POST("products", middlewares.AuthMiddleware(), productController.CreateProduct)
	apiRouter.GET("products/:id", middlewares.AuthMiddleware(), productController.GetProductById)
	apiRouter.PUT("products/:id", middlewares.AuthMiddleware(), controllers.UpdateProduct)
	apiRouter.DELETE("products/:id", middlewares.AuthMiddleware(), productController.DeleteProduct)

	// route category
	apiRouter.GET("categories", controllers.GetCategories)

	// route transaction
	apiRouter.POST("transactions", transactionController.CreateTransaction)
	//apiRouter.GET("transactions/:order_number", transactionController.GetTransaction)
	apiRouter.GET("transactions", middlewares.AuthMiddleware(), transactionController.GetAllTransactions)
	apiRouter.GET("transactions/:id", transactionController.GetTransactionByID)

	apiRouter.GET("notifications", middlewares.AuthMiddleware(), notificationController.GetNotifications)
	apiRouter.GET("notifications/unread-count", middlewares.AuthMiddleware(), notificationController.GetUnreadCount)
	apiRouter.PUT("notifications/:id/read", middlewares.AuthMiddleware(), notificationController.MarkAsRead)
	apiRouter.PUT("notifications/mark-all-read", middlewares.AuthMiddleware(), notificationController.MarkAllAsRead)

	return router
}
