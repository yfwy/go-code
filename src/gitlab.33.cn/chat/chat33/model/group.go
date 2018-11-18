package model

import (
	"encoding/json"
	"regexp"

	l "github.com/inconshreveable/log15"
	"gitlab.33.cn/chat/chat33/db"
	"gitlab.33.cn/chat/chat33/result"
	logic "gitlab.33.cn/chat/chat33/router"
	"gitlab.33.cn/chat/chat33/types"
	"gitlab.33.cn/chat/chat33/utility"
)

var groupModelLog = l.New("module", "chat33/model/group")

type GroupInfo struct {
	GroupId       string `json:"groupId"`
	GroupName     string `json:"groupName"`
	CreateTime    int64  `json:"createTime"`
	Avatar        string `json:"avatar"`
	Description   string `json:"description"`
	OpenTime      int64  `json:"openTime"`
	CloseTime     int64  `json:"closeTime"`
	Status        int    `json:"status"`
	TotalNumber   int    `json:"totalNumber"`
	UserNumber    int    `json:"userNumber"`
	VisitorNumber int    `json:"visitorNumber"`
}

type GroupChatHistory struct {
	LogId       string      `json:"logId"`
	ChannelType string      `json:"channelType"`
	FromId      string      `json:"fromId"`
	TargetId    string      `json:"targetId"`
	MsgType     int         `json:"msgType"`
	Msg         interface{} `json:"msg"`
	Datetime    int64       `json:"datetime"`
	SenderInfo  interface{} `json:"senderInfo"`
}

type GroupUserList struct {
	Id        string `json:"id"`
	Account   string `json:"account"`
	Name      string `json:"name"`
	Uid       string `json:"uid"`
	Avatar    string `json:"avatar"`
	Remark    string `json:"remark"`
	Verified  string `json:"verified"`
	UserLevel int    `json:"userLevel"`
}

type GroupUserInfo struct {
	Id            string `json:"id"`
	Account       string `json:"account"`
	Name          string `json:"name"`
	Uid           string `json:"uid"`
	Avatar        string `json:"avatar"`
	MutedTime     int64  `json:"mutedTime"`
	MutedLastTime int64  `json:"mutedLastTime"`
	Verified      string `json:"verified"`
	UserLevel     int    `json:"userLevel"`
}

type KickOutEvent struct {
	EventType      int   `json:"event_type"`
	MutedTime      int64 `json:"muted_time"`
	MutedLastTime  int64 `json:"muted_last_time"`
	RemoveTime     int64 `json:"remove_time"`
	RemoveLastTime int64 `json:"remove_last_time"`
	Datetime       int64 `json:"datetime"`
}

type GroupEvent struct {
	EventType int    `json:"event_type"`
	GroupId   string `json:"group_id"`
	Datetime  int64  `json:"datetime"`
}

// 将用户踢出聊天室事件
func KickOutUserEvent(userid string) bool {
	devMap, ok := utility.Usermap[userid]
	if ok {
		//发送被踢出消息
		_nowTs := utility.NowMillionSecond()
		mutedAndRemoveInfo, _ := db.GetMuteAndRemoveInfo(userid)
		var ge = &KickOutEvent{EventType: 3, MutedTime: mutedAndRemoveInfo["muted_time"], MutedLastTime: mutedAndRemoveInfo["muted_last_time"], RemoveTime: mutedAndRemoveInfo["remove_time"], RemoveLastTime: mutedAndRemoveInfo["remove_last_time"], Datetime: _nowTs}
		v, _ := json.Marshal(ge)
		for _, _client := range devMap {
			_client.WriteMsg(v)
			_client.ClearGroupId()
		}
		for _, _group := range utility.GroupList {
			delete(_group, userid)
		}
	}
	return true
}

// 将用户禁言事件
func MutedUserEvent(userid string) bool {
	devMap, ok := utility.Usermap[userid]
	if ok {
		//发送被踢出消息
		_nowTs := utility.NowMillionSecond()
		mutedAndRemoveInfo, _ := db.GetMuteAndRemoveInfo(userid)
		var ge = &KickOutEvent{EventType: 3, MutedTime: mutedAndRemoveInfo["muted_time"], MutedLastTime: mutedAndRemoveInfo["muted_last_time"], RemoveTime: mutedAndRemoveInfo["remove_time"], RemoveLastTime: mutedAndRemoveInfo["remove_last_time"], Datetime: _nowTs}
		v, _ := json.Marshal(ge)
		for _, _client := range devMap {
			_client.WriteMsg(v)
		}
	}
	return true
}

