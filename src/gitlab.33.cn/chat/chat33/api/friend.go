package api

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gitlab.33.cn/chat/chat33/model"
	"gitlab.33.cn/chat/chat33/result"
	"gitlab.33.cn/chat/chat33/types"
)

/*
	获取好友列表
*/
func FriendList(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	type friendListParams struct {
		Type int `json:"type"`
		Time int `json:"time"`
	}

	var params friendListParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	friendList, errMsg := model.FriendList(userID.(string), params.Type, params.Time)
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", friendList)
	c.JSON(http.StatusOK, ret)
}

/*
	申请添加好友
*/
func AddFriend(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	type addFriendParams struct {
		Id     string `json:"id" binding:"required"`
		Remark string `json:"remark"`
		Reason string `json:"reason"`
		RoomId string `json:"roomId"`
	}

	var params addFriendParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	errMsg := model.AddFriend(userID.(string), params.Id, params.Remark, params.Reason, params.RoomId)
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	c.JSON(http.StatusOK, ret)
}

/*
	处理好友申请
*/
func HandleFriendRequest(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	type addFriendParams struct {
		Id    string `json:"id" binding:"required"`
		Agree *int   `json:"agree" binding:"exists"`
	}

	var params addFriendParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	if (*params.Agree != types.FriendRequestAccept) && (*params.Agree != types.FriendRequestReject) {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	errMsg := model.HandleFriendRequest(userID.(string), params.Id, *params.Agree)
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	c.JSON(http.StatusOK, ret)
}

//设置备注
func FriendSetRemark(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	type setRemarkParams struct {
		Id     string  `json:"id" binding:"required"`
		Remark *string `json:"remark" binding:"exists"`
	}
	var params setRemarkParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	errMsg := model.SetFriendRemark(userID.(string), params.Id, *params.Remark)
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	c.JSON(http.StatusOK, ret)
}

/*
	设置好友免打扰
*/
func FriendSetDND(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	type setDNDParams struct {
		Id              string `json:"id" binding:"required"`
		SetNoDisturbing *int   `json:"setNoDisturbing" binding:"exists"`
	}
	var params setDNDParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	errMsg := model.SetFriendDND(userID.(string), params.Id, *params.SetNoDisturbing)
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	c.JSON(http.StatusOK, ret)
}

/*
	设置好友置顶
*/
func FriendSetTop(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	type setTopParams struct {
		Id  string `json:"id" binding:"required"`
		Top *int   `json:"stickyOnTop" binding:"exists"`
	}
	var params setTopParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	errMsg := model.SetFriendTop(userID.(string), params.Id, *params.Top)
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	c.JSON(http.StatusOK, ret)
}

/*
	删除好友
*/
func DeleteFriend(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		WriteJson(c, result.ComposeHttpAck(result.SessionError, "缺少user_id", nil))
		return
	}

	type deleteFriendParams struct {
		Id string `json:"id" binding:"required"`
	}
	var params deleteFriendParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, result.ComposeHttpAck(result.ParamsError, "", nil))
		return
	}

	err := model.DeleteFriend(userID.(string), params.Id)
	WriteJson(c, result.ComposeHttpAck(err.ErrorCode, err.Message, nil))
}

//查看好友详情
func FriendInfo(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		WriteJson(c, result.ComposeHttpAck(result.SessionError, "缺少user_id", nil))
		return
	}

	type friendInfoParams struct {
		Id string `json:"id" binding:"required"`
	}
	var params friendInfoParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, result.ComposeHttpAck(result.ParamsError, "", nil))
		return
	}

	data, errMsg := model.FriendInfo(userID.(string), params.Id)
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", data)
	c.JSON(http.StatusOK, ret)
}

//获取好友消息记录
func CatLog(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		WriteJson(c, result.ComposeHttpAck(result.SessionError, "缺少user_id", nil))
	}
	type CatLogPara struct {
		Id      string `json:"id" binding:"required"`
		StartId string `json:"startId"`
		Number  int    `json:"number" binding:"required"`
	}
	var params CatLogPara
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, result.ComposeHttpAck(result.ParamsError, "", nil))
		return
	}

	data, errMsg := model.FindCatLog(userID.(string), params.Id, params.StartId, params.Number)
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	ret := result.ComposeHttpAck(result.CodeOK, "", data)
	c.JSON(http.StatusOK, ret)
}

//删除、撤回消息记录
func DeleteMsg(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		WriteJson(c, result.ComposeHttpAck(result.SessionError, "缺少user_id", nil))
	}
	type DeleteMsgPara struct {
		LogId string `json:"logId" binding:"required"`
		Tp    int    `json:"type" binding:"required"`
	}
	var params DeleteMsgPara
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, result.ComposeHttpAck(result.ParamsError, "", nil))
		return
	}

	errMsg := model.DeleteMsg(userID.(string), params.LogId, params.Tp)
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	ret := result.ComposeHttpAck(result.CodeOK, "", nil)
	c.JSON(http.StatusOK, ret)
}

//GetAllFriendUnreadMsg获取所有好友未读消息统计
func GetAllFriendUnreadMsg(c *gin.Context) {
	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		WriteJson(c, result.ComposeHttpAck(result.SessionError, "缺少user_id", nil))
	}
	data, errMsg := model.GetAllFriendUnreadMsg1(userId.(string))
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", data)
	c.JSON(http.StatusOK, ret)

}
