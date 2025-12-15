package controller

import (
	"net/http"

	"github.com/MhmdEagel/lms-usti-be/lib"
	"github.com/MhmdEagel/lms-usti-be/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func CreateClassroom(c *gin.Context) {
	var body model.CreateClassroomInput
	if err := c.ShouldBindJSON(&body); err != nil {
		msg := lib.GetValidationMessage(err.(validator.ValidationErrors))
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	userId, exist := c.Get("userId")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan"})
		return
	}
	body.DosenId = userId.(string)
	newClassroom := model.NewClassroom(&body)
	if err := model.StoreClassroom(newClassroom); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan."})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Kelas berhasil dibuat."})
}
