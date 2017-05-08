package main

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
