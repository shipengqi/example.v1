package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	done := make(chan error, 1)
	reader, writer, err := os.Pipe()
	if err != nil {
		fmt.Println("Pipe", err)
		return
	}
	go func() {
		defer writer.Close()
		for i := 0; i < 5; i ++ {
			_, err = writer.WriteString(fmt.Sprintf("content %d", i)) // writer 写数据
			if err != nil {
				fmt.Println("Write", err)
				return
			}
			fmt.Println("write", i)
			time.Sleep(1*time.Second)
		}
	}()

	s := bufio.NewScanner(reader)
	go func() {
		for s.Scan() {
			fmt.Println(string(s.Bytes()))
		}
		if err := reader.Close(); err != nil {
			fmt.Println("reader.Close", err)
		}
		fmt.Println("wait read end")
		done <- s.Err()
		close(done)
	}()
	fmt.Println("wait done chan")
	e := <-done
	fmt.Println("done: ", e)
}

// Output
// wait done chan
// write 0
// write 1
// write 2
// write 3
// write 4
// content 0content 1content 2content 3content 4
// wait read end
// done:  <nil>
