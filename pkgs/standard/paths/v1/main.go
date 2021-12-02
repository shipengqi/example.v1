package main

import "path/filepath"

func main() {
	p := "/root/kube.tar.gz"
	b := filepath.Base(p)
	println(b)
}
