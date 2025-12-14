package lib

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func GetValidationMessage(err validator.ValidationErrors) string {
	if len(err) == 0 {
		return ""
	}
	for _, v := range err {
		switch v.Tag() {
		case "required":
			return fmt.Sprintf("The %s field is required.", strings.ToLower(v.Field()))
		case "min":
			return fmt.Sprintf("The %s field must be at least %s characters long.", strings.ToLower(v.Field()), v.Param())
		case "max":
			return fmt.Sprintf("The %s field must be at most %s characters long.", strings.ToLower(v.Field()), v.Param())
		case "email":
			return fmt.Sprintf("The %s field must be a valid email address.", strings.ToLower(v.Field()))
		case "oneof":
			return fmt.Sprintf("The %s field is not a valid role.", strings.ToLower(v.Field()))
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
