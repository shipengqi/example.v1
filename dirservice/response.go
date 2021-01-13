package main

type FileInfo struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
	Size  int64  `json:"size"`
}

type FilesResponse struct {
	Path string     `json:"path"`
	Dirs []FileInfo `json:"dirs"`
}

type StatisticsResponse struct {
	Path      string `json:"path"`
	DirCount  int    `json:"dirCount"`
	FileCount int    `json:"fileCount"`
	TotalSize int    `json:"totalSize"`
}
