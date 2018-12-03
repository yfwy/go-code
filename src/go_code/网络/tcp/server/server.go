package main

import (
	"fmt"
	"net"
)

func process(conn net.Conn)  {
	defer conn.Close()
	for  {
		buf := make([]byte,100)
		n,err:=conn.Read(buf)
		if err != nil {
			fmt.Print("读取失败")
			return
		}
		fmt.Print(buf[0:n])
	}
}
func main() {
	Listen, err := net.Listen("tcp", "0.0.0.0:8888")
	if err != nil {
		fmt.Print("监听失败")
		return
	}
	for {
		conn, err := Listen.Accept()
		if err != nil {
			fmt.Print("连接失败")
			continue
		}
		go process(conn)
	}
}