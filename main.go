package main

import (
	"github.com/akshay-kumart/go-api/controllers"
	"github.com/akshay-kumart/go-api/initializers"
	"github.com/akshay-kumart/go-api/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	initializers.LoadEnv()
	r := gin.Default()

	r.POST("/signUp", controllers.SignUp)
	r.POST("/login", controllers.Login)
	protected := r.Group("/api", middleware.AuthMiddle)
	{
		protected.GET("/validate", controllers.Validate)
		protected.GET("/role", middleware.AdminOnly(), controllers.Role)
	}
	r.Run(":8090")
}
