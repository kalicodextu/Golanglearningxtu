#Golang 学习：方法
##一、方法的定义

　　Go语言的结构体类型（Struct）比函数类型更加灵活。它可以封装属性和操作。前者即是结构体类型中的字段，而后者则是结构体类型所拥有的方法。
　　在关键字func和名称Grow之间的那个圆括号及其包含的内容就是接收者声明。其中的内容由两部分组成。第一部分是代表它依附的那个类型的值的标识符。第二部分是它依附的那个类型的名称。后者表明了依附关系，而前者则使得在该方法中的代码可以使用到该类型的值（也称为当前值）。代表当前值的那个标识符可被称为接收者标识符，或简称为接收者。
　　这里的接收者指代它所依附的那个类型。我们仍以结构体类型Person为例。下面是依附于它的一个名为Grow的方法的声明：
```
func (person *Person) Grow() {
    person.Age++
} 
```
　　如上所示，在关键字func和名称Grow之间的那个圆括号及其包含的内容就是接收者声明。其中的内容由两部分组成。第一部分是代表它依附的那个类型的值的标识符。第二部分是它依附的那个类型的名称。后者表明了依附关系，而前者则使得在该方法中的代码可以使用到该类型的值（也称为当前值）。代表当前值的那个标识符可被称为接收者标识符，或简称为接收者。请看下面的示例：
```
p := Person{"Robert", "Male", 33}
p.Grow()
```   
　　我们可以直接在Person类型的变量p之上应用调用表达式来调用它的方法Grow。注意，此时方法Grow的接收者标识符person指代的正是变量p的值。这也是“当前值”这个词的由来。在Grow方法中，我们通过使用选择表达式选择了当前值的字段Age，并使其自增。因此，在语句p.Grow()被执行之后，p所代表的那个人就又年长了一岁（p的Age字段的值已变为34）。
　　需要注意的是，在Grow方法的接收者声明中的那个类型是*Person，而不是Person。实际上，前者是后者的指针类型。这也使得person指代的是p的指针，而不是它本身。

- 方法总是绑定对象实例，并隐式将实例作为**第一实参**（receiver）。
    - 只能为当前包内命名类型定义方法。
    - 参数receive可任意命名。如方法中未曾使用，可省略参数名。
    - 参数receive类型可以是T或*T。类型T不能是接口指针。
    - 不支持方法重载，receive只是参数签名的组成部分。
    - 可用实例value或pointer调用全部方法，编译器自动转换。

- 没有构造和析构方法，通常用简单工厂模式返回对象实例。
```
type Queue struct {
	elements []interface{}
}

func NewQueue() *Queue { // 创建对象实例。
	return &Queue{make([]interface{}, 10)}
}
func (*Queue) Push(e interface{}) error { // 省略 receiver 参数名。
	panic("not implemented")
}

// func (Queue) Push(e int) error { // Error: method redeclared: Queue.Push
// panic("not implemented")
// }
func (self *Queue) length() int { // receiver 参数名可以是 self、this 或其他。
	return len(self.elements)
}
```
- 方法不过是一种特殊的函数，只需将其还原，就知道receiver T和*T的差别。
```
type Data struct {
	x int
}

func (self Data) ValueTest() { // func ValueTest(self Data);
	fmt.Printf("Value: %p\n", &self)
}
func (self *Data) PointerTest() { // func PointerTest(self *Data);
	fmt.Printf("Pointer: %p\n", self)
}
func main() {
	d := Data{}
	p := &d
	fmt.Printf("Data: %p\n", p)
	d.ValueTest()   // ValueTest(d)
	d.PointerTest() // PointerTest(&d)
	p.ValueTest()   // ValueTest(*p)
	p.PointerTest() // PointerTest(p)
}
```
输出结果：
```
Data    : 0x1167c0ec
Value   : 0x1167c108
Pointer : 0x1167c0ec
Value   : 0x1167c10c
Pointer : 0x1167c0ec
```
结果分析：
　　从上面的结果可以看出当receiver是T类型时，它是Data的一个copy，而当receiver是*T类型时，就指真正的Data。如果一个method的receiver是*T,你可以在一个T类型的实例变量V上面调用这个method，而不需要&V去调用这个method。如果一个method的receiver是*T,你可以在一个T类型的实例变量V上面调用这个method，而不需要&V去调用这个method类似的如果一个method的receiver是T，你可以在一个*T类型的变量P上面调用这个method，而不需要 *P去调用这个method。所以，你不用担心你是调用的指针的method还是不是指针的method，Go知道你要做的一切。**我真是惊呆了，这就是说Go的设计是需要你能理解why，但是不需要你去做what。**

- 从Go version1.4开始，不再支持多级指针查找方法成员。
```
type X struct{}

func (*X) test() {
	println("X.test")
}
func main() {
	p := &X{}
	p.test()
	// Error: calling method with receiver &p (type **X) requires explicit dereference
	// (&p).test()
}
```
##二、匿名字段
- 可以像字段成员那样访问匿名字段方法，编译器负责查找。
```
type User struct {
	id   int
	name string
}
type Manager struct {
	User
}

func (self *User) ToString() string {
	return fmt.Sprintf("User: %p,%v", self, self)
}

func main() {
	m := Manager{User{1, "Tom"}}
	fmt.Println("Manager: %p\n", &m)
	fmt.Println(m.ToString())
}
```
输出结果：
```
Manager: %p
 &{{1 Tom}}
User: 0x115b2170,&{1 Tom}
```
- 通过匿名的手段，可以获得和继承类似的复用能力。依据编译器的查找次序，只需在外层定义同名方法，就可以实现override。
```
type User struct {
	id   int
	name string
}
type Manager struct {
	User
	title string
}

func (self *User) ToString() string {
	return fmt.Sprintf("User: %p, %v", self, self)
}
func (self *Manager) ToString() string {
	return fmt.Sprintf("Manager: %p, %v", self, self)
}
func main() {
	m := Manager{User{1, "Tom"}, "Administrator"}
	fmt.Println(m.ToString())
	fmt.Println(m.User.ToString())
}
```
输出结果：
```
Manager: 0x1162fb20, &{{1 Tom} Administrator}
User: 0x1162fb20, &{1 Tom}
```
##三、方法集
　　每个类型都有与之关联的方法集，这会影响到接口实现规则。

