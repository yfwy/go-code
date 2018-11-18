package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main(){
	fileName:="d:/c/text01.txt"
	file,err:=os.OpenFile(fileName,os.O_RDWR|os.O_APPEND,0666)
	if err!=nil{
		fmt.Print("buxing" )
		return
	}
	reader:=bufio.NewReader(file)
	for{
		str,err:=reader.ReadString('\n')
		if err==io.EOF {
			break
		}
		fmt.Print(str)
	}
	writer:=bufio.NewWriter(file)
	str:="jjjjj"
	for i := 0; i < 5; i++ {
		writer.WriteString(str)
	}
	writer.Flush()
	}