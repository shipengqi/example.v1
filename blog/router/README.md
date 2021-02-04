
静态文件：
```go
// images static file system
r.StaticFS("/images", http.Dir(upload.GetImageFullPath()))
```

源码：

```go
// StaticFS works just like `Static()` but a custom `http.FileSystem` can be used instead.
// Gin by default user: gin.Dir()
func (group *RouterGroup) StaticFS(relativePath string, fs http.FileSystem) IRoutes {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}
	handler := group.createStaticHandler(relativePath, fs)
	// *filepath 将匹配所有文件路径，并且 *filepath 必须在 Pattern 的最后
	urlPattern := path.Join(relativePath, "/*filepath")

	// Register GET and HEAD handlers
	group.GET(urlPattern, handler)
	group.HEAD(urlPattern, handler)
	return group.returnObj()
}
```

在暴露的 URL 中禁止了 `*` 和 `:` 符号的使用，通过 `createStaticHandler` 创建了静态文件服务，实质最终调用的还是 `fileServer.ServeHTTP` 和
一些处理逻辑。

```go
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := group.calculateAbsolutePath(relativePath)
	// http.StripPrefix 主要作用是从请求 URL 的路径中删除给定的前缀，最终返回一个 Handler
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	_, nolisting := fs.(*onlyfilesFS)
	return func(c *Context) {
		if nolisting {
			c.Writer.WriteHeader(404)
		}
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}
```

通常 `http.FileServer` 要与 `http.StripPrefix` 相结合使用，否则当你运行：
```go
http.Handle("/images", http.FileServer(http.Dir("upload/images")))
```

会无法正确的访问到文件目录，因为 `upload/images` 也包含在了 URL 路径中，必须使用：

```go
http.Handle("/images", http.StripPrefix("upload/images", http.FileServer(http.Dir("upload/images"))))
```
