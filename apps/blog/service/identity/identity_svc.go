package identity

import (
	"github.com/shipengqi/example.v1/apps/blog/dao"
	model2 "github.com/shipengqi/example.v1/apps/blog/model"
	e2 "github.com/shipengqi/example.v1/apps/blog/pkg/e"
	"github.com/shipengqi/example.v1/apps/blog/pkg/jwt"
	log "github.com/shipengqi/example.v1/apps/blog/pkg/logger"
	"github.com/shipengqi/example.v1/apps/blog/pkg/setting"
	utils2 "github.com/shipengqi/example.v1/apps/blog/pkg/utils"
	"strconv"
	"strings"

	jwt2 "github.com/dgrijalva/jwt-go"
)

type Interface interface {
	Login(user, pass string) (string, *model2.UserRBAC, error)
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

func (i *identity) Login(user, pass string) (string, *model2.UserRBAC, error) {
	info, err := i.dao.GetUser(user)
	if err != nil {
		return "", nil, e2.Wrap(err, e2.ErrInternalServer.Message())
	}
	if len(info.Username) == 0 {
		return "", nil, e2.ErrUserNotFound
	}
	if info.Deleted {
		return "", nil, e2.ErrUserDeleted
	}
	if info.Locked {
		return "", nil, e2.ErrUserLocked
	}
	if info.Password != utils2.EncodeMD5WithSalt(pass, setting.AppSettings().Salt) {
		return "", nil, e2.ErrPassWrong
	}
	token, err := i.jwt.GenerateToken(user, pass, int(info.ID))
	if err != nil {
		return "", nil, err
	}

	rbac, err := i.dao.GetUserRbac(info.ID)
	if err != nil {
		return "", nil, err
	}

	rbac.U = model2.User{
		Username: info.Username,
		Phone:    info.Phone,
		Email:    info.Email,
	}
	return token, rbac, nil
}

func (i *identity) Authenticate(token string) (claims *jwt.Claims, err error) {
	if token == "" {
		return nil, e2.ErrUnauthorized
	}

	claims, err = i.jwt.ParseToken(token)
	if err != nil {
		if ve, ok := err.(*jwt2.ValidationError); ok {
			log.Error().Err(err).Msgf("Authenticate ValidationError")
			if ve.Errors&jwt2.ValidationErrorMalformed != 0 {
				return nil, e2.ErrTokenMalformed
			} else if ve.Errors&jwt2.ValidationErrorExpired != 0 {
				return nil, e2.ErrTokenExpired
			} else if ve.Errors&jwt2.ValidationErrorNotValidYet != 0 {
				return nil, e2.ErrNotValidYet
			} else {
				return nil, e2.ErrTokenInvalid
			}
		}
		return nil, e2.Wrapf(err, "invalid token")
	}
	return
}

func (i *identity) Authorize(claims *jwt.Claims, url, method string) error {
	userId, err := utils2.DecodeXOR(claims.UID, string(i.jwt.GetSigningKey()))
	if err != nil {
		return e2.ErrTokenInvalid
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
		return e2.ErrForbidden
	}
	return nil
}
