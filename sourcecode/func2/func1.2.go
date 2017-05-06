package main

import "fmt"

func test(s string, n ...int) string {
	var x int
	for _, i := range n {
		x += i
	}
	return fmt.Sprintf(s, x)
}
func main() {
	//fmt.Println(test("sum: %d", 1, 2, 3))
	s := []int{1, 2, 3}
	fmt.Println(test("Count: %d", s...))
}
