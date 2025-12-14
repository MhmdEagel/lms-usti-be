package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)
type VerificationToken struct {
	ID      string    `gorm:"primary_key;not null"`
	Email   string    `json:"email" gorm:"unique;not null"`
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}


type Verification struct {
	Token string `json:"token" binding:"required"`
}


func (verification *VerificationToken) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	verification.ID = id.String()
	return err
}
func NewVerificationToken(email string) *VerificationToken {
	id := uuid.New()
	token := id.String()
	expires := time.Now().Add(10 * time.Minute)
	return &VerificationToken{Email: email, Token: token, Expires: expires}
}
func CreateVerificationToken(token *VerificationToken) error {	
	_, err := GetVerificationTokenByEmail(token.Email)
	if err != nil {
		return err
	}
	DB.Create(&token)
	return nil
}
func GetVerificationTokenByEmail(email string) (*VerificationToken, error) {
	var token VerificationToken
	res := DB.Where("email = ?", email).First(&token)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("email tidak ditemukan")
	}
	return &token, nil
}
func GetVerificationTokenByToken(token string) (*VerificationToken, error) {
	var t VerificationToken
	res := DB.Where("token = ?", token).First(&t)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("email tidak ditemukan")
	}
	return &t, nil
}
