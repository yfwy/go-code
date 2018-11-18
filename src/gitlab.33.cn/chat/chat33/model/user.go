package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	cmn "dev.33.cn/33/common"
	"github.com/astaxie/beego/orm"
	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	dbtools "gitlab.33.cn/chat/chat33/db"
	"gitlab.33.cn/chat/chat33/result"
	"gitlab.33.cn/chat/chat33/types"
	"gitlab.33.cn/chat/chat33/utility"
)

var logUser = log15.New("module", "model/user")

type userStatistics struct {
	Muted_num   int `json:"muted_num"`
	Lastday_num int `json:"lastday_num"`
	Kickout_num int `json:"kickout_num"`
	Active_num  int `json:"active_num"`
}

type ReportInfo struct {
	Id       string      `json:"id"`
	Uid      string      `json:"uid"`
	Name     string      `json:"name"`
	Account  string      `json:"account"`
	MsgType  int         `json:"msg_type"`
	Content  interface{} `json:"content"`
	Datetime int64       `json:"datetime"`
}

type ReportInfoList struct {
	Msg_list []ReportInfo `json:"msg_list"`
}

type Permission struct {
	PermissionId   int    `json:"permission_id"`
	PermissionName string `json:"permission_name"`
}

type PermissionList struct {
	PermissionList []Permission `json:"permission_list"`
}

// 根据uid和appid获取用户id
func GetUserID(uid, appid string) (string, error) {
	rows, err := dbtools.GetUserInfoByUID(uid)
	if err != nil {
		//输出错误日志
		return "", err
	}

	if len(rows) > 0 {
		_id, ok := rows[0]["user_id"]
		if !ok {
			return "", nil
		}
		//
		return _id, nil
	} else {
		return "", nil
	}
}

// 根据id查询用户当前是否被禁言
func CheckUserMutedById(id string) (bool, error) {
	return CheckUserMuted(id, utility.NowMillionSecond())
}

// 根据用户Id和当前查询时间查询用户是否被禁言
func CheckUserMuted(id string, timestamp int64) (bool, error) {
	rows, err := dbtools.GetMutedInfo(id, timestamp)
	if err != nil {
		//输出错误日志
		return true, err
	}
	if len(rows) > 0 {
		return true, nil
	}

	return false, nil
}

func CheckUserKickoutById(id string) (bool, error) {
	return CheckUserKickout(id, utility.NowMillionSecond())
}

func CheckUserKickout(id string, timestamp int64) (bool, error) {
	rows, err := dbtools.GetKickOutInfo(id, timestamp)
	if err != nil {
		return true, err
	}
	if len(rows) > 0 {
		return true, nil
	}
	return false, nil
}

// 根据接收者id 查看接收者是否是客服
func CheckUserIsCS(id, appId string) (bool, error) {
	rows, err := dbtools.GetCSInfoByID(id, appId)
	if err != nil {
		//输出错误日志
		return false, err
	}

	if len(rows) > 0 {
		return true, nil
	}

	return false, nil
}

