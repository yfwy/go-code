package main

import (
	"errors"
	"fmt"
)

type Math struct {
	arr [4] int
	maxSize int
	front int
	rear int
}

func (this *Math) Add (a int)(err error)  {
	if this.rear == this.maxSize-1{
		return  errors.New("cuo wu")
	}
	this.rear++
	this.arr[this.rear]= a
	return
}
func(this*Math) Show(){
	for i:=this.rear+1 ;i<this.front ;i++  {
		fmt.Print(this.arr[i])
	}
}
func (this*Math)get(i int)(val int,err error)  {
	if this.front == this.rear {
		return  1,errors.New("cuo wu" )
	}
	this.front++
	i = this.arr[i]
	return i ,err
}