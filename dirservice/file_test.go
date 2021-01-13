package main

import "testing"

func TestReadDirRecursive(t *testing.T) {
	dirPath := "C:\\code\\example.v1\\dirservice"
	err := ReadDirRecursive(dirPath)
	if err != nil {
		t.Fatalf("read %s, err: %+v", dirPath, err)
	}

	for _, v := range fileList {
		t.Logf("file: %+v", v)
	}
}
