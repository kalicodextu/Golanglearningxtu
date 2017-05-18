#Golang 学习：数据
##一、Array
　　Array 就是数组。它的定义方式：
```
var arr [n]type
```

- 在[n]type中，n表示数组长度，type表示存储元素的类型。对数组的操作和其它语言类似，都是通过[]来进行读取或赋值。
　　但Go 中的Array较其它语言有一些差异。
>- 数组是值类型，赋值和传参会复制整个数组，而不是指针。
- 数组长度必须是常量，且是类型的组成部分。[2]int 和 [3]int 是不同的类型。
- 支持 “==”、“！=”操作符，因为内存总是被初始化过的。
- 指针数组[n]*T，数组指针 *[n]T.



 可以用复合语句初始化。
```
a := [3]int{1, 2} // 未初始化元素值为 0。
b := [...]int{1, 2, 3, 4} // 通过初始化值确定数组长度。
c := [5]int{2: 100, 4:200} // 使用索引号初始化元素。
d := [...]struct {
	name string
	age uint8
}{
	{"user1", 10}, // 可省略元素类型。
	{"user2", 20}, // 别忘了最后一行的逗号。
}
```
- 支持多维数组
```
a := [2][3]int{{1, 2, 3}, {4, 5, 6}}
b := [...][2]int{{1, 1}, {2, 2}, {3, 3}} // 第 2 纬度不能用 "..."。
```
- 值拷贝行为会造成性能问题，通常建议使用slice或者数组指针。
```
func test(x [2]int) {
	fmt.Printf("x: %p\n", &x)
	x[1] = 1000
}
func main() {
	a := [2]int{}
	fmt.Printf("a: %p\n", &a)
	test(a)
	fmt.Println(a)
}
```
输出结果：
```
a: 0x114c20e0
x: 0x114c2120
[0 0]  ///值拷贝
```
- 内置函数 len 和 cap 都会返回数组长度（元素数量）。
```
a := [2]int{}
println(len(a), cap(a)) // 2, 2
```
##二、slice
　　slice 有类似动态数组的功能。注意slice和数组在声明时的区别：声明数组时，方括号内写明了数组的长度或使用...自动计算长度，而声明slice时，方括号内没有任何字符。
```
var fslice []int
```
　　但是slice并不是数组或者数组指针。它通过内部指针和相关属性引用数组片段，以实现变长方案。
```
///runtime.h
struct Slice
{ // must not move anything
	byte* array; // actual data
	uintgo len; // number of elements
	uintgo cap; // allocated number of elements
};
```
> - 引用类型。但自身是结构体，值拷贝传递。
- 属性len表示可用元素数量，读写操作不能超过该限制。
- 属性cap表示最大扩张容量，不能超出数组限制。
- 如果slice==nil，那么len、cap结果都是等于0。

