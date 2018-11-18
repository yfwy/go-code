package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	cmn "dev.33.cn/33/common"
	l "github.com/inconshreveable/log15"
	"gitlab.33.cn/chat/chat33/db"
	m "gitlab.33.cn/chat/chat33/model"
	"gitlab.33.cn/chat/chat33/result"
	"gitlab.33.cn/chat/chat33/types"
	"gitlab.33.cn/chat/chat33/utility"
)

const (
	InviteCode = "888888"
)

var redPacketLog = l.New("module", "chat/api/red_envelope")

// coin ---> red packet currency
// var coinMap = map[int]int{
//	1: 1, //BTY
//	2: 2, //YCC
// }

// // 发红包
// func Send(context *gin.Context) {
// 	var ret m.BaseReturn

// 	session := sessions.Default(context)
// 	sUid := session.Get("user_id")
// 	if sUid == nil {
// 		WriteMsg(context.Writer, []byte("session lack user_id"))
// 		return
// 	}
// 	uid := sUid.(string)
// 	/*session, err := utility.SessionStore.Get(context.Request, utility.SESSION_LOGIN)
// 	if err != nil {
// 		redPacketLog.Error("invalid session", "err", err.Error())
// 		WriteMsg(context.Writer, []byte("invalid session"))
// 		return
// 	}
// 	uid := cmn.ToString(session.Values["user_id"])*/

// 	token := context.GetHeader("FZM-AUTH-TOKEN")
// 	if token == "" {
// 		ret := result.ComposeHttpAck(result.LackParam, "Header: token", result.Empty{})
// 		context.JSON(http.StatusOK, ret)
// 		return
// 	}

// 	device := context.GetHeader("FZM-DEVICE")
// 	if device == "" {
// 		ret := result.ComposeHttpAck(result.LackParam, "Header: device", result.Empty{})
// 		context.JSON(http.StatusOK, ret)
// 		return
// 	}

// 	redPacketLog.Info("user info", "token", token, "user_id", uid)

// 	body, _ := ioutil.ReadAll(context.Request.Body)
// 	var params map[string]interface{}
// 	err := json.Unmarshal(body, &params)
// 	if err != nil {
// 		ret := result.ComposeHttpAck(result.ParamsError, "Header: device", result.Empty{})
// 		context.JSON(http.StatusOK, ret)
// 		return
// 	}

// 	appId := context.GetHeader("FZM-APP-ID")
// 	if appId == "" {
// 		if _, ok := params["app_id"]; !ok {
// 			ret := result.ComposeHttpAck(result.LackParam, "Header: FZM-APP-ID", result.Empty{})
// 			context.JSON(http.StatusOK, ret)
// 			return
// 		} else {
// 			appId = cmn.ToString(params["app_id"])
// 		}
// 	}

// 	if _, ok := params["is_group"]; !ok {
// 		ret := result.ComposeHttpAck(result.LackParam, "is_group", result.Empty{})
// 		context.JSON(http.StatusOK, ret)
// 		return
// 	}
// 	isGroup := cmn.ToInt(params["is_group"])

// 	if _, ok := params["to_id"]; !ok {
// 		ret := result.ComposeHttpAck(result.LackParam, "to_id", result.Empty{})
// 		context.JSON(http.StatusOK, ret)
// 		return
// 	}
// 	toId := cmn.ToString(params["to_id"])

// 	if _, ok := params["amount"]; !ok {
// 		ret := result.ComposeHttpAck(result.LackParam, "amount", result.Empty{})
// 		context.JSON(http.StatusOK, ret)
// 		return
// 	}
// 	amount := cmn.ToString(params["amount"])

// 	if _, ok := params["size"]; !ok {
// 		ret := result.ComposeHttpAck(result.LackParam, "size", result.Empty{})
// 		context.JSON(http.StatusOK, ret)
// 		return
// 	}
// 	size := cmn.ToString(params["size"])

// 	if _, ok := params["remark"]; !ok {
// 		ret := result.ComposeHttpAck(result.LackParam, "remark", result.Empty{})
// 		context.JSON(http.StatusOK, ret)
// 		return
// 	}
// 	remark := cmn.ToString(params["remark"])

