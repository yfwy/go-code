package utils

import (
	"chatroom2/chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)
type Transfer struct {
	//分析它应该有哪些字段
	Conn net.Conn
	Buf [8096]byte //这时传输时，使用缓冲
}
func (this*Transfer)WriterPkg( data []byte) (err error) {
	var pkgLen uint32
	pkgLen =uint32(len(data))
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen)
	n,err:=this.Conn.Write(this.Buf[0:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	n , err =this.Conn.Write(data)
	if err != nil {
		fmt.Print("信息发送失败")
		return 
	}
	return 
}
func (this*Transfer)ReadPkg()(Mes message.Message,err error)  {
	fmt.Println("读取客户端发送的数据...")
	//conn.Read 在conn没有被关闭的情况下，才会阻塞
	//如果客户端关闭了 conn 则，就不会阻塞
	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil {
		//err = errors.New("read pkg header error")
		return
	}
	//根据buf[:4] 转成一个 uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[0:4])
	//根据 pkgLen 读取消息内容
	n, err := this.Conn.Read(this.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		//err = errors.New("read pkg body error")
		return
	}
	//把pkgLen 反序列化成 -> message.Message
	// 技术就是一层窗户纸 &mes！！
	err = json.Unmarshal(this.Buf[:pkgLen], &Mes)
	if err != nil {
		fmt.Println("json.Unmarsha err=", err)
		return
	}
	return
}