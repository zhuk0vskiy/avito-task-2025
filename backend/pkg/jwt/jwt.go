package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtPayload struct {
	Id       string
	// Username string
	// Role     string
}

var (
	ErrJWTInit = errors.New("failed to init jwt key")
)

func GenerateAuthToken(id uuid.UUID, jwtKey string) (tokenString string, err error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":   id.String(),
			"exp":  time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err = token.SignedString([]byte(jwtKey))

	if err != nil {
		return "", ErrJWTInit
	}

	return tokenString, nil
}

func VerifyAuthToken(tokenString, jwtKey string) (payload *JwtPayload, err error) {
	//fmt.Println(tokenString, "===")
	//tokenString = "1"
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	//fmt.Println(token, err)
	if err != nil {
		return nil, fmt.Errorf("парсинг токена: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("токен невалидный")
	}

	payload = new(JwtPayload)
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		payload.Id = fmt.Sprintf("%v", claims["id"])
		// payload.Username = fmt.Sprint(claims["sub"])
		// payload.Role = fmt.Sprint(claims["role"])
	}

	//fmt.Println(payload, "payload")

	return payload, nil
}