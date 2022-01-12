package errorsdemo

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrAbort             = errors.New("abort")
)

type ErrMissingFlag struct {
	names[] string
}

func NewErrMissingFlag(name ...string) ErrMissingFlag {
	return ErrMissingFlag{names: name}
}

func (e ErrMissingFlag) Error() string {
	var b strings.Builder
	b.WriteString("missing flag")
	for i := range e.names {
		name := e.names[i]
		if !strings.HasPrefix(name, "--") {
			name = fmt.Sprintf("--%s", name)
		}
		if i > 0 {
			b.WriteString(" or")
		}
		b.WriteString(" '")
		b.WriteString(name)
		b.WriteString("'")
	}

	return b.String()
}
