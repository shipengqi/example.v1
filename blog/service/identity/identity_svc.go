package identity

import (
	"github.com/dgrijalva/jwt-go"

	"github.com/shipengqi/example.v1/blog/pkg/e"
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
		return "", e.Wrap(err, e.ErrGenTokenFailed.Message())
	}
	// return token and user info
	return token, err
}

func (s *Svc) Authenticate(token string) (claims *jwt2.Claims, err error) {
	if token == "" {
		return nil, e.ErrUnauthorized
	}

	claims, err = s.jwt.ParseToken(token)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			log.Error().Err(err).Msgf("Authenticate ValidationError")
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, e.ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, e.ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, e.ErrNotValidYet
			} else {
				return nil, e.ErrTokenInvalid
			}
		}
		return nil, e.Wrapf(err, "invalid token")
	}
	return
}
