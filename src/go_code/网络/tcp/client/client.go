package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn,err:=net.Dial("tcp" ,"0.0.0.0:8888")
	if err != nil {
		fmt.Print("链接失败")
		return
	}
	defer conn.Close()
	reader:=bufio.NewReader(os.Stdin)
	for {
		line,_:=reader.ReadString('\n')
		conn.Write([]byte(line))
		strings.Trim(line," \n\r")
		if line=="exit" {
			fmt.Print("客户端推出")
			return
		}
		/*buf := make([]byte, 4096)
		for {
			count, err := conn.Read(buf)
			if err != nil {
				break
			}
			fmt.Println(string(buf[0:count]))
		}*/
	}

}
