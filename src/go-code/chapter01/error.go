package main

import "fmt"

func text (a int ,b int ) int {
	var c int
	c =a/b
	return c
}
func main() {
c :=text(100 , 10)
fmt.Println(c) 
}