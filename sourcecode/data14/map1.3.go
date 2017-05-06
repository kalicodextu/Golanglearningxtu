package main

//func main() {
//	type user struct{ name string }
//	m := map[int]user{ // 当 map 因扩张而重新哈希时，各键值项存储位置都会发生改变。 因此，map
//		1: {"user1"}, // 被设计成 not addressable。 类似 m[1].name 这种期望透过原 value
//	} // 指针修改成员的行为自然会被禁止。
//	m[1].name = "user2"
//	println(m[1].name)
//}
//func main() {
//	type user struct{ name string }
//	m := map[int]user{
//		1: {"user1"},
//	}
//	tmp := m[1]
//	tmp.name = "user2"
//	m[1] = tmp
//	println(m[1].name)
//}
//func main() {
//	type user struct{ name string }
//	m := map[int]*user{
//		1: &user{"user1"},
//	}
//	m[1].name = "user2" // 返回的是指针复制品。透过指针修改原对象是允许的。
//	println(m[1].name)
//}

func main() {
	type user struct{ name string }
	m := map[int]user{
		1: {"user1"},
	}
	m[1] = struct {
		name string
	}{"user2"}

	println(m[1].name)
}
