package main

import (
	"os"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	l "github.com/inconshreveable/log15"

	"gitlab.33.cn/chat/chat33/db"
	"gitlab.33.cn/chat/chat33/types"

	"github.com/BurntSushi/toml"
	"gitlab.33.cn/chat/chat33/api"
	"gitlab.33.cn/chat/chat33/model"

	"gitlab.33.cn/chat/chat33/comet"
)

func initLogLevel(cfg *types.Config) {
	var level l.Lvl
	switch cfg.Log.Level {
	case "debug":
		level = l.LvlDebug
	case "info":
		level = l.LvlInfo
	case "warn":
		level = l.LvlWarn
	case "error":
		level = l.LvlError
	case "crit":
		level = l.LvlCrit
	default:
		level = l.LvlWarn
	}
	l.Root().SetHandler(l.LvlFilterHandler(level, l.StreamHandler(os.Stdout, l.TerminalFormat())))
}

func main() {
	os.Chdir(pwd())
	d, _ := os.Getwd()
	l.Info("project info:", "dir", d)
	var cfg types.Config
	if _, err := toml.DecodeFile(d+"/etc/config.toml", &cfg); err != nil {
		panic(err)
	}
	l.Info("config info", "cfg", cfg)
	initLogLevel(&cfg)

	db.InitDB(&cfg)
	model.Init(&cfg)

	r := gin.Default()
	store := cookie.NewStore([]byte("veryHardToGuess"))
	r.Use(sessions.Sessions("session-login", store))
	// websocket
	r.GET("/ws", func(context *gin.Context) {
		comet.ServeWs(context)
	})

	//精确搜索用户或群
	r.POST("/chat33/search", api.ClearlySearch)
	//
	r.POST("/chat33/applyList", api.GetApplyList)
	//创建群
	r.POST("/room/create", api.CreateRoom)
	//删除群
	r.POST("/room/delete", api.RemoveRoom)
	//退出群
	r.POST("/room/loginOut", api.LoginOutRoom)
	//踢出群
	r.POST("/room/kickOut", api.KickOutRoom)
	//获取群列表
	r.POST("/room/list", api.GetRoomList)
	//获取群信息
	r.POST("/room/info", api.GetRoomInfo)
	//获取群成员列表
	r.POST("/room/userList", api.GetRoomUserList)
	//获取群成员信息
	r.POST("/room/userInfo", api.GetRoomUserInfo)
	//管理员设置群
	r.POST("/room/setPermission", api.AdminSetPermission)
	//群内用户身份设置
	r.POST("/room/setLevel", api.SetLevel)
	//群成员设置免打扰
	r.POST("/room/setNoDisturbing", api.SetNoDisturbing)
	//群成员设置群内昵称
	r.POST("/room/setMemberNickname", api.SetMemberNickname)
	//邀请入群
	r.POST("/room/joinRoomInvite", api.JoinRoomInvite)
	//申请入群
	r.POST("/room/joinRoomApply", api.JoinRoomApply)
	//入群申请处理
	r.POST("/room/joinRoomApprove", api.JoinRoomApprove)
	//群消息记录
	r.POST("/room/chatLog", api.GetRoomChatLog)
	//群成员设置群置顶
	r.POST("/room/stickyOnTop", api.SetStickyOnTop)
	//获取所有群未读消息统计
	r.POST("/room/unread", api.GetRoomUnreadStatistics)
	//获取群在线人数
	r.POST("/room/getOnlineNumber", api.GetOnlineNumber)
	//搜索群成员信息
	r.POST("/room/searchMember", api.GetRoomSearchMember)

	// 获取好友列表
	r.POST("/friend/list", api.FriendList)
	// 添加好友
	r.POST("/friend/add", api.AddFriend)
	// 删除好友
	r.POST("friend/delete", api.DeleteFriend)
	// 处理好友请求
	r.POST("/friend/response", api.HandleFriendRequest)
	// 修改好友备注
	r.POST("/friend/setRemark", api.FriendSetRemark)
	// 设置好友免打扰
	r.POST("/friend/setNoDisturbing", api.FriendSetDND)
	// 设置好友消息置顶
	r.POST("/friend/stickyOnTop", api.FriendSetTop)
	//获取好友消息记录
	r.POST("/friend/chatLog", api.CatLog)
	//获取所有好友未读消息统计
	r.POST("/friend/unread", api.GetAllFriendUnreadMsg)

	//删除指定的一条消息
	r.POST("/deleteMsg", api.DeleteMsg)
	//撤回指定的一条消息
	r.POST("/withdrawMsg", api.DeleteMsg)

	// 首页统计信息
	r.POST("/index/statistics", api.IndexStatistics)
	// 获取币种信息
	r.POST("/coin/coinInfo", api.CoinList)
	// 获取权限列表
	r.POST("/permission/permissionInfo", api.PermissionInfo)
	// 获取手机验证码
	r.POST("/send/sms", api.SendSms)
	// 获取app信息
	r.POST("/app/appInfo", api.AppInfo)

	// 用户密码登录
	r.POST("/user/pwdLogin", api.UserPwdLogin)
	// 用户后台token登录
	r.POST("/user/tokenLogin", api.UserTokenLogin)
	// 用户是否已注册找币账户
	r.POST("/user/isreg", api.IsReg)
	// 用户统计信息
	r.POST("/user/userStatistics", api.UserStatistics)
	// 用户修改昵称
	r.POST("/user/editNickname", api.UserEditNickname)
	// 用户自定义头像
	r.POST("/user/editAvatar", api.UserEditAvatar)
	// 用户信息编辑
	r.POST("/user/editInfo", api.UserEditInfo)
	// 用户举报详情
	r.POST("/user/reportInfo", api.UserReportInfo)
	// 举报用户
	r.POST("/user/report", api.UserReport)
	// 禁言用户
	r.POST("/user/muted", api.UserMuted)
	// 移除用户
	r.POST("/user/kickOut", api.UserKickout)
	// 用户信息列表
	r.POST("/user/infoList", api.UserInfoList)
	// 用户信息
	r.POST("/user/info", api.UserInfo)
	//查看用户详情
	r.POST("/user/userInfo", api.FriendInfo)

	// 客服数量统计
	//r.POST("/user/customServiceStatistics", api.CustomServiceStatistics)
	// 删除客服信息
	r.POST("/user/removeCustomService", api.RemoveCustomService)
	// 编辑客服名称
	r.POST("/user/editCSName", api.EditCSName)
	// 编辑客服权限
	r.POST("/user/editCSPermission", api.EditCSPermission)
	// 获取客服信息列表
	r.POST("/user/customServiceInfoList", api.CSInfoList)
	// 添加客服
	r.POST("/user/addCustomService", api.AddCS)
	// 获取客服操作记录
	r.POST("/user/customServiceOperateLog", api.CSOperateLog)

	// 添加聊天室
	r.POST("/group/addGroup", api.AddGroup)
	// 聊天室列表
	r.POST("/group/list", api.GroupInfoList)
	// 聊天室详情
	r.POST("/group/info", api.GetGroupInfo)
	// 聊天室用户列表
	r.POST("/group/userList", api.GroupUserList)
	// 聊天室用户信息
	r.POST("/group/userInfo", api.GroupUserInfo)
	// 群聊天记录
	r.POST("/group/getGroupChatHistory", api.GroupChatHistory)
	// 聊天室头像编辑
	r.POST("/group/editAvatar", api.GroupEditAvatar)
	// 聊天室状态编辑
	r.POST("/group/setStatus", api.SetGroupStatus)
	// 编辑聊天室名称
	r.POST("/group/editGroupName", api.SetGroupName)
	// 获取聊天室总人数
	r.POST("/group/getOnlineNumber", api.GetGroupOnlineNumber)

	// 查询余额
	r.POST("/red-packet/balance", api.Balance)
	// 发红包
	r.POST("/red-packet/send", api.Send)
	// 收红包
	r.POST("/red-packet/receive-entry", api.ReceiveEntry)
	// 收红包
	r.POST("/red-packet/receive", api.Receive)
	// 红包入账
	r.POST("/red-packet/entry", api.Entry)
	// 注册领取红包
	r.POST("/red-packet/register-entry", api.RegisterEntry)
	// 红包统计信息
	r.POST("/red-packet/statistics", api.RedEnvelopeStatistics)
	// 红包信息列表
	r.POST("/red-packet/infoList", api.RedEnvelopeInfoList)
	// 红包详情
	r.POST("/red-packet/detail", api.RedEnvelopeDetail)
	// 获取收红包历史
	r.POST("/red-packet/receive-record", api.RPReceiveHistory)
	// 获取发红包历史
	r.POST("/red-packet/send-record", api.RPSendHistory)
	// 获取收发红包历史
	r.POST("/red-packet/record", api.RPHistory)

	r.Run(cfg.Server.Addr)
}

/*
	---workdir/
		| -- bin/
		|     |-- chat(I am here)
		|
		| -- etc/
			  |-- config.toml
			  |-- config.json
*/
func pwd() string {
	dir, err := filepath.Abs(filepath.Dir(filepath.Dir(os.Args[0])))
	if err != nil {
		panic(err)
	}
	return dir
}