- 类型T方法集包含全部的receive T方法。
- 类型*T方法集包含全部receive T + *T方法。
- 如类型S包含你名字段T，则S方法集包含T方法。
- 如类型S包含匿名字段*T，则S方法集包含T + *T方法。
- 不管嵌入T或*T，*S方法集总是包含T + *T方法。

用实例value和pointer调用方法（含匿名字段）不受方法集约束，编译器总是查找全部方法，并自动转换receive实参。
##四、表达式
- 根据调用者的不同，方法分为两种表现形式：
```
instance.method(args...)  ---> <type>.func(interface, args...)
```
- 前者称为 method value，后者称为method expression。
- 两者都可像普通函数那样赋值和传参，区别在于method value绑定实例，而 method expression则须显式传参。
```
type User struct {
	id   int
	name string
}

func (self *User) Test() {
	fmt.Printf("%p, %v\n", self, self)
}
func main() {
	u := User{1, "Tom"}
	u.Test()
	mValue := u.Test
	mValue() // 隐式传递 receiver
	mExpression := (*User).Test
	mExpression(&u) // 显式传递 receiver
}
```
输出结果：
```
0x11482170, &{1 Tom}
0x11482170, &{1 Tom}
0x11482170, &{1 Tom}
```
- **注意method value会复制receiver**。
```
type User struct {
	id   int
	name string
}

func (self User) Test() {
	fmt.Println(self)
}
func main() {
	u := User{1, "Tom"}
	mValue := u.Test // 立即复制 receiver，因为不是指针类型，不受后续修改影响。
	u.id, u.name = 2, "Jack"
	u.Test()
	mValue()
}
```
输出结果：
```
{2 Jack}
{1 Tom}
```
- 在汇编层面，method value和闭包的实现方式相同，实际返回FuncVal类型对象。
```
FuncVal { method_address, receiver_copy }
```
- 可根据方法集转换 method expression，注意receive类型的差异。
```
type User struct {
	id   int
	name string
}

func (self *User) TestPointer() {
	fmt.Printf("TestPointer: %p, %v\n", self, self)
}
func (self User) TestValue() {
	fmt.Printf("TestValue: %p, %v\n", &self, self)
}
func main() {
	u := User{1, "Tom"}
	fmt.Printf("User: %p, %v\n", &u, u)
	mv := User.TestValue
	mv(u)
	mp := (*User).TestPointer
	mp(&u)
	mp2 := (*User).TestValue // *User 方法集包含 TestValue。
	mp2(&u)                  // 签名变为 func TestValue(self *User)。
} // 实际依然是 receiver value copy。
```
输出结果：
```
User        : 0x116420e0, {1 Tom}
TestValue   : 0x11642140, {1 Tom}
TestPointer : 0x116420e0, &{1 Tom}
TestValue   : 0x11642170, {1 Tom}
```
- 将方法还原成函数。
```
type Data struct{}

func (Data) TestValue()    {}
func (*Data) TestPointer() {}
func main() {
	var p *Data = nil
	p.TestPointer()
	(*Data)(nil).TestPointer() // method value
	(*Data).TestPointer(nil)   // method expression
	// p.TestValue() // invalid memory address or nil pointer dereference
	// (Data)(nil).TestValue() // cannot convert nil to type Data
	// Data.TestValue(nil) // cannot use nil as type Data in function argument
}
```
##五、补充练习示例
　　为源码文件中声明的结构体类型Person添加相应的字段和方法，使得该文件不会导致任何编译错误并能够在标准输出上打印出`Robert moved from Beijing to San Francisco.`。
```
package main

import "fmt"

type Person struct {
    Name    string
	Gender  string
	Age     uint8
}

func main() {
	p := Person{"Robert", "Male", 33, "Beijing"}
	oldAddress := p.Move("San Francisco")
	fmt.Printf("%s moved from %s to %s.\n", p.Name, oldAddress, p.Address)
}
```
参考代码：
```
package main

import "fmt"

type Person struct {
	Name    string
	Gender  string
	Age     uint8
	Address string
}

func (self *Person) Move(addr string) string {
	a := self.Address
	self.Address = addr
	return a
}
func main() {
	p := Person{"Robert", "Male", 33, "Beijing"}
	oldAddress := p.Move("San Francisco")
	fmt.Printf("%s moved from %s to %s.\n", p.Name, oldAddress, p.Address)
}
```

##六、总结
　　本部分主要介绍了Go面向对象 method，学习了这部分就可以设计一些出基本的面向对象的Go程序。Go里面的面向对象是如此的简单，没有任何的私有、公有关键字，通过大小写来实现(大写开头的为共有，小写开头的为私有)，方法也同样适用这个原则。仔细探究这些用法，愈发感觉Go的设计精妙，它能让开发者开发过程中不忽略关键点，但不许要你为其做多余的，繁杂的动作。

##七、参考资料
>1. [Go 学习笔记(雨痕)](https://github.com/qyuhen/book)
>2. [Go Web 编程](https://github.com/astaxie/build-web-application-with-golang)
>3. [Go语言第一课](http://www.imooc.com/code/7746)
