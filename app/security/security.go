package security

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-errors/errors"
	"theAmazingApiGateway/app/config"
)

type JWTCustomClaims struct {
	Id       uint
	Email    string
	Name     string
	LastName string
	jwt.StandardClaims
}

type JWTToken struct {
	Id       uint
	Name     string
	LastName string
	Email    string
}

func GetTokenData(tokenString string) (JWTToken, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().JWT_SECRET), nil
	})
	if err != nil {
		return JWTToken{}, err
	}

	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return JWTToken{
			Id:       claims.Id,
			Name:     claims.Name,
			LastName: claims.LastName,
			Email:    claims.Email,
		}, nil
	} else {
		return JWTToken{}, errors.New("Invalid token")
	}
}
