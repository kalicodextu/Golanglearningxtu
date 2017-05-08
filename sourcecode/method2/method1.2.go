package main

import "fmt"

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
	fmt.Printf("Data: %p\n", p) //Data   :0x1167c0ec
	d.ValueTest()               // ValueTest(d)       //Value  :0x1167c108
	d.PointerTest()             // PointerTest(&d)    //Pointer:0x1167c0ec
	p.ValueTest()               // ValueTest(*p)	  //Value  :0x1167c10c
	p.PointerTest()             // PointerTest(p)	  //Pointer:0x1167c0ec
}
