package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"kama_chat_server/internal/config"
	"kama_chat_server/pkg/zlog"
)

// JWTConfig JWT 配置
type JWTConfig struct {
	SecretKey     string
	ExpireTime    time.Duration
	RefreshExpire time.Duration
}

// Claims JWT 声明
type Claims struct {
	UUID     string `json:"uuid"`
	Account  string `json:"account"`
	Nickname string `json:"nickname"`
	IsAdmin  int8   `json:"is_admin"`
	jwt.RegisteredClaims
}

var jwtConfig *JWTConfig

func init() {
	// 从配置文件自动初始化 JWT
	conf := config.GetConfig()
	jwtConfig = &JWTConfig{
		SecretKey:     conf.JWTConfig.SecretKey,
		ExpireTime:    time.Duration(conf.JWTConfig.ExpireHours) * time.Hour,
		RefreshExpire: time.Duration(conf.JWTConfig.RefreshHours) * time.Hour,
	}
	zlog.Info("JWT 配置已初始化")
}

// GenerateToken 生成 JWT token
func GenerateToken(uuid, account, nickname string, isAdmin int8) (string, error) {
	if jwtConfig == nil {
		return "", errors.New("JWT config not initialized")
	}

	claims := Claims{
		UUID:     uuid,
		Account:  account,
		Nickname: nickname,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtConfig.ExpireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtConfig.SecretKey))
}

// ParseToken 解析 JWT token
func ParseToken(tokenString string) (*Claims, error) {
	if jwtConfig == nil {
		return nil, errors.New("JWT config not initialized")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtConfig.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken 刷新 token
func RefreshToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	// 生成新的 token
	return GenerateToken(claims.UUID, claims.Account, claims.Nickname, claims.IsAdmin)
}

