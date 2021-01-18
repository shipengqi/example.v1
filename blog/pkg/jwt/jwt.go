package jwt

import (
	"time"

	jwt2 "github.com/dgrijalva/jwt-go"
	"github.com/shipengqi/example.v1/blog/pkg/utils"
)

type Interface interface {
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (*Claims, error)
	RefreshToken(token string) (string, error)
}

type jwt struct {
	SigningKey []byte
}

func New(signingKey string) Interface {
	return &jwt{SigningKey: []byte(signingKey)}
}

// Custom jwt claims
type Claims struct {
	jwt2.StandardClaims

	Username string `json:"username"`
	Password string `json:"password"`
}

// GenerateToken generate tokens used for auth
func (j *jwt) GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Hour)

	claims := &Claims{
		jwt2.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "example.v1",
		},
		utils.EncodeMD5(username),
		utils.EncodeMD5(password),
	}
	return j.gen(claims)
}

// ParseToken parsing token
func (j *jwt) ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt2.ParseWithClaims(
		token,
		&Claims{},
		func(token *jwt2.Token) (interface{}, error) {
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
func (j *jwt) RefreshToken(token string) (string, error) {
	claims, err := j.ParseToken(token)
	if err != nil {
		return "", err
	}
	claims.StandardClaims.ExpiresAt = time.Now().Add(time.Hour).Unix()
	return j.gen(claims)
}

func (j *jwt) gen(claims *Claims) (string, error) {
	tokenClaims := jwt2.NewWithClaims(jwt2.SigningMethodHS256, claims)
	return tokenClaims.SignedString(j.SigningKey)
}
