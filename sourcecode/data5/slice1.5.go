//- 可以直接修改 struct array/slice 成员
package main

import "fmt"

func main() {
	d := [5]struct {
		x int
	}{}
	s := d[:]
	fmt.Printf("%p, %p, %p, %p\n", &s, &s[0], &d, &d[0])
	d[1].x = 10
	s[2].x = 20
	fmt.Println(d)
	fmt.Printf("%p, %p, %p\n", &s, &d, &d[0])
}
