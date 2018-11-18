package types

const (
	AppIdZhaobi      = "1001"
	AppIdPWalelt     = "1002"
	AppIdYuanScraper = "1003"

	DeviceWeb     = "Web"
	DeviceAndroid = "Android"
	DeviceIOS     = "iOS"

	LevelVisitor = 0
	LevelMember  = 1
	LevelCs      = 2
	LevelAdmin   = 3

	// 获取聊天历史记录条数上限
	ChatHistoryLimit = 1000

	PacketTypeLucky = 1 //拼手气红包
	PacketTypeAdv   = 2 //推广红包

	RecvForOldCustomer = 1 //老用户领取
	RecvForNewCustomer = 2 //未注册用户领取

	PageNo    = 1
	PageLimit = 15

	TimeForever int64 = 3155727600000 // 2069/12/31 23:00:00
)

const (
	EventJoinRoom         = 20
	EventRemoveRoom       = 22
	EventLogOutRoom       = 21
	EventRoomOnlineNumber = 24

	EventCloseGroup  = 4
	EventRemoveGroup = 5
	EventOpenGroup   = 7
)
