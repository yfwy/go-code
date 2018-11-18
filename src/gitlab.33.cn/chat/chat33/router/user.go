package router

const (
	VISITOR   = 0
	NOMALUSER = 1
	MANAGER   = 2
)

var (
	UserMap = make(map[string]*User)
)

type User struct {
	Id       string
	Level    int
	Clients  map[string]Client   // device --> Client
	Channels map[string]*Channel // route  --> Channel
}

func NewUser() *User {
	return &User{
		Clients:  make(map[string]Client),
		Channels: make(map[string]*Channel),
	}
}

func (u *User) Add(device string, client Client) bool {
	u.Clients[device] = client
	return true
}

// all client regist current channel
func (u *User) Subscribe(channel *Channel) bool {
	u.Channels[channel.route] = channel
	for _, client := range u.Clients {
		channel.Register(u.Id, client)
	}
	return true
}

func (u *User) UnSubscribe(channel *Channel) bool {
	delete(u.Channels, channel.route)
	for _, client := range u.Clients {
		channel.UnRegister(u.Id, client)
	}
	return true
}

// single device regist channel
func (u *User) DeviceRegister(device string, channel *Channel) bool {
	u.Channels[channel.route] = channel
	if dev, ok := u.Clients[device]; ok {
		channel.Register(u.Id, dev)
		return true
	}
	return false
}

func (u *User) DeviceUnRegister(device string, channel *Channel) bool {
	if _, ok := u.Channels[channel.route]; !ok {
		return false
	}
	if dev, ok := u.Clients[device]; ok {
		channel.UnRegister(u.Id, dev)
		removeChannel := true
		for _, c := range u.Clients {
			if channel.hub.IsExist(c) {
				removeChannel = false
			}
		}
		if removeChannel {
			delete(u.Channels, channel.route)
		}
		return true
	}
	return false
}

func (u *User) Disconnect(client Client) {
	for _, channel := range u.Channels {
		u.DeviceUnRegister(client.GetDevice(), channel)
	}
	client.CloseSender()
	client.Close()
}

func (u *User) GetChannel(route string) *Channel {
	return u.Channels[route]
}

func (u *User) GetChannels() map[string]*Channel {
	return u.Channels
}

func (u *User) SendToAllClients(msg interface{}) {
	for _, c := range u.Clients {
		c.Send(msg)
	}
}

func AppendUser(userId, device string, client Client, level int) (*User, bool) {
	u, ok := UserMap[userId]
	if ok {
		client.Binding(u)
		u.Add(device, client)
	} else {
		u = NewUser()
		u.Id = userId
		u.Level = level
		client.Binding(u)
		u.Add(device, client)
		UserMap[userId] = u
	}
	return u, true
}
