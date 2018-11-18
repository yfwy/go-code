package model

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"gitlab.33.cn/chat/chat33/utility"

	cmn "dev.33.cn/33/common"
	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"gitlab.33.cn/chat/chat33/db"
	"gitlab.33.cn/chat/chat33/result"
	"gitlab.33.cn/chat/chat33/types"
)

var redPackLog = log15.New("module", "model/red_envelope")

// red packet currency  --> zhaobi coin
var coinMap = map[int]int{
	1: 1, //BTY
	2: 2, //YCC
}

var CoinToZBCoin = map[int]int{
	1: 1, // BTY
	2: 2, // YCC
}

func Balance(token string) (interface{}, error) {
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + token
	bytes, err := cmn.HTTPRequest("GET", cfg.Api.Zhaobi+"/api/account/asset", headers, nil)
	if err != nil {
		return nil, err
	}

	/*
		{
		   "code" : 200,
		   "data" : {
		      "valuation" : "512390.23",
		      "list" : {
		         "BTY" : {
		            "name" : "BTY",
		            "realactive" : "100996.700000",
		            "valuation" : "498157.13",
		            "active" : "100996.700000",
		            "poundage" : "0.000000",
		            "total" : "100996.700000",
		            "frozen" : "0.000000"
				 },
				 ...
		      }
		   },
		   "error" : "OK",
		   "ecode" : "200",
		   "message" : "OK"
		}
	*/

	var resp map[string]interface{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	if cmn.ToInt(resp["code"]) != 200 {
		return nil, errors.New(resp["message"].(string))
	}

	data, ok := resp["data"]
	if !ok {
		return 0, errors.New("no 'data' info")
	}

	type Balance struct {
		Coin   int     `json:"coin"`
		Amount float32 `json:"amount"`
	}

	clist, ok := data.(map[string]interface{})["list"]
	if !ok {
		return nil, errors.New("no 'list' info")
	}

	ret := make(map[string]interface{})
	balances := []*Balance{}
	for _, v := range clist.(map[string]interface{}) {
		one := v.(map[string]interface{})
		for _, e := range coinList {
			if e.CoinName == one["name"] {
				b := &Balance{
					Coin:   cmn.ToInt(e.CoinId),
					Amount: cmn.ToFloat32(one["active"]),
				}
				balances = append(balances, b)
				break
			}
		}
	}

	ret["balances"] = balances
	return ret, nil
}

// 发送红包
func Send(token string, req *types.ReqSendRedPacket, remark string) (int, *types.RedPacket, error) {
	headers := make(map[string]string)
	headers["auth-token"] = token

	params := url.Values{}
	params.Set("coin", cmn.ToString(req.Coin))
	params.Set("amount", req.Amount)
	params.Set("size", req.Size)
	params.Set("type", req.Type)
	params.Set("remark", req.Remark)
	params.Set("invitation_code", req.InvitationCode)

	byte, err := cmn.HTTPPostForm(cfg.Api.RedPacket+"/red-packet-random/split", headers, strings.NewReader(params.Encode()))
	if err != nil {
		return result.RPError, nil, err
	}

	/*

		{
			"code":10070,
			"message":"暂无红包功能权限"
		}

		{
		   "code" : 200,
		   "coin" : "2",
		   "url" : "http://47.74.190.154/packets/ycc/pages/open.html?id=b84e7d70-94a0-11e8-b452-631c93934368",
		   "message" : "success",
		   "type" : "1",  //拼手气红包
		   "id" : "b84e7d70-94a0-11e8-b452-631c93934368"
		}
	*/
	var resp map[string]interface{}
	err = json.Unmarshal(byte, &resp)
	if err != nil {
		return result.RPError, nil, err
	}

	code := cmn.ToInt(resp["code"])
	if code != 200 {
		msg := cmn.ToString(resp["message"])
		switch msg {
		case "暂无红包功能权限":
			return result.CannotSendRP, nil, errors.New(msg)
		default:
			return result.RPError, nil, errors.New(msg)
		}
	}

	packet := &types.RedPacket{
		PacketId:   cmn.ToString(resp["id"]),
		PacketUrl:  cmn.ToString(resp["url"]),
		PacketType: cmn.ToInt(resp["type"]),
		Coin:       coinMap[cmn.ToInt(resp["coin"])],
		Remark:     remark,
	}

	return 0, packet, nil
}

// 领取红包
func ReceiveEntry(packetId, uid string) (int, interface{}, error) {
	params := url.Values{}
	params.Set("id", packetId)
	params.Set("user_id", uid)
	byte, err := cmn.HTTPPostForm(cfg.Api.RedPacket+"/red-packet-random/receive-entry", nil, strings.NewReader(params.Encode()))
	if err != nil {
		return result.RPError, nil, err
	}

	redPackLog.Debug("red-packet-random", "receive-entry", string(byte))

	/*
		{
			"message" : "仅限新人领取",
			"code" : 10072
		}

		{
			"message" : "红包已被领完",
			"code" : 10072
		}

		{
		   "total" : 2,
		   "code" : 200,
		   "amount" : "3",
		   "message" : "success",
		   "remain" : 1,
		   "coin" : 1
		}
	*/
	var resp map[string]interface{}
	err = json.Unmarshal(byte, &resp)
	if err != nil {
		return result.RPError, nil, err
	}

	code := cmn.ToInt(resp["code"])
	if code != 200 {
		msg := cmn.ToString(resp["message"])
		switch msg {
		case "仅限新人领取":
			return result.OnlyForNewUser, nil, errors.New(msg)
		case "红包已被领完":
			return result.RPEmpty, nil, errors.New(msg)
		case "红包已领取":
			return result.RPReceived, nil, errors.New(msg)
		default:
			return result.RPError, nil, errors.New(msg)
		}
	}

	ret := make(map[string]interface{})
	ret["coin"] = cmn.ToInt(resp["coin"])
	ret["total"] = cmn.ToInt(resp["total"])
	ret["remain"] = cmn.ToInt(resp["remain"])
	ret["amount"] = cmn.ToFloat32(resp["amount"])
	return 0, ret, nil
}

func Receive(packetId string) (int, interface{}, error) {
	byte, err := cmn.HTTPRequest("GET", cfg.Api.RedPacket+"/red-packet-random/receive?id="+packetId, nil, nil)
	if err != nil {
		return result.RPError, nil, err
	}

	redPackLog.Debug("red-packet-random", "receive", string(byte))
	/*
		{
			"message" : "红包已被领完",
			"code" : 10074
		}
			{
			   "message" : "success",
			   "mark" : "2ca28310-94a6-11e8-b452-631c93934368",
			   "total" : 2,
			   "remain" : 0,
			   "code" : 200,
			   "amount" : "7"
			}
	*/
	var resp map[string]interface{}
	err = json.Unmarshal(byte, &resp)
	if err != nil {
		return result.RPError, nil, err
	}

	code := cmn.ToInt(resp["code"])
	if code != 200 {
		msg := cmn.ToString(resp["message"])
		switch msg {
		case "红包已被领完":
			return result.RPEmpty, nil, errors.New(msg)
		default:
			return result.RPError, nil, errors.New(msg)
		}
	}

	ret := make(map[string]interface{})
	ret["total"] = cmn.ToInt(resp["total"])
	ret["remain"] = cmn.ToInt(resp["remain"])
	ret["amount"] = cmn.ToFloat32(resp["amount"])
	ret["mark"] = cmn.ToString(resp["mark"])
	return 0, ret, nil
}

func Entry(mark, account string) (int, interface{}, error) {
	params := url.Values{}
	params.Set("mark", mark)
	params.Set("account", account)
	byte, err := cmn.HTTPPostForm(cfg.Api.RedPacket+"/red-packet-random/entry", nil, strings.NewReader(params.Encode()))
	if err != nil {
		return result.RPError, nil, err
	}

	redPackLog.Debug("red-packet-random", "entry", string(byte))

	/*

		{
			"code" : 10072,
			"message" : "红包已领取"
		}

		{
			"message" : "仅限新人领取",
			"code" : 10072
		}

		{
			"code" : 10075,
			"message" : "红包标识不匹配",
		}

		{
			"code" : 10027,
			"red_packet_id" : "76292a60-9ea2-11e8-851a-b7c93f977214",
			"coin" : 2,
			"amount" : "3",
			"message" : "红包已领取",
			"user_id" : 200093
		}

		{
			"red_packet_id" : "b84e7d70-94a0-11e8-b452-631c93934368",
			"amount" : "7",
			"mark" : "72a15da0-94a6-11e8-b452-631c93934368",
			"code" : 10038,
			"message" : "用户未注册，请注册后继续",
			"coin" : 1
		}

		{
			"message" : "入账信息已提交，红包将于5五分钟内入账",
			"user_id" : "200463",
			"coin" : 1,
			"mobile" : "18668169201",
			"code" : 200,
			"amount" : "7",
			"red_packet_id" : "b84e7d70-94a0-11e8-b452-631c93934368"
		}
	*/
	var resp map[string]interface{}
	err = json.Unmarshal(byte, &resp)
	if err != nil {
		return result.RPError, nil, err
	}
	code := cmn.ToInt(resp["code"])

	ret := make(map[string]interface{})
	switch code {
	case 200:
		ret["recv_type"] = types.RecvForOldCustomer
		ret["mark"] = ""
		ret["coin"] = cmn.ToInt(resp["coin"])
		ret["amount"] = cmn.ToFloat32(resp["amount"])
		ret["packet_id"] = cmn.ToString(resp["red_packet_id"])
		ret["message"] = cmn.ToString(resp["message"])
		return 0, ret, nil
	case 10038:
		ret["recv_type"] = types.RecvForNewCustomer
		ret["mark"] = cmn.ToString(resp["mark"])
		ret["coin"] = cmn.ToInt(resp["coin"])
		ret["amount"] = cmn.ToFloat32(resp["amount"])
		ret["packet_id"] = cmn.ToString(resp["red_packet_id"])
		ret["message"] = cmn.ToString(resp["message"])
		return 0, ret, nil
	default:
		msg := cmn.ToString(resp["message"])
		switch msg {
		case "仅限新人领取":
			return result.OnlyForNewUser, nil, errors.New(msg)
		case "红包已被领完":
			return result.RPEmpty, nil, errors.New(msg)
		case "标识不匹配", "红包标识不匹配":
			return result.RPIdNotMatch, nil, errors.New(msg)
		case "红包已领取":
			return result.RPReceived, nil, errors.New(msg)
		default:
			return result.RPError, nil, errors.New(msg)
		}
	}
}

func RegisterEntry(mark, account, password, captcha string) (int, interface{}, error) {
	params := url.Values{}
	params.Set("mark", mark)
	params.Set("account", account)
	params.Set("types", "sms")
	params.Set("password", password)
	params.Set("captcha", captcha)
	byte, err := cmn.HTTPPostForm(cfg.Api.RedPacket+"/red-packet-random/register-entry", nil, strings.NewReader(params.Encode()))
	if err != nil {
		return result.RPError, nil, err
	}

	redPackLog.Debug("red-packet-random", "register-entry", string(byte))

	/*
		{
			"message" : "标识不匹配",
			"code" : 10071
		}

		{
		   "message" : "入账信息已提交，红包将于5五分钟内入账",
		   "user_id" : "200463",
		   "coin" : 1,
		   "mobile" : "18668169201",
		   "code" : 200,
		   "amount" : "7",
		   "red_packet_id" : "b84e7d70-94a0-11e8-b452-631c93934368"
		}
	*/
	var resp map[string]interface{}
	err = json.Unmarshal(byte, &resp)
	if err != nil {
		return result.RPError, nil, err
	}

	code := cmn.ToInt(resp["code"])
	if code != 200 {
		msg := cmn.ToString(resp["message"])
		switch msg {
		case "该手机号码已经被注册,请更换":
			return result.UserHasReg, nil, errors.New(msg)
		case "仅限新人领取":
			return result.OnlyForNewUser, nil, errors.New(msg)
		case "红包已被领完":
			return result.RPEmpty, nil, errors.New(msg)
		case "标识不匹配", "红包标识不匹配":
			return result.RPIdNotMatch, nil, errors.New(msg)
		case "验证码不正确":
			return result.VerifyCodeError, nil, errors.New(msg)
		case "验证码已经过期或者已使用":
			return result.VerifyCodeExpired, nil, errors.New(msg)
		default:
			return result.RPError, nil, errors.New(msg)
		}
	}

	ret := make(map[string]interface{})
	ret["zb_uid"] = cmn.ToString(resp["user_id"])
	ret["coin"] = cmn.ToInt(resp["coin"])
	ret["amount"] = cmn.ToFloat32(resp["amount"])
	ret["mobile"] = account
	ret["packet_id"] = cmn.ToString(resp["red_packet_id"])
	return 0, ret, nil
}

// 获取红包统计信息
func RedEnvelopeStatistics(appID string) (interface{}, error) {
	rows, err := db.GetAppsPackets(appID)
	if err != nil {
		return nil, err
	}

	var advNum int
	var todayTotalNum int
	var todayAdvNum int

	now := time.Now()
	totayTimeStart := cmn.MillionSecond(cmn.BeginSecOfDay(now))
	totayTimeEnd := cmn.MillionSecond(cmn.EndSecOfDay(now))
	for _, v := range rows {
		packetType := cmn.ToInt(v["type"])
		if packetType == types.PacketTypeAdv {
			advNum++
		}

		packetTimeMs := cmn.ToInt64(v["created_at"])
		if packetTimeMs >= totayTimeStart && packetTimeMs <= totayTimeEnd {
			todayTotalNum++
			if packetType == types.PacketTypeAdv {
				todayAdvNum++
			}
		}
	}

	ret := &types.RedPacketStatistics{
		AdvNum:        advNum,
		TotalNum:      len(rows),
		TodayAdvNum:   todayAdvNum,
		TodayTotalNum: todayTotalNum,
	}
	return ret, nil
}

// 获取红包信息列表
func RedEnvelopeInfoList(params *types.PacketQueryParam) (interface{}, error) {
	infoList, err := db.GetAppsPacketsWithFilter(params)
	if err != nil {
		redPackLog.Error("red packert info list", "err", err.Error())
		return nil, err
	}

	if infoList == nil {
		return nil, nil
	}

	for _, v := range infoList.Packets {
		rows, err := db.GetUserInfoByID(v.SendUid)
		if err != nil {
			redPackLog.Error("get user info", "err", err.Error())
			return nil, err
		}
		if len(rows) > 0 {
			v.SendAccount = rows[0]["account"]
		}

		recvDetails, err := GetPacketRecvDetails(v.PacketId)
		if err != nil {
			redPackLog.Error("get packet receive details", "err", err.Error())
			return nil, err
		}

		v.ReceiveNum = len(recvDetails)
		var newUserNum int
		for _, e := range recvDetails {
			if e.RecvType == types.RecvForNewCustomer {
				newUserNum++
			}
		}
		v.NewUserNum = newUserNum
		// v.BackNum = TODO
	}
	return infoList, nil
}

// 获取红包详情
func RedEnvelopeDetail(packetId string) (interface{}, error) {
	rows, err := db.GetRedPacket(packetId)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, errors.New("no packet log")
	}

	data := make(map[string]interface{})
	data["app_uid"] = rows[0]["uid"]      // zhaobi uid
	data["send_uid"] = rows[0]["user_id"] // chat uid
	data["send_account"] = rows[0]["username"]
	data["packet_type"] = cmn.StringToInt32(rows[0]["type"])
	data["time"] = cmn.ToInt64(rows[0]["created_at"])
	data["currency"] = cmn.StringToInt32(rows[0]["coin"])
	data["amount"] = cmn.StringToInt32(rows[0]["amount"])
	data["size"] = cmn.StringToInt32(rows[0]["size"])
	data["remark"] = rows[0]["remark"]
	data["avatar"] = utility.ToString(rows[0]["avatar"])

	recvDetails, err := GetPacketRecvDetails(packetId)
	if err != nil {
		redPackLog.Error("get packet receive details", "err", err.Error())
		return nil, err
	}

	data["recv_details"] = recvDetails
	return data, nil
}

