package api

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"gitlab.33.cn/chat/chat33/db"
	"gitlab.33.cn/chat/chat33/model"
	"gitlab.33.cn/chat/chat33/result"
	"gitlab.33.cn/chat/chat33/utility"
)

// var slog = log15.New("module", "chat/api/statistics")

// 获取app列表api
func AppInfo(context *gin.Context) {
	ret, _ := model.GetAppInfoList()
	context.JSON(http.StatusOK, ret)
}

// 获取首页统计信息api
func IndexStatistics(c *gin.Context) {
	roomUserList, errMsg := model.GetIndexStatistics()
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", roomUserList)
	c.JSON(http.StatusOK, ret)
}

// 获取币种信息api
func CoinList(context *gin.Context) {
	ret, _ := model.GetCoinList()
	context.JSON(http.StatusOK, ret)
}

func PermissionInfo(context *gin.Context) {
	var data []model.Permission

	perms := db.GetPermissionInfo()
	fmt.Println(perms)

	for _, perm := range perms {
		data = append(data, model.Permission{
			PermissionId:   utility.ToInt(perm["permission_id"]),
			PermissionName: utility.ToString(perm["permission_name"]),
		})
	}

	ret := result.ComposeHttpAck(result.CodeOK, "", model.PermissionList{PermissionList: data})
	context.JSON(http.StatusOK, ret)
}

// 精确搜索用户或群
func ClearlySearch(c *gin.Context) {
	type requestParams struct {
		Id string `json:"markId" binding:"required"`
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

	info, errMsg := model.ClearlySearch(userId.(string), params.Id)
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", info)
	c.JSON(http.StatusOK, ret)
}

// 获取请求列表
func GetApplyList(c *gin.Context) {
	type requestParams struct {
		Id     string `json:"id"`
		Number int    `json:"number" binding:"required"`
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

	info, errMsg := model.GetApplyList(userId.(string), params.Id, params.Number)
	if errMsg.ErrorCode != result.CodeOK {
		ret := result.ComposeHttpAck(errMsg.ErrorCode, "", result.Empty{})
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := result.ComposeHttpAck(result.CodeOK, "", info)
	c.JSON(http.StatusOK, ret)
}
