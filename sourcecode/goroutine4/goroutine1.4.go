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
	//go Afun(c)
}
