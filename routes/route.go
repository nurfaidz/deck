package routes

import (
	"deck/controllers"
	"deck/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	apiRouter := router.Group("/api/")

	apiRouter.POST("login", controllers.Login)

	// router user
	apiRouter.GET("users", middlewares.AuthMiddleware(), controllers.GetUsers)
	apiRouter.POST("users", middlewares.AuthMiddleware(), controllers.CreateUser)
	apiRouter.GET("users/:id", middlewares.AuthMiddleware(), controllers.GetUserById)
	apiRouter.PUT("users/:id", middlewares.AuthMiddleware(), controllers.UpdateUser)
	apiRouter.DELETE("users/:id", middlewares.AuthMiddleware(), controllers.DeleteUser)

	// router product
	apiRouter.GET("products", middlewares.AuthMiddleware(), controllers.GetProducts)
	apiRouter.POST("products", middlewares.AuthMiddleware(), controllers.CreateProduct)
	apiRouter.GET("products/:id", middlewares.AuthMiddleware(), controllers.GetProductById)
	apiRouter.PUT("products/:id", middlewares.AuthMiddleware(), controllers.UpdateProduct)
	apiRouter.DELETE("products/:id", middlewares.AuthMiddleware(), controllers.DeleteProduct)

	// router transaction
	apiRouter.GET("transactions", middlewares.AuthMiddleware(), controllers.GetTransactions)

	return router
}
