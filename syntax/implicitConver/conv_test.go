package conv

import (
	"fmt"
	"testing"
)

func TestUser_String(t *testing.T) {
	u := User{Name: "pooky"}
	fmt.Println(u)
	fmt.Println(&u)
	fmt.Printf("%p\n",&u)
}