// 	if _, ok := params["type"]; !ok {
// 		ret := result.ComposeHttpAck(result.LackParam, "type", result.Empty{})
// 		context.JSON(http.StatusOK, ret)
// 		return
// 	}
// 	packetType := cmn.ToString(params["type"])
// 	if _, ok := params["coin"]; !ok {
// 		ret := result.ComposeHttpAck(result.LackParam, "coin", result.Empty{})
// 		context.JSON(http.StatusOK, ret)
// 		return
// 	}
// 	coin := cmn.ToInt(params["coin"])
// 	if _, ok := coinMap[coin]; !ok {
// 		redPacketLog.Error("invalid coin", "coin", coin)
// 		WriteMsg(context.Writer, []byte("invalid coin"))
// 		return
// 	}

// 	user, err := m.GetZbUserInfo(token)
// 	if err != nil {
// 		redPacketLog.Error("GetUIDByToken", "err", err)
// 		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
// 		context.JSON(http.StatusOK, ret)
// 		return
// 	}
// 	zbUid := cmn.ToString(user["id"])

// 	// TODO
// 	rows, _ := db.GetUserInfoByID(uid)
// 	if len(rows) == 0 {
// 		redPacketLog.Error("invalid uid", "uid", uid)
// 		ret := result.ComposeHttpAck(result.UserNotExists, "", result.Empty{})
// 		context.JSON(http.StatusOK, ret)
// 		return
// 	}
// 	dbZbUid := cmn.ToString(rows[0]["uid"])

// 	if zbUid != dbZbUid {
// 		redPacketLog.Error("token/uid not match", "uid", uid, "token", token, "zbUid", zbUid, "dbZbUid", dbZbUid)
// 		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
// 		context.JSON(http.StatusOK, ret)
// 		return
// 	}

// 	//在包红包之前先检查是否可以发送 红包发送条件 1、用户未被禁言 2、用户在群
// 	if isGroup == 1 {
// 		muted, err := m.CheckUserMuted(uid, utility.NowMillionSecond())
// 		if err != nil {
// 			redPacketLog.Error("errmsg", "err_msg", err)
// 			ret := result.ComposeHttpAck(result.DbConnectFail, "", result.Empty{})
// 			context.JSON(http.StatusOK, ret)
// 			return
// 		}

// 		if muted {
// 			ret := result.ComposeHttpAck(result.UserMuted, "", result.Empty{})
// 			context.JSON(http.StatusOK, ret)
// 			return
// 		}

// 		//查看是否在此群
// 		group, ok := utility.GroupList[cmn.ToString(toId)]
// 		if !ok {
// 			redPacketLog.Debug("group not find", "group id", toId)
// 			ret := result.ComposeHttpAck(result.GroupNotExists, "", result.Empty{})
// 			context.JSON(http.StatusOK, ret)
// 			return
// 		}

// 		if _, ok = group[uid]; !ok {
// 			redPacketLog.Debug("用户不在此群组", "gourp id", toId, "user id", uid)
// 			ret := result.ComposeHttpAck(result.UserNotEnterGroup, "", result.Empty{})
// 			context.JSON(http.StatusOK, ret)
// 			return
// 		}
// 	}

// 	req := &types.ReqSendRedPacket{
// 		Amount:         amount,
// 		Size:           size,
// 		Remark:         remark,
// 		Type:           packetType,
// 		Coin:           coinMap[coin],
// 		InvitationCode: InviteCode,
// 	}
// 	code, packet, err := m.Send(token, req, remark)
// 	if code != 0 {
// 		redPacketLog.Error("send packet", "err", err.Error())
// 		ret := ret.Get(code, err.Error())
// 		context.JSON(http.StatusOK, ret)
// 		return
// 	}

// 	//todo 根据user_id查出user信息
// 	var _username string
// 	var _user_level int

// 	var sendTime = utility.NowMillionSecond()

// 	var reqstr m.WsBaseMessage
// 	reqstr.Event_type = 0
// 	reqstr.Msg_id = ""
// 	reqstr.From_id = uid
// 	if isGroup == 1 {
// 		reqstr.From_gid = toId
// 		reqstr.To_id = ""
// 	} else {
// 		reqstr.To_id = toId
// 		reqstr.From_gid = ""
// 	}
// 	reqstr.Name = _username
// 	reqstr.User_level = _user_level
// 	reqstr.Msg_type = 4
// 	reqstr.Msg = packet
// 	reqstr.Datetime = sendTime

// 	var client utility.Client = nil
// 	// 通过webSocket发出去
// 	if isGroup == 1 {
// 		if group, ok := utility.GroupList[toId]; ok {
// 			if devMap, ok := group[uid]; ok {
// 				if c, ok := devMap[device]; ok {
// 					client = c
// 				}
// 			}
// 		}
// 	} else {
// 		if devMap, ok := utility.Usermap[uid]; ok {
// 			if c, ok := devMap[device]; ok {
// 				client = c
// 			}
// 		}
// 	}

