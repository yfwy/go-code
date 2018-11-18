package process

import (
	"chatroom2/chatroom/common/message"
	"chatroom2/chatroom/utils"
	"encoding/json"
	"fmt"
	"net"
)

type Userprocess struct {

}
func(this*Userprocess)Login(userId int,userPWD int) {
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Print("连接服务器失败")
		return
	}
	defer conn.Close()
	var LoginMes message.LoginMes
	LoginMes.UsePWD = userPWD
	LoginMes.UserID = userId
	var Mes message.Message
	Mes.Type = message.LoginMesType
	data, err := json.Marshal(LoginMes)
	if err != nil {
		fmt.Print("json   false")
	}
	Mes.Data = string(data)
	data, err = json.Marshal(Mes)
	if err != nil {
		fmt.Print("json   false")
		tf := utils.Transfer{}
		tf.WriterPkg(conn, data)
	}
}