package main

//func InitRouter() *gin.Engine {
//	r := gin.New()
//
//	gin.SetMode(gin.DebugMode)
//
//	r.Use(gin.Logger())
//	r.Use(gin.Recovery())
//
//	r.NoRoute(func(c *gin.Context) {
//		c.String(http.StatusNotFound, "incorrect route")
//	})
//
//	v1 := r.Group("/api/v1")
//	{
//		v1.GET("/statistics", GetStatistics)
//		v1.GET("/files", GetFiles)
//	}
//
//	return r
//}

//func GetStatistics(c *gin.Context) {
//	//
//}



//func GetFiles(c *gin.Context) {
//	path := c.Query("path")
//
//	if strings.TrimSpace(path) == "" {
//		c.JSON(http.StatusBadRequest, "parameter 'path' is required")
//		return
//	}
//
//	fullPath := filepath.Join(rootPath, path)
//	exists, err := IsExits(fullPath)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, err)
//		return
//	}
//	if !exists {
//		c.JSON(http.StatusNotFound, fmt.Sprintf("path: %s not found", fullPath))
//		return
//	}
//
//	isDir, err := IsDir(fullPath)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, err)
//		return
//	}
//
//	if !isDir {
//		info, err := GetFileInfo(fullPath)
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, err)
//			return
//		}
//		c.JSON(http.StatusOK, FilesResponse{
//			Path: fullPath,
//			Dirs: []FileInfo{
//				{
//					Name:  info.Name(),
//					Size:  info.Size(),
//					IsDir: false,
//				},
//			},
//		})
//		return
//	}
//
//	err = ReadDirRecursive(fullPath)
//	defer func() {
//		fileList = []FileInfo{}
//	}()
//
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, err)
//		return
//	}
//
//	c.JSON(http.StatusOK, FilesResponse{
//		Path: fullPath,
//		Dirs: fileList,
//	})
//}
