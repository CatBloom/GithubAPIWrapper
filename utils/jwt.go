package utils

import (
	"crypto/rand"
	"encoding/base64"
	"os"

	"github.com/dgrijalva/jwt-go"
)

// JWTに変換
func GenerateJWT(key string, str string) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims[key] = str

	jwtStr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return jwtStr, nil
}

// JWTを解析
func RefreshJWT(key string, jwtStr string) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")

	parsedToken, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}
	claims := parsedToken.Claims.(jwt.MapClaims)

	str := claims[key].(string)

	return str, nil
}

// ランダムなシークレットキーを作成
func GenerateRandomSecretKey() (string, error) {
	keySize := 32

	key := make([]byte, keySize)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(key), nil
}
