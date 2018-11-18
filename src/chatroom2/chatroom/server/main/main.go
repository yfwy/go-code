package main

import (
	"fmt"
	"net"
)

func process(conn net.Conn)  {
	defer conn.Close()
	process := &Process{
		Conn:conn}
	err := process.process2()
	if err != nil {
		fmt.Println("客户端和服务器通讯协程错误=err", err)
		return
		}
	}
func main() {
	listen,err:=net.Listen("tcp","0.0.0.0:8889")
	defer listen.Close()
	if err != nil {
		fmt.Print("创建端口失败")
	}
	for true {
		fmt.Print("等待链接")
		conn, err:= listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=" ,err)
		}
		go process(conn)
	}
}
