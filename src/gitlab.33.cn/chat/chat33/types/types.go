package types

type Config struct {
	Loglevel string
	LogFile  string
	Server   Server
	Mysql    Mysql
	Api      Api
	Service  Service
	Limit    Limit
	Log      Log
}

type Mysql struct {
	Host string
	Port int32
	User string
	Pwd  string
	Db   string
}

type Server struct {
	Addr string
}

type Service struct {
	DockingDeadline int64
}

type Api struct {
	Zhaobi    string
	RedPacket string
}

type Limit struct {
	ChatPool int
	ChatRate int
}

type Log struct {
	Level string
}

type RedPacket struct {
	PacketId   string `json:"packet_id"`
	PacketUrl  string `json:"packet_url"`
	PacketType int    `json:"packet_type"`
	Coin       int    `json:"coin"`
	Remark     string `json:"remark"`
}

type ReqSendRedPacket struct {
	Amount         string `json:"amount"`
	Size           string `json:"size"`
	Remark         string `json:"remark"`
	Type           string `json:"type"`
	Coin           int    `json:"coin"`
	InvitationCode string `json:"invitation_code"`
}

type PacketQueryParam struct {
	PacketId   string
	AppId      string
	PacketType int
	Coin       int
	Uid        string
	StartTime  int64
	EndTime    int64
	Page       int
	Number     int
}

type RedPacketStatistics struct {
	AdvNum        int `json:"adv_num"`
	TotalNum      int `json:"total_num"`
	TodayAdvNum   int `json:"today_adv_num"`
	TodayTotalNum int `json:"today_total_num"`
}

type RedPacketInfoList struct {
	TotalNum int              `json:"total_num"`
	Packets  []*RedPacketItem `json:"packets"`
}

type RedPacketItem struct {
	PacketId    string `json:"packet_id"`
	PacketType  int    `json:"packet_type"`
	SendUid     string `json:"send_uid"`
	AppUid      string `json:"app_uid"`
	SendAccount string `json:"send_account"`
	Coin        int    `json:"coin"`
	Amount      int    `json:"amount"`
	Size        int    `json:"size"`
	ReceiveNum  int    `json:"receive_num"`
	NewUserNum  int    `json:"newuser_num"`
	BackNum     int    `json:"back_num"`
	Time        int64  `json:"time"` // sec
}

type RedPacketDetail struct {
	Base        RedPacketItem
	RecvDetails []RecvDetail
}

type RecvDetail struct {
	RecvUid     int    `json:"recv_uid"` // chat uid
	AppUid      int    `json:"app_uid"`  // Zhaobi uid
	RecvAccount string `json:"recv_account"`
	RecvType    int    `json:"recv_type"`
	RecvTime    int64  `json:"recv_time"` //sec
	Amount      int    `json:"amount"`
	Avatar      string `json:"avatar"`
}

type Coin struct {
	CoinId   int    `json:"coin_id"`
	CoinName string `json:"coin_name"`
}

type App struct {
	AppId   string `json:"app_id"`
	AppName string `json:"app_name"`
	ImgUrl  string `json:"img_url"`
}

type Statistics struct {
	OnlineUserNum   int `json:"user_num"`
	TodayAdvPackets int `json:"red_envelope_num"`
	GroupNum        int `json:"group_num"`
	CsNum           int `json:"cs_num"`
}
