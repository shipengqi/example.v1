package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

var rootPath string
var fileList []FileInfo

func init()  {
	flag.StringVar(&rootPath, "root", "", "Specifies the root path.")
}

func main()  {

	flag.Parse()

	if err := preCheck(); err != nil {
		log.Fatalf("precheck err: %s", err)
	}

	// r := InitRouter()

	http.HandleFunc("/files", getFiles)
	// http.HandleFunc("/statistics", handleRequest)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Listen err: %s", err)
	}
}

func preCheck() error {
	if strings.TrimSpace(rootPath) == "" {
		return errors.New("root path is invalid")
	}

	return nil
}

func getFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	qpath := r.URL.Query().Get("path")
	if strings.TrimSpace(qpath) == "" {
		http.Error(w, "parameter 'path' is required", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(rootPath, qpath)
	exists, err := IsExits(fullPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, fmt.Sprintf("path: %s not found", fullPath), http.StatusNotFound)
		return
	}

	isDir, err := IsDir(fullPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !isDir {
		info, err := GetFileInfo(fullPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		res := FilesResponse{
			Path: fullPath,
			Dirs: []FileInfo{
				{
					Name:  info.Name(),
					Size:  info.Size(),
					IsDir: false,
				},
			},
		}
		data, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(data)
		return
	}

	err = ReadDirRecursive(fullPath)
	defer func() {
		fileList = []FileInfo{}
	}()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	res := FilesResponse{
		Path: fullPath,
		Dirs: fileList,
	}
	data, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(data)
}

// main handler function
func handleRequest(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case http.MethodGet:
		err = handleGet(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) error {
	fmt.Println(r.URL.Path)
	fmt.Println(path.Base(r.URL.Path))
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("OK"))
	return nil
}