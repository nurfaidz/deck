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

	router.SetTrustedProxies(nil)

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	// Initialize services
	transactionService := services.NewTransactionService(database.DB)
	midtransService := services.NewMidtransService()

	// Initialize controllers
	transactionController := controllers.NewTransactionController(database.DB, transactionService)
	paymentController := controllers.NewPaymentController(database.DB, midtransService)

	apiRouter := router.Group("/api/")

	apiRouter.POST("login", controllers.Login)

	// router user
	apiRouter.GET("users", middlewares.AuthMiddleware(), controllers.GetUsers)
	apiRouter.POST("users", middlewares.AuthMiddleware(), controllers.CreateUser)
	apiRouter.GET("users/:id", middlewares.AuthMiddleware(), controllers.GetUserById)
	apiRouter.PUT("users/:id", middlewares.AuthMiddleware(), controllers.UpdateUser)
	apiRouter.DELETE("users/:id", middlewares.AuthMiddleware(), controllers.DeleteUser)

	// router product
	apiRouter.GET("products", controllers.GetProducts)
	apiRouter.POST("products", middlewares.AuthMiddleware(), controllers.CreateProduct)
	apiRouter.GET("products/filter-name", controllers.FilterByName)
	apiRouter.GET("products/filter-category", controllers.FilterByCategory)
	apiRouter.GET("products/:id", middlewares.AuthMiddleware(), controllers.GetProductById)
	apiRouter.PUT("products/:id", middlewares.AuthMiddleware(), controllers.UpdateProduct)
	apiRouter.DELETE("products/:id", middlewares.AuthMiddleware(), controllers.DeleteProduct)

	// Public endpoints untuk customer
	apiRouter.POST("transactions", transactionController.CreateTransaction)
	apiRouter.GET("transactions/:order_number", transactionController.GetTransaction)

	// Protected endpoints untuk admin
	apiRouter.GET("transactions", middlewares.AuthMiddleware(), transactionController.GetAllTransactions)
	apiRouter.GET("transactions/id/:id", middlewares.AuthMiddleware(), transactionController.GetTransactionByID)

	// Payment routes
	apiRouter.POST("payments", paymentController.CreatePayment)
	apiRouter.POST("payments/notification", paymentController.MidtransNotification)
	apiRouter.GET("payments/:order_number", paymentController.GetPaymentStatus)

	return router
}
