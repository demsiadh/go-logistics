package util

import (
	"github.com/golang-jwt/jwt/v5"
	"go_logistics/config"
	"time"
)

type CustomClaims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

// 生成 Token
func createToken(name string, secret []byte) (string, error) {
	claims := CustomClaims{
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // 签发时间
			Issuer:    "logistics",                                        // 签发者
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// 解析 Token
func parseToken(tokenString string, secret []byte) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

// GenerateToken 生成 Token
func GenerateToken(name string) (token string, err error) {
	token, err = createToken(name, []byte(config.SecretKey))
	return
}

// CheckToken 检查 Token
func CheckToken(token string) (claims *CustomClaims, err error) {
	claims, err = parseToken(token, []byte(config.SecretKey))
	return
}
