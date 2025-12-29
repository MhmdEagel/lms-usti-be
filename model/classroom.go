package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Classroom struct {
	ID                    string         `json:"id" gorm:"unique;primaryKey;not null"`
	ClassCode             string         `json:"class_code" gorm:"unique;not null"`
	ClassName             string         `json:"class_name" gorm:"not null"`
	Term                  int            `json:"term" gorm:"not null"`
	RoomNumber            int            `json:"room_number" gorm:"not null"`
	Day                   int            `json:"day" gorm:"not null"`
	ClassStart            time.Time      `json:"class_start" gorm:"not null"`
	ClassEnd              time.Time      `json:"class_end" gorm:"not null"`
	DosenId               string         `json:"-"`
	Dosen                 User           `json:"dosen" gorm:"foreignKey:DosenId"`
	ClassroomMahasiswa    []*User        `json:"mahasiswa,omitempty" gorm:"many2many:classroom_mahasiswas;constraint:OnDelete:CASCADE;"`
	ClassroomAnnouncement []Announcement `json:"announcements,omitempty" gorm:"foreignKey:ClassroomId;constraint:OnDelete:CASCADE;"`
	Materials             []Material     `gorm:"foreignKey:ClassroomId"`
}
type ClassroomMahasiswa struct {
	UserId      string    `gorm:"primaryKey"`
	ClassroomId string    `gorm:"primaryKey"`
	User        User      `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	Classroom   Classroom `gorm:"foreignKey:ClassroomID;constraint:OnDelete:CASCADE"`
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
	DosenId    string
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

func EnrollUser(user *User, classroom *Classroom) error {
	model := ClassroomMahasiswa{
		UserId:      user.ID,
		ClassroomId: classroom.ID,
	}
	return DB.Create(model).Error
}

func DeleteClassroom(id, userId string) error {
	if query := DB.Where("id = ?", id).Where("dosen_id = ?", userId).Delete(&Classroom{}); query.Error != nil {
		return query.Error
	}
	return nil
}

func GetDetailClassroom(id string) (*Classroom, error) {
	var classroom Classroom
	res := DB.Preload("Dosen").First(&classroom, "id = ?", id)
	return &classroom, res.Error
}

func GetDetailClassroomByClassCode(classCode string) (*Classroom, error) {
	var classroom Classroom
	res := DB.First(&classroom, "class_code = ?", classCode)
	if res.Error != nil {
		return nil, res.Error
	}
	return &classroom, nil
}