func GetPacketRecvDetails(packetId string) ([]*types.RecvDetail, error) {
	bytes, err := cmn.HTTPRequest("GET", cfg.Api.RedPacket+"/red-packet-random/record?id="+packetId, nil, nil)
	if err != nil {
		return nil, err
	}

	/*
		{
		   "data" : [
		      {
		         "username" : "861866****201",
		         "amount" : 7,
		         "created_at" : "2018-07-31T10:22:34.000Z",
		         "user_id" : 200463,
		         "coin" : 1,
		         "type" : 2    //新用户
		      },
		      {
		         "type" : 1,  //老用户
		         "user_id" : 200093,
		         "coin" : 1,
		         "amount" : 3,
		         "username" : "861385****274",
		         "created_at" : "2018-07-31T09:07:14.000Z"
		      }
		   ],
		   "message" : "success",
		   "code" : 200
		}
	*/
	var recvResp map[string]interface{}
	err = json.Unmarshal(bytes, &recvResp)
	if err != nil {
		return nil, err
	}

	recvDetails := []*types.RecvDetail{}
	records := recvResp["data"].([]interface{})
	for _, e := range records {
		v := e.(map[string]interface{})

		zbUid := cmn.ToString(v["user_id"]) //zhaobi uid
		rows, err := db.GetUserInfoByUID(zbUid)
		if err != nil {
			redPackLog.Error("get user info", "zbUid", zbUid, "err", err.Error())
		}

		var recvUid int
		if len(rows) > 0 {
			recvUid = cmn.ToInt(rows[0]["user_id"])
		} else {
			// maybe this zhaobi user hasn't used the chat functionality, so we cannot query a record from user tb
			redPackLog.Info("no record in tb user", "zbUid", zbUid)
			continue
		}

		one := &types.RecvDetail{
			RecvUid:     recvUid,
			AppUid:      cmn.ToInt(zbUid),
			RecvAccount: v["username"].(string),
			RecvType:    cmn.ToInt(v["type"]),
			RecvTime:    cmn.CstTime(v["created_at"].(string)).UnixNano() / 1e6,
			Amount:      cmn.ToInt(v["amount"]),
			Avatar:      utility.ToString(rows[0]["avatar"]),
		}
		recvDetails = append(recvDetails, one)
	}

	return recvDetails, nil
}

