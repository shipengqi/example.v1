package utils

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func GetSize(f multipart.File) (int, error) {
	content, err := ioutil.ReadAll(f)

	return len(content)/1024/1024, err
}

func GetExt(fileName string) string {
	return filepath.Ext(fileName)
}

func IsExist(src string) bool {
	_, err := os.Stat(src)

	return os.IsNotExist(err)
}

func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

func IsNotExistMkDir(src string) error {
	if notExist := IsExist(src); notExist == true {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}

func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// MustOpen maximize trying to open the file
func MustOpen(filePath, filename string) (*os.File, error) {
	perm := CheckPermission(filePath)
	if perm == true {
		return nil, errors.Errorf("file.CheckPermission Permission denied file: %s", filePath)
	}

	err := IsNotExistMkDir(filePath)
	if err != nil {
		return nil, errors.Errorf("file.IsNotExistMkDir file: %s, err: %v", filePath, err)
	}

	f, err := Open(fmt.Sprintf("%s/%s", filePath, filename), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, errors.Errorf("file.Open src: %s, err: %v", filename, err)
	}

	return f, nil
}