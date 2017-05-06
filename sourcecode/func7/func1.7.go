package main

func test(x int) {
	defer println("a")
	defer println("b")
	defer func() {
		println(100 / x) // div0 异常未被捕获，逐步往外传递，最终终止进程。
	}()
	defer println("c")
}
func main() {
	test(0)
}
