package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey []byte

func Init(secretKeyIn string) {
	secretKey = []byte(secretKeyIn)
}

func CreateJwt(mapClaims map[string]any) (string, error) {
	if _, ok := mapClaims["exp"]; !ok {
		return "", errors.New("exp not found in map claims")
	}

	jwtMapClaims := jwt.MapClaims{}
	for key := range mapClaims {
		jwtMapClaims[key] = mapClaims[key]
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtMapClaims)

	tokenString, err := token.SignedString(secretKey)

	return tokenString, err
}

func VerifyJwt(tokenString string) (map[string]any, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// check validation errorss
	if err != nil {
		return nil, err
	}

	// check if token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// extract claims
	jwtMapClaims := token.Claims.(jwt.MapClaims)
	mapClaims := make(map[string]any)

	for key := range jwtMapClaims {
		mapClaims[key] = jwtMapClaims[key]
	}

	return mapClaims, nil

}
