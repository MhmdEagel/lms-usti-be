package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Classroom struct {
	ID         string    `gorm:"primary_key;not null"`
	ClassCode  string    `json:"class_code" gorm:"unique;not null"`
	ClassName  string    `json:"class_name" gorm:"not null"`
	Term       int       `json:"term" gorm:"not null"`
	RoomNumber int       `json:"room_number" gorm:"not null"`
	Day        int       `json:"day" gorm:"not null"`
	ClassStart time.Time `json:"class_start" gorm:"not null"`
	ClassEnd   time.Time `json:"class_end" gorm:"not null"`
	DosenId    string
}

func (classroom *Classroom) BeforeCreate(tx *gorm.DB) error {
	classId, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	classCodeId, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	classCode := fmt.Sprintf("KLS-%s", strings.ToUpper(strings.ReplaceAll(classCodeId.String()[:5], "-", "")))
	classroom.ID = classId.String()
	classroom.ClassCode = classCode
	return nil
}

type CreateClassroomInput struct {
	ClassName  string    `json:"class_name" binding:"required,min=8"`
	Term       int       `json:"term" binding:"required"`
	RoomNumber int       `json:"room_number" binding:"required"`
	Day        int       `json:"day" binding:"required"`
	ClassStart time.Time `json:"class_start" binding:"required"`
	ClassEnd   time.Time `json:"class_end" binding:"required"`
	DosenId    string    `binding:"required"`
}

func NewClassroom(classroomInput *CreateClassroomInput) *Classroom {
	return &Classroom{ClassName: classroomInput.ClassName, Term: classroomInput.Term, RoomNumber: classroomInput.RoomNumber, Day: classroomInput.Day, ClassStart: classroomInput.ClassStart, ClassEnd: classroomInput.ClassEnd, DosenId: classroomInput.DosenId}
}
func StoreClassroom(classroom *Classroom) error {
	res := DB.Create(classroom)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
