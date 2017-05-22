#Golang 学习：并发
##一、前言
　　Go语言经常被称为21世纪的C语言，原因一是Go语言设计的简洁优雅，原因二就是Go语言从语言层面原生支持并发。并发的意义，简单通俗来说就是并发的意义就是：你可以同时做多件事！
##二、Goroutine
　　goroutine是Go并行设计的核心。goroutine是通过Go的runtime管理的一个线程管理器。goroutine说到底其实就是线程，但是他比线程更小，十几个goroutine可能体现在底层就是五六个线程，Go语言内部帮你实现了这些goroutine之间的内存共享。执行goroutine只需极少的栈内存(大概是4~5KB)，当然会根据相应的数据伸缩。也正因为如此，可同时运行成千上万个并发任务。goroutine比thread更易用、更高效、更轻便。
- 在函数调用语句前添加`go`关键字就可以创建并发单元了。
```
package main

import (
	"fmt"
	"runtime"
)

func say(s string) {
	for i := 0; i < 5; i++ {
		runtime.Gosched()
		fmt.Println(s)
	}
}
func main() {
	go say("world") //开一个新的Goroutines执行
	say("hello")    //当前Goroutines执行
}
```
输出结果：
```
hello
world
hello
world
hello
world
hello
world
hello
world
```
- 上面的多个goroutine运行在同一个进程里面，共享内存数据，设计要遵循：不要通过共享来通信，而要通过通信来共享。
- runtime.Gosched()表示让CPU把时间片让给别人，下次某时继续恢复执行该goroutine。
- 默认情况下，调度器仅使用单线程，也就是说只实现了并发。想要发挥多核处理器的并行，需要在我们的程序中显示的调用 runtime.GOMAXPROCS(n) 告诉调度器同时使用多个线程。GOMAXPROCS 设置了同时运行逻辑代码的系统线程的最大数量，并返回之前的设置。
- 调用 runtime.Goexit 将立即终止当前 goroutine 执行，调度器确保所有已注册 defer延迟调用被执行。
```
package main

import (
	"runtime"
	"sync"
)

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer println("A.defer")
		func() {
			defer println("B.defer")
			runtime.Goexit() // 终止当前 goroutine
			println("B")     // 不会执行
		}()
		println("A") // 不会执行
	}()
	wg.Wait()
}
```
输出结果：
```
B.defer
A.defe
```
- 和协程 yield 作用类似，Gosched 让出底层线程，将当前 goroutine 暂停，放回队列等待下次被调度执行。
```
package main

import (
	"runtime"
	"sync"
)

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 6; i++ {
			println(i)
			if i == 3 {
				runtime.Gosched()
			}
		}
	}()
	go func() {
		defer wg.Done()
		println("Hello, World!")
	}()
	wg.Wait()
}
```
输出结果：
```

```