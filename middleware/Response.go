package middleware

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func SendResponse(response chan []byte, Response map[string]interface{}) {
	res, err := json.Marshal(Response)
	if err != nil {
		fmt.Println(err)
	}

	response <- res
}

var SampleSecretKey = []byte("Ahmedfawzi made this website")

func GenerateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 8760).Unix(),
	})
	tokenString, err := token.SignedString(SampleSecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) (string, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SampleSecretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email := getEmailFromClaims(claims)
		return email, nil
	} else {
		return "", fmt.Errorf("invalid jwt token")
	}
}

func getEmailFromClaims(claims jwt.MapClaims) string {
	emailValue := claims["email"]

	email := emailValue.(string)

	return email
}