func UserCanSendSysMsg(id, appID string) (bool, error) {
	ok, err := HasPermission(id, appID, SendSystemMsg)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// 查询是否是客服且具有聊天室管理权限
func UserCanManageGroup(id string, appID string) (bool, error) {
	ok, err := HasPermission(id, appID, ManageGroup)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func GetUserIdByUidAndAppid(uid, appid string) (string, bool) {
	o := orm.NewOrm()

	var maps []orm.Params
	num, err := o.Raw("select user_id from `user` where uid= ? and app_id = ?", uid, appid).Values(&maps)
	if err != nil || num < 1 {
		return "", false
	}

	return maps[0]["user_id"].(string), true
}

type AppPem struct {
	AppId       string `json:"app_id"`
	Permissions []int  `json:"permissions"`
}

func ZbTokenLogin(token, deviceType string) (map[string]interface{}, error) {
	data, err := GetZbUserInfo(token)
	if err != nil {
		return nil, err
	}

	uid := cmn.ToString(data["id"])
	//username := cmn.ToString(data["username"])
	email := cmn.ToString(data["email"])
	area := cmn.ToString(data["area"])
	mobile := cmn.ToString(data["mobile"])
	group := cmn.ToString(data["group"]) // user/kf/admin/frozen/broker

	var account string
	if mobile != "" {
		account = area + mobile
	} else {
		account = email
	}

	verified, err := UserVerified(token)
	if err != nil {
		return nil, err
	}

	var id int64
	var userLevel int
	var username string
	var avatar string
	rows, err := dbtools.GetUserInfoByUID(uid)
	if len(rows) == 0 {
		username = utility.RandomUsername()
		userLevel = types.LevelMember
		_, id, err = dbtools.InsertUser(uid, username, account, email, area+mobile, cmn.ToString(userLevel), cmn.ToString(verified))
		if err != nil {
			return nil, err
		}
		avatar = ""
	} else {
		id = cmn.ToInt64(rows[0]["user_id"])
		// TODO:
		userLevel = cmn.ToInt(rows[0]["user_level"])
		username = cmn.ToString(rows[0]["username"])
		avatar = utility.ToString(rows[0]["avatar"])
	}

	userID := cmn.ToString(id)

	logUser.Info("Token Login", "userID", userID, "uid", uid, "group", group, "account", account, "level", userLevel)

	ret := make(map[string]interface{})
	ret["id"] = userID
	ret["uid"] = uid // zhaobi uid
	ret["account"] = account
	ret["verified"] = verified
	ret["username"] = username
	ret["avatar"] = avatar
	if userLevel == types.LevelAdmin {
		ret["user_level"] = types.LevelCs
	} else {
		ret["user_level"] = userLevel
	}

	if deviceType == types.DeviceWeb {
		if userLevel == types.LevelMember {
			return nil, errors.New("common member not allowed to login from Web")
		}

		/*var pemList = make([]*AppPem, 0)

		if userLevel == types.LevelAdmin {
			rows, err = dbtools.GetAdminPermission(cmn.ToInt(userID))
		} else {
			rows, err = dbtools.GetPermission(cmn.ToInt(userID))
		}
		if err != nil {
			return nil, err
		}

		for _, v := range rows {
			one := &AppPem{
				AppId:       cmn.ToString(v["app_id"]),
				Permissions: GetCsPermissionList(cmn.ToInt(v["permission"])),
			}
			if userLevel == types.LevelAdmin {
				one.Permissions = append(one.Permissions, ManageCS)
			}
			pemList = append(pemList, one)
		}
		ret["permission_list"] = pemList*/
	}

	err = dbtools.AddUserLoginLog(userID, deviceType)
	if err != nil {
		logUser.Error("token login log user login failed", "err_msg", err)
	}
	return ret, nil
}

// 根据找币token返回找币uid信息
func GetZbUserInfo(zbToken string) (map[string]interface{}, error) {
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + zbToken
	byte, err := cmn.HTTPPostForm(cfg.Api.Zhaobi+"/api/member/info", headers, nil)
	if err != nil {
		return nil, err
	}

	/*
		{
		   "error" : "OK",
		   "message" : "OK",
		   "code" : 200,
		   "ecode" : "200",
		   "data" : {
		      "base" : {
		         "email" : "",
		         "wfrom" : "fxee",
		         "ispwd" : "1",
		         "mobile" : "13858075274",
		         "username" : "86138****5274",
		         "area" : "86",
		         "regip" : "122.235.232.165",
		         "adddate" : "2018-05-04 14:54:30",
		         "group" : "kf",
		         "id" : "200093",
		         "addtime" : "1525416870"
		      }
		   }
		}
	*/
	var resp map[string]interface{}
	err = json.Unmarshal(byte, &resp)
	if err != nil {
		return nil, err
	}

	if cmn.ToInt(resp["code"]) != 200 {
		return nil, errors.New(resp["message"].(string))
	}

	data, ok := resp["data"]
	if !ok {
		return nil, errors.New("no 'data' info")
	}

	base, ok := data.(map[string]interface{})["base"]
	if !ok {
		return nil, errors.New("no 'base' info")
	}

	return base.(map[string]interface{}), nil
}

func GetTokenViaPwd(mobile, password string) (string, int) {
	// TODO use ccmn.HTTPPostForm(...)
	resp, err := http.PostForm(cfg.Api.Zhaobi+"/api/member/login", url.Values{
		"type":         {"sms"},
		"area":         {"86"},
		"mobile":       {mobile},
		"password":     {password},
		"redirect_uri": {"chat"},
	})

	if err != nil {
		logUser.Info("getTokenViaPwd", "err_msg", err)
		return "", result.ZhaobiInteractFailed
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logUser.Info("getTokenViaPwd", "err_msg", err)
		return "", result.ZhaobiInteractFailed
	}
	defer resp.Body.Close()
	var bodyData map[string]interface{}
	json.Unmarshal(body, &bodyData)

	code, ok := bodyData["code"]
	if !ok {
		return "", result.ZhaobiInteractFailed
	}
	if code.(float64) != 200 {
		logUser.Info("getTokenViaPwd", "errMsg", bodyData["error"], "msg", bodyData["message"])
		return "", result.ZhaobiInteractFailed
	}
	data, ok := bodyData["data"]
	mapData, ok := data.(map[string]interface{})
	if !ok {
		return "", result.ZhaobiInteractFailed
	}
	token, ok := mapData["access_token"]
	return utility.ToString(token), result.CodeOK
}

// TODO return arg1
func UserPwdLogin(mobile, password, deviceType string) (interface{}, *result.Error) {
	token, code := GetTokenViaPwd(mobile, password)
	if code != result.CodeOK {
		return nil, &result.Error{ErrorCode: code, Message: ""}
	}
	// TODO user already login, why login twice?
	ret, err := ZbTokenLogin(token, deviceType)
	if err != nil {
		return nil, &result.Error{ErrorCode: result.ZhaobiInteractFailed, Message: ""}
	}
	ret["token"] = token
	return ret, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

// 用户是否实名认证 2: 未认证 1: 已认证
func UserVerified(token string) (int, error) {
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + token
	byte, err := cmn.HTTPPostForm(cfg.Api.Zhaobi+"/api/certification/identityinfo", headers, nil)
	if err != nil {
		return 2, err
	}

	/*
		{
		   "message" : "OK",
		   "error" : "OK",
		   "ecode" : "200",
		   "code" : 200,
		   "data" : {
		      "cardid" : null,
		      "image" : null,
		      "name" : null,
		      "state" : 0
		   }
		}
	*/
	var resp map[string]interface{}
	err = json.Unmarshal(byte, &resp)
	if err != nil {
		return 2, err
	}

	if cmn.ToInt(resp["code"]) != 200 {
		return 2, errors.New(resp["message"].(string))
	}

	data, ok := resp["data"]
	if !ok {
		return 2, errors.New("no 'data' info")
	}

	var ret int
	switch cmn.ToInt(data.(map[string]interface{})["state"]) {
	case 0:
		ret = 2
	case 1:
		ret = 1
	default:
		ret = 2
	}
	return ret, nil
}

//查询用户统计信息
func UserStatistics(appID string) (interface{}, *result.Error) {
	activeNum := 0
	users := make(map[string]bool)
	for _, group := range utility.GroupList {
		for _, devMap := range group {
			for _, user := range devMap {
				if user.GetAppID() != appID {
					if user.GetAppID() != "" {
						continue
					}
					apps, err := dbtools.GetManageAppMapByID(user.GetID())
					if err != nil {
						logUser.Error("UserStatistics GetManageAppMapByID failed", "err_msg", err)
						continue
					}
					_, ok := apps[appID]
					if !ok {
						continue
					}
				}
				userID := user.GetID()
				users[userID] = true
			}
		}
	}

	activeNum = len(users)

	// LastDay num
	lastDayNum, err := dbtools.GetUserLastDayNumWithAppID(appID)
	if err != nil {
		logUser.Error("UserStatistics count lastday num", "err_msg", err)
		return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}

	// muted_num
	mutedNum, err := dbtools.CountMutedNumWithAppID(appID)
	if err != nil {
		logUser.Error("UserStatistics count muted_num", "err_msg", err)
		return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}

	// kickout_num
	kickOutNum, err := dbtools.CountKickOutNumWithAppID(appID)
	if err != nil {
		logUser.Error("UserStatistics count kickout num", "err_msg", err)
		return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}

	return userStatistics{utility.ToInt(mutedNum), utility.ToInt(lastDayNum),
		utility.ToInt(kickOutNum), utility.ToInt(activeNum)}, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

// 根据用户id查询用户信息
func UserDetail(id string) (map[string]interface{}, error) {
	maps, err := dbtools.GetUserInfoWithID(id)
	if err != nil {
		logUser.Error("Get User Detail Database error", "err", err)
		return nil, err
	}
	if len(maps) < 1 {
		return nil, errors.New("No such user")
	}
	//
	data := make(map[string]interface{})
	data["user_id"] = maps[0]["user_id"]
	data["uid"] = maps[0]["uid"]
	data["user_level"] = maps[0]["user_level"]
	//data["remark"], _ = maps[0]["remark"]
	data["avatar"] = maps[0]["avatar"]
	if username := maps[0]["username"]; username != "" {
		data["username"] = username
	} else {
		data["username"] = maps[0]["account"]
	}

	return data, nil
}

type UsersInfo struct {
	Id             string `json:"id"`
	Uid            string `json:"uid"`
	Account        string `json:"account"`
	Username       string `json:"username"`
	GroupName      string `json:"group_name"`
	LivingTime     int64  `json:"living_time"`
	MutedTime      int64  `json:"muted_time"`
	MutedLastTime  int64  `json:"muted_last_time"`
	RemoveTime     int64  `json:"remove_time"`
	RemoveLastTime int64  `json:"remove_last_time"`
	Remark         string `json:"remark"`
	RemarkBy       string `json:"remark_by"`
	Description    string `json:"description"`
	PicUri         string `json:"pic_uri"`
	Verified       int    `json:"verified"`
	UserLevel      int    `json:"user_level"`
	MutedNum       int    `json:"muted_num"`
	KickoutNum     int    `json:"kickout_num"`
	Avatar         string `json:"avatar"`
}

type UsersInfoList struct {
	Totalnum int         `json:"totalnum"`
	UserList []UsersInfo `json:"user_list"`
}

func UserInfo(id string) (interface{}, *result.Error) {
	o := orm.NewOrm()
	var maps []orm.Params
	var err error
	var num int64

	now := utility.NowMillionSecond()

	//游客
	matchStr := `^[0-9]*$`
	r, _ := regexp.Compile(matchStr)
	if !r.MatchString(id) {
		//数据表中找不到用户 查询游客
		if devMap, ok := utility.Usermap[id]; ok {
			var _appId string
			for _, v := range devMap {
				if v.GetDevice() == "Web" {
					continue
				}
				_appId = v.GetAppID()
			}
			var data = make(map[string]interface{})
			lastPrivMsg, err := dbtools.GetLastPrivMsg(id)
			if err != nil {
				logUser.Error("UserInfoList  select priv msg failed", "err_msg", err)
				return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
			}
			data["id"] = id
			data["uid"] = ""
			data["app_id"] = _appId
			data["account"] = id
			data["username"] = utility.VisitorNameSplit(id)
			data["muted_time"] = 0
			data["muted_last_time"] = 0
			data["remove_time"] = 0
			data["remove_last_time"] = 0
			data["remark"] = ""
			data["remark_by"] = ""
			data["description"] = ""
			data["pic_uri"] = ""
			data["user_level"] = 0
			data["muted_num"] = 0
			data["verified"] = 0
			data["kickout_num"] = 0
			data["avatar"] = ""

			data["last_priv_msg"] = lastPrivMsg
			// construct data
			return data, &result.Error{ErrorCode: result.CodeOK, Message: ""}
		} else {
			logUser.Warn("UserInfoList select from user", "msg", "num < 1", "user_id", id)
			return nil, &result.Error{ErrorCode: result.UserNotExists, Message: ""}
		}
	}

	// select muted
	num, err = o.Raw("select user_id, end_time, operate_time from muted where user_id=?", id).Values(&maps)
	if err != nil {
		logUser.Error("UserInfo select from muted failed", "err_msg", err)
		return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}
	var mutedTime, mutedLastTime int64
	if num < 1 {
		mutedTime = 0
		mutedLastTime = 0
	} else {
		optTime := maps[0]["operate_time"]
		endTime := maps[0]["end_time"]
		optTimeStamp := utility.ToInt64(optTime)
		endTimeStamp := utility.ToInt64(endTime)
		if endTimeStamp <= now {
			mutedLastTime = 0
			mutedTime = 0
		} else {
			mutedTime = optTimeStamp
			mutedLastTime = endTimeStamp - optTimeStamp
		}
	}

	// select kickout
	num, err = o.Raw("select user_id, end_time, operate_time from kickout where user_id=?", id).Values(&maps)
	if err != nil {
		logUser.Error("UserInfo  select from kickout failed", "err_msg", err)
		return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}
	var kickoutTime, kickoutLastTime int64
	if num < 1 {
		kickoutTime = 0
		kickoutLastTime = 0
	} else {
		optTime := maps[0]["operate_time"]
		endTime := maps[0]["end_time"]
		optTimeStamp := utility.ToInt64(optTime)
		endTimeStamp := utility.ToInt64(endTime)
		if endTimeStamp <= now {
			kickoutTime = 0
			kickoutLastTime = 0
		} else {
			kickoutTime = optTimeStamp
			kickoutLastTime = endTimeStamp - optTimeStamp
		}
	}
	// account
	num, err = o.Raw("select * from `user` where user_id=?", id).Values(&maps)
	if err != nil {
		logUser.Error("UserInfoList  select from user failed", "err_msg", err)
		return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}
	if num < 1 {
		logUser.Warn("UserInfoList select from user", "msg", "num < 1", "user_id", id)
		return nil, &result.Error{ErrorCode: result.UserNotExists, Message: ""}
	}
	acc := maps[0]["account"]

	var data = make(map[string]interface{})

	lastPrivMsg, err := dbtools.GetLastPrivMsg(id)
	if err != nil {
		logUser.Error("UserInfoList  select priv msg failed", "err_msg", err)
		return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}

	data["id"] = id
	data["uid"] = utility.ToString(maps[0]["uid"])
	data["app_id"] = utility.ToString(maps[0]["app_id"])
	data["account"] = utility.ToString(acc)
	data["username"] = utility.ToString(maps[0]["username"])
	data["muted_time"] = mutedTime
	data["muted_last_time"] = mutedLastTime
	data["remove_time"] = kickoutTime
	data["remove_last_time"] = kickoutLastTime
	data["remark"] = utility.ToString(maps[0]["remark"])
	data["remark_by"] = utility.ToString(maps[0]["remark_cs"])
	data["description"] = utility.ToString(maps[0]["description"])
	data["pic_uri"] = utility.ToString(maps[0]["pic_uri"])
	data["avatar"] = utility.ToString(maps[0]["avatar"])
	userLevel := utility.ToInt(maps[0]["user_level"])
	if userLevel == 3 {
		userLevel = 2
	}
	data["user_level"] = userLevel
	data["muted_num"] = utility.ToInt(maps[0]["muted_num"])
	data["verified"] = utility.ToInt(maps[0]["verified"])
	data["kickout_num"] = utility.ToInt(maps[0]["kickout_num"])

	data["last_priv_msg"] = lastPrivMsg
	// construct data
	return data, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func UserInfoList(appID string, userType int, uid, account string, page, number int) (interface{}, *result.Error) {
	o := orm.NewOrm()
	var err error
	var num int64
	var users = make(map[string]utility.Client)
	var data = make([]UsersInfo, 0)

	var allUsers []orm.Params
	var maps []orm.Params

	sqlStr := "SELECT u.user_id, u.uid, u.app_id, u.account, u.remark, u.user_level, u.username, u.avatar, " +
		" u.verified,u.remark_cs, u.description, u.pic_uri, m.operate_time as mute_time, m.end_time as mute_end_time," +
		" k.operate_time as kick_time, k.end_time as kick_end_time, IFNULL(muted_num, 0) as muted_num, " +
		" IFNULL(kickout_num, 0) as kickout_num from `user` as u " +
		" LEFT JOIN muted as m on u.user_id=m.user_id LEFT JOIN kickout as k on u.user_id=k.user_id where u.app_id = ?" +
		" and u.user_level < 2"

	if uid != "" {
		sqlStr += " and u.uid like '%" + uid + "%'"
	}

	if account != "" {
		sqlStr += " and u.account like '%" + account + "%'"
	}

	num, err = o.Raw(sqlStr, appID).Values(&allUsers)

	if err != nil {
		logUser.Error("UserInfoList query all user failed", "err_msg", err)
		return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}

	for _, group := range utility.GroupList {
		/*ok := false
		for _, user := range group {
			if user.GetAppID() == appID {
				ok = true
			}
			break
		}*/
		for _, devMap := range group {
			for _, user := range devMap {
				if user.GetAppID() != appID {
					continue
				}
				userID := user.GetID()
				users[userID] = user
			}
		}
	}

	l := utility.ToInt(num)
	for i := 0; i < l; i++ {
		now := utility.NowMillionSecond()

		// mute
		var mutedTime, mutedLastTime int64
		optTime := allUsers[i]["mute_time"]
		if optTime == nil {
			if userType == 1 {
				continue
			}
			mutedTime = 0
			mutedLastTime = 0
		} else {
			optTimeStamp := utility.ToInt64(optTime)
			endTimeStamp := utility.ToInt64(allUsers[i]["mute_end_time"])
			if endTimeStamp <= now {
				if userType == 1 {
					continue
				}
				mutedTime = 0
				mutedLastTime = 0
			} else {
				mutedTime = optTimeStamp
				mutedLastTime = endTimeStamp - optTimeStamp
			}
		}

		// kickout
		var kickoutTime, kickoutLastTime int64
		optTime = allUsers[i]["kick_time"]
		if optTime == nil {
			if userType == 2 {
				continue
			}
			kickoutTime = 0
			kickoutLastTime = 0
		} else {
			optTimeStamp := utility.ToInt64(optTime)
			endTimeStamp := utility.ToInt64(allUsers[i]["kick_end_time"])
			if endTimeStamp <= now {
				if userType == 2 {
					continue
				}
				kickoutTime = 0
				kickoutLastTime = 0
			} else {
				kickoutTime = optTimeStamp
				kickoutLastTime = endTimeStamp - optTimeStamp
			}
		}

		// user id
		userID := utility.ToString(allUsers[i]["user_id"])

		// group name && living time
		var groupName = ""
		var livingTime int64 = 0
		if u, ok := users[userID]; ok {
			num, err = o.Raw("select group_name from `group` where group_id=?", u.GetGroupID()).Values(&maps)
			if err != nil {
				logUser.Error("UserInfoList  select from group failed", "err_msg", err)
				return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
			}

			if num > 0 {
				//logUser.Warn("UserInfoList select from group", "msg", "num < 1", "group_id", u.GetGroupID())
				groupName = utility.ToString(maps[0]["group_name"])
			}
			livingTime = now - u.GetLoginTime()
		}

		// account
		acc := utility.ToString(allUsers[i]["account"])

		userLevel, err := dbtools.GetUserLevelWithAppID(userID, appID)
		if err != nil {
			logUser.Error("UserInfoList get userLevel failed", "err_msg", err, "userID", userID, "appID", appID)
		}
		if userLevel == 3 {
			userLevel = 2
		}
		// construct data
		data = append(data, UsersInfo{
			userID,
			utility.ToString(allUsers[i]["uid"]),
			acc,
			utility.ToString(allUsers[i]["username"]),
			groupName,
			livingTime,
			mutedTime,
			mutedLastTime,
			kickoutTime,
			kickoutLastTime,
			utility.ToString(allUsers[i]["remark"]),
			utility.ToString(allUsers[i]["remark_cs"]),
			utility.ToString(allUsers[i]["description"]),
			utility.ToString(allUsers[i]["pic_uri"]),
			utility.ToInt(allUsers[i]["verified"]),
			userLevel,
			utility.ToInt(allUsers[i]["muted_num"]),
			utility.ToInt(allUsers[i]["kickout_num"]),
			utility.ToString(allUsers[i]["avatar"]),
		})
	}
	totalNum := len(data)
	offset := page * number
	if offset >= totalNum {
		return UsersInfoList{totalNum, make([]UsersInfo, 0)}, &result.Error{ErrorCode: result.CodeOK, Message: ""}
	}
	ed := offset + number
	if ed > totalNum {
		ed = totalNum
	}
	return UsersInfoList{totalNum, data[offset:ed]}, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func UserEditInfo(params map[string]string, csID string) *result.Error {
	/*
		用户信息编辑
		要写入操作客服的信息
	*/
	id := params["id"]
	remark, ok := params["remark"]
	if !ok {
		remark = ""
	}
	description, ok := params["description"]
	if !ok {
		description = ""
	}

	picUri, ok := params["pic_uri"]
	if !ok {
		picUri = ""
	}

	_, err := dbtools.UpdateUserInfo(id, remark, description, picUri, csID)
	if err != nil {
		return &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
	}

	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func UserReportInfo(id string) (interface{}, *result.Error) {
	// 返回数据为举报该用户的举报者信息
	rows, err := dbtools.GetReportInfo(id)
	if err != nil {
		return nil, &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
	}
	var data = make([]ReportInfo, 0)
	for _, report := range rows {
		var tmp ReportInfo
		err = json.Unmarshal([]byte(report["content"]), &tmp.Content)
		if err != nil {
			logUser.Warn("unmarshal report info failed", "err_msg", err)
			continue
		}
		//tmp.Content, _ = report["content"]
		date := report["datetime"]
		tmp.Datetime = utility.ToInt64(date)
		tmp.Id = report["user_id"]
		tmp.Account = report["account"]
		tmp.Name = report["username"]
		tmp.Uid = report["uid"]
		tmp.MsgType = utility.ToInt(report["msg_type"])
		data = append(data, tmp)
	}

	return ReportInfoList{data}, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

/*
	举报用户
	sID:被举报者
	oID:举报者
*/
func Report(oID, sID, msgID string) *result.Error {
	err := dbtools.InsertReportInfo(oID, sID, msgID)
	if err != nil {
		logUser.Error("Report write to db failed", "err_msg", err)
		return &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
	}
	num, err := dbtools.CountReportByID(sID)

	if err != nil {
		logUser.Error("Report write to db failed", "err_msg", err)
	}
	const HalfHour = 60 * 30 * 1000
	if num%10 == 0 {
		appID, _ := dbtools.GetAppIDByUserID(sID)
		MuteUser(sID, "", appID, HalfHour)
	}
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func MuteUser(id, csID, appID string, mutedTime int64) *result.Error {
	o := orm.NewOrm()
	var maps []orm.Params

	num, err := o.Raw("select user_id from `user` where user_id=?", id).Values(&maps)
	if err != nil {
		logUser.Error("muteUser query db failed", "err_msg", err)
		return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}
	if num < 1 {
		return &result.Error{ErrorCode: result.UserNotExists, Message: ""}
	}

	optTime := utility.NowMillionSecond()
	endTime := optTime + mutedTime

	_, err = dbtools.MuteUser(id, csID, appID, optTime, endTime)
	if err != nil {
		logUser.Error("muteUser write to db failed", "err_msg", err)
		return &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
	}

	dbtools.RWMutex.Lock()
	defer dbtools.RWMutex.Unlock()
	if mutedTime != 0 {
		_, err := o.Raw("UPDATE `user` set muted_num = IFNULL(muted_num,0)+1 WHERE user_id = ?", id).Exec()
		if err != nil {
			logUser.Error("muteUser write to db failed", "err_msg", err)
			return &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
		}
		var content string
		if endTime > types.TimeForever {
			content = "永久禁言"
		} else {
			mutedTime /= 1000
			seconds := mutedTime % 60
			minutes := mutedTime / 60 % 60
			hours := mutedTime / 3600
			content = fmt.Sprintf("禁言%02d:%02d:%02d", hours, minutes, seconds)
		}
		_, err = o.Raw("insert into operation_log (user_id, cs_id, app_id, content, operate_time) values (?, ?, ?, ?, ?)",
			id, csID, appID, content, optTime).Exec()
		if err != nil {
			logUser.Error("muteUser write to db failed", "err_msg", err)
			return &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
		}
	} else {
		_, err = o.Raw("insert into operation_log (user_id, cs_id, app_id, content, operate_time) values (?, ?, ?, ?, ?)",
			id, csID, appID, "解除禁言", optTime).Exec()
		if err != nil {
			logUser.Error("muteUser write to db failed", "err_msg", err)
			return &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
		}
	}
	//log.Info("mute user", optTime, endTime)

	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func KickOutUser(id, csID, appID string, kickOutTime int64) *result.Error {
	o := orm.NewOrm()
	var maps []orm.Params

	num, err := o.Raw("select user_id from `user` where user_id=?", id).Values(&maps)
	if err != nil {
		logUser.Error("KickOutUser query db failed", "err_msg", err)
		return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}
	if num < 1 {
		return &result.Error{ErrorCode: result.UserNotExists, Message: ""}
	}

	optTime := utility.NowMillionSecond()
	endTime := optTime + kickOutTime

	_, err = dbtools.KickOutUser(id, csID, appID, optTime, endTime)
	if err != nil {
		logUser.Error("KickOutUser write to db failed", "err_msg", err)
		return &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
	}

	dbtools.RWMutex.Lock()
	defer dbtools.RWMutex.Unlock()
	if kickOutTime != 0 {
		_, err := o.Raw("UPDATE `user` set kickout_num = IFNULL(kickout_num,0)+1 WHERE user_id = ?", id).Exec()
		if err != nil {
			logUser.Error("KickOutUser write to db failed", "err_msg", err)
			return &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
		}
		var content string
		if endTime > types.TimeForever {
			content = "永久移出"
		} else {
			kickOutTime /= 1000
			seconds := kickOutTime % 60
			minutes := kickOutTime / 60 % 60
			hours := kickOutTime / 3600
			content = fmt.Sprintf("移出%02d:%02d:%02d", hours, minutes, seconds)
		}
		_, err = o.Raw("insert into operation_log (user_id, cs_id, app_id, content, operate_time) values (?, ?, ?, ?, ?)",
			id, csID, appID, content, optTime).Exec()
		if err != nil {
			logUser.Error("KickOutUser write to db failed", "err_msg", err)
			return &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
		}
	} else {
		_, err = o.Raw("insert into operation_log (user_id, cs_id, app_id, content, operate_time) values (?, ?, ?, ?, ?)",
			id, csID, appID, "解除移出", optTime).Exec()
		if err != nil {
			logUser.Error("KickOutUser write to db failed", "err_msg", err)
			return &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
		}
	}
	//log.Info("kick out user", optTime, endTime)
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func SendSms(mobile string) (interface{}, error) {
	params := url.Values{}
	params.Set("codetype", "validate")
	params.Set("area", "86")
	params.Set("mobile", mobile)
	params.Set("param", "FzmRandom4")
	byte, err := cmn.HTTPPostForm(cfg.Api.Zhaobi+"/api/send/sms2", nil, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	/*
		{
		   "data" : {
		      "gid" : "20180731182203DfAdLd:86:18668169201",
		      "yys" : "Yunpian",
		      "mobile" : "18668169201"
		   },
		   "message" : "OK",
		   "code" : 200,
		   "error" : "OK",
		   "ecode" : "200"
		} .
	*/
	var resp map[string]interface{}
	err = json.Unmarshal(byte, &resp)
	if err != nil {
		return nil, err
	}

	if cmn.ToInt(resp["code"]) != 200 {
		return nil, errors.New(resp["message"].(string))
	}

	return resp, nil
}

func GetCustomServiceList(appid string) []string {
	ret, err := dbtools.GetCSListWithAppID(appid)

	var rlt = make([]string, 0)
	if err == nil {
		for _, v := range ret {
			rlt = append(rlt, v["cs_id"])
		}
	}

	return rlt
}

func IsReg(mobile string) (interface{}, *result.Error) {
	params := url.Values{}
	params.Set("type", "sms")
	params.Set("area", "86")
	params.Set("mobile", mobile)
	byte, err := cmn.HTTPRequest("GET", cfg.Api.Zhaobi+"/api/member/isreg?"+params.Encode(), nil, nil)
	if err != nil {
		return nil, &result.Error{ErrorCode: result.ServerInterError, Message: err.Error()}
	}

	/*
		{
		    "code":200,
		    "ecode":200,
		    "error":"OK",
		    "message":"OK",
		    "data":{
		        "uid":"200018",
		        "ispwd":"1"
		    }
		}
	*/
	var resp map[string]interface{}
	err = json.Unmarshal(byte, &resp)
	if err != nil {
		return nil, &result.Error{ErrorCode: result.ZhaobiInteractFailed, Message: err.Error()}
	}

	if cmn.ToInt(resp["code"]) != 200 {
		return nil, &result.Error{ErrorCode: result.ZhaobiInteractFailed, Message: resp["message"].(string)}
	}

	data, ok := resp["data"]
	if !ok {
		return nil, &result.Error{ErrorCode: result.ZhaobiInteractFailed, Message: "zhaobi /api/member/isreg error"}
	}

	ret := make(map[string]interface{})
	ret["zb_uid"] = cmn.ToString(data.(map[string]interface{})["uid"])
	return ret, nil
}

func CheckEndPointExist(context *gin.Context, client utility.Client) bool {

	/*session, err := utility.SessionStore.Get(r, utility.SESSION_LOGIN)
	if err == nil {
		id, ok := session.Values["id"]
		if ok {
			if id == client.GetId() {
				return true
			}
		}
	}*/
	session := sessions.Default(context)
	id := session.Get("id")
	if id != nil {
		if id.(string) == client.GetID() {
			return true
		}
	}
	return false
}

func MuteUserAll(id, opID string, mutedTime int64) *result.Error {
	//var ret BaseReturn

	// check user exists

	optTime := utility.NowMillionSecond()
	endTime := optTime + mutedTime

	err := dbtools.MuteUserAll(id, opID, optTime, endTime)
	if err != nil {
		logUser.Error("muteUser write to db failed", "err_msg", err)
		return &result.Error{ErrorCode: result.WriteDbFailed, Message: err.Error()}
	}

	// todo: operation log
	/*if mutedTime != 0 {
		var content string
		if endTime > types.TimeForever {
			content = "永久禁言"
		} else {
			mutedTime /= 1000
			seconds := mutedTime % 60
			minutes := mutedTime / 60 % 60
			hours := mutedTime / 3600
			content = fmt.Sprintf("禁言%02d:%02d:%02d", hours, minutes, seconds)
		}
		_, err = o.Raw("insert into operation_log (user_id, cs_id, app_id, content, operate_time) values (?, ?, ?, ?, ?)",
			id, csID, appID, content, optTime).Exec()
		if err != nil {
			logUser.Error("muteUser write to db failed", "err_msg", err)
			return ret.Get(WriteDbFailed, "", nil), err
		}
	} else {
		_, err = o.Raw("insert into operation_log (user_id, cs_id, app_id, content, operate_time) values (?, ?, ?, ?, ?)",
			id, csID, appID, "解除禁言", optTime).Exec()
		if err != nil {
			logUser.Error("muteUser write to db failed", "err_msg", err)
			return ret.Get(WriteDbFailed, "", nil), err
		}
	}*/
	//log.Info("mute user", optTime, endTime)

	return &result.Error{ErrorCode: result.CodeOK, Message: err.Error()}
}
