package utils

import (
	"errors"
	"time"

	"weave/config"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 对密码进行哈希处理
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash 验证密码哈希
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken 生成JWT访问令牌（包含tenant_id）
func GenerateToken(userID uint, tenantID uint) (string, error) {
	// 创建token
	claims := jwt.MapClaims{
		"user_id":   userID,
		"tenant_id": tenantID,
		"type":      "access",
		"exp":       time.Now().Add(time.Minute * time.Duration(config.Config.JWT.AccessTokenExpiry)).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并获取完整的编码后的字符串token
	tokenString, err := token.SignedString([]byte(config.Config.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRefreshToken 生成JWT刷新令牌（包含tenant_id）
func GenerateRefreshToken(userID uint, tenantID uint) (string, error) {
	// 创建刷新令牌
	claims := jwt.MapClaims{
		"user_id":   userID,
		"tenant_id": tenantID,
		"type":      "refresh",
		"exp":       time.Now().Add(time.Hour * time.Duration(config.Config.JWT.RefreshTokenExpiry)).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并获取完整的编码后的字符串token
	tokenString, err := token.SignedString([]byte(config.Config.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyToken 验证JWT令牌，返回userID、token类型与tenantID
func VerifyToken(tokenString string) (uint, string, uint, error) {
	// 解析token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.Config.JWT.Secret), nil
	})

	if err != nil {
		return 0, "", 0, err
	}

	// 提取claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, "", 0, errors.New("invalid token")
	}

	// 提取userID
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", 0, errors.New("invalid user_id in token")
	}

	// 提取tenantID（可选，默认为0）
	var tenantID uint
	if tid, hasTid := claims["tenant_id"].(float64); hasTid {
		tenantID = uint(tid)
	} else {
		tenantID = 0
	}

	// 提取令牌类型
	tokenType, ok := claims["type"].(string)
	if !ok {
		tokenType = "access" // 默认类型
	}

	return uint(userIDFloat), tokenType, tenantID, nil
}

// VerifyRefreshToken 验证JWT刷新令牌，返回userID与tenantID
func VerifyRefreshToken(tokenString string) (uint, uint, error) {
	userID, tokenType, tenantID, err := VerifyToken(tokenString)
	if err != nil {
		return 0, 0, err
	}

	// 验证是否为刷新令牌
	if tokenType != "refresh" {
		return 0, 0, errors.New("not a refresh token")
	}

	return userID, tenantID, nil
}
