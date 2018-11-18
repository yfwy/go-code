package main

import (
	"fmt"
)

type person struct {
	 Name string
	 Age int
	}
	func(p *person)jisuan(){
		if p.Age > 20{
		fmt.Printf("mzw%d",p.Name)
	}else{
		fmt.Printf("shoufei")
	}
	}
	func main(){
	var p person
	p.Name="kkkkk"
	p.Age=20
	p.jisuan()
	}
