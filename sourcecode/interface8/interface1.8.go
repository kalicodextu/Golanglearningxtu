//
// 超集接口对象可转换为子集接口，反之出错。
//
package main

import "fmt"

type Stringer interface {
	String() string
}
type Printer interface {
	String() string
	Print()
}
type User struct {
	id   int
	name string
}

func (self *User) String() string {
	return fmt.Sprintf("%d, %v", self.id, self.name)
}
func (self *User) Print() {
	fmt.Println(self.String())
}
func main() {
	var o Printer = &User{1, "Tom"}
	var s Stringer = o
	fmt.Println(s.String())
}
