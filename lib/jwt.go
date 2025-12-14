package lib

import (
	"fmt"
	"time"

	"github.com/MhmdEagel/lms-usti-be/env"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
func CreateToken(fullname, email, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"fullname": fullname,
			"email":    email,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
			"role":     role,
		})
	tokenString, err := token.SignedString([]byte(env.SECRET_KEY))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func VerifyToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(env.SECRET_KEY), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims := token.Claims.(*JWTClaims)
	return claims, nil
}