type RecHistory struct {
	RedPacketId string  `json:"red_packet_id"`
	Type        int     `json:"type"`
	Amount      float64 `json:"amount"`
	Coin        int     `json:"coin"`
	ReceiveTime int64   `json:"receive_time"`
	Username    string  `json:"username"`
}

func GetRPReceiveHistoryByMonth(token, date string, coin int) (map[string]interface{}, error) {
	zbCoin, ok := CoinToZBCoin[coin]
	if !ok {
		return nil, errors.New("币种类型错误")
	}

	params := url.Values{}
	params.Set("coin", utility.ToString(zbCoin))
	params.Set("date", date)

	headers := make(map[string]string)
	headers["auth-token"] = token
	bytes, err := cmn.HTTPRequest("GET", cfg.Api.RedPacket+"/history/receive-record?"+params.Encode(), headers, nil)
	if err != nil {
		redPackLog.Error("RPReceiveHistory http request error", "err_msg", err)
		return nil, errors.New("请求错误")
	}
	var resp map[string]interface{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		redPackLog.Error("RPReceiveHistory unmarshal error", "err_msg", err, "token", token, "coin", coin, "date", date)
		return nil, errors.New("")
	}

	if utility.ToInt(resp["code"]) != 200 {
		errMsg := utility.ToString(resp["message"])
		return nil, errors.New(errMsg)
	}

	var recList = make([]RecHistory, 0)
	var amount float64 = 0

	rows := resp["rows"].([]interface{})
	for _, row := range rows {
		r := row.(map[string]interface{})
		entry := r["entry_record"].(map[string]interface{})
		rp := r["red_packet"].(map[string]interface{})
		tmp := RecHistory{}
		tmp.Amount = cmn.ToFloat64(entry["amount"])
		tmp.Coin = coin
		tmp.RedPacketId = utility.ToString(rp["id"])
		tmp.Type = utility.ToInt(entry["type"])
		tmp.ReceiveTime = utility.RFC3339ToTimeStampMillionSecond(utility.ToString(entry["created_at"]))
		tmp.Username = utility.ToString(rp["username"])
		recList = append(recList, tmp)

		amount += tmp.Amount
	}

	var data = make(map[string]interface{})
	data["totalnum"] = resp["count"]
	data["total_amount"] = amount
	data["receive"] = recList
	return data, nil
}

