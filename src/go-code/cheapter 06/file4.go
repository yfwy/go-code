package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type number struct {
	enNum int
	szNum int
	cnNum int
}
func main(){
	fileName:="D:/we/新建文件夹/text01.txt"
	file ,err:=os.Open(fileName)
	if err != nil {
		fmt.Print("youwu")
		return
	}
	defer file.Close()
	reader:=bufio.NewReader(file)
	for{
	a,err:=reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		for _,v:=range a{
			fmt.Println(v)
		}
	}
}
