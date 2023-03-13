package main

import (
	"jwt-split-token/controllers"
	"jwt-split-token/database"
	"jwt-split-token/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {

	// initialize Database
	connectionString, port := SetENVs()
	database.Connect(connectionString)
	database.Migrate()

	// initialize router
	router := routerInit()
	router.Run(port)
}

func routerInit() *gin.Engine {
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
