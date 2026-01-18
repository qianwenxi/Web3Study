package utils

import (
	"time"

	"goStudyHm4/config"

	"github.com/dgrijalva/jwt-go"
)

// Claims JWT负载
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, username string) (string, error) {
	// 设置过期时间
	expireTime := time.Now().Add(time.Hour * time.Duration(config.TokenExpireHour))
	// 构建负载
	claims := Claims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), // 过期时间
			IssuedAt:  time.Now().Unix(), // 签发时间
			Issuer:    "blog-backend",    // 签发者
		},
	}
	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 签名生成令牌字符串
	tokenString, err := token.SignedString(config.JWTSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*Claims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return config.JWTSecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 验证令牌并转换claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}
