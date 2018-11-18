package main

import "fmt"

func main(){
	a:=make(map[int]int)
	a[0]=22
	a[1]=33
	a[2]=44
	for i,j :=range a{
		fmt.Print(i,j)
	}
}
