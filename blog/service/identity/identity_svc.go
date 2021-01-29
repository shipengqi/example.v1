package identity

import (
	jwt2 "github.com/dgrijalva/jwt-go"

	"github.com/shipengqi/example.v1/blog/pkg/e"
	"github.com/shipengqi/example.v1/blog/pkg/jwt"
	log "github.com/shipengqi/example.v1/blog/pkg/logger"
)

type Interface interface {
	Login(user, pass string) (string, error)
	Authenticate(token string) (claims *jwt.Claims, err error)
}

type identity struct {
	jwt jwt.Interface
}

func New(signingKey string) Interface {
	return &identity{jwt: jwt.New(signingKey)}
}

func (i *identity) Login(user, pass string) (string, error) {
	// find user
	// verify user credential
	// get role
	token, err := i.jwt.GenerateToken(user, pass)
	if err != nil {
		return "", e.Wrap(err, e.ErrGenTokenFailed.Message())
	}
	// return token and user info
	return token, err
}

func (i *identity) Authenticate(token string) (claims *jwt.Claims, err error) {
	if token == "" {
		return nil, e.ErrUnauthorized
	}

	claims, err = i.jwt.ParseToken(token)
	if err != nil {
		if ve, ok := err.(*jwt2.ValidationError); ok {
			log.Error().Err(err).Msgf("Authenticate ValidationError")
			if ve.Errors&jwt2.ValidationErrorMalformed != 0 {
				return nil, e.ErrTokenMalformed
			} else if ve.Errors&jwt2.ValidationErrorExpired != 0 {
				return nil, e.ErrTokenExpired
			} else if ve.Errors&jwt2.ValidationErrorNotValidYet != 0 {
				return nil, e.ErrNotValidYet
			} else {
				return nil, e.ErrTokenInvalid
			}
		}
		return nil, e.Wrapf(err, "invalid token")
	}
	return
}