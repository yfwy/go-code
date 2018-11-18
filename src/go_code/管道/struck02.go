package main

import "fmt"

func main()  {
	var a chan int
	a =make(chan int,10)
	for i:=0;i<4 ;i++  {
		a<-i
	}
	for j := 0;j<4;j++  {
		h := <-a
		fmt.Print(h)
	}
}
