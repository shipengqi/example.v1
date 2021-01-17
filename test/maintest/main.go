package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
)

// 有时需要在测试之前或之后进行额外的设置（setup）或拆卸（teardown）；有时，测试还需要控制在主线程上运行的代码。为了支持这些需求，
// testing 包提供了 `TestMain` 函数 :
// func TestMain(m *testing.M)
//
// 如果测试文件中包含该函数，那么生成的测试将调用 TestMain(m)，而不是直接运行测试。TestMain 运行在主 goroutine 中 , 可以在
// 调用 m.Run 前后做任何设置和拆卸。
// 注意，在 TestMain 函数的最后，应该使用 m.Run 的返回值作为参数去调用 os.Exit。
//
// 在调用 TestMain 时 , flag.Parse 并没有被调用。所以，如果 TestMain 依赖于 command-line 标志（包括 testing 包的标志），则
// 应该显式地调用 flag.Parse。注意，这里的依赖是指，若 TestMain 函数内需要用到 command-line 标志，则必须显式地调用 flag.Parse，否则
// 不需要，因为 m.Run 中调用 flag.Parse。

func main() {
	http.HandleFunc("/topic/", handleRequest)

	_ = http.ListenAndServe(":8080", nil)
}

// main handler function
func handleRequest(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case http.MethodGet:
		err = handleGet(w, r)
	case http.MethodPost:
		err = handlePost(w, r)
	default:
		fmt.Println("unknown")
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// 如 GET /topic/1
func handleGet(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return err
	}
	topic, err := FindTopic(id)
	if err != nil {
		return err
	}
	output, err := json.MarshalIndent(&topic, "", "\t\t")
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return nil
}

// POST /topic/
func handlePost(w http.ResponseWriter, r *http.Request) (err error) {
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var topic = new(Topic)
	err = json.Unmarshal(body, &topic)
	if err != nil {
		return
	}

	err = topic.Create()
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
