package errorsdemo

import (
	"testing"

	"github.com/pkg/errors"
)

// errors.Is 递归调用 Unwrap 并判断每一层的 err 是否相等，如果有任何一层 err 和传入的目标错误相等，则返回 true。
// errors.As 和 errors.Is 差不多，区别在于 errors.Is 是严格判断相等，而 errors.As 则是判断类型是否相同。并提取第一
// 个符合目标类型的错误，用来统一处理某一类错误。

func TestErrorsAs(t *testing.T)  {
	t.Run("custom err", func(t *testing.T) {
		err := NewErrMissingFlag("--help")
		got := errors.As(err, &ErrMissingFlag{})
		if !got {
			t.Fatalf("want true, got: %v", got)
		}
	})

	t.Run("custom err with wrap", func(t *testing.T) {
		err := NewErrMissingFlag("--help")
		wrapped := errors.Wrap(err, "wrap")
		got := errors.As(wrapped, &ErrMissingFlag{})
		if !got {
			t.Fatalf("want true, got: %v", got)
		}
	})
}

func TestErrorsIs(t *testing.T)  {
	t.Run("const err", func(t *testing.T) {
		err := ErrAbort
		got := errors.Is(err, ErrAbort)
		if !got {
			t.Fatalf("want true, got: %v", got)
		}
	})

	t.Run("const err with wrap", func(t *testing.T) {
		err := errors.Wrap(ErrAbort, "wrap")
		got := errors.Is(err, ErrAbort)
		if !got {
			t.Fatalf("want true, got: %v", got)
		}
	})
}
