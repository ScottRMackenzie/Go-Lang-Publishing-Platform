package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJWT(ID, username string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}
	var jwtKey = []byte(os.Getenv("JWT_SECRET"))

	claims := &Claims{
		UserID:   ID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func VerifyJWT(token string) (*Claims, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	var jwtKey = []byte(os.Getenv("JWT_SECRET"))

	tkn, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tkn.Claims.(*Claims)
	if !ok || !tkn.Valid {
		return nil, err
	}
	return claims, nil
}