func GetRPReceiveHistoryByYear(token, date string, coin int) (map[string]interface{}, error) {
	var amount float64 = 0
	var totalNum = 0
	var recList = make([]RecHistory, 0)
	month := int(time.Now().Month())
	for i := month; i > 0; i-- {
		data, err := GetRPReceiveHistoryByMonth(token, date+fmt.Sprintf("-%02d", i), coin)
		if err != nil {
			return nil, err
		}
		amount += data["total_amount"].(float64)
		totalNum += utility.ToInt(data["totalnum"])
		recList = append(recList, data["receive"].([]RecHistory)...)
	}
	data := make(map[string]interface{})
	data["totalnum"] = totalNum
	data["total_amount"] = amount
	data["receive"] = recList
	return data, nil
}

func RPReceiveHistory(token, date string, coin int) (interface{}, *result.Error) {
	var data map[string]interface{}
	var err error

	const matchStr = `^[0-9]+$`
	r, _ := regexp.Compile(matchStr)
	if !r.MatchString(date) {
		data, err = GetRPReceiveHistoryByMonth(token, date, coin)
	} else {
		data, err = GetRPReceiveHistoryByYear(token, date, coin)
	}
	if err != nil {
		return nil, &result.Error{ErrorCode: result.RPError, Message: ""}
	}
	return data, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

type SendHistory struct {
	RedPacketId string  `json:"red_packet_id"`
	Type        int     `json:"type"`
	Amount      float64 `json:"amount"`
	Coin        int     `json:"coin"`
	SendTime    int64   `json:"send_time"`
}

func GetRPSendHistoryByMonth(token, date string, coin int) (map[string]interface{}, error) {
	zbCoin, ok := CoinToZBCoin[coin]
	if !ok {
		return nil, errors.New("币种类型错误")
	}

	params := url.Values{}
	params.Set("coin", utility.ToString(zbCoin))
	params.Set("date", date)

	headers := make(map[string]string)
	headers["auth-token"] = token
	bytes, err := cmn.HTTPRequest("GET", cfg.Api.RedPacket+"/history/send-record?"+params.Encode(), headers, nil)
	if err != nil {
		redPackLog.Error("RPSendHistory http request error", "err_msg", err)
		return nil, errors.New("请求错误")
	}
	var resp map[string]interface{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		redPackLog.Error("RPSendHistory unmarshal error", "err_msg", err, "token", token, "coin", coin, "date", date)
		return nil, errors.New("")
	}

	if utility.ToInt(resp["code"]) != 200 {
		errMsg := utility.ToString(resp["message"])
		return nil, errors.New(errMsg)
	}

	var sendList = make([]SendHistory, 0)
	var amount float64 = 0

	rows := resp["rows"].([]interface{})
	for _, row := range rows {
		r := row.(map[string]interface{})
		tmp := SendHistory{}
		tmp.Amount = cmn.ToFloat64(r["amount"])
		tmp.Coin = coin
		tmp.RedPacketId = utility.ToString(r["id"])
		tmp.Type = utility.ToInt(r["type"])
		tmp.SendTime = utility.RFC3339ToTimeStampMillionSecond(utility.ToString(r["created_at"]))
		sendList = append(sendList, tmp)

		amount += tmp.Amount
	}

	var data = make(map[string]interface{})
	data["totalnum"] = resp["count"]
	data["send"] = sendList
	data["total_amount"] = amount
	return data, nil
}

func GetRPSendHistoryByYear(token, date string, coin int) (map[string]interface{}, error) {
	var amount float64 = 0
	var totalNum = 0
	var sendList = make([]SendHistory, 0)
	month := int(time.Now().Month())
	for i := month; i > 0; i-- {
		data, err := GetRPSendHistoryByMonth(token, date+fmt.Sprintf("-%02d", i), coin)
		if err != nil {
			return nil, err
		}
		amount += data["total_amount"].(float64)
		totalNum += utility.ToInt(data["totalnum"])
		sendList = append(sendList, data["send"].([]SendHistory)...)
	}
	data := make(map[string]interface{})
	data["totalnum"] = totalNum
	data["total_amount"] = amount
	data["send"] = sendList
	return data, nil
}

func RPSendHistory(token, date string, coin int) (interface{}, *result.Error) {
	var data map[string]interface{}
	var err error

	const matchStr = `^[0-9]+$`
	r, _ := regexp.Compile(matchStr)
	if !r.MatchString(date) {
		data, err = GetRPSendHistoryByMonth(token, date, coin)
	} else {
		data, err = GetRPSendHistoryByYear(token, date, coin)
	}
	if err != nil {
		return nil, &result.Error{ErrorCode: result.RPError, Message: ""}
	}
	return data, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func RPHistory(token, date string, coin int) (interface{}, *result.Error) {
	var receiveData map[string]interface{}
	var sendData map[string]interface{}
	var err error

	const matchStr = `^[0-9]+$`
	r, _ := regexp.Compile(matchStr)
	if !r.MatchString(date) {
		receiveData, err = GetRPReceiveHistoryByMonth(token, date, coin)
		if err != nil {
			return nil, &result.Error{ErrorCode: result.RPError, Message: ""}
		}
		sendData, err = GetRPSendHistoryByMonth(token, date, coin)
	} else {
		receiveData, err = GetRPReceiveHistoryByYear(token, date, coin)
		if err != nil {
			return nil, &result.Error{ErrorCode: result.RPError, Message: ""}
		}
		sendData, err = GetRPSendHistoryByYear(token, date, coin)
	}

	if err != nil {
		return nil, &result.Error{ErrorCode: result.RPError, Message: ""}
	}

	var data = make(map[string]interface{})
	data["receive_num"] = receiveData["totalnum"]
	data["receive_amount"] = receiveData["total_amount"]
	data["receive"] = receiveData["receive"]

	data["send_num"] = sendData["totalnum"]
	data["send_amount"] = sendData["total_amount"]
	data["send"] = sendData["send"]
	return data, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}
