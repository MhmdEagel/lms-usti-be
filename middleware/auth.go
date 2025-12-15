package middleware

import (
	"fmt"
	"net/http"

	"github.com/MhmdEagel/lms-usti-be/lib"
	"github.com/MhmdEagel/lms-usti-be/model"
	"github.com/gin-gonic/gin"
)

func AuthDosenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		claims, err := lib.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token."})
			c.Abort()
			return
		}
		user, err := model.GetUserByEmail(claims.Email)
		fmt.Println(user.Email)
		if err != nil || user.Role != "DOSEN" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized."})
			c.Abort()
			return
		}
		c.Set("userId", user.ID)
		c.Next()
	}
}
