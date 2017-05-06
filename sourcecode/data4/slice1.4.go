package main

import "fmt"

func main() {
	s := []int{0, 1, 2, 3}
	p := &s[2] // *int, 获取底层数组元素指针。
	*p += 100
	fmt.Println(s)
}