// 获取聊天室在线用户 返回值： 总数 用户数 游客数
func GetGroupOnlineNumber(groupId string) (int, int, int) {
	// get group router
	visitorCount := 0
	userCount := 0
	groupRouter := logic.GetGroupRouteById(groupId)
	if channel, ok := logic.ChannelMap[groupRouter]; ok && channel != nil {
		for k := range channel.UserList {
			if user, ok := logic.UserMap[k]; ok {
				if user.Level == logic.VISITOR {
					visitorCount++
				} else if user.Level == logic.NOMALUSER || user.Level == logic.MANAGER {
					userCount++
				}
			}
		}
	}
	return visitorCount + userCount, userCount, visitorCount
}

//-------------------------new version---------------------------//
// 添加聊天日志
func AppendGroupChatLog(senderId, receiveId, msgType, content, logType string) (string, bool) {
	_, logId, err := db.AppendGroupChatLog(senderId, receiveId, msgType, content, logType)
	if err != nil {
		//返回错误
		return "0", false
	}
	return utility.ToString(logId), true
}

// 获取开放的聊天室列表
func GetEnableGroups() ([]int, error) {
	maps, err := db.GetEnableGroups()
	var rlt []int
	if err != nil {
		return rlt, err
	}
	for _, v := range maps {
		rlt = append(rlt, utility.ToInt(v["group_id"]))
	}
	return rlt, nil
}

// 关闭聊天室事件
func CloseGroupEvent(groupId string) {
	channelId := logic.GetGroupRouteById(groupId)
	if _, ok := logic.ChannelMap[channelId]; !ok {
		return
	} else {
		sender := GroupEvent{EventType: types.EventCloseGroup, GroupId: groupId, Datetime: utility.NowMillionSecond()}
		data, err := json.Marshal(sender)
		if err != nil {
			panic(err)
		}

		//发送给所有人
		if cl, ok := logic.ChannelMap["default"]; ok {
			cl.Broadcast(data)
		}
		delete(logic.ChannelMap, channelId)
	}
}

// 删除聊天室事件
func RemoveGroupEvent(groupId string) {
	channelId := logic.GetGroupRouteById(groupId)
	if _, ok := logic.ChannelMap[channelId]; !ok {
		return
	} else {
		sender := GroupEvent{EventType: types.EventRemoveGroup, GroupId: groupId, Datetime: utility.NowMillionSecond()}
		data, err := json.Marshal(sender)
		if err != nil {
			panic(err)
		}
		//发送给所有人
		if cl, ok := logic.ChannelMap["default"]; ok {
			cl.Broadcast(data)
		}
		delete(logic.ChannelMap, channelId)
	}
}

// 开启或者添加聊天室事件
func OpenOrAddGroupEvent(groupId string) {
	newChannelId := logic.GetGroupRouteById(groupId)
	if _, ok := logic.ChannelMap[newChannelId]; !ok {
		//添加聊天室
		logic.ChannelMap[newChannelId] = logic.NewChannel(newChannelId)
	}

	type AppendGroupInfo struct {
		EventType     int    `json:"eventType"`
		GroupId       string `json:"groupId"`
		GroupName     string `json:"groupName"`
		Avatar        string `json:"avatar"`
		Description   string `json:"description"`
		CreateTime    int64  `json:"createTime"`
		OpenTime      int64  `json:"openTime"`
		CloseTime     int64  `json:"closeTime"`
		Status        int    `json:"status"`
		TotalNumber   int    `json:"totalNum"`
		UserNumber    int    `json:"userNum"`
		VisitorNumber int    `json:"visitorNum"`
	}

	groupInfo, errMsg := GetGroupInfo(groupId)
	if errMsg.ErrorCode == result.CodeOK {
		var _item AppendGroupInfo
		_item.EventType = types.EventOpenGroup
		_item.GroupId = groupInfo.GroupId
		_item.GroupName = groupInfo.GroupName
		_item.Description = groupInfo.Description
		_item.Status = groupInfo.Status
		_item.CreateTime = groupInfo.CreateTime
		_item.Avatar = groupInfo.Avatar
		_item.OpenTime = groupInfo.OpenTime
		_item.CloseTime = groupInfo.CloseTime

		_item.VisitorNumber = groupInfo.VisitorNumber
		_item.UserNumber = groupInfo.UserNumber
		_item.TotalNumber = groupInfo.TotalNumber

		data, err := json.Marshal(_item)
		if err != nil {
			return
		}
		//发送给所有人
		if cl, ok := logic.ChannelMap["default"]; ok {
			cl.Broadcast(data)
		}
	}
}