// 	if client != nil {
// 		var message []byte
// 		if message, err = json.Marshal(reqstr); err != nil {
// 			redPacketLog.Error("prase red envelop message to json err", "err", err)
// 			return
// 		}
// 		var caller *m.BaseReturn = new(m.BaseReturn)
// 		rlt, _ := client.DoBroadcast(message, caller)
// 		if rlt != nil {
// 			WriteMsg(context.Writer, rlt)
// 		}
// 	}

// 	//将红包存到数据库
// 	err = db.InsertPacket(packet.PacketId, cmn.ToString(appId), toId, uid, packetType, cmn.ToString(coin),
// 		size, amount, remark, cmn.ToString(sendTime), isGroup)
// 	if err != nil {
// 		redPacketLog.Error("send red packet", "err2", err)
// 		ret := ret.Get(m.UserHasReg, err.Error(), m.Empty{})
// 		context.JSON(http.StatusOK, ret)
// 		return
// 	}

// 	//WriteMsg(context.Writer, ret.Get(0, "操作成功", packet))
// 	context.JSON(http.StatusOK, ret.Get(0, "", packet))
// }

// 红包统计信息
func RedEnvelopeStatistics(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var params map[string]interface{}
	err := json.Unmarshal(body, &params)
	if err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	appId := context.GetHeader("FZM-APP-ID")
	if appId == "" {
		if _, ok := params["app_id"]; !ok {
			ret := result.ComposeHttpAck(result.LackParam, "app_id", result.Empty{})
			context.JSON(http.StatusOK, ret)
			return
		} else {
			appId = cmn.ToString(params["app_id"])
		}
	}

	resp, err := m.RedEnvelopeStatistics(appId)
	if err != nil {
		redPacketLog.Error("red packet statistics", "err", err.Error())
		ret := result.ComposeHttpAck(result.RPError, err.Error(), result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	ret := result.ComposeHttpAck(result.CodeOK, "", resp)
	context.JSON(http.StatusOK, ret)
}

// 红包信息列表api
func RedEnvelopeInfoList(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var params map[string]interface{}
	err := json.Unmarshal(body, &params)
	if err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	queryParams := &types.PacketQueryParam{}
	appId := context.GetHeader("FZM-APP-ID")
	if appId == "" {
		if _, ok := params["app_id"]; !ok {
			ret := result.ComposeHttpAck(result.LackParam, "app_id", result.Empty{})
			context.JSON(http.StatusOK, ret)
			return
		} else {
			appId = cmn.ToString(params["app_id"])
		}
	}
	queryParams.AppId = appId

	if _, ok := params["packet_id"]; ok {
		queryParams.PacketId = cmn.ToString(params["packet_id"])
	}

	if _, ok := params["packet_type"]; ok {
		if cmn.ToString(params["packet_type"]) != "" {
			queryParams.PacketType = cmn.ToInt(params["packet_type"])
		}
	}

	if _, ok := params["coin"]; ok {
		if cmn.ToString(params["coin"]) != "" {
			coin := cmn.ToInt(params["coin"])
			queryParams.Coin = coin
		}
	}

	if _, ok := params["uid"]; ok {
		queryParams.Uid = cmn.ToString(params["uid"])
	}

	if _, ok := params["start_time"]; ok {
		if ss := cmn.ToString(params["start_time"]); ss != "" {
			st := utility.ToInt64(params["start_time"])
			queryParams.StartTime = st
		}
	}

	if _, ok := params["end_time"]; ok {
		if ss := cmn.ToString(params["end_time"]); ss != "" {
			st := utility.ToInt64(params["end_time"])
			queryParams.EndTime = st
		}
	}

	if queryParams.StartTime < 0 || queryParams.EndTime < 0 ||
		(queryParams.StartTime != 0 && queryParams.EndTime == 0) ||
		(queryParams.StartTime == 0 && queryParams.EndTime != 0) ||
		(queryParams.StartTime > queryParams.EndTime) {
		WriteMsg(context.Writer, []byte("invalid start/end time"))
		return
	}

	if _, ok := params["page"]; ok {
		queryParams.Page = cmn.ToInt(params["page"])
	} else {
		queryParams.Page = types.PageNo
	}

	if _, ok := params["number"]; ok {
		queryParams.Number = cmn.ToInt(params["number"])
	} else {
		queryParams.Number = types.PageLimit
	}

	resp, err := m.RedEnvelopeInfoList(queryParams)
	if err != nil {
		redPacketLog.Error("red packet info list", "err", err)
		ret := result.ComposeHttpAck(result.RPError, err.Error(), result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	ret := result.ComposeHttpAck(result.CodeOK, "", resp)
	context.JSON(http.StatusOK, ret)
}

// 红包查看详情api
func RedEnvelopeDetail(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var params map[string]interface{}
	err := json.Unmarshal(body, &params)
	if err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	if _, ok := params["packet_id"]; !ok {
		ret := result.ComposeHttpAck(result.LackParam, "packet_id", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}
	packetId := cmn.ToString(params["packet_id"])

	resp, err := m.RedEnvelopeDetail(packetId)
	if err != nil {
		redPacketLog.Error("red packet detail", "err", err)
		ret := result.ComposeHttpAck(result.RPEmpty, err.Error(), result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	ret := result.ComposeHttpAck(result.CodeOK, "", resp)
	context.JSON(http.StatusOK, ret)
}

// ================ new ver ==================

// 查询账户余额api
func Balance(context *gin.Context) {
	token := context.GetHeader("FZM-AUTH-TOKEN")
	if token == "" {
		ret := result.ComposeHttpAck(result.LackParam, "Header: FZM-AUTH-TOKEN", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	resp, err := m.Balance(token)
	if err != nil {
		redPacketLog.Error("query balance", "err", err.Error())
		ret := result.ComposeHttpAck(result.RPError, err.Error(), result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", resp)
	context.JSON(http.StatusOK, ret)
}

//（登录用户）领取红包
func ReceiveEntry(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		WriteMsg(c.Writer, []byte("session lack user_id"))
		return
	}

	type receiveEntryParams struct {
		PacketID string `json:"packet_id" binding:"required"`
	}
	var params receiveEntryParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, result.ComposeHttpAck(result.ParamsError, "", nil))
		return
	}

	rows, _ := db.GetUserInfoByID(userID.(string))
	if len(rows) == 0 {
		redPacketLog.Error("invalid uid", "uid", userID.(string))
		ret := result.ComposeHttpAck(result.UserNotExists, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	zbUid := cmn.ToString(rows[0]["uid"])
	code, resp, err := m.ReceiveEntry(params.PacketID, zbUid)
	if code != 0 {
		redPacketLog.Error("receive entry", "err", err.Error())
		ret := result.ComposeHttpAck(code, err.Error(), result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	ret := result.ComposeHttpAck(code, "", resp)
	c.JSON(http.StatusOK, ret)
}

//（未登录用户）领取红包
func Receive(c *gin.Context) {
	type receiveParams struct {
		PacketID string `json:"packet_id" binding:"required"`
	}
	var params receiveParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, result.ComposeHttpAck(result.ParamsError, "", nil))
		return
	}

	code, resp, err := m.Receive(params.PacketID)
	if code != 0 {
		redPacketLog.Error("receive", "err", err.Error())
		ret := result.ComposeHttpAck(code, err.Error(), result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	ret := result.ComposeHttpAck(code, "", resp)
	c.JSON(http.StatusOK, ret)
}

//红包入账
func Entry(c *gin.Context) {
	type entryParams struct {
		Mark    string `json:"mark" binding:"required"`
		Account string `json:"account" binding:"required"`
	}
	var params entryParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, result.ComposeHttpAck(result.ParamsError, "", nil))
		return
	}

	code, resp, err := m.Entry(params.Mark, params.Account)
	if err != nil {
		redPacketLog.Error("entry", "err", err.Error())
		ret := result.ComposeHttpAck(code, err.Error(), result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	ret := result.ComposeHttpAck(code, "", resp)
	c.JSON(http.StatusOK, ret)
}

// （新用户）注册领取红包
func RegisterEntry(c *gin.Context) {
	type registerParams struct {
		Mark     string `json:"mark" binding:"required"`
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binding:"required"`
		Captcha  string `json:"captcha" binding:"required"`
	}
	var params registerParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, result.ComposeHttpAck(result.ParamsError, "", nil))
		return
	}

	code, resp, err := m.RegisterEntry(params.Mark, params.Account, params.Password, params.Captcha)
	if err != nil {
		redPacketLog.Error("register-entry", "err", err.Error())
		ret := result.ComposeHttpAck(code, err.Error(), result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	ret := result.ComposeHttpAck(result.CodeOK, "", resp)
	c.JSON(http.StatusOK, ret)
}

/*
	用户收红包历史
*/
func RPReceiveHistory(c *gin.Context) {
	token := c.GetHeader("FZM-AUTH-TOKEN")
	if token == "" {
		ret := result.ComposeHttpAck(result.LackParam, "Header: FZM-AUTH-TOKEN", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	type recParams struct {
		Date string `json:"date" binding:"required"`
		Coin int    `json:"coin"`
	}
	var params recParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, result.ComposeHttpAck(result.ParamsError, "", nil))
		return
	}

	if params.Coin == 0 {
		params.Coin = 1 // BTY
	}

	ret, _ := m.RPReceiveHistory(token, params.Date, params.Coin)
	c.JSON(http.StatusOK, ret)
}

/*
	用户发红包记录
*/
func RPSendHistory(c *gin.Context) {
	token := c.GetHeader("FZM-AUTH-TOKEN")
	if token == "" {
		ret := result.ComposeHttpAck(result.LackParam, "Header: FZM-AUTH-TOKEN", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	type sendParams struct {
		Date string `json:"date" binding:"required"`
		Coin int    `json:"coin"`
	}
	var params sendParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, result.ComposeHttpAck(result.ParamsError, "", nil))
		return
	}

	if params.Coin == 0 {
		params.Coin = 1 // BTY
	}

	resp, _ := m.RPSendHistory(token, params.Date, params.Coin)
	ret := result.ComposeHttpAck(result.CodeOK, "", resp)
	c.JSON(http.StatusOK, ret)
}

/*
	用户红包记录
*/
func RPHistory(c *gin.Context) {
	token := c.GetHeader("FZM-AUTH-TOKEN")
	if token == "" {
		ret := result.ComposeHttpAck(result.LackParam, "Header: FZM-AUTH-TOKEN", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	type hisParams struct {
		Date string `json:"date" binding:"required"`
		Coin int    `json:"coin"`
	}
	var params hisParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "date", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	if params.Coin == 0 {
		params.Coin = 1 // BTY
	}

	resp, _ := m.RPHistory(token, params.Date, params.Coin)
	ret := result.ComposeHttpAck(result.CodeOK, "", resp)
	c.JSON(http.StatusOK, ret)
}

func Send(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		WriteJson(c, result.ComposeHttpAck(result.SessionError, "缺少user_id", nil))
		return
	}

	token := c.GetHeader("FZM-AUTH-TOKEN")
	if token == "" {
		ret := result.ComposeHttpAck(result.LackParam, "Header: FZM-AUTH-TOKEN", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	device := c.GetHeader("FZM-DEVICE")
	if device == "" {
		ret := result.ComposeHttpAck(result.LackParam, "Header: device", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	type sendParams struct {
		CType  *int   `json:"c_type" binding:"exists"`
		ToId   string `json:"to_id" binding:"required"`
		Uid    string `json:"uid" binding:"required"`
		Coin   *int   `json:"coin" binding:"exists"`
		Type   string `json:"type" binding:"required"`
		Amount string `json:"amount" binding:"required"`
		Size   string `json:"size" binding:"required"`
		Remark string `json:"remark" binding:"required"`
	}

	var params sendParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, result.ComposeHttpAck(result.ParamsError, "", nil))
		return
	}

	// compare tokenUid && dbUid

	rows, _ := db.GetUserInfoByID(userID.(string))
	if len(rows) == 0 {
		redPacketLog.Error("invalid user_id", "user_id", userID)
		ret := result.ComposeHttpAck(result.UserNotExists, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	// uid := cmn.ToString(rows[0]["uid"])

	// TODO: 检查是否可发送红包：可发言(禁言、群禁言)、是群成员、是好友？

	req := &types.ReqSendRedPacket{
		Amount:         params.Amount,
		Size:           params.Size,
		Remark:         params.Remark,
		Type:           params.Type,
		Coin:           *params.Coin,
		InvitationCode: InviteCode,
	}

	// todo
	code, packet, err := m.Send(token, req, params.Remark)
	if code != 0 {
		redPacketLog.Error("send packet", "err", err.Error())
		ret := result.ComposeHttpAck(code, err.Error(), result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	// todo: websocket msg
	//InsertPacket(packetID, userID, toID, tType, size, amount, remark string, cType, coin int, time int64)
	err = db.InsertPacket(packet.PacketId, userID.(string), params.ToId, params.Type, params.Size, params.Amount,
		params.Remark, *params.CType, *params.Coin, utility.NowMillionSecond())
	if err != nil {
		redPacketLog.Error("send red packet", "err2", err)
		ret := result.ComposeHttpAck(result.UserHasReg, err.Error(), result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	c.JSON(http.StatusOK, result.ComposeHttpAck(result.CodeOK, "", packet))
}
