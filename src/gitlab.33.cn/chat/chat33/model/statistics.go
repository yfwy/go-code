package model

import (
	"strings"

	cmn "dev.33.cn/33/common"
	"github.com/inconshreveable/log15"
	"gitlab.33.cn/chat/chat33/db"
	"gitlab.33.cn/chat/chat33/result"
	logic "gitlab.33.cn/chat/chat33/router"
	"gitlab.33.cn/chat/chat33/types"
	"gitlab.33.cn/chat/chat33/utility"
)

var slog = log15.New("module", "model/statistics")

type ApplyList struct {
	SenderInfo  ApplyUserInfo `json:"senderInfo"`
	ReceiveInfo ApplyUserInfo `json:"receiveInfo"`
	Id          string        `json:"id"`
	Type        int           `json:"type"`
	ApplyReason string        `json:"applyReason"`
	Status      int           `json:"status"`
	Datetime    int64         `json:"datetime"`
}

type ApplyUserInfo struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Position string `json:"position"`
}

type SearchInfo struct {
	Type     int         `json:"type"`
	RoomInfo interface{} `json:"roomInfo"`
	UserInfo interface{} `json:"userInfo"`
}

// 获取首页统计信息
func GetIndexStatistics() (interface{}, *result.Error) {
	onlineUsers := make(map[string]bool)
	for _, g := range utility.GroupList {
		for u, devMap := range g {
			if len(devMap) > 0 {
				onlineUsers[u] = true
			}
		}
	}

	csNum, err := db.GetCsNum()
	if err != nil {
		slog.Error("get cs num", "err", err.Error())
	}

	//今日推广红包数量
	rows, err := db.GetTodayPackets()
	if err != nil {
		slog.Error("get today red packet", "err", err.Error())
	}

	var todayAdvPackets int
	for _, row := range rows {
		if cmn.ToInt(row["type"]) == types.PacketTypeAdv {
			todayAdvPackets++
		}
	}

	ret := &types.Statistics{
		OnlineUserNum:   len(onlineUsers),
		TodayAdvPackets: todayAdvPackets,
		GroupNum:        len(utility.GroupList),
		CsNum:           csNum,
	}

	return ret, &result.Error{ErrorCode: result.CodeOK, Message: err.Error()}
}

