package auth

import (
	"github.com/dgrijalva/jwt-go"

	"github.com/shipengqi/example.v1/blog/pkg/errno"
	jwt2 "github.com/shipengqi/example.v1/blog/pkg/jwt"
	log "github.com/shipengqi/example.v1/blog/pkg/logger"
)

type Svc struct {
	jwt *jwt2.Jwt
}

func New(signingKey string) *Svc {
	return &Svc{jwt: jwt2.New(signingKey)}
}

func (s *Svc) Login(user, pass string) (string, error) {
	// find user
	// verify user credential
	// get role
	token, err := s.jwt.GenerateToken(user, pass)
	if err != nil {
		return "", errno.Wrap(err, errno.ErrGenTokenFailed.Message())
	}
	// Todo
	// return user info
	return token, err
}

func (s *Svc) Authenticate(token string) (claims *jwt2.Claims, err error) {
	if token == "" {
		return nil, errno.ErrUnauthorized
	}

	claims, err = s.jwt.ParseToken(token)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			log.Error().Err(err).Msgf("Authenticate ValidationError")
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errno.ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errno.ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errno.ErrNotValidYet
			} else {
				return nil, errno.ErrTokenInvalid
			}
		}
		return nil, errno.Wrapf(err, "invalid token")
	}
	return
}
