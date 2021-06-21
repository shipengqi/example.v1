package command

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"

	"github.com/shipengqi/example.v1/cli/pkg/log"
)

func Exec(command string) (string, string, error) {
	cmd := exec.Command("/bin/sh", "-c", command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", "", err
	}
	return string(stdout.Bytes()), string(stderr.Bytes()), nil
}

func ExecSync(command string) error {
	var stderr bytes.Buffer
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Stderr = &stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	reader := bufio.NewReader(stdout)
	var lines int
	for {
		line, _, err2 := reader.ReadLine()
		if err2 != nil || io.EOF == err2 {
			break
		}
		lines ++
		log.Info(string(line))
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	if lines == 0 {
		log.Info(string(stderr.Bytes()))
	}
	return nil
}
