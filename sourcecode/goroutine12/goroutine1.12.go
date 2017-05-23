package main //package main
import (
	"fmt"
	"time"
)

//import (
//	"math/rand"
//	"time"
//)

//func NewTest() chan int {
//	c := make(chan int)
//	rand.Seed(time.Now().UnixNano())
//	go func() {
//		time.Sleep(time.Second)
//		c <- rand.Int()
//	}()
//	return c
//}
//func main() {
//	t := NewTest()
//	println(<-t) // 等待 goroutine 结束返回。
//}
func main() {
	w := make(chan bool)
	c := make(chan int, 2)
	go func() {
		select {
		case v := <-c:
			fmt.Println(v)
		case <-time.After(time.Second * 3):
			fmt.Println("timeout.")
		}
		w <- true
	}()
	// c <- 1 // 注释掉，引发 timeout。
	<-w
}
