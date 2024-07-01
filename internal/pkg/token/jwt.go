package token

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var secretKey = []byte("secretKey")

func createToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp":     jwt.NumericDate{Time: time.Now().Add(2 * time.Hour)},
			"user_id": id,
		})
	return token.SignedString(secretKey)
}

func verifyToken(jwtStr string) bool {
	token, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return false
	}
	return token.Valid
}
