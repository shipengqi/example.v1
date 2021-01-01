package cmap

import "github.com/shipengqi/example.v1/errors"

const (
	ERROR_PACKAGE_CMAP errors.ErrorPackage = "cmap"
	ERROR_TYPE_ILLEGAL_PARAMS errors.ErrorType = "illegal parameter"
	ERROR_TYPE_ILLEGAL_PAIR errors.ErrorType = "illegal pair type"
)

func newParamError(msg string) errors.V1Error {
	return errors.NewV1ErrorWithPkg(ERROR_PACKAGE_CMAP, ERROR_TYPE_ILLEGAL_PARAMS, msg)
}

func newPairError(pair Pair) errors.V1Error {
	return errors.NewV1ErrorWithPkg(ERROR_PACKAGE_CMAP, ERROR_TYPE_ILLEGAL_PAIR, pair.String())
}