package model

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/inconshreveable/log15"
	"gitlab.33.cn/chat/chat33/db"
	"gitlab.33.cn/chat/chat33/proto"
	"gitlab.33.cn/chat/chat33/result"
	logic "gitlab.33.cn/chat/chat33/router"
	"gitlab.33.cn/chat/chat33/types"
	"gitlab.33.cn/chat/chat33/utility"
)

var logFriend = log15.New("module", "model/friend")

//判断用户是否存在
func boolIsUserExist(id string) *result.Error {
	bool, err := db.UserIsExists(id)
	if err != nil {
		logFriend.Warn("CheckFriend  db failed", "err_msg", err)
		return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}

	if !bool {
		return &result.Error{ErrorCode: result.UserNotExists, Message: ""}
	}

	return nil
}

//判断用户是否存在  是否是好友  对好友的一些操作都需要先判断
func boolIsExistIsFriend(userID, friendID string) *result.Error {
	bool, err := db.UserIsExists(friendID)
	if err != nil {
		logFriend.Warn("CheckFriend  db failed", "err_msg", err)
		return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}

	if !bool {
		return &result.Error{ErrorCode: result.UserNotExists, Message: ""}
	}

	bool, err = db.CheckFriend(userID, friendID, types.FriendIsNotDelete)
	if err != nil {
		logFriend.Warn("CheckFriend  db failed", "err_msg", err)
		return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}

	if !bool {
		return &result.Error{ErrorCode: result.IsNotFriend, Message: ""}
	}

	return nil
}

/*
	好友列表
*/
func FriendList(userID string, tp, time int) (interface{}, *result.Error) {
	if time == 0 {
		friends, err := db.GetFriendList(userID, tp, types.FriendIsNotDelete)
		if err != nil {
			logFriend.Warn("friend list query db failed", "err_msg", err)
			return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
		}
		data := make([]map[string]interface{}, 0)

		for _, friend := range friends {
			one := make(map[string]interface{})

			one["id"] = friend["friend_id"]
			one["name"] = friend["username"]
			one["avatar"] = friend["avatar"]
			one["position"] = friend["position"]
			one["remark"] = friend["remark"]
			dnd := utility.ToInt(friend["DND"])
			one["noDisturbing"] = dnd
			commonlyUsed := utility.ToInt(friend["type"])
			one["commonlyUsed"] = commonlyUsed
			top := utility.ToInt(friend["top"])
			one["onTop"] = top
			isDelete := utility.ToInt(friend["is_delete"])
			one["isDelete"] = isDelete
			addTime, err := strconv.Atoi(friend["add_time"])
			if err != nil {
				return nil, &result.Error{ErrorCode: result.ConvFail, Message: ""}
			}
			one["addTime"] = addTime
			data = append(data, one)
		}
		return map[string]interface{}{"userList": data}, &result.Error{ErrorCode: result.CodeOK, Message: ""}
	} else {
		friends, err := db.GetFriendListByTime(userID, tp, time)
		if err != nil {
			logFriend.Warn("friend list query db failed", "err_msg", err)
			return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
		}
		data := make([]map[string]interface{}, 0)

		for _, friend := range friends {
			one := make(map[string]interface{})

			one["id"] = friend["friend_id"]
			one["name"] = friend["username"]
			one["avatar"] = friend["avatar"]
			one["position"] = friend["position"]
			one["remark"] = friend["remark"]

			//one["noDisturbing"] = friend["DND"]
			//one["commonlyUsed"] = friend["type"]
			//one["onTop"] = friend["top"]
			//one["is_delete"] = friend["is_delete"]
			dnd := utility.ToInt(friend["DND"])
			one["noDisturbing"] = dnd
			commonlyUsed := utility.ToInt(friend["type"])
			one["commonlyUsed"] = commonlyUsed
			top := utility.ToInt(friend["top"])
			one["onTop"] = top
			is_delete := utility.ToInt(friend["is_delete"])
			one["is_delete"] = is_delete
			addTime, err := strconv.Atoi(friend["add_time"])
			if err != nil {
				return nil, &result.Error{ErrorCode: result.ConvFail, Message: ""}
			}
			one["add_time"] = addTime
			data = append(data, one)
		}
		return map[string]interface{}{"userList": data}, &result.Error{ErrorCode: result.CodeOK, Message: ""}
	}
}

