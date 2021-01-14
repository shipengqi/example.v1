package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var pool = make(chan int, 5)

func readDirParallel(dir string, filesChan chan<- os.FileInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, f := range readDir(dir) {
		if f.IsDir() {
			wg.Add(1)
			sub := filepath.Join(dir, f.Name())
			filesChan <- f
			go readDirParallel(sub, filesChan, wg)
		} else {
			filesChan <- f
		}
	}
}

func readDir(path string) []os.FileInfo {
	pool <- 1
	defer func() {
		<- pool
	}()
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Printf("read path: %s, err: %s", path, err)
		return nil
	}
	return files
}

func isExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func isDir(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if !f.IsDir() {
		return false, nil
	}
	return true, nil
}

func getFileInfo(path string) (os.FileInfo, error) {
	f, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}