- slice可以从一个数组或一个已经存在的slice中再次声明。slice通过array[i:j]来获取，其中i是数组的开始位置，j是结束位置，但不包含array[j]，它的长度是j-i。
```
data := [...]int{0, 1, 2, 3, 4, 5, 6}
slice := data[1:4:5] // [low : high : max]
```
结构如下：
```
		 +- low high-+   +- max 						len = high - low
		 | 			 |   | 								cap = max  - low
	 +---+---+---+---+---+---+---+ 			+---------+---------+---------+
data | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 	  slice | pointer | len = 3 | cap = 4 |
     +---+---+---+---+---+---+---+ 			+---------+---------+---------+
		 |<-- len -->|   |                      |
		 | 			 	 | 						|
		 |<---- cap ---->| 						|
		 | 										|
		 +---<<<- slice.array pointer ---<<<----+
```
- 创建表达式使用的是元素索引号，而非数量。
```
data := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
expression 	 slice 					len    cap	   comment
------------+----------------------+------+-------+---------------------
data[:6:8] 	 [0 1 2 3 4 5]  		6 	   8 	   省略 low.
data[5:] 	 [5 6 7 8 9] 			5 	   5 	   省略 high、max。
data[:3] 	 [0 1 2] 			    3 	   10 	   省略 low、max。
data[:] 	 [0 1 2 3 4 5 6 7 8 9]  10 	   10 	   全部省略。
```
- 读写操作实际目标是底层数组，只需注意索引号的差别。
```
func main() {
	data := [...]int{0, 1, 2, 3, 4, 5}
	s := data[2:4]
	s[0] += 100
	s[1] += 200
	fmt.Println(s)
	fmt.Println(data)
}
```
输出结果：
```
[102 203]
[0 1 102 203 4 5]
```
- 可以直接创建slice对象，自动分配底层数组。
```
func main() {
	s1 := []int{0, 1, 2, 3, 8: 100} // 通过初始化表达式构造，可使用索引号。
	fmt.Println(s1, len(s1), cap(s1))
	s2 := make([]int, 6, 8) // 使用 make 创建，指定 len 和 cap 值。
	fmt.Println(s2, len(s2), cap(s2))
	s3 := make([]int, 6) // 省略 cap，相当于 cap = len。
	fmt.Println(s3, len(s3), cap(s3))
}

```
输出结果:
```
[0 1 2 3 0 0 0 0 100] 9 9
[0 0 0 0 0 0] 6 8
[0 0 0 0 0 0] 6 6
```
- 使用make动态创建slice，避免数组必须必须用常量做长度的麻烦。还可以用指针直接访问底层数组，退化成普通数组操作。
```
s := []int{0, 1, 2, 3}
p := &s[2] // *int, 获取底层数组元素指针。
*p += 100
fmt.Println(s)
```
输出结果：
```
[0 1 102 3]
```
- [][]T是指元素类型为[]T
```
data := [][]int{
	[]int{1, 2, 3},
	[]int{100, 200},
	[]int{11, 22, 33, 44},
}
```
- 可以直接修改 struct array/slice 成员
```
func main() {
	d := [5]struct {
		x int
	}{}
	s := d[:]
	d[1].x = 10
	s[2].x = 20
	fmt.Println(d)
	fmt.Printf("%p, %p, %p\n", &s, &d, &d[0])
}
```
输出结果：
```
[{0} {10} {20} {0} {0}]
0x116820e0, 0x1167f800, 0x1167f800
```
上面结果解释：
三个点：引用 、自身是结构体地址变化 、值拷贝传递