//添加好友
func AddFriend(userID, friendID, remark, reason, rId string) *result.Error {

	//不能对自己操作
	if userID == friendID {
		return &result.Error{ErrorCode: result.CanNotOperateSelf, Message: ""}
	}

	//判断用户是否存在
	isExists, err := db.UserIsExists(friendID)
	if !isExists {
		if err != nil {
			logFriend.Warn("UserIsExists db failed", "err_msg", err)
			return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
		}
		return &result.Error{ErrorCode: result.UserNotExists, Message: ""}
	}

	//验证是否是好友
	isFriend, err := db.CheckFriend(userID, friendID, types.FriendIsNotDelete)
	if isFriend {
		if err != nil {
			logFriend.Warn("CheckFriend db dailed", "err_msg", err)
			return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
		}
		return &result.Error{ErrorCode: result.IsFriendAlready, Message: ""}
	}

	roomId, err := strconv.Atoi(rId)
	isFromRoom := true
	//不是通过群添加好友
	if err != nil || roomId == 0 {
		isFromRoom = false
		roomId = 0
	}

	if isFromRoom {
		//判断群是否存在
		roomIsExist, err := db.CheckRoomIsExist(roomId)
		if err != nil {
			logFriend.Error("db CheckRoomIsExist", err)
			return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
		}
		if !roomIsExist {
			return &result.Error{ErrorCode: result.RoomNotExists, Message: ""}
		}

		//判断该群是否允许添加好友
		canAdd, err := db.CheckRoomIsCanAddFriend(roomId)
		if err != nil {
			logFriend.Error("db CheckRoomIsExist", err)
			return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
		}
		if !canAdd {
			return &result.Error{ErrorCode: result.CanNotAddFriendInRoom, Message: ""}
		}

		//判断双方是否都在群里
		userInRoom1, err := db.CheckUserInRoom(userID, roomId)
		if err != nil {
			logFriend.Error("db CheckUserInRoom", err)
			return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
		}
		userInRoom2, err := db.CheckUserInRoom(friendID, roomId)
		if err != nil {
			logFriend.Error("db CheckUserInRoom", err)
			return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
		}
		if !userInRoom1 || !userInRoom2 {
			return &result.Error{ErrorCode: result.UserIsNotInRoom, Message: ""}
		}
	}

	//查询cid
	cid, err := db.FindCid(friendID)
	if err != nil {
		logFriend.Warn("FindCid db dailed", "err_msg", err)
		return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}
	//查询自己的信息
	info, err := db.FindUserInfo(userID)
	if err != nil {
		logFriend.Warn("FindUserInfo db dailed", "err_msg", err)
		return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}

	//查询对方的信息
	infoFriend, err := db.FindUserInfo(friendID)
	if err != nil {
		logFriend.Warn("FindUserInfo db dailed", "err_msg", err)
		return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}

	//判断请求是否存在
	cou, err := db.FindApplyCount(userID, friendID)
	if err != nil {
		logFriend.Warn("FindApplyCount db dailed", "err_msg", err)
		return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}

	//封装websock消息
	var msgMap = make(map[string]interface{})
	msgMap["eventType"] = 31

	var senderInfo = make(map[string]interface{})
	senderInfo["id"] = info[0]["user_id"]
	senderInfo["name"] = info[0]["username"]
	senderInfo["avatar"] = info[0]["avatar"]
	if infoFriend[0]["com_id"] == info[0]["com_id"] {
		senderInfo["position"] = info[0]["position"]
	}

	var receiveInfo = make(map[string]interface{})
	receiveInfo["id"] = infoFriend[0]["user_id"]
	receiveInfo["name"] = infoFriend[0]["username"]
	receiveInfo["avatar"] = infoFriend[0]["avatar"]
	if infoFriend[0]["com_id"] == info[0]["com_id"] {
		receiveInfo["position"] = infoFriend[0]["position"]
	}

	msgMap["senderInfo"] = senderInfo
	msgMap["receiveInfo"] = receiveInfo
	msgMap["type"] = types.FriendApply
	msgMap["applyReason"] = reason
	msgMap["status"] = types.FriendStatusUnHandle
	msgMap["datetime"] = utility.NowMillionSecond()

	text := info[0]["username"] + "请求添加您为好友"

	//如果请求存在，只执行一些更新操作
	if cou > 0 {
		err := db.UpdateApply(reason, remark, userID, friendID, types.FriendStatusUnHandle, types.FriendApply, roomId)
		if err != nil {
			logFriend.Warn("FindFriendRequest db dailed", "err_msg", err)
			return &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
		}
		ids, err := db.FindApplyId(userID, friendID)
		if err != nil {
			logFriend.Warn("FindApplyId db dailed", "err_msg", err)
			return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
		}
		msgMap["id"] = ids[0]["id"]
		msg, err := json.Marshal(msgMap)
		if err != nil {
			return &result.Error{ErrorCode: result.ConvFail, Message: ""}
		}
		msgStr := string(msg)
		fmt.Println(msgStr)
		SendMsg(friendID, msg, text, "系统消息", msgStr, cid[0]["cid"], true, true, 1000*60*60*30)
		SendMsg(userID, msg, text, "系统消息", msgStr, cid[0]["cid"], true, true, 1000*60*60*30)

		return &result.Error{ErrorCode: result.CodeOK, Message: ""}
	} else {
		//添加
		id, err := db.AddApply(types.FriendApply, userID, friendID, reason, remark, roomId)
		if err != nil {
			logFriend.Warn("AddApply db failed", "err_msg", err)
			return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
		}
		msgMap["id"] = strconv.Itoa(id)
		msg, err := json.Marshal(msgMap)
		if err != nil {
			return &result.Error{ErrorCode: result.ConvFail, Message: ""}
		}
		msgStr := string(msg)
		fmt.Println(msgStr)
		SendMsg(friendID, msg, text, "系统消息", msgStr, cid[0]["cid"], true, true, 1000*60*60*30)
		SendMsg(userID, msg, text, "系统消息", msgStr, cid[0]["cid"], true, true, 1000*60*60*30)

		return &result.Error{ErrorCode: result.CodeOK, Message: ""}
	}
}

