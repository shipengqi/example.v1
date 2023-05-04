package yamls

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRead(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pwd)
	Read(filepath.Join(pwd, "test.yaml"))
}