package routes

import (
	"inventory_management/handlers"
	"inventory_management/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.Use(middleware.RateLimiterMiddleware())

	api := router.Group("/api/v1")
	{
		api.POST("/login", handlers.Login)
		items := api.Group("/inventory")
		{
			items.GET("", handlers.GetAllItems)                                       // GET /api/v1/inventory
			items.GET("/:id", handlers.GetItemByID)                                   // GET /api/v1/inventory/:id
			items.POST("", middleware.JWTAuthMiddleware(), handlers.CreateItem)       // POST /api/v1/inventory
			items.PUT("/:id", middleware.JWTAuthMiddleware(), handlers.UpdateItem)    // PUT /api/v1/inventory/:id
			items.DELETE("/:id", middleware.JWTAuthMiddleware(), handlers.DeleteItem) // DELETE /api/v1/inventory/:id
		}
	}

	return router
}