###1. reslice
　　reslice 是基于已有的 slice创建新 slice对象，以便在cap允许范围内调整属性。
```
s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
s1 := s[2:5] // [2 3 4]
s2 := s1[2:6:7] // [4 5 6 7]
s3 := s2[3:6] // Error
```
结构如下：
```
 +---+---+---+---+---+---+---+---+---+---+
data | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 |
	 +---+---+---+---+---+---+---+---+---+---+
	  0       2 		  5
			  +---+---+---+---+---+---+---+---+
s1 			  | 2 | 3 | 4 |   |   |   |   |   | len = 3, cap = 8
			  +---+---+---+---+---+---+---+---+
			   0 	   2 			   6   7
					  +---+---+---+---+---+
s2 					  | 4 | 5 | 6 | 7 |   | len = 4, cap = 5
					  +---+---+---+---+---+
					   0 		   3   4   5
								  +---+---+---+
s3 								  | 7 | 8 | X | error: slice bounds out of range
								  +---+---+---+
```
- 新旧对象依旧指向原底层数组。
```
s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
s1 := s[2:5] // [2 3 4]
s1[2] = 100
s2 := s1[2:6] // [100 5 6 7]
s2[3] = 200
fmt.Println(s)
```
输出结果：
```
[0 1 2 3 100 5 6 200 8 9]
```
###2. append
- 向 slice 尾部添加数据，返回新的 slice 对象。
```
func main() {
	s := make([]int, 0, 5)
	fmt.Printf("%p\n", &s)
	s2 := append(s, 1)
	fmt.Printf("%p\n", &s2)
	fmt.Println(s, s2)
}

```
输出结果：
```
0x11482170
0x114821b0
[] [1]
```
- 简单点说，就是在array[slice.high]写数据。
```
func main() {
	data := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s := data[:3]
	s2 := append(s, 100, 200)
	fmt.Println(data)
	fmt.Println(s)
	fmt.Println(s2)
}
```
输出结果：
```
[0 1 2 100 200 5 6 7 8 9]
[0 1 2]
[0 1 2 100 200]
```
- 一旦超出原slice.cap的限制，就会重新分配底层数组，即便原数组并未填满。
```
func main() {
	data := [...]int{0, 1, 2, 3, 4, 10: 0}
	s := data[:2:3]
	fmt.Println(&s[0], &data[0])
	s = append(s, 100, 200)      // 一次 append 两个值，超出 s.cap 限制。
	fmt.Println(s, data)         // 重新分配底层数组，与原数组无关。
	fmt.Println(&s[0], &data[0]) // 比对底层数组起始指针。
}
```
输出结果：
```
0x11488330 0x11488330
[0 1 100 200] [0 1 2 3 4 0 0 0 0 0 0]
0x1148bb40 0x11488330
```
输出结果分析：

append 后的 s 重新分配了底层数组，并复制了数据。
- 通常以 2 倍的容量重新分配底层数组。在大批量添加数据时，建议一次性分配足够大的空间，以减少内存分配和数据复制开销。或初始化足够长的len属性，改用索引号进行操作。及时释放不再使用的slice对象，避免持有过期数组，造成GC无法回收。
```
func main() {
	s := make([]int, 0, 1)
	c := cap(s)
	for i := 0; i < 50; i++ {
		s = append(s, i)
		if n := cap(s); n > c {
			fmt.Printf("cap: %d -> %d\n", c, n)
			c = n
		}
	}
}
```
输出结果：
```
cap: 1 -> 2
cap: 2 -> 4
cap: 4 -> 8
cap: 8 -> 16
cap: 16 -> 32
cap: 32 -> 64
```
###3. copy
　　函数 copy 在两个slice间赋值数据，复制长度以len小的为准。两个slice可指向同一底层数组，允许元素区间重叠。
```
func main() {
	data := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s := data[8:]
	s2 := data[:5]
	fmt.Println(s, s2)
	copy(s2, s) // dst:s2, src:s
	fmt.Println(s2)
	fmt.Println(data)
}
```
输出结果：
```
[8 9] [0 1 2 3 4]
[8 9 2 3 4]
[8 9 2 3 4 5 6 7 8 9]
```
**应该及时将所需的数据 copy 到较小的slice，以便释放较大的底层数组的内存**

