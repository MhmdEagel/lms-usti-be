package router

import (
	"github.com/MhmdEagel/lms-usti-be/controller"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	
	api := r.Group("/lms-usti-api")
	{
		api.GET("", controller.Test)
		auth := api.Group("/auth")
		{
			auth.POST("/login", controller.Login)
			auth.POST("/register", controller.Register)
			auth.POST("/verify", controller.Verify)
		}
	}
	return r
}