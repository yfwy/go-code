package api

import (
	"net/http"

	"gitlab.33.cn/chat/chat33/types"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gitlab.33.cn/chat/chat33/model"
	"gitlab.33.cn/chat/chat33/result"
	"gitlab.33.cn/chat/chat33/utility"
)

// var group_api_log = l.New("module", "chat/api/group")

func AddGroup(c *gin.Context) {
	type requestParams struct {
		GroupName string `json:"groupName" binding:"required"`
		Avatar    string `json:"avatar"`
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

	// TODO //查询是否有添加权限
	// rlt, err := model.UserCanManageGroup(val, app_id)
	// if err != nil {
	// 	var _ack model.BaseReturn
	// 	ret := _ack.Get(model.DbConnectFail, "", model.Empty{})
	// 	context.JSON(http.StatusOK, ret)
	// 	return
	// }
	// if !rlt {
	// 	var _ack model.BaseReturn
	// 	ret := _ack.Get(model.NoCSPermission, "", model.Empty{})
	// 	context.JSON(http.StatusOK, ret)
	// 	return
	// }

	errMsg := model.AddNewGroup(params.GroupName, params.Avatar, utility.ToString(userId))
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	c.JSON(http.StatusOK, ret)
}

func GroupInfoList(c *gin.Context) {
	type requestParams struct {
		GroupStatus int    `json:"groupStatus" binding:"required"`
		GroupName   string `json:"groupName"`
		StartTime   int    `json:"startTime"`
		EndTime     int    `json:"endTime"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	switch params.GroupStatus {
	case 1:
	case 2:
	case 3:
	default:
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	groupInfoList, errMsg := model.GetGroupInfoList(params.GroupStatus, params.StartTime, params.EndTime, params.GroupName)
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", groupInfoList)
	c.JSON(http.StatusOK, ret)
}

func GetGroupInfo(c *gin.Context) {
	type requestParams struct {
		GroupId string `json:"groupId" binding:"required"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	groupInfo, errMsg := model.GetGroupInfo(params.GroupId)
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", groupInfo)
	c.JSON(http.StatusOK, ret)
}

func GroupUserList(c *gin.Context) {
	type requestParams struct {
		GroupId       string `json:"groupId" binding:"required"`
		QueryUserName string `json:"queryUserName"`
		Page          string `json:"page"`
		Number        string `json:"number"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	var queryParams = make(map[string]interface{})
	queryParams["queryUserName"] = params.QueryUserName
	queryParams["page"] = params.Page
	queryParams["number"] = params.Number

	groupInfoList, errMsg := model.GetGroupUserList(params.GroupId, queryParams)
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", groupInfoList)
	c.JSON(http.StatusOK, ret)
}

func GroupUserInfo(c *gin.Context) {
	type requestParams struct {
		Id string `json:"id" binding:"required"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	groupUserInfo, errMsg := model.GetGroupUserInfo(params.Id)
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", groupUserInfo)
	c.JSON(http.StatusOK, ret)
}

// 编辑聊天室状态api
func SetGroupStatus(c *gin.Context) {
	type requestParams struct {
		GroupId string `json:"groupId" binding:"required"`
		Type    int    `json:"type" binding:"required"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	//获取调用者id
	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	// TODO //查询是否有修改权限
	// rlt, err := model.UserCanManageGroup(val, app_id)
	// if err != nil {
	// 	var _ack model.BaseReturn
	// 	ret := _ack.Get(model.DbConnectFail, "", model.Empty{})
	// 	context.JSON(http.StatusOK, ret)
	// 	return
	// }
	// if !rlt {
	// 	var _ack model.BaseReturn
	// 	ret := _ack.Get(model.NoCSPermission, "", model.Empty{})
	// 	context.JSON(http.StatusOK, ret)
	// 	return
	// }

	switch params.Type {
	case types.GroupOpen:
	case types.GroupClose:
	case types.GroupDelete:
	default:
		ret := result.ComposeHttpAck(result.ParamsError, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	errMsg := model.SetGroupStatus(params.GroupId, params.Type)
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	c.JSON(http.StatusOK, ret)
}

func SetGroupName(c *gin.Context) {
	type requestParams struct {
		GroupId   string `json:"groupId" binding:"required"`
		GroupName string `json:"groupName" binding:"required"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	//获取调用者id
	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	// TODO //查询是否有修改权限 1、是客服 2、是对应的APP
	// app_id := model.GetAppidByGroupid(utility.ToString(group_id))
	// rlt, err := model.UserCanManageGroup(val, app_id)
	// if err != nil {
	// 	var _ack model.BaseReturn
	// 	ret := _ack.Get(model.DbConnectFail, "", model.Empty{})
	// 	context.JSON(http.StatusOK, ret)
	// 	return
	// }
	// if !rlt {
	// 	var _ack model.BaseReturn
	// 	ret := _ack.Get(model.NoCSPermission, "", model.Empty{})
	// 	context.JSON(http.StatusOK, ret)
	// 	return
	// }

	errMsg := model.SetGroupName(params.GroupId, params.GroupName)
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	c.JSON(http.StatusOK, ret)
}

func GroupEditAvatar(c *gin.Context) {
	type requestParams struct {
		GroupId string `json:"groupId" binding:"required"`
		Avatar  string `json:"avatar" binding:"required"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	//获取调用者id
	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	// //TODO 查询是否有修改权限 1、是客服 2、是对应的APP
	// app_id := model.GetAppidByGroupid(utility.ToString(group_id))
	// rlt, err := model.UserCanManageGroup(val, app_id)
	// if err != nil {
	// 	var _ack model.BaseReturn
	// 	ret := _ack.Get(model.DbConnectFail, "", model.Empty{})
	// 	context.JSON(http.StatusOK, ret)
	// 	return
	// }
	// if !rlt {
	// 	var _ack model.BaseReturn
	// 	ret := _ack.Get(model.NoCSPermission, "", model.Empty{})
	// 	context.JSON(http.StatusOK, ret)
	// 	return
	// }

	errMsg := model.EditGroupAvatar(params.GroupId, params.Avatar)
	ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
	c.JSON(http.StatusOK, ret)
}

// 获取群聊记录api
func GroupChatHistory(c *gin.Context) {

	type requestParams struct {
		GroupId string `json:"id" binding:"required"`
		StartId string `json:"startId"`
		Number  int    `json:"number"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	//获取调用者id
	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		ret := result.ComposeHttpAck(result.LoginExpired, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	groupChatLog, errMsg := model.GetGroupChatHistory(utility.ToString(userId), params.GroupId, params.StartId, params.Number)
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", groupChatLog)
	c.JSON(http.StatusOK, ret)
}

func GetGroupOnlineNumber(c *gin.Context) {
	type requestParams struct {
		GroupId string `json:"groupId" binding:"required"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ret := result.ComposeHttpAck(result.LackParam, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}

	totalNumber, _, _ := model.GetGroupOnlineNumber(params.GroupId)
	var info = make(map[string]interface{})
	info["totalNumber"] = totalNumber
	ret := result.ComposeHttpAck(result.CodeOK, "", info)
	c.JSON(http.StatusOK, ret)
}
