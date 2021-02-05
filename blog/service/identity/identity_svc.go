package identity

import (
	"strconv"
	"strings"

	jwt2 "github.com/dgrijalva/jwt-go"
	"github.com/shipengqi/example.v1/blog/dao"
	"github.com/shipengqi/example.v1/blog/model"
	"github.com/shipengqi/example.v1/blog/pkg/utils"

	"github.com/shipengqi/example.v1/blog/pkg/e"
	"github.com/shipengqi/example.v1/blog/pkg/jwt"
	log "github.com/shipengqi/example.v1/blog/pkg/logger"
)

type Interface interface {
	Login(user, pass string) (string, *model.UserRBAC, error)
	Authenticate(token string) (claims *jwt.Claims, err error)
	Authorize(claims *jwt.Claims, url, method string) error
}

type identity struct {
	jwt jwt.Interface
	dao dao.Interface
}

func New(signingKey string, d dao.Interface) Interface {
	return &identity{jwt: jwt.New(signingKey), dao: d}
}

func (i *identity) Login(user, pass string) (string, *model.UserRBAC, error) {
	info, err := i.dao.GetUser(user)
	if err != nil {
		return "", nil, e.Wrap(err, e.ErrInternalServer.Message())
	}
	if len(info.Username) == 0 {
		return "", nil, e.ErrNothingFound
	}
	if info.Deleted {
		return "", nil, e.ErrUserDeleted
	}
	if info.Locked {
		return "", nil, e.ErrUserLocked
	}
	if info.Password != utils.EncodeMD5WithSalt(pass) {
		return "", nil, e.ErrPassWrong
	}
	token, err := i.jwt.GenerateToken(user, pass, int(info.ID))
	if err != nil {
		return "", nil, err
	}

	rbac, err := i.dao.GetUserRbac(info.ID)
	if err != nil {
		return "", nil, err
	}

	rbac.U = model.User{
		Username: info.Username,
		Phone:    info.Phone,
		Email:    info.Email,
	}
	return token, rbac, nil
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

func (i *identity) Authorize(claims *jwt.Claims, url, method string) error {
	userId, err := utils.DecodeXOR(claims.UID, string(i.jwt.GetSigningKey()))
	if err != nil {
		return e.ErrTokenInvalid
	}
	id, _ := strconv.Atoi(userId)
	rbac, err := i.dao.GetUserRbac(uint(id))
	if err != nil {
		return err
	}

	permissions, err := i.dao.GetPermissionsWithRoles(rbac.Roles)
	if err != nil {
		return err
	}
	denied := true
	for k := range permissions {
		if strings.HasPrefix(url, permissions[k].URL) && permissions[k].Method == method {
			denied = false
			break
		}
	}
	if denied {
		return e.ErrForbidden
	}
	return nil
}