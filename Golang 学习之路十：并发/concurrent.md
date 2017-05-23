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
	"fmt"

	"runtime"
)

func say(s string) {

	for i := 0; i < 2; i++ {

		runtime.Gosched()

		fmt.Println(s)

	}

}

func main() {

	go say("world")

	say("hello")

}
```
输出结果：
```
hello
world
hello
world
```
##三、Chanel
　　goroutine是运行在相同的地址空间（即CSP模式），Go提供了通信机制channel来进行同步。channel可以与Unix shell 中的双向管道做类比：可以通过它发送或者接收值。这些值只能是特定的类型：channel类型。定义一个channel时，也需要定义发送到channel的值的类型。因为默认同步模式，需要发送与接收配对，否则会被阻塞,知道另一方准备好后被唤醒。
- 使用`make`创建`channel`
```
ci := make(chan int)
cs := make(chan string)
cf := make(chan interface{})
```
- channel通过操作符`<-`来接收和发送数据
```
ch <- v // 发送v到channel ch.
v := <-ch // 从ch中接收数据，并赋值给v
```
- 默认情况下，channel接收和发送数据都是阻塞的，除非另一端已经准备好，这样就使得Goroutines同步变的更加的简单，而不需要显式的lock。
```
package main

import "fmt"

func Afun(ch chan int) {

	ch <- 1
	ch <- 2

}

func main() {
	c := make(chan int)
	go Afun(c)
	<-c
	<-c

}
```
上面的程序虽无输出但可以运行，而如果是下面的写法是无法通过编译运行的。
```
// 1
package main

import "fmt"

func Afun(ch chan int) {

	ch <- 1
	ch <- 2

}

func main() {
	c := make(chan int)
	<-c
	<-c
	go Afun(c)

}
```
```
// 2
package main

func main() {
	c := make(chan int)
	c <- 1
	c <- 2
	<-c
	<-c
}
```
前一段代码最终会正常结束，但是后一段代码会发生死锁。为什么会出现这种现象呢，咱们把上面两段代码的逻辑跑一下。

第一段代码：

        1. 创建了一个无缓冲channel

        2. 启动了一个goroutine，这个routine中对channel执行放入数据操作，但是因为这时候channel为空，所以这个取出操作发生阻塞，但是主routine可没有发生阻塞，它还在继续运行呢

        3. 主goroutine这时候继续执行下一行，在channel中取出了数据

        4. 这时阻塞的那个routine检测到了channel没有数据了，所以解除阻塞，放入数据到channel，程序就此完毕



第二、三段代码：

        1.  创建了一个无缓冲的channel

        2.  主routine要从channel中取出一个数据，但是因为channel没有缓冲，相当于channel一直都是空的，所以这里会发生阻塞。可是下面的那个goroutine还没有创建呢，主routine在这里一阻塞，整个程序就只能这么一直阻塞下去了，就会发生死锁！

※从这里可以看出，对于无缓冲的channel，放入操作和取出操作不能再同一个routine中，而且应该是先确保有某个routine对它执行取出操作，然后才能在另一个routine中执行放入操作。

- 异步方式通过判断缓冲区来决定是否阻塞。如果缓冲区已满，发送被阻塞；缓冲区为空，接收被阻塞。
```
ch := make(chan type, value)
value == 0 ! 无缓冲（阻塞）
value > 0 ! 缓冲（非阻塞，直到value 个元素）
```
用简单例子证明：
```
package main

import "fmt"

func main() {
	c := make(chan int, 2)//修改2为1就报错，修改2为3可以正常运行
	c <- 1
	c <- 2
	fmt.Println(<-c)
	fmt.Println(<-c)
}
```
- 缓存区是内部属性，并非构成要素
```
var a, b chan int = make(chan int), make(chan int, 3)
```
- 用`range`可以判断`channel`是否关闭，我们还可以用`ok-idiom`模式来判断`channel`是否关闭。
```
package main

import "fmt"

