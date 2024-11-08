package middleware

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	ID           string   `json:"id"`
	Email        string   `json:"email"`
	Organization string   `json:"organization"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	Permissions  []string `json:"permissions"`
	jwt.StandardClaims
}

func GenerateToken(data Claims) (string, error) {
	expirationTime := time.Now().Add(48 * time.Hour)
	claims := &Claims{
		ID:           data.ID,
		Email:        data.Email,
		Organization: data.Organization,
		FirstName:    data.FirstName,
		LastName:     data.LastName,
		Permissions:  data.Permissions,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	return claims, nil
}
