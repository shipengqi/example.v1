package cmap

import "github.com/pkg/errors"

const (
	ERROR_PACKAGE_CMAP  = "cmap"
	ERROR_TYPE_ILLEGAL_PARAMS  = "illegal parameter"
	ERROR_TYPE_ILLEGAL_PAIR  = "illegal pair type"
)

func newParamError(msg string) error {
	return errors.New(ERROR_TYPE_ILLEGAL_PARAMS)
}

func newPairError(pair Pair) error {
	return errors.New(ERROR_TYPE_ILLEGAL_PAIR)
}
