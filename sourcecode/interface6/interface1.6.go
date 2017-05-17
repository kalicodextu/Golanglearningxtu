//
//利用类型推断，可判断接口对象是否某个具体的接口或类型。
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
	if value, ok := o.(fmt.Stringer); ok { // ok-idiom
		fmt.Println(value, ok)
	}
	u := o.(*User)
	// u := o.(User) // panic: interface is *main.User, not main.User
	fmt.Println(u)
}
