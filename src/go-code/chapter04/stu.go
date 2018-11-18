package main

import (
	"fmt"
)

type Person2 struct {
	Name string
	Sex  string
}
type stu struct {
	Person2
}

func (s *stu)function()  {
	fmt.Print("du collage",s.Name,)
}
type nile struct {
	Person2
}
func (n *nile)function2(){
	fmt.Print("du xiaoxue")
}
func main(){
	var c stu
	c.Name="kjjj"
	c.Sex="hhhh"
	c.function()
}