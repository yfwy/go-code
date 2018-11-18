package proto

import (
	"encoding/json"

	l "github.com/inconshreveable/log15"
	"gitlab.33.cn/chat/chat33/db"
	"gitlab.33.cn/chat/chat33/result"
	logic "gitlab.33.cn/chat/chat33/router"
	"gitlab.33.cn/chat/chat33/utility"
)

var proto_log = l.New("module", "chat/proto/proto")

type Proto struct {
	SourceData []byte
	Data       *map[string]interface{}
	Event      int
	MessageId  string
}

func NewProto() *Proto {
	return &Proto{}
}

func (p *Proto) Prase(user *logic.User, device string, msg []byte) *result.Error {
	//返回消息
	var data = make(map[string]interface{})
	err := json.Unmarshal(msg, &data)
	if err != nil {
		//发送的数据格式错误
		proto_log.Error("无法解析接收到的数据", "message", msg)
		return &result.Error{EventCode: 0, ErrorCode: result.MsgFormatError, Message: ""}
	}
	p.SourceData = msg
	p.Data = &data
	proto_log.Debug("receive srouce data:" + string(msg))

	p.Event = utility.ToInt(data["eventType"])
	p.MessageId = utility.ToString(data["msgId"])

	switch p.Event {
	case 0:
	case 1:
		groupId := utility.ToString(data["groupId"])
		//join in group
		channel, errcode := CheckUserIntheGroup(user.Id, groupId, device)
		if errcode == result.GroupNotExists || errcode == result.UserNotExists {
			return &result.Error{EventCode: 1, ErrorCode: errcode, MsgId: p.MessageId, Message: ""}
		}

		if errcode == result.UserNotEnterGroup {
			// // //查询用户是否可加入聊天室
			// canJoin, err := model.CheckUserKickoutById(user.Id)
			// if canJoin || err != nil {
			// 	return &result.Error{EventCode: 1, ErrorCode: model.CannotJoinGroup, MsgId: p.MessageId, Message: ""}
			// }
			user.DeviceRegister(device, channel)
		}
		proto_log.Debug("成功加入聊天室", "group_id", groupId, "user_id", user.Id)
		return &result.Error{EventCode: 1, ErrorCode: result.CodeOK, MsgId: p.MessageId, Message: ""}
	case 2:
		groupId := utility.ToString(data["groupId"])
		//TODO need lock
		//join in group
		channel, errcode := CheckUserIntheGroup(user.Id, groupId, device)
		if errcode != result.CodeOK {
			return &result.Error{EventCode: 2, ErrorCode: result.CodeOK, MsgId: p.MessageId, Message: ""}
		}
		user.DeviceUnRegister(device, channel)
		// sendNoti
		//TODO need unlock
		proto_log.Debug("成功退出聊天室", "group_id", groupId, "user_id", user.Id)
		return &result.Error{EventCode: 2, ErrorCode: result.CodeOK, MsgId: p.MessageId, Message: ""}
	}
	return nil
}

func (p *Proto) GetEvent() int {
	return p.Event
}

//
const (
	TOGROUP = 1
	TOROOM  = 2
	TOUSER  = 3
)

const (
	SYSTEM  = 0
	TEXT    = 1
	AUDIO   = 2
	PHOTP   = 3
	REDPACK = 4
	VIDEO   = 5
)

type ProtoMsg struct {
	LogId            string
	MessageType      int
	Target           int
	TargetId         string
	MessageId        string
	SourceMsgContent *map[string]interface{}
	TargetMsgContent *map[string]interface{}
	TargetMsgData    []byte
	Route            string
}

func NewProtoMsg(data *map[string]interface{}) (*ProtoMsg, *result.Error) {
	var ret = &ProtoMsg{}

	var err *result.Error
	ret.MessageType, err = ret.CheckMessageType(*data)
	channelType := (*data)["channelType"]
	//fromId, ok := (*data)["fromId"]
	targetId, ok := (*data)["targetId"]
	if !ok {
		return ret, &result.Error{EventCode: 0, ErrorCode: result.ParamsError, MsgId: ret.MessageId, Message: ""}
	}
	switch utility.ToInt(channelType) {
	case 1:
		ret.Target = TOGROUP
		ret.TargetId = utility.ToString(targetId)
		ret.Route = logic.GetGroupRouteById(utility.ToString(targetId))
	case 3:
		ret.Target = TOUSER
		ret.TargetId = utility.ToString(targetId)
		ret.Route = ret.TargetId
	case 2:
		ret.Target = TOROOM
		ret.TargetId = utility.ToString(targetId)
		ret.Route = logic.GetRoomRouteById(utility.ToString(targetId))
	}

	return ret, err
}

