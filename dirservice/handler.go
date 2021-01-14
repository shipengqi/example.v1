package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

func getSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("method: %s not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}
	queryPath := r.URL.Query().Get("path")
	if strings.TrimSpace(queryPath) == "" {
		http.Error(w, "parameter 'path' is required", http.StatusBadRequest)
		return
	}

	routePath := path.Base(r.URL.Path)
	log.Println("route path: ", routePath)

	log.Println("get query.path: ", queryPath)
	fullPath := filepath.Join(rootPath, queryPath)
	log.Println("get fullPath: ", fullPath)

	sum, err := readFullPath(fullPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if routePath == "statistics" {
		sendResponse(w, sum.Statistics)
		return
	}
	sendResponse(w, sum.Files)
}

func readFullPath(fullPath string) (*Summary, error) {
	var filesCount, dirCount, totalBytes int64
	var fileList []*FileInfo

	sum := &Summary{}

	exists, err := isExists(fullPath)
	if err != nil {
		return sum, err
	}
	if !exists {
		return sum, err
	}
	isDir, err := isDir(fullPath)
	if err != nil {
		return sum, err
	}

	if !isDir {
		info, err := getFileInfo(fullPath)
		if err != nil {
			return sum, err
		}
		sum.Files = &FilesResponse{
			Path: fullPath,
			Dirs: []*FileInfo{
				{
					Name:  info.Name(),
					Size:  info.Size(),
					IsDir: false,
				},
			},
		}
		sum.Statistics = &StatisticsResponse{
			Path:      fullPath,
			DirCount:  0,
			FileCount: 1,
			TotalSize: info.Size(),
		}
		return sum, err
	}

	filesChan := make(chan os.FileInfo)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go readDirParallel(fullPath, filesChan, wg)

	go func() {
		wg.Wait()
		close(filesChan)
	}()

LOOP:
	for {
		select {
		case f, ok := <-filesChan:
			if !ok {
				break LOOP
			}
			if f.IsDir() {
				dirCount++
			} else {
				filesCount++
			}
			totalBytes += f.Size()
			fileList = append(fileList, &FileInfo{
				Name:  f.Name(),
				IsDir: f.IsDir(),
				Size:  f.Size(),
			})
		default:
			// log.Printf("%d files, %d bytes\n", filesCount, totalBytes)
		}
	}
	log.Printf("%d files, %d bytes\n", filesCount, totalBytes)

	sum.Files = &FilesResponse{
		Path: fullPath,
		Dirs: fileList,
	}
	sum.Statistics = &StatisticsResponse{
		Path:      fullPath,
		DirCount:  dirCount,
		FileCount: filesCount,
		TotalSize: totalBytes,
	}

	return sum, nil
}