func main() {
	data := make(chan int, 3) // 缓冲区可以存储 3 个元素
	exit := make(chan bool)
	data <- 1 // 在缓冲区未满前，不会阻塞。
	data <- 2
	data <- 3
	go func() {
		for d := range data { // 在缓冲区未空前，不会阻塞。
			fmt.Println(d)
		}
		exit <- true
	}()
	data <- 4 // 如果缓冲区已满，阻塞。
	data <- 5
	close(data)
	<-exit
}
```
`ok-idiom`模式
```
for {
	if d, ok := <-data; ok {
		fmt.Println(d)
	} else {
		break
	}
}
```
- 关闭后的channel可以取数据，但是不能放数据。而且，channel在执行了close()后并没有真的关闭，channel中的数据全部取走之后才会真正关闭。
```
package main

func main() {
	ch := make(chan int, 5)
	ch <- 1
	ch <- 1
	close(ch)
	ch <- 1 //不能对关闭的channel执行放入操作
        
        // 会触发panic
}
//// panic: send on closed channel
```
```
package main

func main() {
	ch := make(chan int, 5)
	ch <- 1
	ch <- 1
	close(ch)
	<-ch //只要channel还有数据，就可能执行取出操作

	//正常结束
}
```
再看：
```
package main

import "fmt"

func main() {
	ch := make(chan int, 5)
	ch <- 1
	ch <- 1
	ch <- 1
	ch <- 1
	close(ch) //如果执行了close()就立即关闭channel的话，下面的循环就不会有任何输出了
	for {
		data, ok := <-ch
		if !ok {
			break
		}
		fmt.Println(data)
	}
}
```
输出结果：
```
1
1
1
1
```
##四、channel 单向
　　我们还可以以数据在通道中的传输方向为依据来划分通道。默认情况下，通道都是双向的，即双向通道。如果数据只能在通道中单向传输，那么该通道就被称作单向通道。
- 我们在初始化一个通道值的时候不能指定它为单向。但是，在编写类型声明的时候，我们却是可以这样做：
```
type Receiver <-chan int //代表了一个只可从中接收数据的单向通道类型。
type Sender chan<- int //代表了一个只可从中发送数据的单向通道类型。

var myChannel = make(chan int, 3)
var sender Sender = myChannel
var receiver Receiver = myChannel  
```

　　单向通道的主要作用是约束程序对通道值的使用方式。比如，我们调用一个函数时给予它一个发送通道作为参数，以此来约束它只能向该通道发送数据。又比如，一个函数将一个接收通道作为结果返回，以此来约束调用该函数的代码只能从这个通道中接收数据。
```
package main

import (
    "fmt"
	"time"
)

type Sender chan<- int

type Receiver <-chan int

func main() {
	var myChannel = make(chan int)
	var number = 6
	go func() {
		var sender Sender = myChannel
		sender <- number
		fmt.Println("Sent!")
	}()
	go func() {
		var receiver Receiver = myChannel
		fmt.Println("Received!", <-receiver)
	}()
	// 让main函数执行结束的时间延迟1秒，
	// 以使上面两个代码块有机会被执行。
	time.Sleep(time.Second)
}
```
输出结果：
```
Received! 6
Sent!
```	
- 不能将单向 channel 转换为普通 channel
##五、select
　　如果需要同时处理多个channel，可使用`select`语句监听`channel`上的数据流动。它随机选择一个可用`channel`做收发操作，或执行 default case。`select`是默认阻塞的。
```
package main

import (
	"fmt"
	"os"
)

