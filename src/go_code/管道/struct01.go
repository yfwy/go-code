package main

import (
	"fmt"
	"time"
)

func process()  {
	for i := 0; i < 4;i++  {
		fmt.Println("美国人")
		time.Sleep(time.Second)
	}
}
func map1()  {
	var a map[int]int
	defer func() {
		if err := recover();err!=nil{
			fmt.Print("携程 map 错误")
		}
	}()
	fmt.Println(a)
}
func main()  {
	go process()
	go map1()
	for  i := 0;i<10 ;i++ {
		fmt.Print("中国人")
		time.Sleep(time.Second)
	}
}

