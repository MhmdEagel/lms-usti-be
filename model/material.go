package model

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Material struct {
	ID            string    `gorm:"primaryKey"`
	Title         string    `json:"title" gorm:"not null"`
	Description   string    `json:"description,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ClassroomId   string
	MaterialFiles []MaterialFile `json:"files" gorm:"foreignKey:MaterialId;constraint:OnDelete:CASCADE;"`
	MaterialLinks []MaterialLink `json:"links" gorm:"foreignKey:MaterialId;constraint:OnDelete:CASCADE;"`
}

type MaterialInput struct {
	MaterialName        string                  `form:"material_name" binding:"required"`
	MaterialDescription string                  `form:"description"`
	Files               []*multipart.FileHeader `form:"files"`
	Links               []LinkInput             `form:"links"`
}

type LinkInput struct {
	LinkName string `json:"link_name" form:"link_name" binding:"required"`
	LinkUrl  string `json:"link_url" form:"link_url" binding:"required"`
}

func NewMaterial(title, description, classroomId string) *Material {
	return &Material{Title: title, Description: description, ClassroomId: classroomId}
}

func (material *Material) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	material.ID = id.String()
	return err
}

func StoreMaterial(material *Material) error {
	return DB.Create(material).Error
}

type MaterialFile struct {
	ID               string `gorm:"primaryKey"`
	UniqueFileName   string `json:"fileName" gorm:"not null"`
	OriginalFileName string `json:"originalFileName" gorm:"not null"`
	FileUrl          string `json:"fileUrl" gorm:"not null"`
	MaterialId       string `json:"-"`
}

func (materialFile *MaterialFile) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	materialFile.ID = id.String()
	return err
}

func NewMaterialFile(uniqueFileName, originalFileName, fileUrl, materialId string) *MaterialFile {
	return &MaterialFile{UniqueFileName: uniqueFileName, OriginalFileName: originalFileName, FileUrl: fileUrl, MaterialId: materialId}
}

func StoreMaterialFile(materialFile *MaterialFile) error {
	return DB.Create(materialFile).Error
}

type MaterialLink struct {
	ID         string `gorm:"primaryKey"`
	LinkName   string `json:"linkName" gorm:"not null"`
	LinkUrl    string `json:"linkUrl" gorm:"not null"`
	MaterialId string `json:"-"`
}

func (materialLink *MaterialLink) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	materialLink.ID = id.String()
	return err
}

func NewMaterialLink(linkName, linkUrl, materialId string) *MaterialLink {
	return &MaterialLink{LinkName: linkName, LinkUrl: linkUrl, MaterialId: materialId}
}

func StoreMaterialLink(materialLink *MaterialLink) error {
	return DB.Create(materialLink).Error
}
