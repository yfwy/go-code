package main

import (
	"fmt"
	"net"
)

func process(conn net.Conn)  {
	for {
		fmt.Print("输入")
		buf:=make([]byte,1024)
		n,err:=conn.Read(buf[:])
		if err!=nil {
			fmt.Print("错误")
			return
		}
		fmt.Print(string(buf[:n]))
	}
	defer conn.Close()
}
func main(){
	fmt.Print("开始监听")
	listen,err:=net.Listen("tcp" ,"127.0.0.1:8888")
	if err!=nil{
		fmt.Print("listen err =",err)
		return
	}
	defer listen.Close()
	fmt.Println("listen suc=%v/n",listen)
	for  {
		fmt.Print("d等待连接")
		conn,err:=listen.Accept()
		if err!=nil{
			fmt.Print("shibai =",err)
			return
		}else {
			fmt.Print("conn=",conn )
		}
		go process(conn)
	}
}