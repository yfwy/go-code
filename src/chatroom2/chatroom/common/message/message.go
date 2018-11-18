package message
const (
	LoginMesType 			= "LoginMes"
	LoginResMesType			= "LoginResMes"
	RegisterMesType			= "RegisterMes"
	RegisterResMesType 		= "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType				= "SmsMes"
)

type Message struct {
	Type string
	Data string
}
type ReginMes struct {
	UserID int
	UsePWD int
	UserName string
}
type LoginMes struct {
	UserID int
	UsePWD int
}