//--------------------------------返回--------------------------------------//
// 添加新的聊天室 group_name:聊天室名称 ;operater:操作者id  []byte 操作结果 error 错误
func AddNewGroup(groupName, avatar, operater string) *result.Error {
	_, num, err := db.AddGroup(groupName, avatar)
	if err != nil {
		//返回错误
		return &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}
	if err == nil && num > 0 {
		OpenOrAddGroupEvent(utility.ToString(num))
	}
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

// 获取聊天室信息列表
func GetGroupInfoList(groupStatus, startTime, endTime int, groupName string) (interface{}, *result.Error) {
	maps, err := db.GetGroupList(groupStatus, startTime, endTime, groupName)
	if err != nil {
		return nil, &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}
	var dataList = make([]GroupInfo, 0)
	for _, val := range maps {
		var _item GroupInfo
		_item.GroupId = utility.ToString(val["group_id"])
		_item.GroupName = utility.ToString(val["group_name"])
		_item.Description = utility.ToString(val["description"])
		_item.Status = utility.ToInt(val["status"])
		_item.CreateTime = utility.ToInt64(val["create_time"])
		_item.Avatar = utility.ToString(val["avatar"])
		_item.OpenTime = utility.ToInt64(val["open_time"])
		_item.CloseTime = utility.ToInt64(val["close_time"])

		totoalCount, userCount, visitorCount := GetGroupOnlineNumber(_item.GroupId)
		_item.VisitorNumber = visitorCount
		_item.UserNumber = userCount
		_item.TotalNumber = totoalCount

		dataList = append(dataList, _item)
	}
	var ret = make(map[string]interface{})
	ret["groups"] = dataList
	return &ret, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

// 获取聊天室详情 group_id:聊天室id
func GetGroupInfo(groupId string) (*GroupInfo, *result.Error) {
	maps, err := db.GetGroupInfo(groupId)
	if err != nil {
		return nil, &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}

	var _item GroupInfo
	if len(maps) > 0 {
		val := maps[0]
		_item.GroupId = utility.ToString(val["group_id"])
		_item.GroupName = utility.ToString(val["group_name"])
		_item.Description = utility.ToString(val["description"])
		_item.Status = utility.ToInt(val["status"])
		_item.CreateTime = utility.ToInt64(val["create_time"])
		_item.Avatar = utility.ToString(val["avatar"])
		_item.OpenTime = utility.ToInt64(val["open_time"])
		_item.CloseTime = utility.ToInt64(val["close_time"])

		totoalCount, userCount, visitorCount := GetGroupOnlineNumber(_item.GroupId)
		_item.VisitorNumber = visitorCount
		_item.UserNumber = userCount
		_item.TotalNumber = totoalCount
	}
	return &_item, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

// 获取聊天室用户列表
func GetGroupUserList(groupId string, params map[string]interface{}) (interface{}, *result.Error) {
	queryUserName := utility.ToString(params["queryUserName"])
	page := utility.ToInt(params["page"])
	number := utility.ToInt(params["number"])

	var userList = make([]GroupUserList, 0)
	//从缓存中找出所有用户
	channelId := logic.GetGroupRouteById(groupId)
	if channel, ok := logic.ChannelMap[channelId]; !ok {
		return nil, &result.Error{ErrorCode: result.GroupNotExists, Message: ""}
	} else {
		for id := range channel.UserList {
			if user, ok := logic.UserMap[id]; !ok {
				continue
			} else {
				switch user.Level {
				case logic.VISITOR:

					var _item GroupUserList
					_item.Id = user.Id
					_item.Uid = ""
					_item.Account = user.Id
					_item.Avatar = ""
					_item.UserLevel = 1
					_item.Verified = "1" //1 否

					if queryUserName != "" {
						r, _ := regexp.Compile("([.]*)" + queryUserName + "([.]*)")
						if !r.MatchString(_item.Name) {
							break
						}
					}
					userList = append(userList, _item)
				case logic.NOMALUSER:
					fallthrough
				case logic.MANAGER:
					// 查出用户详情
					infos, err := db.FindUserInfo(user.Id)
					if err != nil {
						logFriend.Info("FriendInfo query db failed", "err_msg", err)
						return nil, &result.Error{ErrorCode: result.QueryDbFailed}
					}
					if len(infos) >= 1 {
						v := infos[0]
						var _item GroupUserList
						_item.Id = user.Id
						_item.Uid = v["uid"]
						_item.Account = v["account"]
						_item.Name = v["username"]
						_item.Avatar = v["avatar"]
						_item.UserLevel = utility.ToInt(v["user_level"])
						_item.Verified = v["verified"] //1 否 2 是

						if queryUserName != "" {
							r, _ := regexp.Compile("([.]*)" + queryUserName + "([.]*)")
							if !r.MatchString(_item.Name) {
								break
							}
						}
						userList = append(userList, _item)
					}
				default:
				}
			}
		}
	}

	var startIndex int
	if page == 0 {
		startIndex = 0
	} else {
		startIndex = (page - 1) * number
	}

	if number == 0 {
		number = len(userList)
	}

	type _rlt struct {
		UserList interface{} `json:"userList"`
		Totalnum interface{} `json:"totalnum"`
	}
	totalNumb := len(userList[startIndex:number])
	return &_rlt{UserList: userList[startIndex:number], Totalnum: totalNumb}, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

// 获取聊天室用户信息
func GetGroupUserInfo(userId string) (interface{}, *result.Error) {
	var _item GroupUserInfo

	var user *logic.User
	var ok bool
	if user, ok = logic.UserMap[userId]; !ok {
		return nil, &result.Error{ErrorCode: result.UserNotExists, Message: ""}
	}

	switch user.Level {
	case logic.VISITOR:
		_item.Id = user.Id
		_item.Avatar = ""
		_item.Name = user.Id
		_item.Uid = ""
		_item.Account = user.Id
		_item.MutedTime = 0
		_item.MutedLastTime = 0
		_item.UserLevel = 1
		_item.Verified = "1" //1 否
	case logic.NOMALUSER:
		fallthrough
	case logic.MANAGER:
		// 查出用户详情
		infos, err := db.FindUserInfo(user.Id)
		if err != nil {
			logFriend.Info("FriendInfo query db failed", "err_msg", err)
			return nil, &result.Error{ErrorCode: result.QueryDbFailed}
		}
		if len(infos) >= 1 {
			v := infos[0]
			_item.Id = user.Id
			_item.Uid = v["uid"]
			_item.Account = v["account"]
			_item.Name = v["username"]
			_item.Avatar = v["avatar"]
			_item.UserLevel = utility.ToInt(v["user_level"])
			_item.Verified = v["verified"] //1 否 2 是
			_item.MutedTime = 0
			_item.MutedLastTime = 0
		}
	default:
	}
	return _item, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

// 编辑某个聊天室的状态 group_id:聊天室id operate_type:操作类型 1：开启，2：关闭，3：删除
func SetGroupStatus(groupId string, operateType int) *result.Error {
	_, _, err := db.AlterGroupState(groupId, operateType)
	if err != nil {
		//返回错误
		return &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}

	switch operateType {
	case types.GroupOpen:
		OpenOrAddGroupEvent(groupId)
	case types.GroupClose:
		CloseGroupEvent(groupId)
	case types.GroupDelete:
		RemoveGroupEvent(groupId)
	default:
	}

	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

// 编辑某个聊天室的名称
func SetGroupName(groupId, groupName string) *result.Error {
	_, _, err := db.AlterGroupName(groupId, groupName)
	if err != nil {
		//返回错误
		return &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

// 编辑聊天室头像 groupId:聊天室id
func EditGroupAvatar(groupId, avatar string) *result.Error {
	_, _, err := db.UpdateGroupAvatar(groupId, avatar)
	if err != nil {
		//返回错误
		return &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

// 获取聊天室聊天记录 groupId:聊天室id
func GetGroupChatHistory(callerId, groupId string, startid string, number int) (interface{}, *result.Error) {
	chatList := make([]*GroupChatHistory, 0)
	startId := utility.ToInt(startid)
	rows, err := db.GetGroupChatLog(groupId, startId, number)
	if err != nil {
		return nil, &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}
	for _, item := range rows {
		ct := &GroupChatHistory{}
		ct.LogId = utility.ToString(item["id"])
		ct.ChannelType = "1"
		ct.FromId = utility.ToString(item["sender_id"])
		ct.TargetId = utility.ToString(item["receive_id"])
		ct.MsgType = utility.ToInt(item["msg_type"])
		con := utility.StringToJobj(item["content"])
		if ct.MsgType == 4 {
			// 是否已领取
			packetId := utility.ToString(con["packet_id"])
			data, err := RedEnvelopeDetail(packetId)
			if err != nil {
				groupModelLog.Error("red packet detail", "packetId", packetId, "err", err.Error())
				continue
			}

			packetDetail := data.(map[string]interface{})
			var isOpened bool
			con["remark"] = utility.ToString(packetDetail["remark"])
			for _, v := range packetDetail["recv_details"].([]*types.RecvDetail) {
				if callerId == utility.ToString(v.RecvUid) {
					isOpened = true
					break
				}
			}
			con["is_opened"] = isOpened
			con["type"] = utility.ToInt(packetDetail["packet_type"])
		}
		ct.Msg = con
		ct.Datetime = utility.ToInt64(item["send_time"])

		var senderInfo = make(map[string]interface{})
		senderInfo["nickname"] = utility.ToString(item["username"])
		senderInfo["avatar"] = utility.ToString(item["avatar"])
		ct.SenderInfo = senderInfo

		chatList = append(chatList, ct)
	}

	nextLog := "0"
	ret := make(map[string]interface{})
	if len(rows) > number {
		nextLog = chatList[len(chatList)-1].LogId
		chatList = chatList[:len(chatList)-1]
	}
	ret["logs"] = chatList
	ret["nextLog"] = nextLog
	return ret, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}
