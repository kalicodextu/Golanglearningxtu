package main

import "fmt"

type Queue struct {
	elements []interface{}
}

func NewQueue() *Queue { // 创建对象实例。
	return &Queue{make([]interface{}, 10)}
}
func (*Queue) Push(e interface{}) error { // 省略 receiver 参数名。
	panic("not implemented")
}

// func (Queue) Push(e int) error { // Error: method redeclared: Queue.Push
// panic("not implemented")
// }
func (self *Queue) length() int { // receiver 参数名可以是 self、this 或其他。
	return len(self.elements)
}
func main() {
	fmt.Println(" ")
}
