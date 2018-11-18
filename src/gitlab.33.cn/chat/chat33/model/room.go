package model

import (
	"encoding/json"

	"gitlab.33.cn/chat/chat33/db"
	"gitlab.33.cn/chat/chat33/proto"
	"gitlab.33.cn/chat/chat33/result"
	logic "gitlab.33.cn/chat/chat33/router"
	"gitlab.33.cn/chat/chat33/types"
	"gitlab.33.cn/chat/chat33/utility"
)

type RoomListInfo struct {
	Id           string `json:"id"`
	MarkId       string `json:"markId"`
	Name         string `json:"name"`
	Avatar       string `json:"avatar"`
	NoDisturbing int    `json:"noDisturbing"`
	CommonlyUsed int    `json:"commonlyUsed"`
	OnTop        int    `json:"onTop"`
}

type RoomInfo struct {
	Id             string         `json:"id"`
	MarkId         string         `json:"markId"`
	Name           string         `json:"name"`
	Avatar         string         `json:"avatar"`
	OnlineNumber   int            `json:"onlineNumber"`
	MemberNumber   int            `json:"memberNumber"`
	NoDisturbing   int            `json:"noDisturbing"`
	MemberLevel    int            `json:"memberLevel"`
	CanAddFriend   int            `json:"canAddFriend"`
	JoinPermission int            `json:"joinPermission"`
	Users          []RoomUserList `json:"users"`
}

type RoomUserList struct {
	Id           string `json:"id"`
	Nickname     string `json:"nickname"`
	RoomNickname string `json:"roomNickname"`
	Avatar       string `json:"avatar"`
	MemberLevel  int    `json:"memberLevel"`
}

type RoomJoinApplyList struct {
	RoomId      string `json:"roomId"`
	RoomName    string `json:"roomName"`
	UserId      string `json:"userId"`
	UserName    string `json:"userName"`
	UserAvatar  string `json:"userAvatar"`
	ApplyReason string `json:"applyReason"`
	Status      int    `json:"status"`
}

type RoomChatLogList struct {
	LogId       string      `json:"logId"`
	ChannelType int         `json:"channelType"`
	FromId      string      `json:"fromId"`
	TargetId    string      `json:"targetId"`
	MsgType     int         `json:"msgType"`
	Msg         interface{} `json:"msg"`
	Datetime    int64       `json:"datetime"`
	SenderInfo  interface{} `json:"senderInfo"`
}

type RoomUnreadStatistics struct {
	Id      string          `json:"id"`
	Number  int             `json:"number"`
	LastLog RoomChatLogList `json:"lastLog"`
}

// 被拉入群通知
type JoinRoomNotification struct {
	EventType int    `json:"eventType"`
	RoomId    string `json:"roomId"`
	UserId    string `json:"userId"`
	Datetime  int64  `json:"datetime"`
}

// 被解散群通知
type RemoveRoomNotification struct {
	EventType int    `json:"eventType"`
	RoomId    string `json:"roomId"`
	Datetime  int64  `json:"datetime"`
}

// 退群通知
type LogOutRoomNotification struct {
	EventType int    `json:"eventType"`
	RoomId    string `json:"roomId"`
	UserId    string `json:"userId"`
	Type      int    `json:"type"`
	Content   string `json:"content"`
}

// 群在线人数通知
type RoomOnlineNumberNotification struct {
	EventType int    `json:"eventType"`
	RoomId    string `json:"roomId"`
	Number    string `json:"number"`
	Datetime  int64  `json:"datetime"`
}

// 发送入群通知
func SendJoinRoomNotification(roomId, uesrId string, members []string) {
	var ret JoinRoomNotification
	ret.EventType = types.EventJoinRoom
	ret.RoomId = roomId
	ret.UserId = uesrId
	ret.Datetime = utility.NowMillionSecond()
	data, _ := json.Marshal(ret)

	for _, memId := range members {
		client := logic.UserMap[memId]
		if client != nil {
			client.SendToAllClients(data)
		}
	}
}

