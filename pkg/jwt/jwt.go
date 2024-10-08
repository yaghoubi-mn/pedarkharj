package jwt

import (
	"errors"
	"fmt"
	"strconv"
	"time"

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

func VerifyJwt(tokenString string) (mapClaims map[string]any, err error) {
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
	mapClaims = make(map[string]any)

	for key := range jwtMapClaims {
		mapClaims[key] = jwtMapClaims[key]
	}

	return mapClaims, nil

}

func CreateRefreshAndAccessFromUser(refreshExpireTime time.Duration, accessExpireTime time.Duration, id uint64, name string, number string, isRegistered bool) (refresh string, access string, err error) {
	refresh, err = CreateJwt(map[string]any{
		"exp": time.Now().Add(refreshExpireTime).Unix(),
	})

	if err != nil {
		return "", "", err
	}

	access, err = CreateAccessFromUser(accessExpireTime, id, name, number, isRegistered)

	return refresh, access, err

}

func CreateAccessFromUser(accessExpireTime time.Duration, id uint64, name string, number string, isRegistered bool) (access string, err error) {

	if id == 0 {
		return "", errors.New("cannot create jwt: id is zero")
	}
	fmt.Println(id, "----------------")
	access, err = CreateJwt(map[string]any{
		"exp":          time.Now().Add(accessExpireTime).Unix(),
		"id":           id,
		"name":         name,
		"number":       number,
		"isRegistered": isRegistered,
	})

	return access, err

}

func GetUserFromAccess(access string) (id uint64, name string, number string, isRegistered bool, err error) {

	mapClaims, err := VerifyJwt(access)
	if err != nil {
		return 0, "", "", false, err
	}
	fmt.Println(mapClaims, "====================")
	id = uint64(mapClaims["id"].(float64))
	name = mapClaims["name"].(string)
	number = mapClaims["number"].(string)
	isRegistered = mapClaims["isRegistered"].(bool)

	return id, name, number, isRegistered, nil
}

func CreateRefreshAndAccessFromUserWithMap(refreshExpireMinutes time.Duration, accessExpireMinutes time.Duration, id uint64, name string, number string, isRegistered bool) (tokens map[string]string, err error) {
	tokens = make(map[string]string)

	refresh, access, err := CreateRefreshAndAccessFromUser(refreshExpireMinutes, accessExpireMinutes, id, name, number, isRegistered)
	tokens["refresh"] = refresh
	tokens["access"] = access
	tokens["accessExpireSeconds"] = strconv.Itoa(int(accessExpireMinutes.Seconds()))

	return
}
