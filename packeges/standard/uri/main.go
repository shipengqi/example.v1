package main

import (
	"fmt"
	"net/url"
)

func main() {
	tests := []struct {
		raw    string
		change string
	}{
		{
			raw:    "http://127.0.0.1:8080/url/123",
			change: "192.168.1.1",
		},
		{
			raw: "https://127.0.0.1:8080/url/123",
		},
		{
			raw: "https://www.baidu.com/url/123",
		},
		{
			raw: "www.baidu.com/url/123",
		},
		{
			raw: "www.baidu.com:443/url/123",
		},
		{
			raw: "https://www.baidu.com:443/url/123",
		},
	}
	for _, v := range tests {
		parse(v)
	}

}

func parse(uri struct {
	raw    string
	change string
}) {
	parsed, _ := url.Parse(uri.raw)

	// Host
	parsedHost := parsed.Hostname()
	// Port
	parsedPort := parsed.Port()
	// Path
	parsedPath := parsed.Path
	parsedUrl := parsed.String()

	fmt.Println("Parsed Host: ", parsedHost)
	fmt.Println("Parsed Port: ", parsedPort)
	fmt.Println("Parsed Path: ", parsedPath)
	fmt.Println("Parsed Url: ", parsedUrl)

	// Change URl Host
	if uri.change != "" {
		parsed.Host = uri.change + ":" + parsed.Port()
		fmt.Println("Changed Host: ", parsed.String())
	}
	fmt.Println("")
}
