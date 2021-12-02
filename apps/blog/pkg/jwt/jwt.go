package jwt

import (
	utils2 "github.com/shipengqi/example.v1/apps/blog/pkg/utils"
	"strconv"
	"time"

	jwt2 "github.com/dgrijalva/jwt-go"
)

type Interface interface {
	GenerateToken(username, password string, id int) (string, error)
	ParseToken(token string) (*Claims, error)
	RefreshToken(token string) (string, error)
	GetSigningKey() []byte
}

type jwt struct {
	signingKey []byte
}

func New(signingKey string) Interface {
	return &jwt{signingKey: []byte(signingKey)}
}

// Custom jwt claims
type Claims struct {
	jwt2.StandardClaims

	UID      string `json:"uid"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (j *jwt) GetSigningKey() []byte {
	return j.signingKey
}

// GenerateToken generate tokens used for auth
func (j *jwt) GenerateToken(username, password string, id int) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Hour)

	claims := &Claims{
		StandardClaims: jwt2.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "example.v1",
		},
		UID:      utils2.EncodeXOR(strconv.Itoa(id), string(j.signingKey)),
		Username: utils2.EncodeMD5(username),
		Password: utils2.EncodeMD5(password),
	}
	return j.gen(claims)
}

// ParseToken parsing token
func (j *jwt) ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt2.ParseWithClaims(
		token,
		&Claims{},
		func(token *jwt2.Token) (interface{}, error) {
			return j.signingKey, nil
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
	return tokenClaims.SignedString(j.signingKey)
}
