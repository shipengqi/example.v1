package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func generate(writer *io.PipeWriter) {
	for {
		s := fmt.Sprintf("random: %d", r.Uint32())
		_, err := writer.Write([]byte(s))
		if nil != err {
			log.Fatal(err)
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func main() {
	rp, wp := io.Pipe()
	for i := 0; i < 20; i++ {
		go generate(wp)
	}
	time.Sleep(1 * time.Second)
	data := make([]byte, 64)
	for {
		n, err := rp.Read(data)
		if nil != err {
			log.Fatal(err)
		}
		if 0 != n {
			log.Println("main loop", n, string(data))
		}
		time.Sleep(1 * time.Second)
	}
}
