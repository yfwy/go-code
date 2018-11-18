package result

import (
	"encoding/json"
	"strings"

	"github.com/inconshreveable/log15"
	"gitlab.33.cn/chat/chat33/utility"
)

// TODO
const (
	CodeOK                 = 0
	DbConnectFail          = -1000
	ParamsError            = -1001
	LackParam              = -1002
	SessionError           = -1003
	LoginExpired           = -1004
	MsgFormatError         = -1005
	UnknowDeviceType       = -1006
	UserLoginOtherDevice   = -1007
	QueryDbFailed          = -1010
	WriteDbFailed          = -1011
	TooManyRequests        = -1012
	UserNotExists          = -2001
	UserExists             = -2002
	CSExists               = -2003
	ZhaobiTokenLoginFailed = -2004
	UserSendSysMsg         = -2006
	UserMuted              = -2007
	VisitorSendMsg         = -2008
	SendPrivMsg            = -2009
	UserIsServed           = -2010
	JoinListenFailed       = -2011
	ZhaobiInteractFailed   = -2013
	SendPrivMsgBtwCS       = -2014
	VisitorOffline         = -2015
	NoCSOnline             = -2016
	IsFriendAlready        = -2017
	IsNotFriend            = -2018
	FriendRequestHadDeal   = -2019
	NotExistFriendRequest  = -2020
	CanNotOperateSelf      = -2021
	ConvFail               = -2022
	DeleteMsgFailed        = -2023
	PermissionDeny         = -3000
	CannotJoinGroup        = -3001
	NoCSPermission         = -3002
	NoEditGroupPermission  = -3003
	RPError                = -4000
	RPEmpty                = -4001
	UserNotReg             = -4002
	UserHasReg             = -4003
	OnlyForNewUser         = -4004
	RPIdNotMatch           = -4005
	RPIdIllegal            = -4006
	VerifyCodeError        = -4007
	VerifyCodeExpired      = -4008
	RPReceived             = -4009
	CannotSendRP           = -4010
	GroupNotExists         = -5000
	UserNotEnterGroup      = -5001
	QueryChatLogFailed     = -6000
	IsRoomMemberAlready    = -6100
	RoomNotExists          = -6101
	CanNotInvite           = -6102
	CanNotJoinRoom         = -6103
	CanNotLoginOut         = -6104
	UserIsNotInRoom        = -6105
	CanNotAddFriendInRoom  = -6106
	WSMsgFormatError       = -7000
	GetVerifyCodeFailed    = -8000
	ServerInterError       = -9000
	NetWorkError           = -9001
	//....
)

var errorCode = map[int]string{
	CodeOK:                 "操作成功",
	DbConnectFail:          "数据库连接失败",
	ParamsError:            "参数错误",
	LackParam:              "缺少参数",
	SessionError:           "Session错误",
	LoginExpired:           "登录过期",
	MsgFormatError:         "消息格式错误",
	UnknowDeviceType:       "未知的设备类型",
	UserLoginOtherDevice:   "账号已经在其他终端登录",
	QueryDbFailed:          "查询数据库失败",
	WriteDbFailed:          "写入数据库失败",
	TooManyRequests:        "发送频率过快，请稍后再试",
	UserNotExists:          "用户不存在",
	UserExists:             "用户已存在",
	CSExists:               "要添加的客服已存在",
	ZhaobiTokenLoginFailed: "找币token登录失败",
	UserSendSysMsg:         "用户没有发系统消息权限",
	UserMuted:              "用户被禁言",
	VisitorSendMsg:         "游客没有发消息权限",
	SendPrivMsg:            "没有给普通用户发私聊权限",
	UserIsServed:           "当前用户已经被其他客服接待",
	JoinListenFailed:       "加入旁听失败",
	ZhaobiInteractFailed:   "找币交互失败",
	SendPrivMsgBtwCS:       "客服间不能发送私聊消息",
	VisitorOffline:         "游客已离线",
	NoCSOnline:             "暂无客服在线，请稍后再试，或可登录账号给客服留言!",
	PermissionDeny:         "权限不足",
	CannotJoinGroup:        "用户没有加入聊天群权限",
	NoCSPermission:         "没有客服权限",
	NoEditGroupPermission:  "没有修改聊天室权限",
	CanNotInvite:           "不可邀请好友",
	CanNotJoinRoom:         "不可加入该群",
	CanNotLoginOut:         "群主不可退出群",
	RoomNotExists:          "群不存在",
	IsRoomMemberAlready:    "已经是群成员",
	RPError:                "红包内部错误",
	RPEmpty:                "红包已被领完",
	UserNotReg:             "用户未注册",
	UserHasReg:             "用户已注册",
	OnlyForNewUser:         "仅限新人领取",
	RPIdNotMatch:           "红包标识不匹配",
	RPIdIllegal:            "非法的红包ID",
	VerifyCodeError:        "验证码不正确",
	VerifyCodeExpired:      "验证码已经过期或者已使用",
	RPReceived:             "红包已领取",
	CannotSendRP:           "用户无发红包权限",
	GroupNotExists:         "聊天室不存在",
	UserNotEnterGroup:      "用户未进入此聊天室",
	QueryChatLogFailed:     "查询聊天记录失败",
	WSMsgFormatError:       "消息格式错误",
	GetVerifyCodeFailed:    "获取手机验证码失败",
	ServerInterError:       "服务端内部错误",

	IsFriendAlready:       "对方已经是您的好友",
	IsNotFriend:           "对方不是您的好友",
	FriendRequestHadDeal:  "好友请求已经被处理",
	NotExistFriendRequest: "好友请求不存在",
	CanNotOperateSelf:     "不能对自己进行操作",
	ConvFail:              "数据转换异常",
	DeleteMsgFailed:       "删除消息失败",

	NetWorkError:          "网络错误，请重试",
	UserIsNotInRoom:       "用户不在群中",
	CanNotAddFriendInRoom: "该群不允许添加好友",
}

type Empty struct {
}

type Error struct {
	EventCode int
	ErrorCode int
	MsgId     string
	Message   string
}

func ParseError(errcode int, msg string) string {
	//errcode int, msg string
	errStr, ok := errorCode[errcode]
	if !ok {
		log15.Warn("ParseError error code not exists", "errcode", errcode)
		return msg
	}
	return strings.Trim(utility.ParseString(errStr+": %v", msg), " :")
}

func ComposeWsError(errmsg *Error) []byte {
	type WsAck struct {
		EventType int    `json:"eventType"`
		MsgId     string `json:"msgId"`
		Code      int    `json:"code"`
		Content   string `json:"content"`
	}

	var ret WsAck
	ret.EventType = errmsg.EventCode
	ret.MsgId = errmsg.MsgId
	ret.Code = errmsg.ErrorCode
	ret.Content = ParseError(errmsg.ErrorCode, errmsg.Message)

	v, err := json.Marshal(ret)
	if err != nil {
		//log
	}
	return v
}

func ComposeHttpAck(errcode int, errmsg string, data interface{}) interface{} {
	type HttpAck struct {
		Result  int         `json:"result"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
	var ret HttpAck
	ret.Result = errcode
	ret.Message = ParseError(errcode, errmsg)
	ret.Data = data
	return ret
}
