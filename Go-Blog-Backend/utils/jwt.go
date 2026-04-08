package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const TokenExpireDuration = time.Hour * 24

var CustomSecret = []byte("")

type CustomClaims struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(UserId int64, username string) (string, error) {
	//create a claim
	claims := CustomClaims{UserId, username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "Go-blog",
		},
	}

	//use HS256 algorithm to sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(CustomSecret)
}

func ParseToken(tokenString string) (*CustomClaims, error) {
	var claims CustomClaims

	//get token
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return CustomSecret, nil
	})
	if err != nil {
		return nil, err
	}

	//Verify if the token is valid
	if token.Valid {
		return &claims, nil
	}
	return nil, errors.New("invalid token")
}
