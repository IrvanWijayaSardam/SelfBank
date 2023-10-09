package service

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService interface {
	GenerateToken(UserID string, Email string, Jk string, Telephone string, Name string, IdRole uint64) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
	secretKey string
	issuer    string
}

// NewJWTService creates a new instance of JWTService
func NewJWTService() JWTService {
	return &jwtService{
		issuer:    "aminivan",
		secretKey: getSecretKey(),
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	return secretKey
}

func (j *jwtService) GenerateToken(UserID string, Email string, Jk string, Telephone string, Name string, IdRole uint64) (string, error) {
	claims := jwt.MapClaims{
		"userid": UserID,
		"name":   Name,
		"email":  Email,
		"telp":   Telephone,
		"jk":     Jk,
		"idrole": IdRole,
		"iss":    j.issuer,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().AddDate(0, 0, 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}
	return t, nil
}

func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	return parsedToken, nil
}
