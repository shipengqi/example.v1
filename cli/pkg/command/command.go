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
	cmd := exec.Command("/bin/sh", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	reader := bufio.NewReader(stdout)
	for {
		line, _, err2 := reader.ReadLine()
		if err2 != nil || io.EOF == err2 {
			break
		}
		log.Info(string(line))
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}