##三、Map
　　map 就是 python 中的字典的概念相同。在Go中Map是引用类型，哈希表。键必须是支持相等运算符（==，！=）类型，比如：number、string、pointer、array、struct，以及对应的interface。值可以是任意类型，没有限制。
```
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
```
输出结果：
```
user1
```
- 预先给make函数一个合理的元素数量参数，有助于提升性能。可以避免以后的扩张操作。
```
m := make(map[string]int, 1000)
```
常见的操作：
```
func main() {
	m := map[string]int{
		"a": 1,
	}
	if v, ok := m["a"]; ok { // 判断 key 是否存在。
		println(v, ok)
	}
	println(m["c"])       // 对于不存在的 key，直接返回 \0，不会出错。
	m["b"] = 2            // 新增或修改。
	delete(m, "c")        // 删除。如果 key 不存在，不会出错。
	println(len(m))       // 获取键值对数量。cap 无效。
	for k, v := range m { // 迭代，可仅返回 key。随机顺序返回，每次都不相同。
		println(k, v)
	}
}
```
输出结果：
```
1 true
0
2
a 1
b 2
```
不能保证迭代顺序，这些适合 python 特性相同
- 从map中取回的是一个value临时复制品,对其成员修改是没有任何意义的。
```
func main() {
	type user struct{ name string }
	m := map[int]user{ // 当 map 因扩张而重新哈希时，各键值项存储位置都会发生改变。 因此，map
		1: {"user1"}, // 被设计成 not addressable。 类似 m[1].name 这种期望透过原 value
	} // 指针修改成员的行为自然会被禁止。
	m[1].name = "user2"
	println(m[1].name)
}
```
运行结果：
```
cannot assign to struct field m[1].name in map
```
- 可以完整替换value或使用指针达到目的。
```
// 完整代替 value
func main() {
	type user struct{ name string }
	m := map[int]user{
		1: {"user1"},
	}
	tmp := m[1]
	tmp.name = "user2"
	m[1] = tmp
	println(m[1].name)
}
```
```
// 使用指针
func main() {
	type user struct{ name string }
	m := map[int]*user{
		1: &user{"user1"},
	}
	m[1].name = "user2" // 返回的是指针复制品。透过指针修改原对象是允许的。
	println(m[1].name)
}
```
- 允许在迭代期间安全删除键值。但是如果有新增操作,意外不可知预料。
```
func main() {
	for i := 0; i < 5; i++ {
		m := map[int]string{
			0: "a", 1: "a", 2: "a", 3: "a", 4: "a",
			5: "a", 6: "a", 7: "a", 8: "a", 9: "a",
		}
		for k := range m {
			m[k+k] = "x"
			delete(m, k)
		}
		fmt.Println(m)
	}
}
```
输出结果：
```
map[14:x 8:x 18:x 32:x 12:x 20:x]
map[10:x 2:x 14:x 16:x 18:x 24:x 8:x]
map[12:x 14:x 10:x 18:x 32:x 4:x]
map[4:x 6:x 12:x 20:x 16:x 28:x 18:x]
map[8:x 12:x 18:x 10:x 14:x 16:x 2:x]
```
##四、Struct
　　Go语言中，也和C或者其他语言一样，我们可以声明新的类型，作为其它类型的属性或字段的容器。struct是值类型,赋值和传参会复制全部内容.可以用"_"定义补位字段,支持指向自身的类型的指针成员。
```
type Node struct {
	_    int
	id   int
	data *byte
	next *Node
}

func main() {
	n1 := Node{
		id:   1,
		data: nil,
	}
	n2 := Node{
		id:   2,
		data: nil,
		next: &n1,
	}
}
```
- 顺序初始化必须包含全部字段，否则会报错。
```
type User struct {
    name string
    age int
}
u1 := User{"Tom", 20}
u2 := User{"Tom"} // Error: too few values in struct initializer
```
- 支持匿名结构,可用作结构成员或定义变量。
```
type File struct {
    name string
    size int
    attr struct {
        perm int
        owner int
    }
}
f := File{
    name: "test.txt",
    size: 1025,
    // attr: {0755, 1}, // Error: missing type in composite literal
}
f.attr.owner = 1
f.attr.perm = 0755
var attr = struct {
    perm int
    owner int
}{2, 0755}
f.attr = attr
```
- 支持"=="、"!="相等操作符,可用作map键类型。
```
func main() {
	type User struct {
		id   int
		name string
	}
	m := map[User]int{
		User{1, "Tom"}: 100,
	}
	println(m[User{1,"Tom"}]) // 输出 100
}
```
- **可以定义字段标签,用反射读取。标签是类型的组成部分**。
```
var u1 struct { name string "username" } //"username" 标签
var u2 struct { name string }
u2 = u1 // Error: cannot use u1 (type struct { name string "username" }) as
        //        type struct { name string } in assignment
```
- 空结构”节省“内存：实现 set 数据结构,或者实现没有状态只有方法的“静态类”。
```
var null struct{}
set := make(map[string]struct{})
set["a"] = null
```
###1. 匿名字段
- 匿名字段不过是一种语法糖，从根本上说，就是一个与成员类型同名（不包含包名）的字段。被匿名嵌入的可以是任意类型，也可以是指针。
```
type User struct {
	name string
}
type Manager struct {
	User
	title string
}
m := Manager{
	User: User{"Tom"}, // 匿名字段的显式字段名，和类型名相同。
	title: "Administrator",
}
```
- 可以像普通字段那样访问匿名字段成员，编译器从外向内逐级查找所有的匿名字段，直到发现目标或出错。
```
ype Resource struct {
	id int
}
type User struct {
	Resource
	name string
}
type Manager struct {
	User
	title string
}
var m Manager
m.id = 1
m.name = "Jack"
m.title = "Administrator"
```
- 外层同名字段会遮蔽嵌入字段成员，相同层次的同名字段也会让编译器⽆所适从。解决⽅法是使用显式字段名。
```
type Resource struct {
	id int
	name string
}
type Classify struct {
	id int
}
type User struct {
	Resource // Resource.id 与 Classify.id 处于同⼀一层次。
	Classify
	name string // 遮蔽 Resource.name。
}
u := User{
	Resource{1, "people"},
	Classify{100},
	"Jack",
}
println(u.name) // User.name: Jack
println(u.Resource.name) // people
// println(u.id) // Error: ambiguous selector u.id
println(u.Classify.id) // 100
```
- 不能同时嵌入某一类型和其指针类型，因为它们名字相同。
```
type Resource struct {
	id int
}
type User struct {
	*Resource
	// Resource // Error: duplicate field Resource
	name string
}
u := User{
	&Resource{1},
	"Administrator",
}
println(u.id)
println(u.Resource.id)
```
###2. 面向对象
- 面向对象三大特征里，Go 仅支封装，尽管匿名字段的内存布局和行为类似继承。没有class 关键字，没有继承、多态等等。
```
type User struct {
	id int
	name string
}
type Manager struct {
	User
	title string
}
m := Manager{User{1, "Tom"}, "Administrator"}
// var u User = m // Error: cannot use m (type Manager) as type User in assignment
// 没有继承，自然也不会有多态。
var u User = m.User // 同类型拷贝。
```
- 内存布局和 C struct 相同，没有任何附加的 object 信息。
```
	|<-------- User:24 ------->|<-- title:16 -->|
	+--------+-----------+------------+ 		+---------------+
m 	|    1   |   string  |   string   | 		| Administrator | 	[n]byte
	+--------+-----------+------------+ 		+---------------+
				  | 			| 						|
				  | +--->>>------------------>>>--------+
				  |
				  +--->>>-------------------->>>-----+
													 |
				  +--->>>-------------------->>>-+   |
				  | 							 | 	 	|
	+--------+-----------+ 						+---------+
u   |   1    |   string  | 						|   Tom   | [n]byte
	+--------+-----------+ 						+---------+
	|<-id:8->|<-name:16->|
```
- 可用 unsafe 包相关函数输出内存地址信息。
```
m : 0x2102271b0, size: 40, align: 8
m.id : 0x2102271b0, offset: 0
m.name : 0x2102271b8, offset: 8
m.title: 0x2102271c8, offset: 24
```
##五、总结
　　本部分主要介绍了 Go 数据：array、slice、map、struct。从介绍中，基本可以看到Go与C语言的一些相同之处，也可以看到它从现代语言中提取的优点。学习、记录、灵活使用！
##六、参考资料
>1. [Go 学习笔记(雨痕)](https://github.com/qyuhen/book)
>2. [Go Web 编程](https://github.com/astaxie/build-web-application-with-golang)