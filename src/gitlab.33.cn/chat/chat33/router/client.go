package router

type Client interface {
	CloseSender()
	Close()
	Send(interface{}) error
	GetSender() chan interface{}
	GetDevice() string
	Binding(*User) bool
}
