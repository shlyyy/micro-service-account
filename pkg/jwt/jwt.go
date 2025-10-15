package jwtutil

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shlyyy/micro-service-account/pkg/config"
	"github.com/shlyyy/micro-service-account/pkg/logger"
)

// 业务错误常量
var (
	ErrTokenExpired     = errors.New("token has expired")
	ErrTokenNotValidYet = errors.New("token is not yet valid")
	ErrTokenMalformed   = errors.New("that's not even a token")
	ErrTokenInvalid     = errors.New("invalid token") // 涵盖签名无效等其他所有无效情况
)

type CustomClaims struct {
	jwt.RegisteredClaims
	AccountId   int32  `json:"account_id"`
	Nickname    string `json:"nickname"`
	AuthorityId int32  `json:"authority_id"`
}

func GenerateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.Cfg.JWT.Secret))
	if err != nil {
		logger.Error("Generate JWT Error!")
	}
	return tokenStr, err
}

func ParseToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Cfg.JWT.Secret), nil
		})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			// Token 过期
			return nil, ErrTokenExpired
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			// Token 尚未生效
			return nil, ErrTokenNotValidYet
		case errors.Is(err, jwt.ErrTokenMalformed):
			// Token 格式错误
			return nil, ErrTokenMalformed
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			// 签名无效
			return nil, ErrTokenInvalid
		default:
			// 其他所有无法识别的错误，都归为无效 Token
			return nil, ErrTokenInvalid
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
	}
	// 如果 token.Valid 为 false，或者类型断言失败，也视为无效 Token
	return nil, ErrTokenInvalid
}

func RefreshToken(refreshTokenString string) (string, error) {
	claims, err := ParseToken(refreshTokenString)
	if err != nil {
		return "", err
	}

	// 更新过期时间
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour))
	claims.IssuedAt = jwt.NewNumericDate(time.Now())

	return GenerateToken(*claims)
}
