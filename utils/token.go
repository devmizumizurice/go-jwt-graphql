package utils

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userId string, isRefresh bool) (*string, error) {
	var exp int64
	var secret []byte

	if !isRefresh {
		accessTokenLifeTime, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_VALIDATE_MINUTES"))
		if err != nil {
			return nil, err
		}
		exp = time.Now().Add(time.Minute * time.Duration(accessTokenLifeTime)).Unix()
		secret = []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
	} else {
		refreshTokenLifeTime, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_VALIDATE_DAYS"))
		if err != nil {
			return nil, err
		}
		exp = time.Now().Add(time.Hour * 24 * time.Duration(refreshTokenLifeTime)).Unix()
		secret = []byte(os.Getenv("REFRESH_TOKEN_SECRET"))

	}

	claims := jwt.MapClaims{
		"sub": userId,
		"exp": exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)

	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func VerifyToken(tokenString string, isRefresh bool) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("UNMATHED_SIGNING_METHOD")
		}

		var secret []byte

		if !isRefresh {
			secret = []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
		} else {
			secret = []byte(os.Getenv("REFRESH_TOKEN_SECRET"))
		}

		return secret, nil
	})
}
