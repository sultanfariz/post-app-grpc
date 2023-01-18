package commons

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTClaims struct {
	Id                 int
	jwt.StandardClaims
}

type ConfigJWT struct {
	SecretJWT       string
	ExpiresDuration int
}

// func (jwtConf *ConfigJWT) Init() middleware.JWTConfig {
// 	return middleware.JWTConfig{
// 		Claims:     &JWTClaims{},
// 		SigningKey: []byte(jwtConf.SecretJWT),
// 	}
// }

func (jwtConf *ConfigJWT) GenerateToken(Id int, Email string) (string, error) {
	claims := JWTClaims{
		Id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(jwtConf.ExpiresDuration)).Unix(),
			IssuedAt:  time.Now().Local().Unix(),
			Issuer:    "auth.service",
			Subject:   Email,
		},
	}

	// Create token with claims
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(jwtConf.SecretJWT))

	return token, err
}

// func GetUser(ctx context.Context) *JWTMyClaims {
// 	user := c.Get("user").(*jwt.Token)
// 	claims := user.Claims.(*JWTMyClaims)
// 	return claims
// }