func main() {
	a, b := make(chan int, 3), make(chan int)
	go func() {
		v, ok, s := 0, false, ""
		for {
			select { // 随机选择可用 channel，接收数据。
			case v, ok = <-a:
				s = "a"
			case v, ok = <-b:
				s = "b"
			}
			if ok {
				fmt.Println(s, v)
			} else {
				os.Exit(0)
			}
		}
	}()
	for i := 0; i < 5; i++ {
		select { // 随机选择可用 channel，发送数据。
		case a <- i:
		case b <- i:
		}
	}
	close(a)
	select {} // 没有可用 channel，阻塞 main goroutine。
}
```
输出结果：
```
// 随机，结果可能不同
b 3
b 4
a 0
a 1
a 2
```
- 在select里面还有default语法，select其实就是类似switch的功能，default就是当监听的channel都没有准备好的时候，默认执行的（select不再阻塞等待channel）。
```
select {
	case i := <-c:
	// use i
	default:
	// 当c阻塞的时候执行这里
}
```
- 用`select`设置超时，避免整个程序进入阻塞
```
package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)
	o := make(chan bool)
	go func() {
		for {
			select {
			case v := <-c:
				println(v)
			case <-time.After(5 * time.Second):
				println("timeout")
				o <- true
				break
			}
		}
	}()
	c <- 1
	fmt.Println(<-o)
}
```
输出结果：
```
1
timeout
true
```
##六、模式
- 用简单工厂模式打包并发任务和`channel`。
```
package main

import (
	"math/rand"
	"time"
)

func NewTest() chan int {
	c := make(chan int)
	rand.Seed(time.Now().UnixNano())
	go func() {
		time.Sleep(time.Second)
		c <- rand.Int()
	}()
	return c
}
func main() {
	t := NewTest()
	println(<-t) // 等待 goroutine 结束返回。
}
```
- 用 channel 实现信号量 (semaphore)
```
package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(3)
	sem := make(chan int, 1)
	for i := 0; i < 3; i++ {
		go func(id int) {
			defer wg.Done()
			sem <- 1 // 向 sem 发送数据，阻塞或者成功。
			for x := 0; x < 3; x++ {
				fmt.Println(id, x)
			}
			<-sem // 接收数据，使得其他阻塞 goroutine 可以发送数据。
		}(i)
	}
	wg.Wait()
}
```
输出结果：
```
0 0
0 1
0 2
1 0
1 1
1 2
2 0
2 1
2 2
```
- 用 closed channel 发出退出通知。
```
package main

import (
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	quit := make(chan bool)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			task := func() {
				println(id, time.Now().Nanosecond())
				time.Sleep(time.Second)
			}
			for {
				select {
				case <-quit: // closed channel 不会阻塞，因此可用作退出通知。
					return
				default: // 执行正常任务。
					task()
				}
			}
		}(i)
	}
	time.Sleep(time.Second * 5) // 让测试 goroutine 运行一会。
	close(quit)                 // 发出退出通知。
	wg.Wait()
}
```
- channel 是第一类对象，可传参 (内部实现为指针) 或者作为结构成员。
```
import "fmt"

type Request struct {
	data []int
	ret  chan int
}

func NewRequest(data ...int) *Request {
	return &Request{data, make(chan int, 1)}
}
func Process(req *Request) {
	x := 0
	for _, i := range req.data {
		x += i
	}
	req.ret <- x
}
func main() {
	req := NewRequest(10, 20, 30)
	Process(req)
	fmt.Println(<-req.ret)
}
```
输出结果：
```
60
```
##七、补充 `runtime`包处理`goroutine`的函数

- Goexit

	退出当前执行的goroutine，但是defer函数还会继续调用
- Gosched
	
	让出当前goroutine的执行权限，调度器安排其他等待的任务运行，并在下次某个时候从该位置恢复执行。
- NumCPU

	返回 CPU 核数量
- NumGoroutine

	返回正在执行和排队的任务总数
- GOMAXPROCS

	用来设置可以运行的CPU核数

##八、参考
>1. [Go 学习笔记(雨痕)](https://github.com/qyuhen/book)
>2. [Go Web 编程](https://github.com/astaxie/build-web-application-with-golang)
>3. [Go 语言第一课](http://www.imooc.com/learn/345)
>4.  博文：[进一步认识golang中的并发](http://blog.csdn.net/gophers/article/details/24665419)
https://mikespook.com/2011/09/%E5%8F%8D%E5%B0%84%E7%9A%84%E8%A7%84%E5%88%99/