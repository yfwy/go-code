package main

import "fmt"

func main(){
	var b interface{}
	var a float64 =2.1
	b = a
	 c,ok :=b.(float64)
	if ok {
		fmt.Print(c)

	}
}