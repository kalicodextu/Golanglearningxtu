package main

func main() {
	type User struct {
		id   int
		name string
	}
	m := map[User]int{
		User{1, "Tom"}: 100,
	}
	println(m[User{1, "Tom"}])
}

//func main() {
//	var u1 struct {
//		name string "username"
//	}
//	//var u2 struct{ name string }
//	//u2 = u1 // Error: cannot use u1 (type struct { name string "username" }) as
//	//        type struct { name string } in assignment
//	println(u1.name)
//}
