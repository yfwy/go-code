package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	cmn "dev.33.cn/33/common"
	"github.com/inconshreveable/log15"
	dbtools "gitlab.33.cn/chat/chat33/db"
	m "gitlab.33.cn/chat/chat33/model"
	"gitlab.33.cn/chat/chat33/result"
	"gitlab.33.cn/chat/chat33/types"
	"gitlab.33.cn/chat/chat33/utility"
)

var logUser = log15.New("module", "api/user")

/*
	检查当前请求用户是否有某一权限
*/
func CheckPermission(context *gin.Context, appID string, permission int) bool {
	session := sessions.Default(context)
	idRaw := session.Get("user_id")
	if idRaw == nil {
		return false
	}
	/*session, err := utility.SessionStore.Get(r, utility.SESSION_LOGIN)
	if err != nil {
		return false
	}
	idRaw, ok := session.Values["user_id"]
	if !ok {
		return false
	}*/
	id, ok := idRaw.(string)
	if !ok {
		return false
	}
	ok, _ = m.HasPermission(id, appID, permission)
	return ok
}

/*
	检查当前用户是否具有管理员权限
*/
func CheckAdmin(context *gin.Context) bool {
	/*var err error
	session, err := utility.SessionStore.Get(r, utility.SESSION_LOGIN)
	if err != nil {
		return false
	}
	idRaw, ok := session.Values["user_id"]
	if !ok {
		return false
	}*/
	session := sessions.Default(context)
	idRaw := session.Get("user_id")
	if idRaw == nil {
		return false
	}
	id, ok := idRaw.(string)
	if !ok {
		return false
	}
	ok, _ = m.CheckAdmin(id)
	return ok
}

