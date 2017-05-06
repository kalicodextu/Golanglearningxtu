package main

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
