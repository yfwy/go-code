package main

import "fmt"

type A interface {
	Test01()
	Test03()
}
type B interface {
	Test01()
	Test02()
}
type c struct {
	A
	B
}

func main()  {
	var d c
	var A= d
	var B= d
	fmt.Print(A,B)
}