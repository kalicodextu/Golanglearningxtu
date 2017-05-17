//
//switch测试（不支持fallthrough）
//
package main

import "fmt"

type User struct {
	id   int
	name string
}

func (self *User) String() string {
	return fmt.Sprintf("%d, %s", self.id, self.name)
}

func main() {
	var o interface{} = &User{1, "Tom"}
	switch v := o.(type) {
	case nil: // o == nil
		fmt.Println("nil")
	case fmt.Stringer: // interface
		fmt.Println(v)
	case func() string: // func
		fmt.Println(v())
	case *User: // *struct
		fmt.Printf("%d, %s\n", v.id, v.name)
	default:
		fmt.Println("unknown")
	}
}
