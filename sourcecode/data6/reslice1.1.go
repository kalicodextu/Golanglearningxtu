package main

import "fmt"

func main() {
	s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s1 := s[2:5] // [2 3 4]
	s1[2] = 100
	s2 := s1[2:6] // [100 5 6 7]
	s2[3] = 200
	fmt.Println(s)
}
