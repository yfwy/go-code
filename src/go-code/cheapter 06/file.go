package main

import (
	"bufio"
	"fmt"
	"os"
)

func main(){
	file,ok:=os.Open("d:/jj")
	if ok!=nil{
		fmt.Print("kkk")
	}
	defer file.Close()
	read:=bufio.NewReader(file)
	for{
		str,err:=read.ReadString('\n')
		if err!=nil {
			fmt.Print()
			break
		}
		fmt.Print(str)
	}
}