// 发送解散群通知
func SendRemoveRoomNotification(roomId string, members []string) {
	var ret RemoveRoomNotification
	ret.EventType = types.EventRemoveRoom
	ret.RoomId = roomId
	ret.Datetime = utility.NowMillionSecond()
	data, _ := json.Marshal(ret)

	for _, memId := range members {
		client := logic.UserMap[memId]
		if client != nil {
			client.SendToAllClients(data)
		}
	}
}

// 发送退出群通知
func SendLogOutRoomNotification(logOutType int, roomId, userId string, content string, members []string) {
	var ret LogOutRoomNotification
	ret.EventType = types.EventLogOutRoom
	ret.RoomId = roomId
	ret.UserId = userId
	ret.Type = logOutType
	ret.Content = content
	data, _ := json.Marshal(ret)

	for _, memId := range members {
		client := logic.UserMap[memId]
		if client != nil {
			client.SendToAllClients(data)
		}
	}
}

// 发送入群请求和回复通知
func SendApplyNotification(userId string, logId int64, members []string) {
	var ret = make(map[string]interface{})
	ret["eventType"] = types.EventLogOutRoom
	applyInfo, errMsg := GetApplyList(userId, utility.ToString(logId), 1)
	if errMsg.ErrorCode != result.CodeOK {
		return
	}
	if info, ok := applyInfo.(map[string]interface{}); ok {
		if applyList, ok := info["applyList"].([]ApplyList); ok {
			if len(applyList) > 0 && applyList[0].Type == 1 {
				roomApply := applyList[0]
				ret["senderInfo"] = roomApply.SenderInfo
				ret["receiveInfo"] = roomApply.ReceiveInfo
				ret["id"] = roomApply.Id
				ret["type"] = roomApply.Type
				ret["applyReason"] = roomApply.ApplyReason
				ret["status"] = roomApply.Status
				ret["datetime"] = roomApply.Datetime
			}
		}
	}

	data, _ := json.Marshal(ret)

	for _, memId := range members {
		client := logic.UserMap[memId]
		if client != nil {
			client.SendToAllClients(data)
		}
	}
}

func SendAlert(caller, targetId, msg string, target int, members []string) {
	var ret = &proto.ProtoMsg{}
	ret.MessageId = ""
	ret.Target = target
	ret.TargetId = targetId
	ret.MessageType = 6
	var content = make(map[string]interface{})
	content["content"] = msg
	ret.SourceMsgContent = &content

	switch target {
	case proto.TOGROUP:
		ret.Route = logic.GetGroupRouteById(targetId)
		proto.FormatTargetMsg(caller, ret)
		msgErr := proto.AppendGroupChatLog(caller, ret)
		if msgErr != nil {
			return
		}
	case proto.TOROOM:
		ret.Route = logic.GetRoomRouteById(targetId)
		proto.FormatTargetMsg(caller, ret)
		msgErr := proto.AppendRoomChatLog(caller, ret)
		if msgErr != nil {
			return
		}
	case proto.TOUSER:
		ret.Route = targetId
		proto.FormatTargetMsg(caller, ret)
		msgErr := proto.AppendFriendChatLog(caller, ret, 2)
		if msgErr != nil {
			return
		}
	}
	ret.ComposeTargetData()

	for _, memId := range members {
		client := logic.UserMap[memId]
		if client != nil {
			client.SendToAllClients(ret)
		}
	}
}

// 初始化群的channel
func InitRoomChannel(roomId string) {
	channelId := logic.GetGroupRouteById(roomId)
	logic.ChannelMap[channelId] = logic.NewChannel(channelId)
}

func RoomMemberSubscribe(roomId string, members []string) {
	roomChannelId := logic.GetRoomRouteById(roomId)
	cl, ok := logic.ChannelMap[roomChannelId]
	if ok {
		for _, v := range members {
			if user, ok := logic.UserMap[v]; ok {
				user.Subscribe(cl)
			}
		}
	}
}