//处理好友请求
func HandleFriendRequest(userID, friendID string, agree int) *result.Error {
	fmt.Println(userID, friendID, agree)
	var err error

	//先确保该请求是否存在
	res, err := db.FindFriendRequestInfo(userID, friendID)
	if err != nil {
		logFriend.Warn("FindFriendRequest failed", "err_msg", err)
		return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}
	if len(res) == 0 {
		return &result.Error{ErrorCode: result.NotExistFriendRequest, Message: ""}
	}

	//确保请求状态为未处理
	statusStr := res[0]["state"]
	if statusStr != strconv.Itoa(types.FriendStatusUnHandle) {
		return &result.Error{ErrorCode: result.FriendRequestHadDeal, Message: ""}
	}

	//查询自己的信息
	info, err := db.FindUserInfo(userID)
	if err != nil {
		logFriend.Warn("FindUserInfo db dailed", "err_msg", err)
		return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}

	//查询对方的信息
	infoFriend, err := db.FindUserInfo(friendID)
	if err != nil {
		logFriend.Warn("FindUserInfo db dailed", "err_msg", err)
		return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}

	//封装websock消息
	var msgMap = make(map[string]interface{})
	msgMap["eventType"] = 31

	var receiveInfo = make(map[string]interface{})
	receiveInfo["id"] = info[0]["user_id"]
	receiveInfo["name"] = info[0]["username"]
	receiveInfo["avatar"] = info[0]["avatar"]
	if infoFriend[0]["com_id"] == info[0]["com_id"] {
		receiveInfo["position"] = info[0]["position"]
	}

	var senderInfo = make(map[string]interface{})
	senderInfo["id"] = infoFriend[0]["user_id"]
	senderInfo["name"] = infoFriend[0]["username"]
	senderInfo["avatar"] = infoFriend[0]["avatar"]
	if infoFriend[0]["com_id"] == info[0]["com_id"] {
		senderInfo["position"] = infoFriend[0]["position"]
	}

	msgMap["id"] = res[0]["id"]
	msgMap["senderInfo"] = senderInfo
	msgMap["receiveInfo"] = receiveInfo
	msgMap["type"] = types.FriendApply
	msgMap["applyReason"] = res[0]["apply_reason"]
	msgMap["datetime"] = utility.NowMillionSecond()

	//查询cid
	cid, err := db.FindCid(friendID)
	if err != nil {
		logFriend.Warn("FindCid db dailed", "err_msg", err)
		return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}

	if agree == types.FriendRequestReject {
		msgMap["status"] = types.FriendStatusReject
		msg, err := json.Marshal(msgMap)
		if err != nil {
			return &result.Error{ErrorCode: result.ConvFail, Message: ""}
		}
		msgStr := string(msg)
		err = db.RejectFriend(userID, friendID)
		text := info[0]["username"] + "拒绝了您的好友请求"
		SendMsg(friendID, msg, text, "系统消息", msgStr, cid[0]["cid"], true, true, 1000*60*60*30)
		SendMsg(userID, msg, text, "系统消息", msgStr, cid[0]["cid"], true, true, 1000*60*60*30)

	} else {
		msgMap["status"] = types.FriendStatusAccept
		msg, err := json.Marshal(msgMap)
		if err != nil {
			return &result.Error{ErrorCode: result.ConvFail, Message: ""}
		}
		msgStr := string(msg)
		err = db.AcceptFriend(userID, friendID)
		text := info[0]["username"] + "已同意您的好友请求"
		SendMsg(friendID, msg, text, "系统消息", msgStr, cid[0]["cid"], true, true, 1000*60*60*30)
		SendMsg(userID, msg, text, "系统消息", msgStr, cid[0]["cid"], true, true, 1000*60*60*30)

		//如果是群里添加的好友，则需要向群成员发送websocket
		roomId, err := strconv.Atoi(res[0]["room_id"])
		if err != nil {
			return &result.Error{ErrorCode: result.ConvFail, Message: ""}
		}
		if roomId > 0 {
			//todo
			//格式：XX1添加XX1为好友
			//查找群成员
			results, err := db.FindRoomMemberIds(roomId)
			if err != nil {
				logFriend.Error("db FindRoomMemberIds", err)
				return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
			}
			var members = make([]string, 0)
			for _, member := range results {
				members = append(members, member["user_id"])
			}
			//封装消息
			//查找该显示的名称
			nameSM, err := db.FindRoomMemberName(roomId, userID)
			if err != nil {
				return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
			}
			receiveName := nameSM[0]["name"]

			nameSM, err = db.FindRoomMemberName(roomId, friendID)
			if err != nil {
				return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
			}
			senderName := nameSM[0]["name"]

			msgContext := senderName + "添加" + receiveName + "为好友"
			SendAlert(userID, strconv.Itoa(roomId), msgContext, proto.TOROOM, members)
		}

	}
	if err != nil {
		logFriend.Warn("Handle FriendRequest db failed", "err_msg", err)
		return &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
	}
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

//修改备注
func SetFriendRemark(userID, friendID, remark string) *result.Error {
	//不能对自己操作
	if userID == friendID {
		return &result.Error{ErrorCode: result.CanNotOperateSelf, Message: ""}
	}

	//判断用户是否存在  是否是好友
	results := boolIsExistIsFriend(userID, friendID)
	if results != nil {
		return results
	}

	err := db.SetFriendRemark(userID, friendID, remark)
	if err != nil {
		logFriend.Warn("SetFriendRemark update db failed", "err_msg", err)
		return &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
	}
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

/*
	设置好友免打扰
*/
func SetFriendDND(userID, friendID string, DND int) *result.Error {
	//不能对自己操作
	if userID == friendID {
		return &result.Error{ErrorCode: result.CanNotOperateSelf, Message: ""}
	}

	//判断用户是否存在  是否是好友
	results := boolIsExistIsFriend(userID, friendID)
	if results != nil {
		return results
	}
	err := db.SetFriendDND(userID, friendID, DND)
	if err != nil {
		logFriend.Warn("SetFriendDND update db failed", "err_msg", err)
		return &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
	}
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

/*
	设置好友置顶
*/
func SetFriendTop(userID, friendID string, top int) *result.Error {
	//不能对自己操作
	if userID == friendID {
		return &result.Error{ErrorCode: result.CanNotOperateSelf, Message: ""}
	}

	//判断用户是否存在  是否是好友
	results := boolIsExistIsFriend(userID, friendID)
	if results != nil {
		return results
	}
	err := db.SetFriendTop(userID, friendID, top)
	if err != nil {
		logFriend.Warn("SetFriendTop update db failed", "err_msg", err)
		return &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
	}
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

/*
	删除好友
*/
func DeleteFriend(userID, friendID string) *result.Error {
	//不能对自己操作
	if userID == friendID {
		return &result.Error{ErrorCode: result.CanNotOperateSelf, Message: ""}
	}

	//判断用户是否存在  是否是好友
	results := boolIsExistIsFriend(userID, friendID)
	if results != nil {
		return results
	}

	err := db.DeleteFriend(userID, friendID)
	if err != nil {
		logFriend.Warn("Delete Friend failed", "err_msg", err)
		return &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
	}
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func CheckFriendStayConnected(targetId string) bool {
	_, ok := logic.UserMap[targetId]
	return ok
}

//查看好友详情
func FriendInfo(userID, friendID string) (interface{}, *result.Error) {
	//用户是否存在
	resp := boolIsUserExist(friendID)
	if resp != nil {
		return nil, resp
	}

	//好友的信息
	infos, err := db.FindUserInfo(friendID)
	if err != nil {
		return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}
	info := infos[0]

	//用户是否存在
	resp = boolIsUserExist(userID)
	if resp != nil {
		return nil, resp
	}
	//自己的信息
	infos, err = db.FindUserInfo(userID)
	if err != nil {
		logFriend.Info("FriendInfo query db failed", "err_msg", err)
		return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}
	myInfo := infos[0]

	//是否是好友
	isFriend, err := db.CheckFriend(userID, friendID, types.FriendIsNotDelete)
	if err != nil {
		logFriend.Info("CheckFriend query db failed", "err_msg", err)
		return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}

	var returnInfo = make(map[string]interface{})

	//所有人都可以查看的信息
	sex := utility.ToInt(info["sex"])
	returnInfo["sex"] = sex
	returnInfo["avatar"] = info["avatar"]
	returnInfo["id"] = info["user_id"]
	returnInfo["name"] = info["username"]
	returnInfo["mark_id"] = info["mark_id"]
	//同一个公司可以查看职位
	if info["com_id"] == myInfo["com_id"] {
		returnInfo["com_id"] = info["com_id"]
		returnInfo["position"] = info["position"]
	}

	//好友可以查看的信息
	if isFriend {
		//好友之间的信息
		friendInfos, err := db.FindFriend(userID, friendID)
		if err != nil {
			logFriend.Info("FindFreind query db failed", "err_msg", err)
			return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
		}
		friendInfo := friendInfos[0]
		noDisturbing := utility.ToInt(friendInfo["DND"])
		returnInfo["noDisturbing"] = noDisturbing
		stickyOnTop := utility.ToInt(friendInfo["top"])
		returnInfo["stickyOnTop"] = stickyOnTop
		returnInfo["remark"] = friendInfo["remark"]
		addTime, err := strconv.Atoi(friendInfo["add_time"])
		if err != nil {
			return nil, &result.Error{ErrorCode: result.ConvFail, Message: ""}
		}
		returnInfo["addTime"] = addTime
		//是好友就返回1
		returnInfo["isFriend"] = 1
	} else {
		//不是好友就返回2
		returnInfo["isFriend"] = 2
	}

	//管理员客服可以查看的信息  或者自己看自己的信息
	if (myInfo["user_level"] == "2") || (myInfo["user_level"] == "3") || (userID == friendID) {
		returnInfo["uid"] = info["uid"]
		returnInfo["account"] = info["account"]
		returnInfo["phone"] = info["phone"]
		returnInfo["email"] = info["email"]
		verified := utility.ToInt(info["verified"])
		returnInfo["verified"] = verified
		returnInfo["description"] = info["description"]
		userLevel := utility.ToInt(info["user_level"])
		returnInfo["userLevel"] = userLevel
		returnInfo["com_id"] = info["com_id"]
		returnInfo["position"] = info["position"]
	}

	return returnInfo, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

//获取好友消息记录
func FindCatLog(userID, friendID string, startId string, number int) (interface{}, *result.Error) {
	//判断用户是否存在
	resp := boolIsUserExist(friendID)
	if resp != nil {
		return nil, resp
	}
	resp = boolIsUserExist(userID)
	if resp != nil {
		return nil, resp
	}
	if userID == friendID {
		return nil, &result.Error{ErrorCode: result.CanNotOperateSelf, Message: ""}
	}
	//判断是否是好友
	resp = boolIsExistIsFriend(userID, friendID)
	if resp != nil {
		return nil, resp
	}

	start, err := strconv.Atoi(startId)
	if err != nil {
		//获取最新消息id
		res, err := db.FindLastCatLogId(userID, friendID)
		if err != nil {
			return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
		}
		if len(res) == 0 {
			return nil, &result.Error{ErrorCode: result.CodeOK, Message: ""}
		}
		start, err = strconv.Atoi(res[0]["MAX(`id`)"])
		start += 1
		if err != nil {
			return nil, &result.Error{ErrorCode: result.ConvFail, Message: ""}
		}
	}

	//查询消息记录
	logs, err := db.FindCatLog(userID, friendID, start, number)
	if err != nil {
		return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}

	if len(logs) == 0 {
		return nil, &result.Error{ErrorCode: result.CodeOK, Message: ""}
	}

	//查找发送者的信息
	userInfo, err := db.SenderInfo(userID)
	if err != nil {
		return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}
	friInfo, err := db.SenderInfo(friendID)
	if err != nil {
		return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}

	//查找备注
	info, err := db.FindFriend(userID, friendID)
	if err != nil {
		return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}
	remark := info[0]["remark"]

	//封装返回结果
	var next string
	var chatLogs = make([]map[string]interface{}, 0)
	for _, log := range logs {
		var one = make(map[string]interface{})
		one["logId"] = log["id"]
		one["channelType"] = 3
		one["fromId"] = log["sender_id"]
		one["targetId"] = log["receive_id"]
		msgType, err := strconv.Atoi(log["msg_type"])
		if err != nil {
			return nil, &result.Error{ErrorCode: result.ConvFail, Message: ""}
		}
		one["msgType"] = msgType
		con := utility.StringToJobj(log["content"])
		if msgType == 4 {
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
				if userID == utility.ToString(v.RecvUid) {
					isOpened = true
					break
				}
			}
			con["is_opened"] = isOpened
			con["type"] = utility.ToInt(packetDetail["packet_type"])
		}
		one["msg"] = con
		sendTime, err := strconv.Atoi(log["send_time"])
		if err != nil {
			return nil, &result.Error{ErrorCode: result.ConvFail, Message: ""}
		}
		one["datetime"] = sendTime
		if log["sender_id"] == userID {
			one["senderInfo"] = userInfo[0]
		} else {
			one["senderInfo"] = friInfo[0]
			one["remark"] = remark
		}
		chatLogs = append(chatLogs, one)
		next = log["id"]
	}
	var data = make(map[string]interface{})
	data["logs"] = chatLogs
	nextInt, err := strconv.Atoi(next)
	if err != nil {
		return nil, &result.Error{ErrorCode: result.ConvFail, Message: ""}
	}
	//把对方发给我所有的消息改为已读.
	_, _, err = db.ChangePrivateChatLogStstusByUserAndFriendId(userID, friendID)
	if err != nil {
		logFriend.Error("ChangePrivateChatLogStstusByUserAndFriendId failed", err)
		return nil, &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
	}

	//下一次的startId
	nextLog := strconv.Itoa(nextInt)
	data["nextLog"] = nextLog
	return data, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

//删除指定的一条消息
func DeleteMsg(userId, logId string, tp int) *result.Error {
	if tp == types.RoomMsg {
		ret, err := db.CheckRoomMsgContentIsUser(userId, logId)
		if err != nil {
			return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
		}
		if !ret {
			return &result.Error{ErrorCode: result.DeleteMsgFailed, Message: ""}
		}
		count1, err := db.DeleteRoomMsgContent(logId)
		if count1 != 1 {
			return &result.Error{ErrorCode: result.DeleteMsgFailed, Message: ""}
		}
		return &result.Error{ErrorCode: result.CodeOK, Message: ""}
	}
	if tp == types.FriendMsg {
		ret, err := db.CheckCatLogIsUser(userId, logId)
		if err != nil {
			return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
		}
		if !ret {
			return &result.Error{ErrorCode: result.DeleteMsgFailed, Message: ""}
		}
		count1, err := db.DeleteCatLog(logId)
		if count1 != 1 {
			return &result.Error{ErrorCode: result.DeleteMsgFailed, Message: ""}
		}
		return &result.Error{ErrorCode: result.CodeOK, Message: ""}
	}
	return &result.Error{ErrorCode: result.ParamsError, Message: ""}
}

//获取所有好友未读消息统计
func GetAllFriendUnreadMsg1(userId string) (interface{}, *result.Error) {
	//查询所有好友
	friends, err := db.FindFriendIdByUserId(userId)
	if err != nil {
		logFriend.Warn("FindFriendByUserId query db failed", "err_msg", err)
		return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}
	var data = make(map[string]interface{})
	var infos = make([]map[string]interface{}, 0)
	for _, fid := range friends {
		//查询好友的未读聊天记录数
		count, err := db.FindUnReadNum(userId, fid["friend_id"], types.NotRead)
		if err != nil {
			logFriend.Warn("FindUnReadNum query db failed", "err_msg", err)
			return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
		}

		if count > 0 {
			var one = make(map[string]interface{})

			//查询好友信息
			info, err := db.FindFriendInfoByUserId(userId, fid["friend_id"])
			if err != nil {
				logFriend.Warn("FindFriendInfoByUserId query db failed", "err_msg", err)
				return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
			}
			var senderInfo = make(map[string]interface{})
			senderInfo["nickname"] = info[0]["username"]
			senderInfo["avatar"] = info[0]["avatar"]

			//查询好友第一条未读聊天记录
			msg, err := db.FindFirstMsg(userId, fid["friend_id"])
			if err != nil {
				logFriend.Warn("FindFirstMsg query db failed", "err_msg", err)
				return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
			}

			var lastLog = make(map[string]interface{})
			//msgId, err := strconv.Atoi(msg[0]["id"])
			//if err != nil {
			//	return nil, &result.Error{ErrorCode: result.ConvFail}
			//}
			if len(msg) > 0 {
				lastLog["logId"] = msg[0]["id"]
				lastLog["channelType"] = 3
				lastLog["fromId"] = msg[0]["sender_id"]
				lastLog["targetId"] = msg[0]["receive_id"]
				msgType, err := strconv.Atoi(msg[0]["msg_type"])
				if err != nil {
					return nil, &result.Error{ErrorCode: result.ConvFail, Message: ""}
				}
				lastLog["msgType"] = msgType
				con := utility.StringToJobj(msg[0]["content"])
				if msgType == 4 {
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
						if userId == utility.ToString(v.RecvUid) {
							isOpened = true
							break
						}
					}
					con["is_opened"] = isOpened
					con["type"] = utility.ToInt(packetDetail["packet_type"])
				}
				lastLog["msg"] = con
				sendTime, err := strconv.Atoi(msg[0]["send_time"])
				if err != nil {
					return nil, &result.Error{ErrorCode: result.ConvFail, Message: ""}
				}
				lastLog["datetime"] = sendTime

				lastLog["senderInfo"] = senderInfo
			}

			one["id"] = fid["friend_id"]
			one["number"] = count
			one["lastLog"] = lastLog
			infos = append(infos, one)
		}
	}
	data["infos"] = infos
	return data, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}
