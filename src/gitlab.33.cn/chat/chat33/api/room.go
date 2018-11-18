package api

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gitlab.33.cn/chat/chat33/model"
	"gitlab.33.cn/chat/chat33/result"
	"gitlab.33.cn/chat/chat33/types"
	"gitlab.33.cn/chat/chat33/utility"
)

// 创建群
func CreateRoom(c *gin.Context) {
	type requestParams struct {
		RoomName   string   `json:"roomName"`
		RoomAvatar string   `json:"roomAvatar"`
		Users      []string `json:"users" binding:"required"`
	}

	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	errMsg := model.CreateRoom(userId.(string), params.RoomName, params.RoomAvatar, types.CanAddFriend, types.ShouldNotApproval, types.AdminNotMuted, types.MasterNotMuted, params.Users)
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	c.JSON(http.StatusOK, ret)
}

// 删除群
func RemoveRoom(c *gin.Context) {
	type requestParams struct {
		RoomId string `json:"roomid" binding:"required"`
	}

	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	errMsg := model.RemoveRoom(userId.(string), params.RoomId)
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	c.JSON(http.StatusOK, ret)
}

// 退出群
func LoginOutRoom(c *gin.Context) {
	type requestParams struct {
		RoomId string `json:"roomid" binding:"required"`
	}

	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	//check user is not master of this room
	if model.CheckUserIsMaster(params.RoomId, userId.(string)) {
		//master cannot logout
		ret := result.ComposeHttpAck(result.CanNotLoginOut, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	errMsg := model.LoginOutRoom(userId.(string), params.RoomId)
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	c.JSON(http.StatusOK, ret)
}

// 踢出群
func KickOutRoom(c *gin.Context) {
	type requestParams struct {
		RoomId string   `json:"roomId" binding:"required"`
		Users  []string `json:"users" binding:"required"`
	}

	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	if len(params.Users) < 1 {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	//check user is not master of this room
	if !model.CheckUserIsMaster(params.RoomId, userId.(string)) {
		//master cannot logout
		ret := result.ComposeHttpAck(result.PermissionDeny, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	var ret interface{}
	for _, v := range params.Users {
		errMsg := model.KickOutRoom(userId.(string), params.RoomId, v)
		if ret == nil || errMsg.ErrorCode == result.CodeOK {
			ret = result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		}
	}
	c.JSON(http.StatusOK, ret)
}

// 获取群列表
func GetRoomList(c *gin.Context) {
	type requestParams struct {
		Type int `json:"type" binding:"required"`
	}

	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	roomList, errMsg := model.GetRoomList(userId.(string), params.Type)
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", roomList)
	c.JSON(http.StatusOK, ret)
}

// 获取群信息
func GetRoomInfo(c *gin.Context) {
	type requestParams struct {
		RoomId string `json:"roomId" binding:"required"`
	}

	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	roomInfo, errMsg := model.GetRoomInfo(userId.(string), params.RoomId)
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", roomInfo)
	c.JSON(http.StatusOK, ret)
}

// 获取群成员列表
func GetRoomUserList(c *gin.Context) {
	type requestParams struct {
		RoomId string `json:"roomId" binding:"required"`
	}

	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	roomUserList, errMsg := model.GetRoomUserList(params.RoomId)
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", roomUserList)
	c.JSON(http.StatusOK, ret)
}

// 获取群成员信息
func GetRoomUserInfo(c *gin.Context) {
	type requestParams struct {
		RoomId string `json:"roomId" binding:"required"`
		UserId string `json:"userId" binding:"required"`
	}

	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	roomUserInfo, errMsg := model.GetRoomUserInfo(params.RoomId, params.UserId)
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", roomUserInfo)
	c.JSON(http.StatusOK, ret)
}

// 搜索群成员信息
func GetRoomSearchMember(c *gin.Context) {
	type requestParams struct {
		RoomId string `json:"roomId" binding:"required"`
		Query  string `json:"query" binding:"required"`
	}

	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	roomUserInfo, errMsg := model.GetRoomSearchMember(params.RoomId, params.Query)
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", roomUserInfo)
	c.JSON(http.StatusOK, ret)
}

// 管理员设置群
func AdminSetPermission(c *gin.Context) {
	type requestParams struct {
		RoomId         string `json:"roomId" binding:"required"`
		CanAddFriend   int    `json:"canAddFriend"`
		JoinPermission int    `json:"joinPermission"`
	}

	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	//check caller is manager
	if !model.CheckUserIsMamnagerOrMaster(params.RoomId, utility.ToString(userId)) {
		ret := result.ComposeHttpAck(result.PermissionDeny, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	errMsg := model.AdminSetPermission(params.RoomId, params.CanAddFriend, params.JoinPermission)
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	c.JSON(http.StatusOK, ret)
}

// 群内用户身份设置
func SetLevel(c *gin.Context) {
	type requestParams struct {
		RoomId string `json:"roomId" binding:"required"`
		UserId string `json:"userId" binding:"required"`
		Level  int    `json:"level" binding:"required"`
	}

	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	// check userId is admin
	if !model.CheckUserIsMaster(params.RoomId, utility.ToString(userId)) {
		ret := result.ComposeHttpAck(result.PermissionDeny, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	errMsg := model.SetLevel(utility.ToString(userId), params.UserId, params.RoomId, params.Level)
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	c.JSON(http.StatusOK, ret)
}

// 群成员设置免打扰
func SetNoDisturbing(c *gin.Context) {
	type requestParams struct {
		RoomId          string `json:"roomId" binding:"required"`
		SetNoDisturbing int    `json:"setNoDisturbing" binding:"required"`
	}

	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	switch params.SetNoDisturbing {
	case 1:
	case 2:
	default:
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	errMsg := model.SetNoDisturbing(userId.(string), params.RoomId, params.SetNoDisturbing)
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	c.JSON(http.StatusOK, ret)
}

// 群成员设置消息置顶
func SetStickyOnTop(c *gin.Context) {
	type requestParams struct {
		RoomId      string `json:"roomId" binding:"required"`
		StickyOnTop int    `json:"stickyOnTop" binding:"required"`
	}

	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	switch params.StickyOnTop {
	case 1:
	case 2:
	default:
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	errMsg := model.SetStickyOnTop(userId.(string), params.RoomId, params.StickyOnTop)
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	c.JSON(http.StatusOK, ret)
}

// 群成员设置群内昵称
func SetMemberNickname(c *gin.Context) {
	type requestParams struct {
		RoomId   string `json:"roomId" binding:"required"`
		Nickname string `json:"nickname" binding:"required"`
	}

	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	errMsg := model.SetMemberNickname(userId.(string), params.RoomId, params.Nickname)
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	c.JSON(http.StatusOK, ret)
}

// 邀请入群
func JoinRoomInvite(c *gin.Context) {
	type requestParams struct {
		RoomId string   `json:"roomId" binding:"required"`
		Users  []string `json:"users" binding:"required"`
	}

	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	errMsg := model.JoinRoomInvite(userId.(string), params.RoomId, params.Users)
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	c.JSON(http.StatusOK, ret)
}

// 入群申请
func JoinRoomApply(c *gin.Context) {
	type requestParams struct {
		RoomId      string `json:"roomId" binding:"required"`
		ApplyReason string `json:"applyReason"`
	}

	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	errMsg := model.JoinRoomApply(userId.(string), params.RoomId, params.ApplyReason)
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	c.JSON(http.StatusOK, ret)
}

// 入群申请处理
func JoinRoomApprove(c *gin.Context) {
	type requestParams struct {
		RoomId string `json:"roomId" binding:"required"`
		UserId string `json:"userId" binding:"required"`
		Agree  int    `json:"agree" binding:"required"`
	}

	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	errMsg := model.JoinRoomApprove(utility.ToString(userId), params.RoomId, params.UserId, params.Agree)
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	c.JSON(http.StatusOK, ret)
}

// 获取消息记录
func GetRoomChatLog(c *gin.Context) {
	type requestParams struct {
		RoomId  string `json:"id" binding:"required"`
		StartId string `json:"startId"`
		Number  int    `json:"number"`
	}

	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	roomUserList, errMsg := model.GetRoomChatLog(userId.(string), params.RoomId, params.StartId, params.Number)
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", roomUserList)
	c.JSON(http.StatusOK, ret)
}

// 获取群未读消息统计
func GetRoomUnreadStatistics(c *gin.Context) {
	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	list, errMsg := model.GetRoomUnReadStatistics(userId.(string))
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", list)
	c.JSON(http.StatusOK, ret)
}

// 获取群在线人数
func GetOnlineNumber(c *gin.Context) {
	type requestParams struct {
		RoomId string `json:"roomId" binding:"required"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	info, errMsg := model.GetRoomOnlineNumber(params.RoomId)
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", info)
	c.JSON(http.StatusOK, ret)
}
