package jwt

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
	// fmt.Println(token)
	tokenString, err = token.SignedString([]byte(m.jwtKey))

	if err != nil {
		return "", ErrJWTInit
	}

	return tokenString, nil
}

// func (m *JwtManager) GetStringClaimFromJWT(ctx *gin.Context, claim string) (strVal string, err error) {
// 	token, err := parseAuthHeader(ctx)
// 	if err != nil {
// 		return "", err
// 	}

// 	strVal, err = parseToken(token, claim, m.jwtKey)

// 	id, ok := claims[claim]
// 	if !ok {
// 		return "", fmt.Errorf("failed getting claim '%s' from JWT token", claim)
// 	}

// 	strVal, ok = id.(string)
// 	if !ok {
// 		return "", fmt.Errorf("converting interface to string")
// 	}

// 	return strVal, err
// }

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

	// token, err := parseAuthHeader(ctx)
	// if err != nil {
	// 	return "", err
	// }

	// strVal, err = parseToken(token, claim, m.jwtKey)

	// id, ok := claims[claim]
	// if !ok {
	// 	return "", fmt.Errorf("failed getting claim '%s' from JWT token", claim)
	// }

	// strVal, ok = id.(string)
	// if !ok {
	// 	return "", fmt.Errorf("converting interface to string")
	// }

	// return strVal, err
}

func parseAuthHeader(ctx context.Context) (string, error) {
	// header := c.GetHeader("Authorization")
	header, _ := ctx.Value("Authorization").(string)
	// fmt.Println(ok)
	if header == "" {
		return "", fmt.Errorf("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", fmt.Errorf("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", fmt.Errorf("token is empty")
	}

	return headerParts[1], nil
}

func parseToken(accessToken string, claim string, signingKey string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get claims from token")
	}

	strVal, ok := claims[claim].(string)
	if !ok {
		return "", fmt.Errorf("failed getting claim '%s' from JWT token", claim)
	}
	// fmt.Println(strVal)
	// strVal, ok := id.(string)
	// if !ok {
	// 	return "", fmt.Errorf("converting interface to string")
	// }

	exp, err := claims.GetExpirationTime()
	if err != nil {
		return "", fmt.Errorf("error getting expiration time: %w", err)
	}
	if exp.Before(time.Now()) {
		return "", fmt.Errorf("token is expired")
	}

	return strVal, nil
}
