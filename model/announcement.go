package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Announcement struct {
	ID          string `json:"id" gorm:"primary_key;not null"`
	Title       string `json:"title" gorm:"not null;"`
	Content     string `json:"content" gorm:"not null;type:text"`
	CreatedAt   time.Time
	ClassroomId string
	DosenId     string
}

func (announcement *Announcement) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	announcement.ID = id.String()
	return err
}

type AnnouncementInput struct {
	Title   string `binding:"required"`
	Content string `binding:"required"`
}

type AnnouncementDeleteInput struct {
	Id string `binding:"required"`
}

func NewAnnouncement(announce *AnnouncementInput, userId, classroomId string) *Announcement {
	return &Announcement{Title: announce.Title, Content: announce.Content, DosenId: userId, ClassroomId: classroomId}
}

func StoreAnnouncement(announce *Announcement) error {
	res := DB.Create(&announce)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func DestroyAnnouncement(announcementId, classroomId, userId string) error {
	if query := DB.Where("id = ?", announcementId).Where("classroom_id = ?", classroomId).Where("dosen_id = ?", userId).Delete(&Announcement{}); query.Error != nil {
		return query.Error
	}
	return nil
}
