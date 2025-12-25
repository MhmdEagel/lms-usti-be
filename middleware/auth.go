package middleware

import (
	"net/http"

	"github.com/MhmdEagel/lms-usti-be/lib"
	"github.com/MhmdEagel/lms-usti-be/model"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		claims, err := lib.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized."})
			c.Abort()
			return
		}
		c.Set("user", model.Me{Email: claims.Email, Role: claims.Role, UserId: claims.UserId})
		c.Next()
	}
}


func AuthDosenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exist := c.Get("user")
		if !exist {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan"})
			return
		}
		user := val.(model.Me)
		foundUser, err := model.GetUserByEmail(user.Email)
		if err != nil || foundUser.Role != "DOSEN" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized."})
			c.Abort()
			return
		}
		c.Next()
	}
}
func AuthMahasiswaMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exist := c.Get("user")
		if !exist {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan"})
			return
		}
		user := val.(model.Me)
		foundUser, err := model.GetUserByEmail(user.Email)
		if err != nil || foundUser.Role != "MAHASISWA" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized."})
			c.Abort()
			return
		}
		c.Next()
	}
}
