package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtPayload struct {
	Id string
	// Username string
	// Role     string
}

type ManagerIntf interface {
	GenerateAuthToken(id uuid.UUID) (tokenString string, err error)
	VerifyAuthToken(tokenString, jwtKey string) (payload *JwtPayload, err error)
	GetStringClaimFromJWT(ctx context.Context, claim string) (strVal string, err error)
}

type JwtManager struct {
	jwtKey      string
	expTimeHour int
}

func NewJwtManager(jwtKey string, expTimeHour int) ManagerIntf {
	return &JwtManager{
		jwtKey:      jwtKey,
		expTimeHour: expTimeHour,
	}
}

var (
	ErrJWTInit = errors.New("failed to init jwt key")
)

func (m *JwtManager) GenerateAuthToken(id uuid.UUID) (tokenString string, err error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  id.String(),
			"exp": time.Now().Add(time.Duration(m.expTimeHour) * time.Hour).Unix(),
		})

	tokenString, err = token.SignedString([]byte(m.jwtKey))

	if err != nil {
		return "", ErrJWTInit
	}

	return tokenString, nil
}

func (m *JwtManager) VerifyAuthToken(tokenString, jwtKey string) (payload *JwtPayload, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("парсинг токена: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("токен невалидный")
	}

	payload = new(JwtPayload)
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		payload.Id = fmt.Sprintf("%v", claims["id"])
	}

	return payload, nil
}

func (m *JwtManager) GetStringClaimFromJWT(ctx context.Context, claim string) (strVal string, err error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return "", err
	}

	id, ok := claims[claim]
	if !ok {
		return "", fmt.Errorf("failed getting claim '%s' from JWT token", claim)
	}

	strVal, ok = id.(string)
	if !ok {
		return "", fmt.Errorf("converting interface to string")
	}

	return strVal, nil
}
