package main

import (
	"fmt"
	"net"
)

func main()  {
	conn,err:=net.Dial("tcp","127.0.0.1:8888")
	if err != nil {
		fmt.Print("cuowu",err)
		return
	}else {
		fmt.Printf("成功 ip=%v\n",conn.RemoteAddr().String())
	}
	for true {

		buf:="777"
		n,err:=conn.Write([]byte(buf))
		if err != nil {
			fmt.Print("错误",err)
		}
		fmt.Printf("客户端 %d",n)
		return
	}

}