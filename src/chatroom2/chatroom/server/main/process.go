package main

import (
	"chatroom2/chatroom/common/message"
	"chatroom2/chatroom/utils"
	"fmt"
	"net"
)

type Process struct {
	Conn net.Conn
}

func (this*Process)serverProcessMes(Mes message.Message)  {
	switch Mes.Type {
	case message.LoginMesType:

	}
}
func (this*Process)process2()(err error){
	for{
		tf:=utils.Transfer{
			Conn:this.Conn,
		}
		var Mes message.Message
		Mes,err=tf.ReadPkg()
		if err != nil {
			fmt.Print()
			this.serverProcessMes(Mes)
		}
	}
}