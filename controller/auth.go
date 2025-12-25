package controller

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/MhmdEagel/lms-usti-be/lib"
	"github.com/MhmdEagel/lms-usti-be/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func Login(c *gin.Context) {
	var body model.Login
	if err := c.ShouldBindJSON(&body); err != nil {
		msg := lib.GetValidationMessage(err.(validator.ValidationErrors))
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	res, err := model.GetUserByEmail(body.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email tidak ditemukan"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan."})
		return
	}
	if !res.EmailVerified.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email belum diverifikasi"})
		return
	}
	if lib.IsPasswordMatch(res.Password, body.Password) {
		token, err := lib.CreateToken(res.Fullname, res.Email, res.Role, res.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan."})
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{"error": "email atau password salah."})

}
func Register(c *gin.Context) {
	var body model.Register
	if err := c.ShouldBindJSON(&body); err != nil {
		msg := lib.GetValidationMessage(err.(validator.ValidationErrors))
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	user := model.NewUser(&body)
	if err := model.CreateUser(user); err != nil {
		log.Println(err.Error())
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah digunakan."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan."})
		return
	}
	token := model.NewVerificationToken(user.Email)
	if err := model.CreateVerificationToken(token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan."})
		log.Println(err.Error())
		return
	}
	if err := lib.SendVerificationEmail(user.Email, token.Token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan."})
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User successfully created."})
}
func VerifyEmail(c *gin.Context) {
	var body model.Verification
	if err := c.ShouldBindJSON(&body); err != nil {
		msg := lib.GetValidationMessage(err.(validator.ValidationErrors))
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	t, err := model.GetVerificationTokenByToken(body.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if lib.IsExpired(&t.Expires) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token sudah tidak berlaku"})
		model.DB.Where("id = ?", t.ID).Delete(&model.VerificationToken{})
		return
	}
	res := model.DB.Model(&model.User{}).Where("email = ?", t.Email).Update("email_verified", time.Now())
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan"})
		return
	}
	if res := model.DB.Where("id = ?", t.ID).Delete(&model.VerificationToken{}); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User successfully verified"})
}

func ResendVerification(c *gin.Context) {
	var body model.ResendVerificationInput
	if err := c.ShouldBindJSON(&body); err != nil {
		msg := lib.GetValidationMessage(err.(validator.ValidationErrors))
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	user, err := model.GetUserByEmail(body.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email tidak ditemukan"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "terjadi kesalahan"})
		return
	}
	token := model.NewVerificationToken(user.Email)
	if err := model.CreateVerificationToken(token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan."})
		return
	}
	if err := lib.SendVerificationEmail(token.Email, token.Token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengirimkan email."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email verifikasi berhasil dikirim."})
}
