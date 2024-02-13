package routes

import (
	"project_first/controllers"

	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.RouterGroup) {
	router.GET("/insert-fields", controllers.DisplayInsertionModule)
	router.GET("/display-items", controllers.DisplayItems)
	router.GET("/admin-settings", controllers.GetAdminSettings)
	api := router.Group("/api")
	{
		api.POST("/insert-items", controllers.InsertItems)
		api.GET("/display-all-items", controllers.DisplayAllItems)
		api.GET("/exchange-currency", controllers.CurrencyExchange)
		api.PUT("/update-admin-settings",controllers.UpdateAdminSettings)
	}
}

func SetupRouter() *gin.Engine {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	AddRoutes(&router.RouterGroup)
	return router
}
