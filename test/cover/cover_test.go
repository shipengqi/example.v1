package cover

import (
	"testing"
)

func TestTag(t *testing.T) {
	Tag(1)
	Tag(2)
}

// go test -cover 可以显示覆盖率，例如上面的示例输出：
//Android
//Go
//PASS
//coverage: 60.0% of statements
//ok      github.com/shipengqi/example.v1/test/cover      0.151s

// go test -coverprofile=coverage.out 可以手机更多信息，输出到 coverage.out 文件中
// go tool cover -func=coverage.out 可以查看覆盖率报告，输出：
//github.com/shipengqi/example.v1/test/cover/cover.go:5:  Tag             60.0%
//total:                                                  (statements)    60.0%

// go tool cover -html=coverage.out 可以在浏览器打开覆盖率报告，可以查看具体没有被覆盖的代码
// go tool cover -html=coverage.out -o=coverage..html 可以生成 html 文件

