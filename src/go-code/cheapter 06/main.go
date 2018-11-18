package main

import (
	"bufio"
	"fmt"
	"os"
)

func che() {
	fileName := "d:/c/text01.txt"
	file, ok := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if ok != nil {
		fmt.Print()
		return
	}
	defer file.Close()
	str:="hello world\n"
	write:=bufio.NewWriter(file)
	for i:=0;i<5 ;i++  {
		write.WriteString(str)
	}
	write.Flush()
}