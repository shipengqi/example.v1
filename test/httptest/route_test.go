package route

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 模拟 http 网络调用

func init() {
	// 初始化路由
	Routes()
}

// 利用 httptest.NewRecorder() 创建一个h ttp.ResponseWriter，模拟了真实服务端的响应，
// 这种响应时通过调用 http.DefaultServeMux.ServeHTTP 方法触发的。
func TestSendJson(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/sendJson", nil)

	rw := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rw, req)

	t.Log("code: ", rw.Code)
	t.Log("body: ", rw.Body.String())
}


// 另一种模拟调用的方式，httptest.NewServer 函数模拟服务器的创建，接收一个 http.Handler 处理 API 请求的接口。
// 例如 httptest.NewServer 函数可以传入 Gin 的 engine 对象。
func mockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(SendJson))
}

func TestSendJson2(t *testing.T) {
	// 创建一个模拟的服务器
	server := mockServer()
	defer server.Close()

	// Get 请求发往模拟服务器的地址
	res, err := http.Get(server.URL)
	if err != nil {
		t.Fatal("create get request", err)
	}
	defer res.Body.Close()

	t.Log("code:", res.StatusCode)
	json, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal("read response", err)
	}

	t.Log("body: ", string(json))
}
