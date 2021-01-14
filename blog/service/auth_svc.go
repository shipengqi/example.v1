package service

import "github.com/shipengqi/example.v1/blog/pkg/errno"

func (s *Service) Login(user, pass string) (string, error) {
	return s.jwt.GenerateToken(user, pass)
}

func (s *Service) Authenticate(token string) error {
	if token == "" {
		return errno.ErrUnauthorized
	}

	claims, err := s.jwt.ParseToken(token)
	if err != nil {
		return errno.Wrapf(err, "invalid token")
	}
	claims.
	return nil
}
