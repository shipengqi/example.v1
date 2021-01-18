package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/shipengqi/example.v1/blog/pkg/utils"
)

type Jwt struct {
	SigningKey []byte
}

func New(signingKey string) *Jwt {
	return &Jwt{SigningKey: []byte(signingKey)}
}

// Custom jwt claims
type Claims struct {
	jwt.StandardClaims

	Username string `json:"username"`
	Password string `json:"password"`
}

// GenerateToken generate tokens used for auth
func (j *Jwt) GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Hour)

	claims := &Claims{
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "example.v1",
		},
		utils.EncodeMD5(username),
		utils.EncodeMD5(password),
	}
	return j.gen(claims)
}

// ParseToken parsing token
func (j *Jwt) ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(
		token,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return j.SigningKey, nil
		},
	)

	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

// RefreshToken refresh token
func (j *Jwt) RefreshToken(token string) (string, error) {
	claims, err := j.ParseToken(token)
	if err != nil {
		return "", err
	}
	claims.StandardClaims.ExpiresAt = time.Now().Add(time.Hour).Unix()
	return j.gen(claims)
}

func (j *Jwt) gen(claims *Claims) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString(j.SigningKey)
}
