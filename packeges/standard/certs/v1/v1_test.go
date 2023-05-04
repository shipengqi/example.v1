package v1

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTransport(t *testing.T) {
	t.Run("CA1", func(t *testing.T) {
		trans, err := CreateTransport("testdata/ca.crt")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(trans.TLSClientConfig.RootCAs.Subjects()))
		for _, v := range trans.TLSClientConfig.RootCAs.Subjects() {
			t.Log(string(v))
		}
	})

	t.Run("CA2", func(t *testing.T) {
		trans, err := CreateTransport("testdata/ca2.crt")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(trans.TLSClientConfig.RootCAs.Subjects()))
		for _, v := range trans.TLSClientConfig.RootCAs.Subjects() {
			ns := bytes.Split(v, []byte(","))
			t.Logf("%s", ns)
		}
	})
}
