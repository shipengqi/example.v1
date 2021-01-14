package main

import (
	"encoding/json"
	"net/http"
)

type FileInfo struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
	Size  int64  `json:"size"`
}

type FilesResponse struct {
	Path string      `json:"path"`
	Dirs []*FileInfo `json:"dirs"`
}

type StatisticsResponse struct {
	Path      string `json:"path"`
	DirCount  int64  `json:"dirCount"`
	FileCount int64  `json:"fileCount"`
	TotalSize int64  `json:"totalSize"`
}

type Summary struct {
	Statistics *StatisticsResponse
	Files      *FilesResponse
}

func sendResponse(w http.ResponseWriter, res interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(data)
}
