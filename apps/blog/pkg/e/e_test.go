package e

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const TEST_CODE_BASE = 1000000

type testPair struct {
	key int
	val Errno
}

func TestErrnoAllMethods(t *testing.T) {
	table := genErrorTestTable(30)
	t.Run("New", func(t *testing.T) {
		for _, pair := range table {
			v := pair.val
			t.Run(fmt.Sprintf("Code=%d", pair.key), func(t *testing.T) {
				assert.Equal(t, pair.key, v.Code())
				assert.Equal(t, fmt.Sprintf("code: %d, msg: %s", v.Code(), v.Message()), v.Error())
				assert.Equal(t, nil, v.RawError())
				assert.Equal(t, "", v.Stack())
			})
		}

		t.Run("New duplicate error", func(t *testing.T) {
			defer func() {
				if e := recover(); e != nil {
					t.Logf("duplicate error %#v", e)
				}
			}()
			_ = add(0, "OK")
		})

		t.Run("New code less than 0 error", func(t *testing.T) {
			defer func() {
				if e := recover(); e != nil {
					t.Logf("code less than 0 error error %#v", e)
				}
			}()
			_ = New(-1, "less than 0")
		})
	})
}

func TestEqualError(t *testing.T) {
	t.Run("EqualError OK", func(t *testing.T) {
		actual := Is(&Code{
			code:    0,
			message: "OK",
			err:     nil,
		}, OK)
		assert.Equal(t, true, actual)
	})

	t.Run("EqualError ErrInternalServer", func(t *testing.T) {
		actual := Is(errors.New("internal err"), ErrInternalServer)
		assert.Equal(t, true, actual)
	})
}

func TestString(t *testing.T) {
	t.Run("String OK", func(t *testing.T) {
		actual := String("")
		assert.Equal(t, OK, actual)
	})

	t.Run("String test err", func(t *testing.T) {
		actual := String("test err")
		assert.Equal(t, ErrInternalServer.Code(), actual.Code())
		assert.Equal(t, "test err", actual.Message())
	})

	t.Run("String 10", func(t *testing.T) {
		actual := String("10")
		expected := &Code{
			code: 10,
		}
		assert.Equal(t, expected, actual)
	})
}

func TestCause(t *testing.T) {
	t.Run("Cause OK", func(t *testing.T) {
		actual := Cause(nil)
		assert.Equal(t, OK, actual)
	})

	t.Run("Cause test err", func(t *testing.T) {
		actual := Cause(errors.New("test err"))
		assert.Equal(t, ErrInternalServer.Code(), actual.Code())
		assert.Equal(t, "test err", actual.Message())
	})

	t.Run("Cause internal server err", func(t *testing.T) {
		actual := Cause(ErrInternalServer)
		assert.Equal(t, ErrInternalServer, actual)
	})

	t.Run("Cause wrapped err", func(t *testing.T) {
		actual := Cause(errors.Wrap(ErrForbidden, "wrapped"))
		assert.Equal(t, ErrForbidden, actual)
	})

	t.Run("Cause e.Wrap err", func(t *testing.T) {
		actual := Cause(Wrap(ErrForbidden, "wrapped"))
		assert.Equal(t, ErrForbidden, actual)
	})

	t.Run("Cause wrapped native err", func(t *testing.T) {
		err := errors.Wrap(errors.New("native err"), "wrapped")
		actual := Cause(err)
		assert.Equal(t, &Code{
			code:    ErrInternalServer.Code(),
			message: "wrapped: native err",
			err:     nil,
		}, actual)
		is := Is(err, ErrInternalServer)
		assert.Equal(t, true, is)
	})

	t.Run("errors.Cause wrapped err", func(t *testing.T) {
		actual := errors.Cause(errors.Wrap(ErrForbidden, "wrapped"))
		assert.Equal(t, ErrForbidden, actual)
	})
}

func TestWrap(t *testing.T) {
	t.Run("Wrap nil", func(t *testing.T) {
		actual := Wrap(nil, "test msg")
		assert.Equal(t, ErrInternalServer.Code(), actual.Code())
		assert.Equal(t, "test msg", actual.Message())
		assert.Equal(t, "test msg", actual.RawError().Error())
	})

	t.Run("Wrap test err", func(t *testing.T) {
		actual := Wrap(errors.New("test err"), "test msg")
		assert.Equal(t, ErrInternalServer.Code(), actual.Code())
		assert.Equal(t, "test msg", actual.Message())
		assert.Equal(t, "test err", actual.RawError().Error())
	})

	t.Run("Wrap test err with format", func(t *testing.T) {
		actual := Wrapf(errors.New("test err"), "test msg %s", "format")
		assert.Equal(t, ErrInternalServer.Code(), actual.Code())
		assert.Equal(t, "test msg format", actual.Message())
		assert.Equal(t, "test err", actual.RawError().Error())
	})
}

func TestWrapf(t *testing.T) {
	t.Cleanup(func() {
		t.Log("cleanup env")
		SetPrintStack(false)
		SetErrorStackSkip(3)
	})
	SetPrintStack(true)
	t.Run("Wrapf nil with stack", func(t *testing.T) {
		actual := Wrap(nil, "test msg")
		assert.Equal(t, ErrInternalServer.Code(), actual.Code())
		assert.Equal(t, "test msg", actual.Message())
		assert.Equal(t, "test msg", actual.RawError().Error())
		assert.Contains(t, actual.Stack(), "call stack")
		assert.NotContains(t, actual.Stack(), "errno.(*Code).genStackTrace")
	})

	t.Run("Wrap test err with stack", func(t *testing.T) {
		actual := Wrap(errors.New("test err"), "test msg")
		assert.Equal(t, ErrInternalServer.Code(), actual.Code())
		assert.Equal(t, "test msg", actual.Message())
		assert.Equal(t, "test err", actual.RawError().Error())
		assert.Contains(t, actual.Stack(), "call stack")
		assert.NotContains(t, actual.Stack(), "errno.(*Code).genStackTrace")
	})

	t.Run("Wrap test err with format with stack", func(t *testing.T) {
		SetErrorStackSkip(0)
		actual := Wrapf(errors.New("test err"), "test msg %s", "format")
		assert.Equal(t, ErrInternalServer.Code(), actual.Code())
		assert.Equal(t, "test msg format", actual.Message())
		assert.Equal(t, "test err", actual.RawError().Error())
		assert.Contains(t, actual.Stack(), "call stack")
		assert.Contains(t, actual.Stack(), ".(*Code).genStackTrace")
	})
}

func genErrorTestTable(num int) []testPair {
	var tp []testPair
	for i := 0; i < num; i++ {
		key := randInt()
		tp = append(tp, testPair{
			key,
			New(key, randString()),
		})
	}
	return tp
}

// randInt 会生成并返回一个伪随机 int
func randInt() int {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Nanosecond * 2)
	return rand.Int() + TEST_CODE_BASE
}

// randString 会生成并返回一个伪随机字符串
func randString() string {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, rand.Int31())
	return hex.EncodeToString(buf.Bytes())
}