// 获取app列表
func GetAppInfoList() (interface{}, *result.Error) {
	type rlt struct {
		AppList interface{} `json:"app_list"`
	}

	return rlt{AppList: appList}, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

// 获取币种信息
func GetCoinList() (interface{}, *result.Error) {
	type rlt struct {
		CoinList interface{} `json:"coin_list"`
	}
	return rlt{CoinList: coinList}, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func ClearlySearch(caller, searchId string) (interface{}, *result.Error) {
	maps, err := db.GetRoomsInfoByMarkId(searchId)
	if err != nil {
		return nil, &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}
	var ret SearchInfo
	if len(maps) > 0 {
		var room RoomInfo
		r := maps[0]
		room.Id = utility.ToString(r["id"])
		room.MarkId = utility.ToString(r["mark_id"])
		room.Name = utility.ToString(r["name"])
		room.Avatar = utility.ToString(r["avatar"])
		room.CanAddFriend = utility.ToInt(r["can_add_friend"])
		room.JoinPermission = utility.ToInt(r["join_permission"])
		room.NoDisturbing = 0
		room.OnlineNumber = 0
		room.Users = make([]RoomUserList, 0)
		ret.Type = 1
		ret.RoomInfo = &room
		var empty struct{}
		ret.UserInfo = empty
		return &ret, &result.Error{ErrorCode: result.CodeOK, Message: ""}
	}

	//find friend
	//find user by uid
	infos, err := db.FindUserByMarkId(searchId)
	if err != nil {
		return nil, &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}
	if len(infos) > 0 {
		userInfo := infos[0]
		userId := userInfo["user_id"]
		friendInfo, errMsg := FriendInfo(caller, userId)
		if errMsg.ErrorCode != result.CodeOK {
			return nil, errMsg
		}
		if friendInfo != nil {
			ret.Type = 2
			ret.UserInfo = friendInfo
			var empty struct{}
			ret.RoomInfo = empty
			return &ret, &result.Error{ErrorCode: result.CodeOK, Message: ""}
		}
	}
	var empty struct{}
	return empty, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func GetApplyList(caller, id string, number int) (interface{}, *result.Error) {
	// get all managaer room
	rooms, _ := GetUserManageRooms(caller)
	sql := ""
	for i, v := range rooms {
		if i != 0 {
			sql += ","
		}
		sql += v
	}
	maps, err := db.GetApplyList(caller, utility.ToInt64(id), number+1, sql)
	if err != nil {
		return nil, &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
	}

	var applyInfoList = make([]ApplyList, 0)
	for _, v := range maps {
		var applyInfo ApplyList
		applyInfo.Id = utility.ToString(v["id"])
		applyInfo.Type = utility.ToInt(v["type"])
		applyInfo.ApplyReason = utility.ToString(v["apply_reason"])
		applyInfo.Status = utility.ToInt(v["state"])
		applyInfo.Datetime = utility.ToInt64(v["datetime"])

		targetId := utility.ToString(v["target"])
		applyUser := utility.ToString(v["apply_user"])
		//get sender Info
		var senderInfo ApplyUserInfo
		userInfo, errMsg := FriendInfo(caller, applyUser)
		if errMsg.ErrorCode != result.CodeOK {
			return nil, errMsg
		}
		if userInfo != nil {
			userInfoMap := userInfo.(map[string]interface{})
			senderInfo.Id = userInfoMap["id"].(string)
			senderInfo.Name = userInfoMap["name"].(string)
			senderInfo.Avatar = userInfoMap["avatar"].(string)
			if userInfoMap["position"] != nil {
				senderInfo.Position = userInfoMap["position"].(string)
			}
		}

		//get target Info
		var targetInfo ApplyUserInfo
		if applyInfo.Type == 1 {
			roomInfos, err := db.GetRoomsInfo(targetId)
			if err != nil {
				return nil, &result.Error{ErrorCode: result.DbConnectFail, Message: ""}
			}
			if len(roomInfos) < 1 {
				continue
			}
			info := roomInfos[0]
			targetInfo.Id = info["id"]
			targetInfo.Name = info["name"]
			targetInfo.Avatar = info["avatar"]
			targetInfo.Position = ""
		} else if applyInfo.Type == 2 {
			userInfo, errMsg := FriendInfo(caller, targetId)
			if errMsg.ErrorCode != result.CodeOK {
				return nil, errMsg
			}
			if userInfo != nil {
				userInfoMap := userInfo.(map[string]interface{})
				targetInfo.Id = userInfoMap["id"].(string)
				targetInfo.Name = userInfoMap["name"].(string)
				targetInfo.Avatar = userInfoMap["avatar"].(string)
				if userInfoMap["position"] != nil {
					targetInfo.Position = userInfoMap["position"].(string)
				}
			}
		}
		applyInfo.SenderInfo = senderInfo
		applyInfo.ReceiveInfo = targetInfo

		applyInfoList = append(applyInfoList, applyInfo)
	}

	var ret = make(map[string]interface{})
	count, _ := db.GetApplyListNumber()
	nextId := "0"
	if len(maps) > number {
		nextId = applyInfoList[len(applyInfoList)-1].Id
		applyInfoList = applyInfoList[0 : len(applyInfoList)-1]
	}
	ret["applyList"] = applyInfoList
	ret["nextId"] = nextId
	ret["totalNumber"] = count
	return &ret, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func SendMsgWithGT(userId, geTuiHeader, geTuiConent string, msg interface{}) *result.Error {
	//查询cid
	cid, err := db.FindCid(userId)
	if err != nil {
		logFriend.Warn("FindCid db dailed", "err_msg", err)
		return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}
	//TODO if not find cid
	if len(cid) > 0 {
		SendMsg(userId, msg, geTuiHeader, "系统消息", geTuiConent, cid[0]["cid"], true, true, 1000*60*60*30)
	}
	return &result.Error{ErrorCode: result.CodeOK}
}

func SendMsg(uid string, msg interface{}, text, title, transmissionContent, cid string, isOffline, transmissionType bool, offlineExpireTime int) error {
	user, ok := logic.UserMap[uid]
	if ok {
		count := 0
		for _, cli := range user.Clients {
			cli.Send(msg)
			count++
		}
		if count <= 0 {
			//个推
			cids := strings.Split(cid, "#")
			for _, c := range cids {
				err := utility.GTPushSingle(text, title, transmissionContent, c, isOffline, transmissionType, offlineExpireTime)
				if err != nil {
					return err
				}
			}
		}
		return nil
	} else {
		//个推
		cids := strings.Split(cid, "#")
		for _, c := range cids {
			err := utility.GTPushSingle(text, title, transmissionContent, c, isOffline, transmissionType, offlineExpireTime)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