func UserPwdLogin(c *gin.Context) {
	deviceType := c.GetHeader("FZM-DEVICE")
	if deviceType == "" {
		ret := result.ComposeHttpAck(result.LackParam, "Header: device", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	type pwdParams struct {
		Mobile   string `json:"mobile" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var params pwdParams
	if err := c.ShouldBindJSON(&params); err != nil {
		WriteJson(c, result.ComposeHttpAck(result.ParamsError, "", result.Empty{}))
		return
	}

	data, errs := m.UserPwdLogin(params.Mobile, params.Password, deviceType)
	if errs.ErrorCode != result.CodeOK {
		WriteJson(c, result.ComposeHttpAck(errs.ErrorCode, errs.Message, result.Empty{}))
		return
	}

	session := sessions.Default(c)
	session.Set("id", utility.RandomID())
	session.Set("user_id", data.(map[string]interface{})["id"])
	if deviceType == types.DeviceWeb {
		session.Set("ismanager", true)
	}
	session.Set("token", data.(map[string]interface{})["token"])
	session.Set("devtype", deviceType)
	err := session.Save()
	if err != nil {
		logUser.Error("pwdLogin save session failed", "err", err)
	}

	WriteJson(c, result.ComposeHttpAck(result.CodeOK, "", data))
}

// 根据token登录
func UserTokenLogin(context *gin.Context) {
	token := context.GetHeader("FZM-AUTH-TOKEN")
	if token == "" {
		ret := result.ComposeHttpAck(result.LackParam, "Header: token", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	deviceType := context.GetHeader("FZM-DEVICE")
	if deviceType == "" {
		ret := result.ComposeHttpAck(result.LackParam, "Header: device", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	data, err := m.ZbTokenLogin(token, deviceType)
	if err != nil {
		logUser.Error("zhaobi login", "err", err.Error())
		ret := result.ComposeHttpAck(result.ZhaobiTokenLoginFailed, err.Error(), result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	// //判断设备重复登录
	// //检查设备是否登录
	// if devMap, ok := utility.Usermap[utility.ToString(data["id"])]; ok {
	// 	if deviceType == "Web" {
	// 		if client, ok := devMap[deviceType]; ok && client != nil {
	// 			var c = client.(*chatsocket.Client)
	// 			if !m.CheckEndPointExist(context, c) {
	// 				logUser.Debug("user already login")
	// 				//context.Writer.Write(ret.Get(m.UserLoginOtherDevice, "", m.Empty{}))
	// 				context.JSON(http.StatusOK, ret.Get(m.UserLoginOtherDevice, "", m.Empty{}))
	// 				return
	// 			}
	// 		}
	// 	} else {
	// 		for k, _ := range devMap {
	// 			if k != "Web" {
	// 				if client, ok := devMap[k]; ok && client != nil {
	// 					var c = client.(*chatsocket.Client)
	// 					if !m.CheckEndPointExist(context, c) {
	// 						logUser.Debug("user already login")
	// 						//context.Writer.Write(ret.Get(m.UserLoginOtherDevice, "", m.Empty{}))
	// 						context.JSON(http.StatusOK, ret.Get(m.UserLoginOtherDevice, "", m.Empty{}))
	// 						return
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	// //清除旁听状态
	// var key string
	// var ok bool
	// if key, ok = chatsocket.FormatListenerKey(utility.ToString(data["id"]), deviceType); ok {
	// 	for _, docking_info_inf := range utility.DockingMap {
	// 		var docking_info *chatsocket.DockingInfo
	// 		docking_info = docking_info_inf.(*chatsocket.DockingInfo)
	// 		delete(docking_info.ListenCsList, key)
	// 	}
	// }

	session := sessions.Default(context)
	session.Set("id", utility.RandomID())
	session.Set("user_id", data["id"])
	if deviceType == types.DeviceWeb {
		/*var appList = make([]string, 0)
		pemList, _ := data["permission_list"].([]*m.AppPem)
		for _, pem := range pemList {
			appList = append(appList, pem.AppId)
		}
		session.Set("app_list", appList)*/
		session.Set("ismanager", true)
	}
	// TODO: remove app_id
	session.Set("app_id", "1001")
	session.Set("token", token)
	session.Set("devtype", deviceType)
	err = session.Save()
	if err != nil {
		logUser.Error("tokenLogin save session failed", "err", err)
	}

	ret := result.ComposeHttpAck(result.CodeOK, "", data)
	context.JSON(http.StatusOK, ret)
}

// 3、 UserLogin 用户统计信息API
func UserStatistics(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var params map[string]string
	err := json.Unmarshal(body, &params)
	if err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	app_id, ok := params["app_id"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "app_id", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	logUser.Info("user statistics", "app_id", app_id)
	resp, errMsg := m.UserStatistics(app_id)
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", resp)
	context.JSON(http.StatusOK, ret)
}

/*
	获取用户信息
	入参：用户id
*/
func UserInfo(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var params map[string]string
	err := json.Unmarshal(body, &params)
	if err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	userID, ok := params["id"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "id", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	resp, errMsg := m.UserInfo(userID)
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", resp)
	context.JSON(http.StatusOK, ret)
}

// 4、 UserInfoList 用户信息列表API
func UserInfoList(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var params map[string]interface{}
	err := json.Unmarshal(body, &params)
	if err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	appID, ok := params["app_id"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "app_id", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	userType, ok := params["user_type"]
	if !ok {
		userType = 0
	}

	uid, ok := params["uid"]
	if !ok || uid == nil {
		uid = ""
	}

	account, ok := params["account"]
	if !ok || account == nil {
		account = ""
	}

	page, ok := params["page"]
	if !ok {
		page = 1
	}

	number, ok := params["number"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "number", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	resp, errMsg := m.UserInfoList(appID.(string), utility.ToInt(userType), uid.(string),
		account.(string), utility.ToInt(page)-1, utility.ToInt(number))
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", resp)
	context.JSON(http.StatusOK, ret)
}

/*
	编辑用户信息
*/
func UserEditInfo(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var params map[string]string
	err := json.Unmarshal(body, &params)
	if err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	userID, ok := params["id"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "id", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	/*session, err := utility.SessionStore.Get(context.Request, utility.SESSION_LOGIN)
	csID, ok := session.Values["user_id"]

	if !ok {
		var _ack m.BaseReturn
		ret := _ack.Get(m.PermissionDeny, "", m.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}*/
	session := sessions.Default(context)
	idRaw := session.Get("user_id")
	if idRaw == nil {
		ret := result.ComposeHttpAck(result.SessionError, "session error", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	csID := idRaw.(string)

	appID, err := dbtools.GetAppIDByUserID(utility.ToString(userID))
	if err != nil {
		ret := result.ComposeHttpAck(result.QueryDbFailed, "session error", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	ok, err = m.HasPermission(csID, appID, m.Remark)
	if err != nil || !ok {
		ret := result.ComposeHttpAck(result.PermissionDeny, "session error", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	errMsg := m.UserEditInfo(params, csID)
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	context.JSON(http.StatusOK, ret)
}

/*
	用户举报详情
*/
func UserReportInfo(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	//logUser.Info("post /user/report", "body", string(body))
	var params map[string]string
	err := json.Unmarshal(body, &params)
	if err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	id, ok := params["id"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "id", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	resp, errMsg := m.UserReportInfo(id)
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", resp)
	context.JSON(http.StatusOK, ret)
}

/*
	举报用户
	入参：被举报用户id
		举报内容
	返回：举报成功与否
*/
func UserReport(context *gin.Context) {
	var err error
	/*session, _ := utility.SessionStore.Get(context.Request, utility.SESSION_LOGIN)
	oID, ok := session.Values["user_id"]
	if !ok {
		result = ret.Get(m.PermissionDeny, "登录用户才可举报")
		context.JSON(http.StatusOK, ret)
		return
	}*/
	session := sessions.Default(context)
	idRaw := session.Get("user_id")
	if idRaw == nil {
		ret := result.ComposeHttpAck(result.SessionError, "session error", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	oID := idRaw.(string)

	body, _ := ioutil.ReadAll(context.Request.Body)
	var params map[string]interface{}
	err = json.Unmarshal(body, &params)
	if err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	id, ok := params["id"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "id", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}
	msgID, ok := params["msg_id"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "msg_id", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	errMsg := m.Report(cmn.ToString(oID), cmn.ToString(id), cmn.ToString(msgID))
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	context.JSON(http.StatusOK, ret)
}

/*
	移除用户
*/
func UserKickout(context *gin.Context) {
	var err error
	body, _ := ioutil.ReadAll(context.Request.Body)
	var params map[string]interface{}
	err = json.Unmarshal(body, &params)
	if err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	id, ok := params["id"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "id", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}
	kickOutTime, ok := params["kickout_time"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "kickout_time", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	// check type
	id, ok = id.(string)
	if !ok {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	appID, err := dbtools.GetAppIDByUserID(utility.ToString(id))
	if err != nil {
		ret := result.ComposeHttpAck(result.UserNotExists, err.Error(), result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}
	if !CheckPermission(context, appID, m.KickOut) {
		ret := result.ComposeHttpAck(result.PermissionDeny, err.Error(), result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	/*session, _ := utility.SessionStore.Get(context.Request, utility.SESSION_LOGIN)
	csID, ok := session.Values["user_id"]*/
	session := sessions.Default(context)
	csID := session.Get("user_id")
	if csID == nil {
		ret := result.ComposeHttpAck(result.SessionError, "session error", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	errMsg := m.KickOutUser(id.(string), csID.(string), appID, utility.ToInt64(kickOutTime))

	// 调用socket接口
	m.KickOutUserEvent(id.(string))

	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	context.JSON(http.StatusOK, ret)
}

/*
	客服信息列表
	根据appID、添加时间（可选）筛选客服
*/
func CSInfoList(context *gin.Context) {
	var err error
	body, _ := ioutil.ReadAll(context.Request.Body)
	//logUser.Info("post /user/customServiceInfoList", "body", string(body))

	var params map[string]interface{}
	err = json.Unmarshal(body, &params)
	if err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	_, ok := params["app_id"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "app_id", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	resp, errMsg := m.CSInfoList(params)
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", resp)
	context.JSON(http.StatusOK, ret)
}

/*
	添加客服
	客服UID为找币UID
*/
func AddCS(context *gin.Context) {
	var err error
	if !CheckAdmin(context) {
		ret := result.ComposeHttpAck(result.PermissionDeny, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	body, _ := ioutil.ReadAll(context.Request.Body)
	//logUser.Info("post /user/addCustomService", "body", string(body))

	var params map[string]interface{}
	err = json.Unmarshal(body, &params)
	if err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	appID, ok := params["app_id"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "app_id", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}
	csUID, ok := params["cs_uid"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "cs_uid", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}
	csName, ok := params["cs_name"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "cs_name", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}
	pList, ok := params["permission_list"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "permission_list", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	errMsg := m.AddCS(appID.(string), csUID.(string), csName.(string), pList.([]interface{}))
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	context.JSON(http.StatusOK, ret)
}

/*
	删除客服信息
*/
func RemoveCustomService(context *gin.Context) {
	var err error
	if !CheckAdmin(context) {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	body, _ := ioutil.ReadAll(context.Request.Body)
	var params map[string]interface{}
	err = json.Unmarshal(body, &params)
	if err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	csID, ok := params["cs_id"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "cs_id", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}
	id, ok := csID.(string)
	if !ok {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	appID, ok := params["app_id"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "app_id", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	errMsg := m.RemoveCustomService(id, utility.ToString(appID))
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	context.JSON(http.StatusOK, ret)
}

/*
	编辑客服名称
*/
func EditCSName(context *gin.Context) {
	var err error
	if !CheckAdmin(context) {
		ret := result.ComposeHttpAck(result.PermissionDeny, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	body, _ := ioutil.ReadAll(context.Request.Body)
	var params map[string]interface{}
	err = json.Unmarshal(body, &params)
	if err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	csID, ok := params["cs_id"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "cs_id", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}
	id, ok := csID.(string)
	if !ok {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	csName, ok := params["cs_name"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "cs_name", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}
	name, ok := csName.(string)
	if !ok {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	appID, ok := params["app_id"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "app_id", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	errMsg := m.EditCSName(id, utility.ToString(appID), name)
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	context.JSON(http.StatusOK, ret)
}

/*
	编辑客服权限
*/
func EditCSPermission(context *gin.Context) {
	if !CheckAdmin(context) {
		ret := result.ComposeHttpAck(result.PermissionDeny, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}
	body, _ := ioutil.ReadAll(context.Request.Body)
	logUser.Info("post /user/EditCSPermission", "body", string(body))
	var params map[string]interface{}
	err := json.Unmarshal(body, &params)
	if err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}
	appID, ok := params["app_id"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "app_id", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	csID, ok := params["cs_id"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "cs_id", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	pList, ok := params["permission_list"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "permission_list", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	errMsg := m.EditCSPermission(appID.(string), csID.(string), pList.([]interface{}))
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	context.JSON(http.StatusOK, ret)
}

/*
	客服操作记录
*/
func CSOperateLog(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	//logUser.Info("post /user/customServiceOperateLog", "body", string(body))

	var params map[string]interface{}
	err := json.Unmarshal(body, &params)
	if err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	csID, ok := params["cs_id"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "cs_id", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	appID, ok := params["app_id"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "app_id", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	page, ok := params["page"]
	if !ok {
		page = float64(1)
	}

	number, ok := params["number"]
	if !ok {
		ret := result.ComposeHttpAck(result.LackParam, "number", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}

	resp, errMsg := m.CSOperateLog(csID.(string), utility.ToString(appID), int(page.(float64))-1, utility.ToInt(number))
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		context.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", resp)
	context.JSON(http.StatusOK, ret)
}

// ==================== rewrite =================

/*
	用户自定义头像
*/
func UserEditAvatar(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		ret := result.ComposeHttpAck(result.SessionError, "lack user_id", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	type editAvatarParams struct {
		Avatar string `json:"avatar" binding:"required"`
	}
	var params editAvatarParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	err := dbtools.UpdateUserAvatar(userID.(string), params.Avatar)
	if err != nil {
		logUser.Error("UserEditAvatar", "err_msg", err)
		ret := result.ComposeHttpAck(result.WriteDbFailed, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	} else {
		ret := result.ComposeHttpAck(result.CodeOK, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
}

/*
	用户修改昵称
*/
func UserEditNickname(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		ret := result.ComposeHttpAck(result.SessionError, "lack user_id", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	type editNickNameParams struct {
		Nickname string `json:"nickname" binding:"required"`
	}
	var params editNickNameParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	err := dbtools.UpdateUserName(userID.(string), params.Nickname)

	if err != nil {
		logUser.Error("UserEditNickName", "err_msg", err)
		ret := result.ComposeHttpAck(result.WriteDbFailed, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	} else {
		ret := result.ComposeHttpAck(result.CodeOK, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
}

/*
	判断用户是否在zb注册
*/
func IsReg(c *gin.Context) {
	type isRegParams struct {
		Mobile string `json:"mobile" binding:"required"`
	}
	var params isRegParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	data, err := m.IsReg(params.Mobile)
	if err != nil {
		logUser.Error("is reg", "err", err.Message)
		ret := result.ComposeHttpAck(err.ErrorCode, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", data)
	c.JSON(http.StatusOK, ret)
}

/*
	管理端-禁言用户-全局禁言
*/
func UserMuted(c *gin.Context) {
	// todo: check permission
	// todo: websocket msg
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		ret := result.ComposeHttpAck(result.SessionError, "lack user_id", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	type muteParams struct {
		Id   string `json:"id" binding:"required"`
		Time *int64 `json:"time" binding:"exists"`
	}
	var params muteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	errMsg := m.MuteUserAll(params.Id, userID.(string), *params.Time)
	// 调用socket接口
	m.MutedUserEvent(params.Id)

	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	c.JSON(http.StatusOK, ret)
}

func SendSms(c *gin.Context) {
	type smsParams struct {
		Mobile string `json:"mobile" binding:"required"`
	}

	var params smsParams
	if err := c.ShouldBindJSON(&params); err != nil || len(params.Mobile) != 11 {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	resp, err := m.SendSms(params.Mobile)
	if err != nil {
		redPacketLog.Error("send sms", "err", err.Error())
		ret := result.ComposeHttpAck(result.GetVerifyCodeFailed, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	ret := result.ComposeHttpAck(result.CodeOK, "", resp)
	c.JSON(http.StatusOK, ret)
}
