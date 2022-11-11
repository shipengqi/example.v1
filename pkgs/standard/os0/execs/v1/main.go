package v1

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func base() error {
	c := exec.Command("ls")
	err := c.Run()
	if err != nil {
		return err
	}

	return nil
}

// The content is output to the standard output
func output2stdout() error {
	c := exec.Command("ls")
	c.Stdout = os.Stdout
	c.Stderr = os.Stdin
	err := c.Run()
	if err != nil {
		return err
	}

	return nil
}

// The content is output to a file
func output2file() error {
	f, err := os.Open("out.txt")
	if err != nil {
		return err
	}

	c := exec.Command("ls")
	c.Stdout = f
	c.Stderr = f
	err = c.Run()
	if err != nil {
		return err
	}

	return nil
}

// The content is output to the net
func cal(w http.ResponseWriter, r *http.Request) {
	c := exec.Command("ls")
	c.Stdout = w
	c.Stderr = w
	err := c.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func serve()  {
	http.HandleFunc("/cal", cal)
	_ = http.ListenAndServe(":8080", nil)
}

// Save the content to the buffer
func output2buf() error {
	var buf bytes.Buffer

	c := exec.Command("ls")
	c.Stdout = &buf
	c.Stderr = &buf
	err := c.Run()
	if err != nil {
		return err
	}

	return nil
}

// The content is output to a multiple writer
func output2multi(w http.ResponseWriter, r *http.Request) {
	f, _ := os.Open("out.txt")
	var buf bytes.Buffer

	mw := io.MultiWriter(w, f, &buf)
	c := exec.Command("ls")
	c.Stdout = mw
	c.Stderr = mw
	err := c.Run()
	if err != nil {
		log.Fatal(err)
	}
}


// CombinedOutput
func combined() error {
	c := exec.Command("ls")
	output, err := c.CombinedOutput()
	if err != nil {
		return err
	}
	log.Println(output)
	return nil
}

// stdin
func stdinfromreader() {
	c := exec.Command("cat")
	// command will read content from the given Reader first
	c.Stdin = bytes.NewBufferString("hello, world")
	c.Stdout = os.Stdout
	err := c.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// compress
func zipcompress(d []byte) {
	var out bytes.Buffer
	// -c 压缩, -9 最高压缩等级
	c := exec.Command("bizp2", "-c", "-9")
	// bizp2 会从这个 buffer 中读取数据并压缩
	c.Stdin = bytes.NewBuffer(d)
	c.Stdout = &out
	err := c.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// env
func setenvs(d []byte) {
	c := exec.Command("bash", "-c", "./test.sh")

	env1 := "TEST_ENV1=test1"
	env2 := "TEST_ENV2=test2"

	envs := append(os.Environ(), env1, env2)
	c.Env = envs
	_, err := c.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
}
