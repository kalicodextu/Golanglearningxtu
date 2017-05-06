package main

func test() {
	x, y := 10, 20
	defer func(i int) {
		println("defer:", i, y) // y 闭包引用
	}(x) // x 被复制
	x += 10
	y += 100
	println("x =", x, "y =", y)
}
func main() {
	test()
}
