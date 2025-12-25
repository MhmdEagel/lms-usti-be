package lib

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetValidationMessage(err validator.ValidationErrors) string {
	if len(err) == 0 {
		return ""
	}
	for _, v := range err {
		switch v.Tag() {
		case "required":
			return fmt.Sprintf("%s harus diisi.", strings.ToLower(v.Field()))
		case "min":
			return fmt.Sprintf("%s minimal %s karakter", strings.ToLower(v.Field()), v.Param())
		case "max":
			return fmt.Sprintf("%s maksimal %s karakter.", strings.ToLower(v.Field()), v.Param())
		case "email":
			return fmt.Sprintf("%s bukan email yang valid.", strings.ToLower(v.Field()))
		case "oneof":
			return fmt.Sprintf("%s bukan pilihan yang valid.", strings.ToLower(v.Field()))
		default:
			fmt.Println(v.Tag())
		}
	}
	return ""
}

func HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 15)
	return string(hash)
}

func IsPasswordMatch(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func OmitFields(columns ...string) func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Omit(columns...)
	}
}