func RemoveRoomChannel(roomId string) {
	channelId := logic.GetGroupRouteById(roomId)
	delete(logic.ChannelMap, channelId)
}

func RoomMemberUnSubscribe(roomId string, members []string) {
	roomChannelId := logic.GetRoomRouteById(roomId)
	cl, ok := logic.ChannelMap[roomChannelId]
	if ok {
		for _, v := range members {
			if user, ok := logic.UserMap[v]; ok {
				user.UnSubscribe(cl)
			}
		}
	}
}

func JoinInRoom(userId, roomId string) {
	if user, ok := logic.UserMap[userId]; ok && user != nil {
		roomChannelId := logic.GetRoomRouteById(roomId)
		if channel, ok := logic.ChannelMap[roomChannelId]; ok && channel != nil {
			user.Subscribe(channel)
		}
	}
}

// 根据群号获取在线人数
func GetOnlineNumber(roomId string) int {
	channelId := logic.GetRoomRouteById(roomId)
	if channel, ok := logic.ChannelMap[channelId]; ok {
		return channel.GetRegisterNumber()
	}
	return 0
}

func GetEnableRoomIds() ([]string, error) {
	var rlt []string
	list, err := db.GetEnabledRooms()
	if err != nil {
		return rlt, err
	}
	for _, v := range list {
		rlt = append(rlt, utility.ToString(v["id"]))
	}
	return rlt, nil
}

// 获取用户加入的所有群id
func GetUserJoinedRooms(userId string) ([]string, error) {
	rooms, err := db.GetRoomsById(userId)
	if err != nil {
		return nil, err
	}
	data := make([]string, 0)
	for _, v := range rooms {
		roomId := v["room_id"]
		data = append(data, roomId)
	}
	return data, nil
}

// 获取用户所管理的所有群id
func GetUserManageRooms(userId string) ([]string, error) {
	rooms, err := db.GetManageRoomsById(userId)
	if err != nil {
		return nil, err
	}
	data := make([]string, 0)
	for _, v := range rooms {
		roomId := v["room_id"]
		data = append(data, roomId)
	}
	return data, nil
}

// 获取群中所有管理员 返回：key 用户ID value 个推cid
func GetRoomManagerAndMaster(roomId string) (map[string]string, error) {
	var ret = make(map[string]string)
	maps, err := db.GetRoomManagerAndMaster(roomId)
	if err != nil {
		return nil, err
	}
	for _, v := range maps {
		userId := utility.ToString(v["user_id"])
		ret[userId] = utility.ToString(v["level"])
	}
	return ret, nil
}

// 获取群中所有用户 返回：key 用户ID value 个推cid
func GetRoomUsers(roomId string) (map[string]string, error) {
	var ret = make(map[string]string)
	maps, err := db.GetRoomMembers(roomId, -1)
	if err != nil {
		return nil, err
	}
	for _, v := range maps {
		userId := utility.ToString(v["user_id"])
		ret[userId] = utility.ToString(v["getui_cid"])
	}
	return ret, nil
}

// 获取未连接的群成员 返回：key 用户ID value 个推cid
func GetNotStayConnectedRoomUsers(roomId string) map[string]string {
	users, _ := GetRoomUsers(roomId)
	var ret = make(map[string]string)
	for k, v := range users {
		if _, ok := logic.UserMap[k]; !ok {
			ret[k] = v
		}
	}
	return ret
}

// 判断用户是否为管理员或者群主
func CheckUserIsMamnagerOrMaster(roomId, userId string) bool {
	maps, err := db.GetRoomMemberInfo(roomId, userId)
	if err != nil || len(maps) < 1 {
		return false
	}
	userInfo := maps[0]
	return utility.ToInt(userInfo["level"]) > 1
}

