package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)



var fileList []FileInfo

func ReadDirRecursive(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			fileList = append(fileList, FileInfo{
				Name:  f.Name(),
				Size:  f.Size(),
				IsDir: true,
			})
			err := ReadDirRecursive(filepath.Join(path, f.Name()))
			if err != nil {
				return err
			}
		} else {
			fileList = append(fileList, FileInfo{
				Name:  f.Name(),
				Size:  f.Size(),
				IsDir: false,
			})
		}
	}

	return nil
}

func IsExits(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func IsDir(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if !f.IsDir() {
		return false, nil
	}
	return true, nil
}

func GetFileInfo(path string) (os.FileInfo, error) {
	f, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}
