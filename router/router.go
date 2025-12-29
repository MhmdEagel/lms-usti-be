package router

import (
	"github.com/MhmdEagel/lms-usti-be/controller"
	"github.com/MhmdEagel/lms-usti-be/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20
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
		classroom.Use(middleware.AuthMiddleware())
		{
			classroom.GET("", controller.ReadClassrooms)
			classroom.POST("/join", controller.JoinClassroom)
			classroom.POST("/create", middleware.AuthDosenMiddleware(), controller.CreateClassroom)

			classroom.GET("/:id", controller.ReadDetailClassroom)
			classroom.DELETE("/:id/delete", middleware.AuthDosenMiddleware(), controller.DestroyClassroom)
			classroom.PATCH("/:id/update", middleware.AuthDosenMiddleware(), controller.UpdateClassroom)
			
			classroom.POST("/:id/announcement/create", middleware.AuthDosenMiddleware(), controller.CreateAnnouncement)
			classroom.POST("/:id/announcement/delete", middleware.AuthDosenMiddleware(), controller.DeleteAnnouncement)

			classroom.POST("/:id/material/upload", middleware.AuthDosenMiddleware(), controller.CreateMaterial)
				
		}
	}
	return r
}