// 判断用户是否为群主
func CheckUserIsMaster(roomId, userId string) bool {
	maps, err := db.GetRoomMemberInfo(roomId, userId)
	if err != nil || len(maps) < 1 {
		return false
	}
	userInfo := maps[0]
	return utility.ToInt(userInfo["level"]) > 2
}

// 判断用户是否为群成员
func CheckUserIsMember(roomId, userId string) bool {
	maps, err := db.GetRoomMemberInfo(roomId, userId)
	if err != nil {
		return false
	}
	if len(maps) > 0 {
		return true
	}
	return false
}

//-------------------------api about-------------------------------//
func GetRoomList(queryUser string, Type int) (interface{}, *result.Error) {
	maps, err := db.GetRoomList(queryUser, Type)
	if err != nil {
		return nil, &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}

	var roomList = make([]RoomListInfo, 0)
	for _, r := range maps {
		var room RoomListInfo
		room.Id = utility.ToString(r["id"])
		room.MarkId = utility.ToString(r["mark_id"])
		room.Name = utility.ToString(r["name"])
		room.Avatar = utility.ToString(r["avatar"])
		room.NoDisturbing = utility.ToInt(r["no_disturbing"])
		room.CommonlyUsed = utility.ToInt(r["common_use"])
		room.OnTop = utility.ToInt(r["room_top"])

		roomList = append(roomList, room)
	}

	var ret = make(map[string]interface{})
	ret["roomList"] = roomList
	return ret, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

// 获取群信息
func GetRoomInfo(queryUser, roomId string) (interface{}, *result.Error) {
	var room RoomInfo
	maps, err := db.GetRoomsInfoAsUser(roomId, queryUser)
	if err != nil {
		return nil, &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}

	if len(maps) > 0 {
		r := maps[0]
		room.Id = utility.ToString(r["room_id"])
		room.MarkId = utility.ToString(r["mark_id"])
		room.Name = utility.ToString(r["name"])
		room.Avatar = utility.ToString(r["avatar"])
		room.CanAddFriend = utility.ToInt(r["can_add_friend"])
		room.JoinPermission = utility.ToInt(r["join_permission"])
		room.NoDisturbing = utility.ToInt(r["no_disturbing"])
		room.MemberLevel = utility.ToInt(r["level"])
		room.OnlineNumber = GetOnlineNumber(roomId)
		room.MemberNumber, _ = db.GetRoomMemberNumber(roomId)

		var users = make([]RoomUserList, 0)
		maps, err := db.GetRoomMembers(roomId, 16)
		if err != nil {
			return nil, &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
		}
		for _, r := range maps {
			var user RoomUserList
			user.Id = utility.ToString(r["user_id"])
			groupNickname := utility.ToString(r["user_nickname"])
			if groupNickname == "" {
				user.RoomNickname = utility.ToString(r["username"])
			} else {
				user.RoomNickname = groupNickname
			}
			user.Nickname = utility.ToString(r["username"])
			user.MemberLevel = utility.ToInt(r["level"])
			user.Avatar = utility.ToString(r["avatar"])
			users = append(users, user)
		}
		room.Users = users
		return &room, &result.Error{ErrorCode: result.CodeOK, Message: ""}
	}
	var empty struct{}
	return empty, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func GetRoomUserList(roomId string) (interface{}, *result.Error) {
	var users = make([]RoomUserList, 0)
	maps, err := db.GetRoomMembers(roomId, -1)
	if err != nil {
		return nil, &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}
	for _, r := range maps {
		var user RoomUserList
		user.Id = utility.ToString(r["user_id"])
		groupNickname := utility.ToString(r["user_nickname"])
		if groupNickname == "" {
			user.RoomNickname = utility.ToString(r["username"])
		} else {
			user.RoomNickname = groupNickname
		}
		user.Nickname = utility.ToString(r["username"])
		user.MemberLevel = utility.ToInt(r["level"])
		user.Avatar = utility.ToString(r["avatar"])

		users = append(users, user)
	}

	var ret = make(map[string]interface{})
	ret["userList"] = users
	return ret, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func GetRoomUserInfo(roomId, userId string) (interface{}, *result.Error) {
	maps, err := db.GetRoomMemberInfo(roomId, userId)
	if err != nil || len(maps) < 1 {
		return nil, &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}

	var user RoomUserList
	if len(maps) == 1 {
		r := maps[0]
		user.Id = utility.ToString(r["user_id"])
		groupNickname := utility.ToString(r["user_nickname"])
		if groupNickname == "" {
			user.RoomNickname = utility.ToString(r["username"])
		} else {
			user.RoomNickname = groupNickname
		}
		user.Nickname = utility.ToString(r["username"])
		user.MemberLevel = utility.ToInt(r["level"])
		user.Avatar = utility.ToString(r["avatar"])

		return &user, &result.Error{ErrorCode: result.CodeOK, Message: ""}
	}
	var empty struct{}
	return empty, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func GetRoomSearchMember(roomId, name string) (interface{}, *result.Error) {
	maps, err := db.GetRoomMemberInfoByName(roomId, name)
	if err != nil {
		return nil, &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}

	var infoList = make([]RoomUserList, 0)
	for _, r := range maps {
		var user RoomUserList
		user.Id = utility.ToString(r["user_id"])
		groupNickname := utility.ToString(r["user_nickname"])
		if groupNickname == "" {
			user.RoomNickname = utility.ToString(r["username"])
		} else {
			user.RoomNickname = groupNickname
		}
		user.Nickname = utility.ToString(r["username"])
		user.MemberLevel = utility.ToInt(r["level"])
		user.Avatar = utility.ToString(r["avatar"])

		infoList = append(infoList, user)
	}
	var ret = make(map[string]interface{})
	ret["data"] = infoList
	return &ret, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func GetRoomOnlineNumber(roomId string) (interface{}, *result.Error) {
	var ret = make(map[string]int)
	ret["onlineNumber"] = GetOnlineNumber(roomId)
	return &ret, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func CreateRoom(creater, roomName, roomAvatar string, canAddFriend, joinPermission, adminMuted, masterMuted int, members []string) *result.Error {
	count := 10
	randomId := utility.RandomRoomId()
	for {
		count--
		isExist, err := db.CheckRoomMarkIdExist(randomId)
		if err != nil {
			return &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
		}
		if !isExist {
			break
		} else if count <= 0 {
			return &result.Error{ErrorCode: result.NetWorkError, Message: ""}
		}
	}

	createTime := utility.NowMillionSecond()
	roomId, err := db.CreateRoom(creater, roomName, roomAvatar, canAddFriend, joinPermission, adminMuted, masterMuted, members, randomId, createTime)
	if err != nil {
		return &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}

	//init room channel
	InitRoomChannel(utility.ToString(roomId))
	//member subscribe
	RoomMemberSubscribe(utility.ToString(roomId), members)
	//send notification to all member
	SendJoinRoomNotification(utility.ToString(roomId), creater, members)

	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func RemoveRoom(operator, roomId string) *result.Error {
	//get all member
	maps, _ := GetRoomUsers(roomId)
	var members = make([]string, 0)
	for k := range maps {
		members = append(members, k)
	}

	err := db.DeleteRoomById(roomId)
	if err != nil {
		return &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}
	//send remove room notification to all member
	SendRemoveRoomNotification(roomId, members)
	//member unsubscribe
	RoomMemberUnSubscribe(roomId, members)
	//delete room channel
	RemoveRoomChannel(roomId)
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func LoginOutRoom(operator, roomId string) *result.Error {
	//get all member
	maps, _ := GetRoomUsers(roomId)
	var members = make([]string, 0)
	for k := range maps {
		members = append(members, k)
	}

	err := db.DeleteRoomMemberById(operator, roomId)
	if err != nil {
		return &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}
	//send logOut room notification to all member
	SendLogOutRoomNotification(1, roomId, operator, "", members)
	//current member unsubscribe
	RoomMemberUnSubscribe(roomId, []string{operator})

	var managers = make([]string, 0)
	managerInfos, _ := db.GetRoomManagerAndMaster(roomId)
	for _, k := range managerInfos {
		managers = append(managers, utility.ToString(k["id"]))
	}
	// get user info in the room
	var userRoomNickname string
	user, _ := db.GetRoomMemberInfo(roomId, operator)
	if len(user) > 0 && user[0]["user_nickname"] != "" {
		userRoomNickname = user[0]["user_nickname"]
	} else if len(user) > 0 && user[0]["username"] != "" {
		userRoomNickname = user[0]["username"]
	}
	var msg = userRoomNickname + "退出群聊"
	SendAlert(operator, roomId, msg, proto.TOROOM, managers)
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func KickOutRoom(caller, roomId string, userId string) *result.Error {
	//get all member
	maps, _ := GetRoomUsers(roomId)
	var members = make([]string, 0)
	for k := range maps {
		members = append(members, k)
	}

	// get user info in the room
	var userRoomNickname string
	user, _ := db.GetRoomMemberInfo(roomId, userId)
	if len(user) > 0 && user[0]["user_nickname"] != "" {
		userRoomNickname = user[0]["user_nickname"]
	} else if len(user) > 0 && user[0]["username"] != "" {
		userRoomNickname = user[0]["username"]
	}

	err := db.DeleteRoomMemberById(userId, roomId)
	if err != nil {
		return &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}
	//send logOut room notification to all member
	SendLogOutRoomNotification(1, roomId, userId, "", members)
	//current member unsubscribe
	RoomMemberUnSubscribe(roomId, []string{userId})

	var msg = userRoomNickname + "被移出群聊"
	SendAlert(caller, roomId, msg, proto.TOROOM, members)
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func AdminSetPermission(roomId string, canAddFriend, joinPermission int) *result.Error {
	err := db.AlterRoomCanAddFriendPermission(roomId, canAddFriend)
	if err != nil {
		return &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}

	err = db.AlterRoomJoinPermission(roomId, joinPermission)
	if err != nil {
		return &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func SetLevel(master, userId, roomId string, level int) *result.Error {
	switch level {
	case types.RoomLevelNomal:
		fallthrough
	case types.RoomLevelManager:
		_, _, err := db.SetRoomMemberLevel(userId, roomId, level)
		if err != nil {
			return &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
		}
	case types.RoomLevelMaster:
		err := db.SetNewMaster(master, userId, roomId, level)
		if err != nil {
			return &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
		}
	default:
		return &result.Error{ErrorCode: result.ParamsError, Message: ""}
	}

	//get all member
	maps, _ := GetRoomUsers(roomId)
	var members = make([]string, 0)
	for k := range maps {
		members = append(members, k)
	}
	// get user info in the room
	var userRoomNickname string
	user, _ := db.GetRoomMemberInfo(roomId, userId)
	if len(user) > 0 && user[0]["user_nickname"] != "" {
		userRoomNickname = user[0]["user_nickname"]
	} else if len(user) > 0 && user[0]["username"] != "" {
		userRoomNickname = user[0]["username"]
	}

	var msg string
	if level == types.RoomLevelMaster {
		msg = userRoomNickname + "已被群主设为管理员"
	}
	if level == types.RoomLevelManager {
		msg = userRoomNickname + "成为新的群主"
	}
	SendAlert(master, roomId, msg, proto.TOROOM, members)

	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func SetNoDisturbing(caller, roomId string, noDisturbing int) *result.Error {
	_, _, err := db.SetRoomNoDisturbing(caller, roomId, noDisturbing)
	if err != nil {
		return &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func SetStickyOnTop(caller, roomId string, onTop int) *result.Error {
	_, _, err := db.SetRoomOnTop(caller, roomId, onTop)
	if err != nil {
		return &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func SetMemberNickname(caller, roomId string, nickname string) *result.Error {
	_, _, err := db.SetMemberNickname(caller, roomId, nickname)
	if err != nil {
		return &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func JoinInRoomImiditly(caller, userId, roomId string) *result.Error {
	createTime := utility.NowMillionSecond()
	_, _, err := db.RoomAddMember(userId, roomId, createTime)
	if err != nil {
		return &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}
	JoinInRoom(userId, roomId)

	//get all member
	maps, _ := GetRoomUsers(roomId)
	var members = make([]string, 0)
	for k := range maps {
		members = append(members, k)
	}
	//send notification to all member
	SendJoinRoomNotification(roomId, userId, members)

	// get user info in the room
	var userRoomNickname string
	user, _ := db.GetRoomMemberInfo(roomId, userId)
	if len(user) > 0 && user[0]["user_nickname"] != "" {
		userRoomNickname = user[0]["user_nickname"]
	} else if len(user) > 0 && user[0]["username"] != "" {
		userRoomNickname = user[0]["username"]
	}
	msg := userRoomNickname + "加入群聊"
	SendAlert(caller, roomId, msg, proto.TOROOM, members)
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func JoinRoomInvite(operator, roomId string, users []string) *result.Error {
	//check operator is master or admin
	if CheckUserIsMamnagerOrMaster(roomId, operator) {
		//join in room Imiditly
		for _, userId := range users {
			JoinInRoomImiditly(operator, userId, roomId)
		}
		return &result.Error{ErrorCode: result.CodeOK, Message: ""}
	} else {
		return &result.Error{ErrorCode: result.PermissionDeny, Message: ""}
	}
}

func JoinRoomApply(operator, roomId, applyReason string) *result.Error {
	//get room configuration
	maps, err := db.GetRoomsInfo(roomId)
	if err != nil {
		return &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}
	var joinPermission int
	if len(maps) > 0 {
		r := maps[0]
		joinPermission = utility.ToInt(r["join_permission"])
	} else {
		return &result.Error{ErrorCode: result.RoomNotExists, Message: ""}
	}

	//check operator is room member
	if CheckUserIsMember(roomId, operator) {
		return &result.Error{ErrorCode: result.IsRoomMemberAlready, Message: ""}
	}

	switch joinPermission {
	case types.CanNotJoinRoom:
		return &result.Error{ErrorCode: result.CanNotJoinRoom, Message: ""}
	case types.ShouldApproval:
		_, logid, err := db.AppendJoinRoomApplyLog(roomId, operator, applyReason, types.JoinApplyWait)
		if err != nil {
			return &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
		}
		// get all manager and master
		maps, err := GetRoomManagerAndMaster(roomId)
		if err == nil {
			delete(maps, operator)
			var arry []string
			for k := range maps {
				arry = append(arry, k)
			}
			SendApplyNotification(operator, logid, arry)
		}
	case types.ShouldNotApproval:
		//join in room Imiditly
		errMsg := JoinInRoomImiditly(operator, operator, roomId)
		return errMsg
	}
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func JoinRoomApprove(operator, roomId, userId string, aggre int) *result.Error {
	tx, err := db.GetNewTx()
	if err != nil {
		return &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}
	var members = make([]string, 0)
	var status int
	if aggre == 1 {
		err := db.JoinRoomApproveStepInsert(tx, roomId, userId)
		if err != nil {
			return &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
		}
		JoinInRoom(userId, roomId)
		//get all member
		maps, _ := GetRoomUsers(roomId)
		for k := range maps {
			members = append(members, k)
		}
		//send notification to all member
		SendJoinRoomNotification(roomId, userId, members)
		status = types.JoinApplyOk
	} else {
		status = types.JoinApplyNo
	}
	id, err := db.JoinRoomApproveStepChangeState(tx, roomId, userId, status)
	if err != nil {
		return &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}
	// get all manager and master
	maps, err := GetRoomManagerAndMaster(roomId)
	if err == nil {
		var array []string
		for k := range maps {
			array = append(array, k)
		}
		SendApplyNotification(operator, id, array)
	}
	//
	if status == types.JoinApplyNo {
		// get user info in the room
		var userRoomNickname string
		user, _ := db.GetRoomMemberInfo(roomId, userId)
		if len(user) > 0 && user[0]["user_nickname"] != "" {
			userRoomNickname = user[0]["user_nickname"]
		} else if len(user) > 0 && user[0]["username"] != "" {
			userRoomNickname = user[0]["username"]
		}
		msg := userRoomNickname + "加入群聊"
		SendAlert(operator, roomId, msg, proto.TOROOM, members)
	}
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func GetRoomChatLog(callerId, roomId, startId string, number int) (interface{}, *result.Error) {
	if startId != "" && number < 1 {
		return nil, &result.Error{ErrorCode: result.ParamsError, Message: ""}
	}
	startLogId := utility.ToInt64(startId)
	if startLogId == 0 {
		number = 20
	}
	// 清除未读状态
	_, _, err := db.UpdateReceiveStateReaded(roomId, callerId)
	if err != nil {
		return nil, &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}
	maps, nextLog, err := db.GetChatlog(roomId, startLogId, number)
	if err != nil {
		return nil, &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}

	var list = make([]RoomChatLogList, 0)
	for _, r := range maps {
		var info RoomChatLogList
		info.LogId = utility.ToString(r["id"])
		info.ChannelType = 2
		info.FromId = utility.ToString(r["sender_id"])
		info.TargetId = utility.ToString(r["room_id"])
		info.MsgType = utility.ToInt(r["msg_type"])
		con := utility.StringToJobj(r["content"])
		if info.MsgType == 4 {
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
		info.Msg = con
		info.Datetime = utility.ToInt64(r["datetime"])

		var senderInfo = make(map[string]interface{})
		senderInfo["nickname"] = utility.ToString(r["username"])
		senderInfo["avatar"] = utility.ToString(r["avatar"])
		info.SenderInfo = senderInfo

		list = append(list, info)
	}

	var ret = make(map[string]interface{})
	ret["logs"] = list
	ret["nextLog"] = nextLog
	return ret, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

//群中未读消息统计
func GetRoomUnReadStatistics(callerId string) (interface{}, *result.Error) {
	maps, err := db.GetRoomsUnreadNumber(callerId)
	if err != nil {
		return nil, &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}
	var list = make([]RoomUnreadStatistics, 0)
	for _, r := range maps {
		roomId := r["room_id"]
		count := r["count"]

		var roomStatisticsInfo RoomUnreadStatistics
		roomStatisticsInfo.Id = roomId
		roomStatisticsInfo.Number = utility.ToInt(count)

		logs, _, err := db.GetChatlog(roomId, 0, 1)
		if err != nil {
			return nil, &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
		}

		if len(logs) >= 1 {
			log := logs[0]
			var info RoomChatLogList
			info.LogId = utility.ToString(log["id"])
			info.ChannelType = 2
			info.FromId = utility.ToString(log["sender_id"])
			info.TargetId = utility.ToString(log["room_id"])
			info.MsgType = utility.ToInt(log["msg_type"])
			con := utility.StringToJobj(log["content"])
			if info.MsgType == 4 {
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
			info.Msg = con
			info.Datetime = utility.ToInt64(log["datetime"])

			var senderInfo = make(map[string]interface{})
			senderInfo["nickname"] = utility.ToString(log["username"])
			senderInfo["avatar"] = utility.ToString(log["avatar"])
			info.SenderInfo = senderInfo

			roomStatisticsInfo.LastLog = info
		}
		list = append(list, roomStatisticsInfo)
	}

	var ret = make(map[string]interface{})
	ret["infos"] = list
	return ret, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}
