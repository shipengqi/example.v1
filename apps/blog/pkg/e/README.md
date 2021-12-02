# golang error 处理

Golang 开发中经常需要检查返回的错误值并作相应处理：

```go
package main

import (
   "database/sql"
   "fmt"
)

func foo() error {
   return sql.ErrNoRows
}

func bar() error {
   return foo()
}

func main() {
   err := bar()
   if err != nil {
      fmt.Printf("got err, %+v\n", err) // got err, sql: no rows in result set
   }
}
```

有时需要根据返回的 error 类型作不同处理，例如：

```go
package main

import (
   "database/sql"
   "fmt"
)

func foo() error {
   return sql.ErrNoRows
}

func bar() error {
   return foo()
}

func main() {
   err := bar()
   if err == sql.ErrNoRows {
      fmt.Printf("data not found, %+v\n", err) // data not found, sql: no rows in result set
      return
   }
   if err != nil {
      // Unknown error
   }
}
```

经常需要为错误增加上下文信息后再返回，以方便调用者了解错误场景：

```go
func foo() error {
   return fmt.Errorf("foo err, %v", sql.ErrNoRows)
}
```

这时 `err == sql.ErrNoRows` 便不再成立。而且上述写法都在**返回错误时都丢掉了调用栈这个重要的诊断信息**。


## 解决方案

### github.com/pkg/errors

`github.com/pkg/errors` 中包含三个关键的方法：
1. `Wrap` 方法用来包装底层错误，增加上下文文本信息并附加调用栈。 一般用于包装对第三方代码（标准库或第三方库）的调用。
2. `WithMessage` 方法仅增加上下文文本信息，不附加调用栈。 如果确定错误已被 `Wrap` 过或不关心调用栈，可以使用此方法。
注意：**不要反复 `Wrap`，会导致调用栈重复**。
3. `Cause` 方法用来判断底层错误 。

```go
package main

import (
   "database/sql"
   "fmt"

   "github.com/pkg/errors"
)

func foo() error {
   return errors.Wrap(sql.ErrNoRows, "foo failed")
}

func bar() error {
   return errors.WithMessage(foo(), "bar failed")
}

func main() {
   err := bar()
   if errors.Cause(err) == sql.ErrNoRows {
      fmt.Printf("data not found, %v\n", err)
      fmt.Printf("%+v\n", err)
      return
   }
   if err != nil {
      // unknown error
   }
}

/*Output:
data not found, bar failed: foo failed: sql: no rows in result set
sql: no rows in result set
foo failed
main.foo
    /usr/three/main.go:11
main.bar
    /usr/three/main.go:15
main.main
    /usr/three/main.go:19
runtime.main
    ...
*/
```

使用 `%v` 作为格式化参数，那么错误信息会保持一行， 其中依次包含调用栈的上下文文本。 使用 `%+v`，则会输出完整的调用栈详情。

如果不需要增加额外上下文信息，仅附加调用栈后返回，可以使用 `WithStack` 方法：

```go
func foo() error {
   return errors.WithStack(sql.ErrNoRows)
}
```

无论是 `Wrap`， `WithMessage` 还是 `WithStack`，当传入的 err 参数为 `nil` 时， 都会返回 `nil`。

#### Cause

`Cause` 方法的源码：

```go
// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}
```

可以看出 `Cause` 方法的入参 err 需要实现 causer 接口，才可以得到底层的错误。


### 其他方案

`golang.org/x/xerrors` 和 1.13 后的 `errors` 标准库，目前的实现并不完善。
