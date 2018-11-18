package router

var (
	ChannelMap = make(map[string]*Channel) // channelId --> Channel
)

// struct of room or personal chan dec
type Channel struct {
	route    string
	UserList map[string]int // userId --> endpoint count
	hub      *Hub
}

func init() {
	ChannelMap["default"] = NewChannel("default")
}

func NewChannel(route string) *Channel {
	return &Channel{
		route:    route,
		hub:      NewHub(),
		UserList: make(map[string]int),
	}
}

func (cl *Channel) Register(clientId string, client Client) {
	// TODO
	cl.UserList[clientId] += 1
	cl.hub.register <- client
}

func (cl *Channel) UnRegister(clientId string, client Client) {
	cl.UserList[clientId] -= 1
	if cl.UserList[clientId] == 0 {
		delete(cl.UserList, clientId)
	}
	cl.hub.unregister <- client
}

func (cl *Channel) Broadcast(msg interface{}) {
	cl.hub.broadcast <- msg
}

func (cl *Channel) GetRegisterNumber() int {
	return len(cl.UserList)
}

func (cl *Channel) DeviceIsRegisted(dev Client) bool {
	return cl.hub.IsExist(dev)
}
