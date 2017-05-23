//package main

//func main() {
//	ch := make(chan int, 5)
//	ch <- 1
//	ch <- 1
//	close(ch)
//	ch <- 1 //不能对关闭的channel执行放入操作

//	// 会触发panic
//}

package main

func main() {
	ch := make(chan int, 5)
	ch <- 1
	ch <- 1
	close(ch)
	<-ch //只要channel还有数据，就可能执行取出操作

	//正常结束
}
