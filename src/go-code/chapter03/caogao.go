package main

import "fmt"

func main(){
	a :=make([]int,5,5)
	fmt.Print(a)
	var b map[string]string
	b =make(map[string]string)
	b["1"]="2"
	fmt.Print(b)
	c:=make(map[string]string)
	c["1"]="11"
	c["2"]="22"
	fmt.Print(c)
}
