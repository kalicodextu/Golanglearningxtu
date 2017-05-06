package main

func main() {
	m := map[int]struct {
		name string
		age  int
	}{
		1: {"user1", 10}, // 可省略元素类型。
		2: {"user2", 20},
	}
	println(m[1].name)
}