func (p *ProtoMsg) CheckMessageType(data map[string]interface{}) (int, *result.Error) {
	var msg_id = utility.ToString(data["msgId"])
	p.MessageId = msg_id
	var _msg_type = data["msgType"]
	msg_type := utility.ToInt(_msg_type)

	var _mgs map[string]interface{}
	var ok bool
	if _mgs, ok = data["msg"].(map[string]interface{}); !ok {
		proto_log.Debug("param `msg` error")
		return msg_type, &result.Error{EventCode: 0, ErrorCode: result.ParamsError, MsgId: msg_id, Message: ""}
	}

	switch msg_type {
	case SYSTEM:
		if _, ok = _mgs["content"]; !ok {
			proto_log.Debug("param error")
			return msg_type, &result.Error{EventCode: 0, ErrorCode: result.ParamsError, MsgId: msg_id, Message: ""}
		}
	case TEXT:
		if _, ok = _mgs["content"]; !ok {
			proto_log.Debug("param error")
			return msg_type, &result.Error{EventCode: 0, ErrorCode: result.ParamsError, MsgId: msg_id, Message: ""}
		}
	case AUDIO:
		if _, ok = _mgs["mediaUrl"]; !ok {
			proto_log.Debug("param error")
			return msg_type, &result.Error{EventCode: 0, ErrorCode: result.ParamsError, MsgId: msg_id, Message: ""}
		}
		if _, ok = _mgs["time"]; !ok {
			proto_log.Debug("param error")
			return msg_type, &result.Error{EventCode: 0, ErrorCode: result.ParamsError, MsgId: msg_id, Message: ""}
		}
	case PHOTP:
		if _, ok = _mgs["imageUrl"]; !ok {
			proto_log.Debug("param error")
			return msg_type, &result.Error{EventCode: 0, ErrorCode: result.ParamsError, MsgId: msg_id, Message: ""}
		}
		if _, ok = _mgs["width"]; !ok {
			proto_log.Debug("param error")
			return msg_type, &result.Error{EventCode: 0, ErrorCode: result.ParamsError, MsgId: msg_id, Message: ""}
		}
		if _, ok = _mgs["height"]; !ok {
			proto_log.Debug("param error")
			return msg_type, &result.Error{EventCode: 0, ErrorCode: result.ParamsError, MsgId: msg_id, Message: ""}
		}
	case REDPACK:
	case VIDEO:
	default:
		proto_log.Debug("param error")
		return msg_type, &result.Error{EventCode: 0, ErrorCode: result.ParamsError, MsgId: msg_id, Message: ""}
	}
	p.SourceMsgContent = &data
	return msg_type, nil
}

func (p *ProtoMsg) GetRouter() string {
	return p.Route
}

func (p *ProtoMsg) CheckGroupPush(userId, device string, userLevel int) *result.Error {
	if userLevel == logic.VISITOR {
		//
		proto_log.Debug("visitor can not send message to group")
		return &result.Error{EventCode: 0, ErrorCode: result.VisitorSendMsg, MsgId: p.MessageId, Message: ""}
	}

	_, errcode := CheckUserIntheGroup(userId, p.TargetId, device)
	if errcode != result.CodeOK {
		proto_log.Debug("检查用户是否在群未通过", "gourp id", p.TargetId, "user id", userId)
		return &result.Error{EventCode: 0, ErrorCode: errcode, MsgId: p.MessageId, Message: ""}
	}

	// //查询是否禁言
	// rlt, err := model.CheckUserMuted(userId, utility.NowMillionSecond())
	// if err != nil {
	// 	//数据库错误
	// 	proto_log.Error("查询数据库失败", "err", err.Error())
	// 	return &result.Error{EventCode: 0, ErrorCode: model.DbConnectFail, MsgId: p.MessageId, Message: ""}
	// } else if rlt {
	// 	//用户被禁言
	// 	return &result.Error{EventCode: 0, ErrorCode: model.UserMuted, MsgId: p.MessageId, Message: ""}
	// }
	return nil
}

/**
*
* 返回：用户在聊天室中 true，不在false
**/
func CheckUserIntheGroup(userId, groupId, deviceKey string) (*logic.Channel, int) {
	//get group key
	routeGid := logic.GetGroupRouteById(groupId)
	channel, ok := logic.ChannelMap[routeGid]
	if !ok {
		return nil, result.GroupNotExists
	}
	user, ok := logic.UserMap[userId]
	if !ok || user == nil {
		return channel, result.UserNotExists
	}
	dev, ok := user.Clients[deviceKey]
	if !ok || user == nil {
		return channel, result.NetWorkError
	}
	ok = channel.DeviceIsRegisted(dev)
	if !ok {
		return channel, result.UserNotEnterGroup
	}
	group := user.GetChannel(routeGid)
	if group == nil {
		return channel, result.UserNotEnterGroup
	}
	return group, result.CodeOK
}

