package main

import "fmt"

func main() {
	data := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s := data[:3]
	s2 := append(s, 100, 200)
	fmt.Println(data)
	fmt.Println(s)
	fmt.Println(s2)
}
