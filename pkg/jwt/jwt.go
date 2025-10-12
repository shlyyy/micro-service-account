package jwt

import (
	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	jwt.StandardClaims
	Id          int32
	NikeName    string
	AuthorityId int32
}

var jwtSecret []byte

func InitJWT(cfg JWTConfig) {
	jwtSecret = []byte(cfg.SingingKey)
}

func GenerateJWT(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
