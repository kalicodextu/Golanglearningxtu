#Golang学习之路：接口（interface）
##一、前言
　　Go 语言和传统的OO语言概念思想上不同，它在语法上不支持类与集成的概念。但是为了实现类似于C++等语言中多态行为，Go语言引入了`interface`类型，实现了类似于多态的功能。注意`interface`与`method`关系密切，在学习`interface`之前需要将`method`的概念理解清楚。可以参考前面的博文[Golang 学习之路七：面向对象-方法（Method）](http://blog.csdn.net/sevensevensevenday/article/details/71424815),[language specification](http://golang.org/ref/spec#Method_expressions)和[Effective Go](http://golang.org/doc/effective_go.html#methods)
##二、接口的定义
　　接口是一个或多个方法签名的集合，任何类型的方法集中只要拥有与之对应的全部方法，它就表示“实现”了该接口，无需在该类型上显式添加接口声明。而所谓的对应方法，是指有相同的名称、参数列表（不包括参数名）以及返回值。当然，该类型还可以有其它方法。

- 接口的命名通常使用`er`结尾，结构体。
- 接口只有方法签名，没有实现。
- 接口没有数据字段。
- 可以在接口中嵌入其它接口。
- 类型可以实现多个接口
```
package main

import "fmt"

type Stringer interface {
	String() string
}
type Printer interface {
	Stringer // 接口嵌入。
	Print()
}
type User struct {
	id   int
	name string
}

func (self *User) String() string {
	return fmt.Sprintf("user %d, %s", self.id, self.name)
}
func (self *User) Print() {
	fmt.Println(self.String())
}
func main() {
	var t Printer = &User{1, "Tom"} // *User 方法集包含 String、Print。
	t.Print()
}
```
输出结果：
```
user 1, Tom
```
- 空interface(interface{})不包含任何的method，正因为如此，所有的类型都实现了空interface。空interface对于描述起不到任何的作用(因为它不包含任何的method），但是空interface在我们需要存储任意类型的数值的时候相当有用，因为它可以存储任意类型的数值。它有点类似于C语言的void*类型。
```
package main

import "fmt"

func Print(v interface{}) {
	fmt.Printf("%T: %v\n", v, v)
}
func main() {
	Print(1)
	Print("Hello World")
}
```
输出结果：
```
int: 1
string: Hello World
```
- 匿名接口可用作变量类型，或结构成员。
```
package main

import "fmt"

type Tester struct {
	s interface {
		String() string
	}
}
type User struct {
	id   int
	name string
}

func (self *User) String() string {
	return fmt.Sprintf("user %d, %s", self.id, self.name)
}
func main() {
	t := Tester{&User{1, "Tom"}}
	fmt.Println(t.s.String())
}
```
输出结果：
```
user 1, Tom
```
##三、接口的执行机制
　　Go新版使用Go重写，《Go学习笔记》使用了原来C/C++。下面贴出的是分别贴出两个interface部分源码。
```
// @file /Go/src/runtime/runtime2.go (Go 1.8)
type iface struct {
	tab  *itab
	data unsafe.Pointer
}

type itab struct {
	inter  *interfacetype
	_type  *_type
	link   *itab
	bad    int32
	inhash int32      // has this itab been added to hash?
	fun    [1]uintptr // variable sized
}
```

```
// Go 旧版本源码：
struct Iface
{
    Itab* tab;
    void* data;
};
struct Itab
{
    InterfaceType* inter;
    Type* type;
    void (*fun[])(void);
};
```
- 接口表存储元数据信息，包括接口类型、动态类型，以及实现接口的方法指针。

- 数据指针持有的是目标对象的只读复制品，复制完整对象或指针。
```
package main

import "fmt"

type User struct {
	id   int
	name string
}

func main() {
	u := User{1, "Tom"}
	var i interface{} = u
	u.id = 2
	u.name = "Jack"
	fmt.Printf("%v\n", u)
	fmt.Printf("%v\n", i.(User))
}
```
输出结果：
```
{2 Jack}
{1 Tom}
```
- 接口转型返回临时对象，只有使用指针才能修改其状态。
```
package main

import "fmt"

type User struct {
	id   int
	name string
}

func main() {
	u := User{1, "Tom"}
	var i interface{} = &u
	u.id = 2
	u.name = "Jack"
	fmt.Printf("%v\n", u)
	fmt.Printf("%v\n", i.(*User))
}
```
输出结果：
```
{2 Jack}
&{2 Jack}
```
- 只有`tab`和`data`都为`nil`时，接口才等于`nil`。
```
package main // tab = nil, data = nil
import (
	"fmt"
	"reflect"
	"unsafe"
)

var a interface{} = nil
var b interface{} = (*int)(nil) // tab 包含 *int 类型信息, data = nil
type iface struct {
	itab, data uintptr
}

func main() {
	ia := *(*iface)(unsafe.Pointer(&a))
	ib := *(*iface)(unsafe.Pointer(&b))
	fmt.Println(a == nil, ia)
	fmt.Println(b == nil, ib, reflect.ValueOf(b).IsNil())
}
```
输出结果：
```
true {0 0}
false {4701888 0} true
```
##四、接口变量存储的类型
　　`interface`变量里面可以存储任意类型的数值(非接口),只要这个值实现了接口的方法。但是，如何知道变量里面实际保存的是哪种类型的对象？
###1. Comma-ok断言
　　Go语言里面有一个语法，可以直接判断是否是该类型的变量： value, ok =element.(T)，这里value就是变量的值，ok是一个bool类型，element是interface变量，T是断言的类型。
- 如果element里面确实存储了T类型的数值，那么ok返回true，否则返回false。
```
package main

import "fmt"

type User struct {
	id   int
	name string
}

func (self *User) String() string {
	return fmt.Sprintf("%d, %s", self.id, self.name)
}
func main() {
	var o interface{} = &User{1, "Tom"}
	if value, ok := o.(fmt.Stringer); ok { // ok-idiom
		fmt.Println(value, ok)
	}
	u := o.(*User)
	// u := o.(User) // panic: interface is *main.User, not main.User
	fmt.Println(u)
}
```
输出结果：
```
1, Tom true
1, Tom
```
##2. switch测试（不支持fallthrough）
- element.(type)语法不能在switch外的任何逻里面使用，如果你要在switch外面判断一个类型就使用comma-ok。
```
package main

import "fmt"

type User struct {
	id   int
	name string
}

func (self *User) String() string {
	return fmt.Sprintf("%d, %s", self.id, self.name)
}

func main() {
	var element interface{} = &User{1, "Tom"}
	switch v := element.(type) {
	case nil: // element == nil
		fmt.Println("nil")
	case fmt.Stringer: // interface
		fmt.Println(v)
	case func() string: // func
		fmt.Println(v())
	case *User: // *struct
		fmt.Printf("%d, %s\n", v.id, v.name)
	default:
		fmt.Println("unknown")
	}
}
```
###3. 补充
- 超集接口对象可转换为子集接口，反之出错。
```
package main

import "fmt"

type Stringer interface {
	String() string
}
type Printer interface {
	String() string
	Print()
}
type User struct {
	id   int
	name string
}

func (self *User) String() string {
	return fmt.Sprintf("%d, %v", self.id, self.name)
}
func (self *User) Print() {
	fmt.Println(self.String())
}
func main() {
	var o Printer = &User{1, "Tom"}
	var s Stringer = o
	fmt.Println(s.String())
}
```
输出结果：
```
1, Tom
```
##五、嵌入接口
　　Go里面真正吸引人的是他内置的逻辑语法，就像我们在学习Struct时学习的匿名字段，非常的优雅，那么相同的逻辑引入到interface里面，那不是更加完美了。如果一个interface1作为interface2的一个嵌入字段，那么
interface2隐式的包含了interface1里面的method。
- 查看源码包中的一个例子
```
//@file /Go/src/container/heap/heap.go
type Interface interface {
	sort.Interface
	Push(x interface{}) // add x as element Len()
	Pop() interface{}   // remove and return element Len() - 1.
}
```
- 找到sort.Interface的源码如下：
```
//@file /Go/src/sort/sort.go
type Interface interface {
	// Len is the number of elements in the collection.
	Len() int
	// Less reports whether the element with
	// index i should sort before the element with index j.
	Less(i, j int) bool
	// Swap swaps the elements with indexes i and j.
	Swap(i, j int)
}
```
- heap.goInterface嵌入sort。Interface就把sort.Interface的所有method给隐式的包含进来了。

##六、反射简单介绍

　　反射是在运行时反射是程序检查其所拥有的结构，尤其是类型的一种能力；这是元编程的一种形式。它同时也是造成混淆的重要来源。我们一般需要使用reflect包。
使用reflect一般分成三步
- 要去反射是一个类型的值(这些值都实现了空interface)，首先需要把它转化成reflect对象(reflect.Type或者reflect.Value，根据不同的情况调用不同的函数)。这两种获取方式如下：
```
t := reflect.TypeOf(i) //得到类型的元数据,通过t我们能获取类型定义里面的所有元素。
v := reflect.ValueOf(i) //得到实际的值，通过v我们获取存储在里面的值，还可以去改变值
```
- 转化为reflect对象之后我们就可以进行一些操作了，也就是将reflect对象转化成相应的值，例如
```
tag := t.Elem().Field(0).Tag //获取定义在struct里面的标签
name := v.Elem().Field(0).String() //获取存储在第一个字段里面的值
```
获取反射值能返回相应的类型和数值
```
var x float64 = 3.4
v := reflect.ValueOf(x)
fmt.Println("type:", v.Type())
fmt.Println("kind is float64:", v.Kind() == reflect.Float64)
fmt.Println("value:", v.Float())
```
- 最后，反射的话，那么反射的字段必须是可修改的，我们前面学习过传值和传引用，这个里面也是一样的道理，反射的字段必须是可读写的意思是，如果下面这样写，那么会发生错误。
```
var x float64 = 3.4
v := reflect.ValueOf(x)
v.SetFloat(7.1)
```
如果要修改相应的值，需要进行下面的修改
```
var x float64 = 3.4
p := reflect.ValueOf(&x)
v := p.Elem()
v.SetFloat(7.1)
```
上面只是摘录《Go Web编程》对反射的简单介绍，更深入的理解可以查看[laws of reflection](http://golang.org/doc/articles/laws_of_reflection.html),可以查看转载的中文翻译版本[[翻译]反射的规则](http://blog.csdn.net/sevensevensevenday/article/details/72401393)或[原始翻译版本](https://mikespook.com/2011/09/%E5%8F%8D%E5%B0%84%E7%9A%84%E8%A7%84%E5%88%99/)。

##七、总结
　　Go语言的interface是它比较具有特色的地方。Go语言的主要设计者之一罗布·派克（ Rob Pike）曾经说过，如果只能选择一个Go语言的特 性移植到其他语言中，他会选择接口。它在Go开发中无处不在，从Go源码中，你可以看到这些。
##八、参考资料
>1. [Go 学习笔记(雨痕)](https://github.com/qyuhen/book)
>2. [Go Web 编程](https://github.com/astaxie/build-web-application-with-golang)
>3. 博文：[[翻译]反射的规则](https://mikespook.com/2011/09/%E5%8F%8D%E5%B0%84%E7%9A%84%E8%A7%84%E5%88%99/)
https://mikespook.com/2011/09/%E5%8F%8D%E5%B0%84%E7%9A%84%E8%A7%84%E5%88%99/
>4. 博文：[golang技术随笔（一）深入理解interface](http://blog.csdn.net/justaipanda/article/details/43155949)http://blog.csdn.net/justaipanda/article/details/43155949
>5. 博文：[【GoLang笔记】浅析Go语言Interface类型的语法行为及用法](http://blog.csdn.net/slvher/article/details/44492223)
http://blog.csdn.net/slvher/article/details/44492223