package router

import (
	"github.com/MhmdEagel/lms-usti-be/controller"
	"github.com/MhmdEagel/lms-usti-be/middleware"
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
			auth.POST("/verify-email", controller.VerifyEmail)
			auth.POST("/verify-email/resend", controller.ResendVerification)
		}
		classroom := api.Group("/classroom")
		classroom.Use(middleware.AuthDosenMiddleware())
		{
			classroom.POST("/create", controller.CreateClassroom)
		}
	}
	return r
}