func AppendGroupChatLog(userId string, p *ProtoMsg) *result.Error {
	//添加聊天日志 聊天室
	msg := (*p.TargetMsgContent)["msg"]
	var msgStr = utility.StructToString(msg)
	if msgStr == "" {
		proto_log.Error("组成消息体失败", "user id", userId, "group id", p.TargetId)
		return &result.Error{EventCode: 0, ErrorCode: result.ParamsError, MsgId: p.MessageId, Message: ""}
	}
	_, logId, err := db.AppendGroupChatLog(userId, p.TargetId, utility.ToString(p.MessageType), utility.ToString(msgStr), "1")
	if err != nil {
		proto_log.Error("添加群聊日志失败", "user id", userId, "group id", p.TargetId)
		return &result.Error{EventCode: 0, ErrorCode: result.DbConnectFail, MsgId: p.MessageId, Message: ""}
	}
	p.LogId = utility.ToString(logId)
	return nil
}

func AppendRoomChatLog(userId string, p *ProtoMsg) *result.Error {
	//添加聊天日志 群组
	msg := (*p.TargetMsgContent)["msg"]
	var msgStr = utility.StructToString(msg)
	if msgStr == "" {
		proto_log.Error("组成消息体失败", "user id", userId, "group id", p.TargetId)
		return &result.Error{EventCode: 0, ErrorCode: result.ParamsError, MsgId: p.MessageId, Message: ""}
	}
	_, log_id, err := db.AppendRoomChatLog(userId, p.TargetId, utility.ToString(p.MessageType), utility.ToString(msgStr))
	if err != nil {
		proto_log.Error("添加群聊日志失败", "user id", userId, "group id", p.TargetId)
		return &result.Error{EventCode: 0, ErrorCode: result.DbConnectFail, MsgId: p.MessageId, Message: ""}
	}
	p.LogId = utility.ToString(log_id)
	return nil
}

func AppendFriendChatLog(userId string, p *ProtoMsg, state int) *result.Error {
	//添加与好友聊天日志
	msg := (*p.TargetMsgContent)["msg"]
	var msgStr = utility.StructToString(msg)
	if msgStr == "" {
		proto_log.Error("组成消息体失败", "user id", userId, "group id", p.TargetId)
		return &result.Error{EventCode: 0, ErrorCode: result.ParamsError, MsgId: p.MessageId, Message: ""}
	}
	_, log_id, err := db.AddPrivateChatLog(userId, p.TargetId, p.MessageType, utility.ToString(msgStr), state)
	if err != nil {
		proto_log.Error("添加群聊日志失败", "user id", userId, "group id", p.TargetId)
		return &result.Error{EventCode: 0, ErrorCode: result.DbConnectFail, MsgId: p.MessageId, Message: ""}
	}
	p.LogId = utility.ToString(log_id)
	return nil
}

func FormatTargetMsg(userId string, p *ProtoMsg) bool {
	var data = make(map[string]interface{})
	//获取用户信息
	data["eventType"] = 0
	data["msgId"] = p.MessageId
	data["fromId"] = userId
	data["channelType"] = p.Target
	data["targetId"] = p.TargetId
	data["msgType"] = p.MessageType
	data["msg"] = (*p.SourceMsgContent)["msg"]
	data["datetime"] = utility.NowMillionSecond()
	var userInfo map[string]string
	maps, err := db.GetUserInfoWithID(userId)
	if err == nil && len(maps) > 0 {
		userInfo = maps[0]
	}
	var senderInfo = make(map[string]interface{})
	if userInfo == nil {
		//未找到用户
		senderInfo["nickname"] = utility.VisitorNameSplit(userId)
		senderInfo["avatar"] = ""
	} else {
		senderInfo["nickname"] = utility.ToString(userInfo["username"])
		senderInfo["avatar"] = utility.ToString(userInfo["avatar"])
	}
	data["senderInfo"] = senderInfo
	p.TargetMsgContent = &data

	return true
}

func (p *ProtoMsg) ComposeTargetData() bool {
	(*p.TargetMsgContent)["logId"] = p.LogId
	p.TargetMsgData, _ = json.Marshal(*p.TargetMsgContent)
	return true
}
