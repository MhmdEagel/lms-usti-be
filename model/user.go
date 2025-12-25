package model

import (
	"database/sql"

	"github.com/MhmdEagel/lms-usti-be/lib"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                     string         `json:"-" gorm:"primary_key;not null"`
	Fullname               string         `json:"fullname" gorm:"not null"`
	Email                  string         `json:"email" gorm:"unique;not null"`
	EmailVerified          sql.NullTime   `json:"-"`
	Image                  string         `json:",omitempty"`
	Password               string         `json:"-" gorm:"not null"`
	Role                   string         `json:"role" gorm:"not null"`
	DosenClassrooms        []Classroom    `json:"dosenClassrooms,omitempty" gorm:"foreignKey:DosenId;"`
	MahasiswaClassrooms    []Classroom    `json:"mahasiswaClassrooms,omitempty" gorm:"many2many:classroom_mahasiswas;constraint:OnDelete:CASCADE;"`
	ClassroomAnnouncements []Announcement `json:"dosenAnnouncements,omitempty" gorm:"foreignKey:DosenId;constraint:OnDelete:CASCADE;"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	user.ID = id.String()
	return err
}

type Register struct {
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"required,oneof=MAHASISWA DOSEN"`
}
type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type ResendVerificationInput struct {
	Email string `json:"email" binding:"required,email"`
}

type JoinClassroomInput struct {
	Code string `json:"code" binding:"required"`
}

type Me struct {
	UserId string
	Email  string
	Role   string
}

func NewUser(register *Register) *User {
	return &User{Fullname: register.Fullname, Email: register.Email, Password: lib.HashPassword(register.Password), Role: register.Role}
}
func CreateUser(user *User) error {
	result := DB.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func GetUserByEmail(email string) (*User, error) {
	var user User
	result := DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetUserById(id string) (*User, error) {
	var user User
	result := DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
