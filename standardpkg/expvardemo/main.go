package main

import (
	"expvar"
	"fmt"
	"net/http"
)

// expvar 为公共变量提供了一个标准化的接口，如服务器中的操作计数器。它以 JSON 格式通过 `/debug/vars` 接口以 HTTP 的方式公开这些公共变量。
// 设置或修改这些公共变量的操作是原子的。
// 除了为程序增加 HTTP handler，此包还注册以下变量：
// cmdline   os.Args
// memstats  runtime.Memstats
//
// 通过如下形式引入到程序中：
// import _ "expvar"
// 只能注册其 HTTP handler 和上述变量。

var visits = expvar.NewInt("visits")

func handler(w http.ResponseWriter, r *http.Request) {
	visits.Add(1)
	_, _ = fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

// 导入 expvar 包后，为 http.DefaultServeMux 上的 PATH /debug/vars 注册一个处理函数
func main()  {
	http.HandleFunc("/visits", handler)
	_ = http.ListenAndServe(":8080", nil)
}

// 运行代码并访问 http://localhost:8080/debug/vars
// 默认情况下该包注册了 os.Args 和 runtime.Memstats 两个指标。expvar.NewInt("visits")，增加了一个新的指标，
// 通过访问 http://localhost:8080/visits 来增加 visits 计数器。
// Output:
//{
//	"cmdline": [
//		"C:\\Users\\shipengqi\\AppData\\Local\\Temp\\___go_build_main_go.exe"
//	],
//	"memstats": {
//		"Alloc": 297984,
//		"TotalAlloc": 297984,
//		"Sys": 6705288,
//		"Lookups": 0,
//		"Mallocs": 971,
//		"Frees": 20,
//		"HeapAlloc": 297984,
//		"HeapSys": 4030464,
//		"HeapIdle": 3039232,
//		"HeapInuse": 991232,
//		"HeapReleased": 3039232,
//		"HeapObjects": 951,
//		"StackInuse": 163840,
//		"StackSys": 163840,
//		"MSpanInuse": 28288,
//		"MSpanSys": 32768,
//		"MCacheInuse": 20448,
//		"MCacheSys": 32768,
//		"BuckHashSys": 3675,
//		"GCSys": 1407368,
//		"OtherSys": 1034405,
//		"NextGC": 4473924,
//		"LastGC": 0,
//		"PauseTotalNs": 0,
//		"PauseNs": [
//			0,
//			0
//		],
//		"PauseEnd": [
//			0
//		],
//		"NumGC": 0,
//		"NumForcedGC": 0,
//		"GCCPUFraction": 0,
//		"EnableGC": true,
//		"DebugGC": false,
//		"BySize": [
//			{
//				"Size": 64,
//				"Mallocs": 45,
//				"Frees": 0
//			},
//			{
//				"Size": 80,
//				"Mallocs": 13,
//				"Frees": 0
//			},
//			{
//				"Size": 96,
//				"Mallocs": 14,
//				"Frees": 0
//			},
//			{
//				"Size": 112,
//				"Mallocs": 3,
//				"Frees": 0
//			},
//			{
//				"Size": 128,
//				"Mallocs": 6,
//				"Frees": 0
//			},
//			{
//				"Size": 144,
//				"Mallocs": 2,
//				"Frees": 0
//			},
//			{
//				"Size": 10240,
//				"Mallocs": 12,
//				"Frees": 0
//			}
//		]
//	},
//	"visits": 0
//}