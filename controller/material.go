package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/MhmdEagel/lms-usti-be/env"
	"github.com/MhmdEagel/lms-usti-be/lib"
	"github.com/MhmdEagel/lms-usti-be/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateMaterial(c *gin.Context) {
	var body model.MaterialInput
	classroomId := c.Param("id")
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	classroom, err := model.GetDetailClassroom(classroomId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "kelas tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "terjadi kesalahan"})
	}

	newMaterial := model.NewMaterial(body.MaterialName, body.MaterialDescription, classroom.ID)
	if err := model.StoreMaterial(newMaterial); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "gagal menyimpan materi"})
		return
	}

	for _, link := range body.Links {
		if !lib.IsUrl(link.LinkUrl) {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid link"})
			return
		}
		newMaterialLink := model.NewMaterialLink(link.LinkName, link.LinkUrl, newMaterial.ID)
		if err := model.StoreMaterialLink(newMaterialLink); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "gagal menyimpan materi"})
			return
		}
	}


	os.Mkdir(fmt.Sprintf("./storage/%s", classroom.ID), 0755)
	for _, file := range body.Files {
		fileType := lib.DetectFileType(file.Filename)
		file.Filename = filepath.Base(file.Filename)
		if !lib.IsAllowedFileType(file.Filename, fileType) {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid file type"})
			return
		}
		uniqueFileName := lib.GenerateUniqueFilename(file.Filename)
		filePath := fmt.Sprintf("%s/classroom/%s/serve/file/%s", env.BASE_URL, classroom.ID, uniqueFileName)
		newMaterialFile := model.NewMaterialFile(uniqueFileName, file.Filename, filePath, newMaterial.ID)
		if err := model.StoreMaterialFile(newMaterialFile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan file ke database"})
			return
		}
		if err := c.SaveUploadedFile(file, fmt.Sprintf("./storage/%s/%s", classroom.ID, uniqueFileName)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan file ke server"})
			return
		}
}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "materi berhasil dibuat"})

}
