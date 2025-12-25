package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/MhmdEagel/lms-usti-be/lib"
	"github.com/MhmdEagel/lms-usti-be/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func CreateClassroom(c *gin.Context) {
	var body model.CreateClassroomInput
	if err := c.ShouldBindJSON(&body); err != nil {
		msg := lib.GetValidationMessage(err.(validator.ValidationErrors))
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	val, exist := c.Get("user")
	user := val.(model.Me)
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan"})
		return
	}
	body.DosenId = user.UserId
	log.Println(body)
	newClassroom := model.NewClassroom(&body)
	if err := model.StoreClassroom(newClassroom); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "kelas berhasil dibuat"})
}
func ReadClassrooms(c *gin.Context) {
	var user model.User
	val, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan"})
		return
	}
	loginUser := val.(model.Me)
	log.Println(user)
	if err := model.DB.Preload("DosenClassrooms.Dosen").Preload("MahasiswaClassrooms.Dosen").Where("id = ?", loginUser.UserId).Preload("DosenClassrooms").Preload("MahasiswaClassrooms").First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan"})
		return
	}
	log.Println(user)
	switch loginUser.Role {
	case "DOSEN":
		c.JSON(http.StatusOK, gin.H{"data": user.DosenClassrooms})
		return
	case "MAHASISWA":
		c.JSON(http.StatusOK, gin.H{"data": user.MahasiswaClassrooms})
		return
	}
}

func ReadDetailClassroom(c *gin.Context) {
	id := c.Param("id")
	classroom, err := model.GetDetailClassroom(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": classroom})
}

func DestroyClassroom(c *gin.Context) {
	id := c.Param("id")
	val, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan"})
		return
	}
	user := val.(model.Me)
	if err := model.DeleteClassroom(id, user.UserId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "kelas tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "kelas berhasil dihapus"})
}

func UpdateClassroom(c *gin.Context) {

}

func JoinClassroom(c *gin.Context) {
	var body model.JoinClassroomInput
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		msg := lib.GetValidationMessage(err.(validator.ValidationErrors))
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	val, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan"})
		return
	}
	me := val.(model.Me)

	user, err := model.GetUserById(me.UserId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user tidak ditemukan"})
			return
		}
		return
	}

	classroom, err := model.GetDetailClassroomByClassCode(body.Code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "kelas tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan"})
		return
	}

	if err := model.EnrollUser(user, classroom); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "berhasil gabung ke kelas"})
}

func CreateAnnouncement(c *gin.Context) {
	var body model.AnnouncementInput
	if err := c.ShouldBindJSON(&body); err != nil {
		msg := lib.GetValidationMessage(err.(validator.ValidationErrors))
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	val, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan"})
		return
	}
	me := val.(model.Me)

	classroomId := c.Param("id")
	classroom, err := model.GetDetailClassroom(classroomId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "kelas tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan"})
		return
	}

	newAnnouncement := model.NewAnnouncement(&body, me.UserId, classroom.ID)
	if err := model.StoreAnnouncement(newAnnouncement); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat pengumuman"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "pengumuman berhasil dibuat"})
}

func DeleteAnnouncement(c *gin.Context) {
	var body model.AnnouncementDeleteInput
	if err := c.ShouldBindJSON(&body); err != nil {
		msg := lib.GetValidationMessage(err.(validator.ValidationErrors))
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	val, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan"})
		return
	}
	classroomId := c.Param("id")
	user := val.(model.Me)
	if err := model.DestroyAnnouncement(body.Id, classroomId, user.UserId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "pengumuman tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "pengumuman berhasil dihapus"})
}
