package helpers

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/andrey67895/go_diplom_second/internal/logger"
)

var log = logger.Logger()
var sampleSecretKey = []byte(EncodeHashSha512("GoDiplomSecondKey"))

type JwtToken struct {
	Login *string
}

func DecodeJWT(tokenString string) (string, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", t.Header["alg"])
		}
		return sampleSecretKey, nil
	}
	var claims = jwt.RegisteredClaims{}

	parsedToken, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc)
	if err != nil {
		log.Error("Ошибка разбора: ", err)
		return "", fmt.Errorf("ошибка разбора: %s", err.Error())
	}

	if !parsedToken.Valid {
		err := fmt.Errorf("недействительный токен")
		log.Error(err)
		return "", err
	}
	return claims.Subject, nil

}

func GenerateJWTAndCheck(username string) (string, error) {
	var claims jwt.RegisteredClaims
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 30).UTC())
	claims.Subject = username
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		log.Error("Ошибка генерации токена: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}
