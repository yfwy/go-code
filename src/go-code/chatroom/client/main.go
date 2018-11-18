package login

import "fmt"

func main() {
	var key int
	var loop bool
	var userID int
	var userPWD int
	for loop {
		fmt.Println("........欢迎来到聊天室........")
		fmt.Println("\t\t\t  1.用户登陆")
		fmt.Println("\t\t\t  2.注册用户")
		fmt.Println("\t\t\t  3.退出系统")
		fmt.Println("\t\t\t  请选择")
		fmt.Scanf("%d\n",&key)
		switch key {
		case 1:
			fmt.Print("输入用户密码")
			loop = false
		case 2:
			fmt.Print("注册" )
			loop = false
		case 3:
			fmt.Print("退出")
			loop = false
		default:
			fmt.Print("输入有误")
		}
		if key == 1 {
			fmt.Print("用户名")
			fmt.Scanln("%d",&userID)
			fmt.Print("密码")
			fmt.Scanln("%d",&userPWD)
			err:=login(userID,userPWD)
			if err!=nil {
				fmt.Print("登陆失败" )
			}else {

			}
		}
	}
}