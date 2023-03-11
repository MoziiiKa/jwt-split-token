package main

import (
	"fmt"
	"jwt-split-token/controllers"
	"jwt-split-token/database"
	"jwt-split-token/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {

	// Load Configurations from config.json using Viper
	LoadAppConfig()

	fmt.Println("TokenMaxAge: ", AppConfig.TokenMaxAge)

	// Initialize Database
	database.Connect(AppConfig.ConnectionString)
	database.Migrate()
	// Initialize Router
	router := initRouter()
	err := router.Run(AppConfig.Port)
	if err != nil {
		return
	}

}

func initRouter() *gin.Engine {
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		api.POST("/token-management/token", controllers.GenerateToken)
		api.POST("/user-management/registration", controllers.RegisterUser)
		secured := api.Group("/access-management").Use(middlewares.Auth())
		{
			secured.GET("/time-ir", controllers.Access)
		}
	}
	return router
}
