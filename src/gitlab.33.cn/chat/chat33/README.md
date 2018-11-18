# 		 新版聊天平台
## 测试机ip地址: 172.16.103.31:8090

## 目录
<!-- TOC -->

- [新版聊天平台](#新版聊天平台)
    - [测试机ip地址: 172.16.103.31:8090](#测试机ip地址-17216103318090)
    - [目录](#目录)
    - [统一错误代码](#统一错误代码)
    - [约定](#约定)
    - [HTTP 接口](#http-接口)
        - [获取手机验证码  /send/sms](#获取手机验证码--sendsms)
    - [用户相关接口](#用户相关接口)
        - [用户密码登录 /user/pwdLogin](#用户密码登录-userpwdlogin)
        - [用户登录 /user/tokenLogin](#用户登录-usertokenlogin)
        - [查看用户详情 /user/userInfo](#查看用户详情-useruserInfo)
        - [用户编辑头像 /user/editAvatar](#用户编辑头像-usereditavatar)
        - [用户修改昵称 /user/editNickname](#用户修改昵称-usereditnickname)
        - [用户是否已注册找币账户 /user/isreg](#用户是否已注册找币账户-userisreg)
        - [管理端禁言用户（全局禁言） /user/muted](#管理端禁言用户全局禁言-usermuted)
        - [获取权限信息/permission/permissionInfo](#获取权限信息permissionpermissioninfo)
    - [聊天室有关接口](#聊天室有关接口)
        - [聊天室列表 /group/list](#聊天室列表-grouplist)
        - [聊天室详情 /group/info](#聊天室详情-groupinfo)
        - [聊天室用户列表 /group/userList](#聊天室用户列表-groupuserlist)
        - [聊天室用户信息 /group/userInfo](#聊天室用户信息-groupuserinfo)
        - [获取聊天室聊天记录 /group/getGroupChatHistory](#获取聊天室聊天记录-groupgetgroupchathistory)
        - [添加聊天室 /group/addGroup](#添加聊天室-groupaddgroup)
        - [聊天室头像编辑 /group/editAvatar](#聊天室头像编辑-groupeditavatar)
        - [聊天室状态编辑 /group/setStatus](#聊天室状态编辑-groupsetstatus)
        - [编辑聊天室名称 /group/editGroupName](#编辑聊天室名称-groupeditgroupname)
        - [禁言聊天室用户 /group/muted](#禁言聊天室用户-groupmuted)
        - [举报聊天室用户 /group/report](#举报聊天室用户-groupreport)
    - [群有关接口](#群有关接口)
        - [获取群列表 /room/list](#获取群列表-roomlist)
        - [获取群信息 /room/info](#获取群信息-roominfo)
        - [获取群成员列表 /room/userList](#获取群成员列表-roomuserlist)
        - [创建群 /room/create](#创建群-roomcreate)
        - [删除群 /room/delete](#删除群-roomdelete)
        - [管理员设置群 /room/setPermission](#管理员设置群-roomsetpermission)
        - [群内用户身份设置 /room/setLevel](#群内用户身份设置-roomsetlevel)
        - [群成员设置免打扰 /room/setNoDisturbing](#群成员设置免打扰-roomsetnodisturbing)
        - [群成员设置昵称 /room/setMemberNickname](#群成员设置昵称-roomsetmembernickname)
        - [直接入群 /room/joinInRoom](#直接入群-roomjoininroom)
        - [入群申请 /room/joinInRoomApply](#入群申请-roomjoininroomapply)
        - [邀请入群申请 /room/inviteJoinInRoomApply](#邀请入群申请-roominvitejoininroomapply)
        - [入群申请回复 /room/joinInRoomApprove](#入群申请回复-roomjoininroomapprove)
        - [获取入群申请列表 /room/joinInRoomApplyList](#获取入群申请列表-roomjoininroomapplylist)
        - [获取群消息记录 /room/chatLog](#获取群消息记录-roomchatlog)
    - [好友有关接口](#好友有关接口)
        - [获取好友列表 /friend/list](#获取好友列表-friendlist)
        - [获取所有好友未读消息统计 /friend/unread](#获取好友列表-friendunread)
        - [添加好友申请/friend/add](#添加好友申请friendadd)
        - [删除好友/friend/delete](#删除好友frienddelete)
        - [好友申请处理 /friend/response](#好友申请处理-friendresponse)
        - [获取好友请求列表 /friend/requestList](#获取好友请求列表-friendrerequestList)
        - [获取好友申请列表 /friend/requestList](#获取好友申请列表-friendrequestlist)
        - [获取好友消息记录 /friend/chatLog](#获取好友消息记录-friendchatlog)
        - [修改好友备注 /friend/setRemark](#修改好友备注-friendsetremark)
        - [设置免打扰 /friend/setDND](#设置免打扰-friendsetdnd)
        - [设置消息置顶 /friend/stickyOnTop](#设置消息置顶-friendstickyontop)
        - [删除指定的一条消息(websocket 通知) /](#删除指定的一条消息websocket-通知-)
        - [撤回指定的一条消息 (websocket 通知)/](#撤回指定的一条消息-websocket-通知)
    - [红包接口类](#红包接口类)
        - [查询账号余额 /red-packet/balance](#查询账号余额-red-packetbalance)
        - [(已登录用户)收红包 /red-packet/receive-entry](#已登录用户收红包-red-packetreceive-entry)
        - [(未登录用户)收红包 /red-packet/receive](#未登录用户收红包-red-packetreceive)
        - [(未登录用户)收红包入账 /red-packet/entry](#未登录用户收红包入账-red-packetentry)
        - [(未登录用户)注册领取红包 /red-packet/register-entry](#未登录用户注册领取红包-red-packetregister-entry)
        - [查询用户红包收发记录 /red-packet/record](#查询用户红包收发记录-red-packetrecord)
        - [查询用户红包领取记录 /red-packet/receive-record](#查询用户红包领取记录-red-packetreceive-record)
        - [查询用户发放红包记录 /red-packet/send-record](#查询用户发放红包记录-red-packetsend-record)
    - [webSocket接口：](#websocket接口)
        - [普通消息](#普通消息)
                - [Msg 格式说明：](#msg-格式说明)
        - [登录聊天室 C->S](#登录聊天室-c-s)
        - [用户禁言通知 S->C](#用户禁言通知-s-c)
        - [关闭聊天室通知 S->C](#关闭聊天室通知-s-c)
        - [删除聊天室通知 S->C](#删除聊天室通知-s-c)
        - [开启聊天室通知 S->C](#开启聊天室通知-s-c)
        - [消息有关通知 S->C](#消息有关通知-s-c)
        - [群中添加好友通知 S->C](#群中添加好友通知-s-c)
        - [入群通知 S->C](#入群通知-s-c)
        - [解散群通知 S->C](#解散群通知-s-c)
        - [被拉入群通知 S->C](#被拉入群通知-s-c)
        - [被踢出群通知 S->C](#被踢出群通知-s-c)

<!-- /TOC -->

## 统一错误代码

```
-1000 数据库连接失败
-1001 参数错误
-1002 缺少参数
-1003 Session错误
-1004 登录过期
-1005 消息格式错误
-1006 未知的设备类型
-1007 账号已经在其他终端登录
-1010 查询数据库失败
-1011 写入数据库失败
-1012 发送频率过快，请稍后再试
-2001 用户不存在
-2002 用户已存在
-2003 要添加的客服已存在
-2004 找币token登录失败
-2006 用户没有发系统消息权限
-2007 用户被禁言
-2008 游客没有发消息权限
-2009 没有给普通用户发私聊权限
-2010 当前用户已经被其他客服接待
-2011 加入旁听失败
-2013 找币交互失败
-2014 客服间不能发送私聊消息
-2015 游客已离线
-2016 暂无客服在线，请稍后再试，或可登录账号给客服留言!
-2017 已经是好友关系
-2018 对方不是您的好友
-2019 好友请求已经被处理
-2020 好友请求不存在
-2021 不能对自己进行操作
-2022 数据转换异常
-2023 删除消息失败
-3000 权限不足
-3001 用户没有加入聊天群权限
-3002 没有客服权限
-4000 红包内部错误
-4001 红包已被领完
-4002 用户未注册
-4003 用户已注册
-4004 仅限新人领取
-4005 红包标识不匹配
-4006 非法的红包ID
-4007 验证码不正确
-4008 验证码已经过期或者已使用
-4009 红包已领取
-4010 用户无发红包权限
-5000 聊天室不存在
-5001 用户未进入此聊天室
-6000 查询聊天记录失败
-7000 消息格式错误
-8000 获取手机验证码失败
-9000 服务端内部错误
```

## 约定

- 所有 `http` 接口请求都是 `POST` 方法， 数据请求和返回结果 都是 `JSON` 格式
- 接口里需要的 `FZM-AUTH-TOKEN`（找币）， `FZM-DEVICE`（登录设备类型）从 `Header` 里面读取
- 返回结果的 `result` 为 `0` 代表成功，`data` 里面是具体数据，`result` 不为 `0` 代表失败，`message` 为失败信息
- 游客 `Header` 中带上`FZM-UUID`（设备mac）
- 传递的时间相关的参数均以毫秒为单位

## HTTP 接口

### 获取手机验证码  /send/sms

请求参数：

| **参数** | **名字** | **类型** | **说明**   |
| -------- | -------- | -------- | ---------- |
| mobile   | 手机号   | string   | 11位手机号 |

返回参数：

```json
{
    "result": 0,
    "message": "操作成功",
    "data": {}
}
```

## 接口

### 获取入群/好友申请列表/chat33/applyList

`post`请求参数：

| **参数** | **名字**       | **类型** | **约束** | **说明** |
| -------- | -------------- | -------- | -------- | -------- |
| id       | 最新一条记录id | string   | 必填     |          |
| number   | 数量           | int      | 必填     |          |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
        "applyList":[
            {
                senderInfo:{
                    "id":"1123",
                    "name":"用户1",
                    "avatar":"http://...../***.jpg",
                    "position":"产品"
                },
                receiveInfo:{
                    "id":"1123",
                    "name":"用户1",
                    "avatar":"http://...../***.jpg",
                    "position":"产品"
                },
                "id":123,
                "type":1, //1 群 2 好友
                "applyReason":"申请理由",
                "status":1, //1:等待验证 2:已拒绝 3:已同意 
                "datetime":1676764266167
            }
        ],
        "totalNumber": 1231, // 总数量
        "nextId":"123"
	}
}
```

### 精确搜索用户/群 /chat33/search

`post`请求参数：

| **参数** | **名字**             | **类型** | **约束** | **说明**           |
| -------- | -------------------- | -------- | -------- | ------------------ |
| markId   | 用户markId或群markId | string   | 必填     | 群和用户外部用的id |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
        "type":1, //1：群 2：用户
        "roomInfo":{
            //群
            "id":"123", //群组id
            "name":"群名称",
            "avatar":"群头像",
            "canAddFriend":1, //1：可以添加好友，2：不可以
            "joinPermission":1 //1：需要审批，2：不需要
        },
        "userInfo":{
            //用户
            "id":"1123",
            "name":"用户1",
            "avatar":"http://...../***.jpg",
            "position":"产品",
            "remark":"好友备注，不是好友为空"
        }
	}
}
```

### 删除指定的一条消息(websocket 通知) /chat33/

`post`
请求参数：
**Body**:

| **参数** | **名字** | 类型   | 约束 | **说明**               |
| -------- | -------- | ------ | ---- | ---------------------- |
| logId    | 消息id   | string | 必填 |                        |
| type     | 消息类型 | int    | 必填 | 1：群消息；2：好友消息 |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
	}
}
```

### 撤回指定的一条消息 (websocket 通知)/chat33/

`post`
请求参数：
**Body**:

| **参数** | **名字** | 类型   | 约束 | **说明**               |
| -------- | -------- | ------ | ---- | ---------------------- |
| logId    | 消息id   | string | 必填 |                        |
| type     | 消息类型 | int    | 必填 | 1：群消息；2：好友消息 |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
	}
}
```

## 用户相关接口

### 用户密码登录 /user/pwdLogin
`post`

请求参数：

| **参数** | **名字** | **类型** | **约束** | **说明** |
| -------- | -------- | -------- | -------- | -------- |
| mobile |  手机号    | string   | 必填     |          |
| password |  密码    | string   | 必填     |          |

Headers:

带上device

返回参数：

| **参数**   | **名字** | **类型** | **说明** |
| --------- | -------- | -------- | -------- |
| account   | 账号     | string   |         |
| id        | 用户id   | string   |       |
| uid       | 用户uid | string    |         |
| user_level   | 用户级别 | int    |        |
| avatar    | 用户头像 | string  |          |
| username  | 用户名  | string   |         |
| verified  | int    | 是否实名 | 1: 已实名 2: 未实名 |
| token     | string | 用户token |              |

```json
{
    "result": 0,
    "message": "操作成功",
    "data": {
        "account": "8612345678",
        "avatar": "",
        "id": "3306",
        "uid": "200884",
        "user_level": 1,
        "username": "dafdaf",
        "verified": 0,
        "token": "fdasfdgsadfwerf"
    }
}
```

### 用户登录 /user/tokenLogin
`post`

请求参数：

无

Headers:

带上token和device

返回参数：

| **参数**   | **名字** | **类型** | **说明** |
| --------- | -------- | -------- | -------- |
| account   | 账号     | string   |         |
| id        | 用户id   | string   |       |
| uid       | 用户uid | string    |         |
| user_level   | 用户级别 | int    |        |
| avatar    | 用户头像 | string  |          |
| username  | 用户名  | string   |         |
| verified  | int    | 是否实名 | 1: 已实名  2: 未实名 |

```json
{
    "result": 0,
    "message": "操作成功",
    "data": {
        "account": "8612345678",
        "avatar": "",
        "id": "3306",
        "uid": "200884",
        "user_level": 1,
        "username": "dafdaf",
        "verified": 0
    }
}
```

### 查看用户详情 /user/userInfo

`post`

请求参数：

**Body**:

| **参数** | **名字** | 类型   | 约束 | **说明** |
| -------- | -------- | ------ | ---- | -------- |
| id       | 对方id   | string | 必填 |          |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
        "id": "001",
        "sex":1,	// 0 未设定 1 男 2 女
        "name": "chatmRrdIsLTxT",
        "avatar": "https://pic.qqtn.com/up/2016-10/14762726301049405.jpg",
        "position":"职位",	// 互为内部员工时可看
        "isFriend", 1,		 // 1是 2否
        //仅好友可看
        "noDisturbing":1, // 1:开启消息免打扰 1:关闭 
        "stickyOnTop":1,		// 1:开启消息置顶 1:关闭
        "addTime": 1539247447457,
        "remark":"备注",	//好友的备注
        // 以下部分仅《管理员》查看时返回
        "com_id": "1", //公司id
        "uid": "101",
        "account": "12323",
        "phone":"",
        "email":"",
        "verified": 1,
        "userLevel": 1,
        "description": "描述"
	}
}
```

### 用户编辑头像 /user/editAvatar

`post`

请求参数

| **参数** | **名字** | **类型** | **约束** | **说明** |
| -------- | -------- | -------- | -------- | -------- |
| avatar |  头像URL    | string   | 必填     |          |

 返回参数

| **参数** | **名字** | **类型** | **说明**        |
| -------- | -------- | -------- | --------------- |
| result   | 结果     | int      | 0:成功，-1:失败 |

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {}
}
```

### 用户修改昵称 /user/editNickname

`post`

请求参数

| **参数** | **名字** | **类型** | **约束** | **说明** |
| -------- | -------- | -------- | -------- | -------- |
| nickname | 昵称     | string   | 必填     |          |

 返回参数

| **参数** | **名字** | **类型** | **说明**        |
| -------- | -------- | -------- | --------------- |
| result   | 结果     | int      | 0:成功，-1:失败 |

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {}
}
```

### 用户是否已注册找币账户 /user/isreg
`post`

请求参数：

| **参数** | **名字** | **类型** | **约束** | **说明** |
| -------- | ---------- | -------- | -------- | -------- |
| mobile    | 手机号 | string   | 必填      |   11位    |

返回参数

| **参数** | **名字** | **类型** | **说明**     |
| -------- | -------- | -------- | -----------|
| zb_uid      | 找币uid   | string   |  “0” 代表未注册      |

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
		"zb_uid": "3003",  //找币 uid
	}
}
```

### 管理端禁言用户（全局禁言） /user/muted

`post`

请求参数

| **参数**  | **名字** | **类型** | **约束** | **说明** |
| ---------- | -------- | -------- | -------- | -------- |
| id         | 用户id   | string   | 必填     |          |
| muted_time | 禁言时长 | int64 | 必填     |          |

返回参数

| **参数**| **名字** | **类型** | **说明**        |
| -------- | -------- | -------- | --------------- |
| result   | 结果     | int      | 0:成功，-1:失败 |

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {}
}
```

### (已弃用)获取权限信息/permission/permissionInfo
`post`

请求参数：无

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
		"permissionList": [
            {
			"permissionId": 1,
			"permissionName": "发系统消息"
            },{
                ...
            }
        ]
	}
}
```


## 聊天室有关接口

### 聊天室列表 /group/list

`post`
请求参数

| **参数**     | **名字**   | **类型** | **约束** | **说明**    |
| ------------ | ---------- | -------- | -------- | ---------|
| groupStatus | 聊天室状态 | int      | 必填  | 1：开启中，2：关闭中， 3：全部 |
| groupName | 聊天室名称 | string | 非必填 | 模糊查询 |
| startTime | 开始时间   | int64 | 非必填   | 模糊查询 |
| endTime    | 结束时间   | int64 | 非必填   | 模糊查询 |

返回参数

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
		"groups": [
            {
                "groupId": "333333",
                "groupName": "BTY交流一号",
                "avatar":"聊天室头像",
                "description":"聊天室描述",	//预留
                "createTime": 1530238075000,
                "openTime": 1530238075000,		//预留
                "closeTime": 1530238075000,	//预留
                "status": 1,	//1 开启 2 关闭
                "totalNumber": 300,
                "userNumber": 100,
                "visitorNumber": 200
            },{
                ...
            }
        ]
	}
}
```

### 聊天室详情 /group/info

`post`
请求参数

**Body**:

| **参数** | **名字** | **类型** | **约束** | **说明** |
| -------- | -------- | -------- | -------- | -------- |
| groupId  | 聊天室id | string   | 必填     |          |

返回参数
```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
        "groupId": "333333",
        "groupName": "BTY交流一号",
        "avatar":"聊天室头像",
        "description":"聊天室描述",	//预留
        "createTime": 1530238075000,
        "openTime": 1530238075000,		//预留
        "closeTime": 1530238075000,	//预留
        "status": 1,	//1 开启 2关闭
        "totalNumber": 300,
        "userNumber": 100,
        "visitorNumber": 200
	}
}
```


### （已弃用）聊天室用户列表 /group/userList

`post`

请求参数

| **参数** | **名字**    | **类型** | **约束** | **说明** |
| -------- | ------------ | -------- | -------- | -------- |
| groupId | 聊天室ID     | string   | 必填     |          |
| queryUserName | 用户名 | string | 非必填 | |
| page     | 第几页       | int      | 非必填   | 默认首页 |
| number   | 每页数据条数 | int      |          |          |

返回参数

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
		"totalnum": 9,
		"userList": [
			{
                 "id":"用户id",
				"account": "1231231",
				"name": "游客编号123",
                  "uid":"123",
				"avatar": "",
                  "remark":"备注",
				"verified": 1, 	//1 是 ；2 否
				"userLevel": 1, //0 游客 ；1 普通用户 ；2: 客服
            },{
                ...
            }
		]
	}
}
```

### （已弃用）聊天室用户信息 /group/userInfo

`post`

请求参数

| **参数** | **名字** | **类型** | **约束** | **说明** |
| -------- | -------- | -------- | -------- | -------- |
| id       | 用户id   | string   | 必填     |          |

返回参数

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
        "id": "001",
        "name": "haha",
        "avatar": "https://pic.qqtn.com/up/2016-10/14762726301049405.jpg",
        "uid": "101",
        "account": "12323",
        "mutedTime": 1516297800000,
        "mutedLastTime":1234000,
        "userLevel": 1	//0:游客，1：普通用户，2：客服
	}
}
```

### 获取聊天室聊天记录 /group/getGroupChatHistory

`post`

请求参数

| **参数** | **名字**      | **类型** | **约束** | **说明**    |
| -------- | ------------- | -------- | -------- | ---------|
| id | 聊天室id       | string   | 必填     |           |
| startId | 当前起始记录id | string | 非必填   | 从start_id的消息开始往前拉取记录，若id为空则拉取最近的消息 |
| number   | 获取记录数     | int      | 必填     |           |

返回参数

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
        "logs":[
            {
                "logId": "11",
             	"channelType": 1,
                "fromId": "1",
                "targetId": "9",
                "msgType": 1, //0：系统消息（只能是文字），1:文字，2:音频，3：图片，4：红包，5：视频
                "msg": {	
                    "content":"文本消息"
                },
                "datetime": 1530238075000,
                "senderInfo":{
                    "nickname": "客服1",
                    "avatar":"头像"
                }
            }
        ],
        "nextLog":"123"	//上一条消息id
	}
}
```

### 获取聊天室总人数 /group/getOnlineNumber

`post`

请求参数

| **参数** | **名字**   | 类型   | **约束** | **说明** |
| -------- | ---------- | ------ | -------- | -------- |
| groupId  | 聊天室名称 | string |          |          |

返回参数

```json
{
	"result": 0,
	"message": "操作成功",
    "data": {
        "totalNumber": 300
    }
}
```

### 添加聊天室 /group/addGroup

`post`

请求参数

| **参数**   | **名字**   | 类型 | **约束** | **说明**|
| ---------- | ---------- | -------- | -------- | -------- |
| groupName | 聊天室名称 | string |          |          |
| avatar | 聊天室图像url | string | | |

返回参数

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {}
}
```

### 聊天室头像编辑 /group/editAvatar

`post`

请求参数

| **参数** | **名字** | **类型** | **约束** | **说明** |
| -------- | -------- | -------- | -------- | -------- |
| groupId  | 聊天室id | string   | 必填     |          |
| avatar   | 头像URL  | string   | 必填     |          |

 返回参数

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {}
}
```

### 聊天室状态编辑 /group/setStatus

`post`

请求参数

| **参数**     | **名字** | **类型** | **约束** | **说明**       |
| ------------ | -------- | -------- | -------- | ------------|
| groupId    | 聊天室id | string   | 必填     |               |
| type | 操作类型 | int      | 必填     | 0：开启，1：关闭，2：删除 |

返回参数

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {}
}
```

### 编辑聊天室名称 /group/editGroupName

`post`

请求参数

| **参数**   | **名字**   | **类型** | **约束** | **说明**|
| ---------- | ---------- | -------- | -------- | -------- |
| groupId  | 聊天室id   | string   | 必填     |          |
| groupName | 聊天室名称 | string   | 必填     |          |

返回参数

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {}
}
```

### （已弃用）禁言聊天室用户 /group/muted

`post`

请求参数

| **参数**  | **名字** | **类型** | **约束** | **说明** |
| --------- | -------- | -------- | -------- | -------- |
| id        | 用户id   | string   | 必填     |          |
| mutedTime | 禁言时长 | int64    | 必填     |          |

返回参数

| **参数** | **名字** | **类型** | **说明**        |
| -------- | -------- | -------- | --------------- |
| result   | 结果     | int      | 0:成功，-1:失败 |

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {}
}
```

### （已弃用）举报聊天室用户 /group/report

`post`

请求参数

| **参数** | **名字**       | **类型** | **约束** | **说明** |
| -------- | -------------- | -------- | -------- | -------- |
| id       | 被举报用户的id | string   | 必填     |          |
| msgId    | 举报消息id     | string   | 必填     |          |

返回参数

| **参数** | **名字** | **类型** | **说明**            |
| -------- | -------- | -------- | ------------------- |
| result   | 结果     | int      | 0:成功，其他值:失败 |

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {}
}
```

## 群有关接口

### 获取群列表 /room/list

`post`
请求参数：
**Body**:
| **参数** | **名字** | 类型 | 约束 | **说明** |
| -------- | ---------- | -------- | -------- | -------- |
| type | 是否常用 | int | 必填 | 1：普通，2：常用 , 3：全部 |
返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
        "roomList": [
            {
                "id":"1123",
                "markId":"1231313", //用于显示的群id
                "name":"群聊1",
                "avatar":"http://...../***.jpg",
                "noDisturbing":1, 	//1：开启了免打扰，2：关闭
                "commonlyUsed": 1,   //1：普通 2 常用
                "onTop": 1 	//1：置顶 2 不置顶
            }, {
                "id":"1124",
                "markId":"1231313", //用于显示的群id
                "name":"群聊2",
                "avatar":"http://...../***.jpg",
                "noDisturbing":1,	//1：开启了免打扰，2：关闭
                "commonlyUsed": 1,   //1：普通 2 常用
                "onTop": 2 	//1：置顶 2 不置顶
            }
        ]
	}
}
```

### 获取所有群未读消息统计 /room/unread

`post`
请求参数：

无

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
        "infos": [
            {
                "id":"1123",	//群id
                "number":213,	//未读消息数
                "lastLog": {
                    "logId": "11",
                    "channelType": 2,
                    "fromId": "1",
                    "targetId": "9",
                    "msgType": 1,
                    //聊天记录
                    "msg": {	
                        "content":"文本消息"
                    },
                    "datetime": 1530238075000,
                    "senderInfo":{
                         "nickname":"昵称",
                         "avatar":"http://xxx/xxx/xxx.jpg"
                    }
                }
            }
        ]
	}
}
```

### 获取群信息 /room/info

`post`
请求参数：
**Body**:

| **参数** | **名字** | 类型   | 约束 | **说明** |
| -------- | -------- | ------ | ---- | -------- |
| roomId   | 群Id     | string | 必填 |          |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
        "id":"123", //群组id
        "markId":"1231313", //用于显示的群id
        "name":"群名称",
        "avatar":"群头像",
        "onlineNumber":123,	//在线人数
	    "memberNumber":123,	//总成员人数
        "noDisturbing":1, 	//1：开启了免打扰，2：关闭
        "memberLevel":1,	//1.普通用户;2.管理员;3.群主
        "canAddFriend":1, //1：可以添加好友，2：不可以
        "joinPermission":1, //1：需要审批，2：不需要
        "users":[
            {
                "id":"1123",
                "nickname":"群聊1",	    //用户昵称
                "roomNickname":"群聊2",	//在群中的昵称
                "avatar":"http://...../***.jpg",
                "memberLevel":1	//1.普通用户;2.管理员;3.群主
            }
        ]
	}
}
```
### 获取群成员列表 /room/userList

`post`
请求参数：
**Body**:

| **参数** | **名字** | 类型   | 约束 | **说明** |
| -------- | -------- | ------ | ---- | -------- |
| roomId   | 群Id     | string | 必填 |          |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
        "userList": [
            {
                "id":"1123",
                "nickname":"群聊1",	    //用户昵称
                "roomNickname":"群聊2",	//在群中的昵称
                "avatar":"http://...../***.jpg",
                "memberLevel":1	//1.普通用户;2.管理员;3.群主
            }, {
                "id":"1124",
                "nickname":"群聊2",
                "roomNickname":"群聊2",
                "avatar":"http://...../***.jpg",
                "memberLevel":1	//1.普通用户;2.管理员;3.群主
            }
        ]
	}
}
```

### 获取群成员信息 /room/userInfo

`post`
请求参数：
**Body**:

| **参数** | **名字** | 类型   | 约束 | **说明** |
| -------- | -------- | ------ | ---- | -------- |
| roomId   | 群Id     | string | 必填 |          |
| userId   | 群成员id | string | 必填 |          |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
        "id":"1123",
        "nickname":"群聊1",	    //用户昵称
        "roomNickname":"群聊2",	//在群中的昵称
        "avatar":"http://...../***.jpg",
        "memberLevel":1	//1.普通用户;2.管理员;3.群主
	}
}
```

### 搜索群成员信息 /room/searchMember

`post`
请求参数：
**Body**:

| **参数** | **名字**   | 类型   | 约束 | **说明**             |
| -------- | ---------- | ------ | ---- | -------------------- |
| roomId   | 群Id       | string | 必填 |                      |
| query    | 群成员名称 | string | 必填 | username或者群内昵称 |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
        "userList": [
            {
                "id":"1123",
                "nickname":"群聊1",	    //用户昵称
                "roomNickname":"群聊2",	//在群中的昵称
                "avatar":"http://...../***.jpg",
                "memberLevel":1	//1.普通用户;2.管理员;3.群主
            }, {
                "id":"1124",
                "nickname":"群聊2",
                "roomNickname":"群聊2",
                "avatar":"http://...../***.jpg",
                "memberLevel":1	//1.普通用户;2.管理员;3.群主
            }
        ]
	}
}
```

### 获取群在线人数 /room/getOnlineNumber

`post`

请求参数

| **参数** | **名字**   | 类型   | **约束** | **说明** |
| -------- | ---------- | ------ | -------- | -------- |
| roomId   | 聊天室名称 | string |          |          |

返回参数

```json
{
	"result": 0,
	"message": "操作成功",
    "data": {
        "onlineNumber":123	//在线人数
    }
}
```

### 创建群 /room/create

`post`
请求参数：
**Body**:
| **参数** | **名字** | 类型 | 约束 | **说明** |
| -------- | ---------- | -------- | -------- | -------- |
| roomName | 群名 | string | 非必填 |  |
| roomAvatar | 群头像 | string | 非必填 | 不填则为默认头像 |
| users | 群成员列表 | string[] | 必填 |  |
返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
	}
}
```

### 删除群 /room/delete

`post`
请求参数：
**Body**:

| **参数** | **名字** | 类型   | 约束 | **说明** |
| -------- | -------- | ------ | ---- | -------- |
| roomId   | 群Id     | string | 必填 |          |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
	}
}
```

### 踢出群 /room/kickOut

`post`
请求参数：
**Body**:

| **参数** | **名字**   | 类型     | 约束 | **说明** |
| -------- | ---------- | -------- | ---- | -------- |
| rootId   | 群Id       | string   | 必填 |          |
| users    | 成员id数组 | string[] | 必填 |          |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
	}
}
```

### 成员退出群 /room/loginOut

`post`
请求参数：
**Body**:

| **参数** | **名字** | 类型   | 约束 | **说明** |
| -------- | -------- | ------ | ---- | -------- |
| rootId   | 群Id     | string | 必填 |          |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
	}
}
```

### 管理员设置群 /room/setPermission

`post`
请求参数：
**Body**:
| **参数** | **名字** | 类型 | 约束 | **说明** |
| -------- | ---------- | -------- | -------- | -------- |
| roomId | 群id | string | 必填 |  |
| canAddFriend | 可否添加好友 | int | 非必填 | 1：可以，2：不可以 |
| joinPermission | 进群权限设置 | int | 非必填 | 1：需要审批，2：不需要审批，3：禁止加群 |
返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
	}
}
```
### 群内用户身份设置 /room/setLevel

`post`
请求参数：
**Body**:

| **参数** | **名字** | 类型   | 约束 | **说明**                   |
| -------- | -------- | ------ | ---- | -------------------------- |
| roomId   | 群id     | string | 必填 |                            |
| userId   | 用户id   | string | 必填 |                            |
| level    | 等级     | int    | 必填 | 1.普通用户;2.管理员;3.群主 |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
	}
}
```

### 群成员设置免打扰 /room/setNoDisturbing

`post`
请求参数：
**Body**:
| **参数** | **名字** | 类型 | 约束 | **说明** |
| -------- | ---------- | -------- | -------- | -------- |
| roomId | 群id | string | 必填 |  |
| setNoDisturbing | 设置消息免打扰 | int | 必填 | 1：开启，2：关闭 |
返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
	}
}
```

### 群成员设置群置顶  /room/stickyOnTop

`post`
请求参数：
**Body**:

| **参数**    | **名字** | 类型   | 约束 | **说明**        |
| ----------- | -------- | ------ | ---- | --------------- |
| roomId      | 群id     | string | 必填 |                 |
| stickyOnTop | 是否置顶 | int    | 必填 | 1 置顶 2 不置顶 |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
	}
}
```

### 群成员设置群内昵称  /room/setMemberNickname

`post`
请求参数：
**Body**:

| **参数** | **名字** | 类型   | 约束 | **说明** |
| -------- | -------- | ------ | ---- | -------- |
| roomId   | 群id     | string | 必填 |          |
| nickname | 昵称     | string | 必填 |          |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
	}
}
```

### 邀请入群 /room/joinRoomInvite

`post`
请求参数：
**Body**:

| **参数** | **名字**     | 类型     | 约束 | **说明** |
| -------- | ------------ | -------- | ---- | -------- |
| roomId   | 群Id         | string   | 必填 |          |
| users    | 受邀人id数组 | string[] | 必填 |          |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
	}
}
```

### 申请入群 /room/joinRoomApply

`post`
请求参数：
**Body**:

| **参数**    | **名字** | 类型   | 约束   | **说明** |
| ----------- | -------- | ------ | ------ | -------- |
| roomId      | 群Id     | string | 必填   |          |
| applyReason | 申请理由 | string | 非必填 |          |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
	}
}
```

### 入群申请处理 /room/joinRoomApprove 

`post`
请求参数：
**Body**:

| **参数** | **名字** | 类型   | 约束 | **说明**         |
| -------- | -------- | ------ | ---- | ---------------- |
| userId   | 申请人id | string | 必填 |                  |
| roomId   | 房间号   | string | 必填 |                  |
| agree    | 是否同意 | int    | 必填 | 1：同意  2：拒绝 |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
	}
}
```

### 获取群消息记录 /room/chatLog

`post`
请求参数：
**Body**:

| **参数** | **名字**       | 类型   | 约束   | **说明**                                                   |
| -------- | -------------- | ------ | ------ | ---------------------------------------------------------- |
| id       | 群id           | string | 必填   |                                                            |
| startId  | 当前起始记录id | string | 非必填 | 从start_id的消息开始往前拉取记录，若id为空则拉取最近的消息 |
| number   | 获取记录数     | int    | 必填   |                                                            |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
        "logs":[
            {
                "logId": "11",
                "channelType": 2,
                "fromId": "1",
                "targetId": "9",
                "msgType": 1,
                //聊天记录
                "msg": {	
                    "content":"文本消息"
                },
                "datetime": 1530238075000,
                "senderInfo":{
                     "nickname":"昵称",
                     "avatar":"http://xxx/xxx/xxx.jpg"
                }
            }
        ],
        "nextLog":"123" //上一条消息id
	}
}
```

## 好友有关接口

###  获取好友列表 /friend/list

`post`

请求参数：

**Body**:

| **参数** | **名字** | **类型** | **约束** | **说明** |
| -------- | ---------- | ---- | -------- | -------- |
| type | 好友类型（是否常用） | int | 必填 | 1：普通，2：常用  3：全部 |
| time | 最新时间 | int | 非必填 | 向后查询 |


返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
        "userList": [
            {
                "addTime": 1539832902773,
                "id":"1123",
                "name":"用户1",
                "avatar":"http://...../***.jpg",
                "position":"产品",
                "remark":"备注",
                "noDisturbing":1, 	//1：开启了免打扰，2：关闭
                "commonlyUsed": 1,   //1：普通 2 常用
              	"onTop": 1, 	 //1：置顶 2 不置顶
              	"isDelete": 1 //1未删除 2删除
            }, {
                ...
            }
        ]
	}
}
```

### 获取所有好友未读消息统计 /friend/unread

`post`
请求参数：

无

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
        "infos": [
            {
                "id":"1123",	//好友id
                "number":213,	//未读消息数
                "lastLog": {      
                    "logId": "11",
                    "channelType": 3,
                    "fromId": "1",
                    "targetId": "9",
                    "msgType": 1,
                    //聊天记录
                    "msg": {	
                        "content":"文本消息"
                    },
                    "datetime": 1530238075000,
                    "senderInfo":{
                         "nickname":"昵称",
                         "avatar":"http://xxx/xxx/xxx.jpg"
                    }
                }
            }
        ]
	}
}
```

### 添加好友申请/friend/add

`post`

请求参数：

**Body**:

| **参数** | **名字** | 类型 | 约束 | **说明** |
| -------- | ---------- | -------- | -------- | -------- |
| id | 对方id | string | 必填 |  |
| remark | 备注 | string | 非必填 |  |
| reason | 申请理由 | string | 非必填 |  |
| roomId | 群id | string | 非必填 | 有群id就表示是通过群加的好友 |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
	}
}
```

(websocket 通知) 

### 删除好友/friend/delete

`post`

请求参数：

**Body**:

| **参数** | **名字** | 类型 | 约束 | **说明** |
| -------- | ---------- | -------- | -------- | -------- |
| id | 对方id | string | 必填 |  |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
	}
}
```


### 好友申请处理 /friend/response

`post`

请求参数：

**Body**:

| **参数** | **名字** | 类型 | 约束 | **说明** |
| -------- | ---------- | -------- | -------- | -------- |
| id | 对方id | string | 必填 |  |
| agree | 是否同意 | int | 必填 | 1：同意 2：拒绝 |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
	}
}
```

(成功则 websocket 通知)
### 获取好友请求列表 /friend/requestList

`post`

请求参数：
**Body**:
返回参数：

```json
{
    "result": 0,
    "message": "操作成功",
    "data": [
        {
            "avatar": "https://img2.woyaogexing.com/2018/10/03/266b6ca89bba410b83958fd46af04a6b!400x400.jpeg",
            "id": "8",
            "reason": "快同意",
            "sex": "",
            "status": 1,
            "time": "1539740934915",
            "username": "你的素颜如水"
        }
    ]
}
```

### 获取好友消息记录 /friend/chatLog

`post`

请求参数：

**Body**:

| **参数** | **名字**       | 类型   | 约束   | **说明**                                                   |
| -------- | -------------- | ------ | ------ | ---------------------------------------------------------- |
| id | 好友id         | string | 必填   |                                                            |
| startId  | 当前起始记录id | string | 非必填 | 从start_id的消息开始往前拉取记录，若id为空则拉取最近的消息 |
| number   | 获取记录数     | int    | 必填   |                                                            |

返回参数：

```json
{
	"result": 0,
    "message": "操作成功",
    "data": {
        "logs": [
            {
                "channelType": 3,
                "datetime": "1537955497824",
                "fromId": "4",
                "logId": "4",
                "msg": {
                   "content": "他"
                },
                "msgType": 1,
                "remark": "亚索",
                "senderInfo": [
                    {
                        "avatar": "头像地址",
                        "name": "chatvWFYGWSHvD"
                    }
                ],
                "targetId": "1"
            }
        ],

        "nextLog": "3"
    }
}
```

### 修改好友备注 /friend/setRemark

`post`

请求参数：

**Body**:

| **参数** | **名字** | 类型   | 约束 | **说明** |
| -------- | -------- | ------ | ---- | -------- |
| id | 好友id   | string | 必填 |          |
| remark   | 备注     | string | 必填 |          |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
	}
}
```

### 设置免打扰 /friend/setNoDisturbing 

`post`

请求参数：

**Body**:

| **参数**        | **名字** | 类型   | 约束 | **说明**        |
| --------------- | -------- | ------ | ---- | --------------- |
| id        | 好友id   | string | 必填 |                 |
| setNoDisturbing | 免打扰   | int    | 必填 | 1：开启 2:关闭 |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
	}
}
```

### 设置消息置顶 /friend/stickyOnTop

`post`

请求参数：

**Body**:

| **参数**   | **名字** | 类型   | 约束 | **说明**        |
| ---------- | -------- | ------ | ---- | --------------- |
| id   | 好友id   | string | 必填 |                 |
| stickyOnTop | 消息置顶 | int    | 必填 | 1:开启 2:关闭 |

返回参数：

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
	}
}
```

## 红包接口类

### 查询账号余额 /red-packet/balance

`post`

请求参数：
无

Header里带上用户token

返回参数

| **参数** | **名字** | **类型** | **说明** |
| -------- | -------- | -------- | -------- |
| balances   | 余额     |  []list  |          |
| coin   | 币种     |  int  | 1: BTY  2: YCC   |
| amount   | 数量     |  float  |          |


```json
{
    "result": 0,
    "message": "操作成功",
    "data": {
        "balances": [
            {
                "coin": 2,
                "amount": 10
            },
            {
                "coin": 1,
                "amount": 9980
            }
        ]
    }
}
```

### (已登录用户)收红包 /red-packet/receive-entry

`post`

请求参数

| **参数**        | **名字** | **类型** | **约束** | **说明** |
| --------------- | -------- | -------- | -------- | -------- |
| packet_id | 红包id   | string   | 必填     |          |

返回参数

| **参数** | **名字** | **类型** | **说明** |
| -------- | -------- | -------- | -------- |
| amount   | 金额     | float   |          |
| total    | 红包总数 | int      |          |
| remain   | 剩余数量 | int      |          |
| coin     | 币种     | int      |  1:BTY 2:YCC        |

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
	    "amount": 3,
        "total": 10,
        "remain": 9,
        "coin": 1
	}
}
```

### (未登录用户)收红包 /red-packet/receive

`post`

请求参数

| **参数**        | **名字** | **类型** | **约束** | **说明** |
| --------------- | -------- | -------- | -------- | -------- |
| packet_id | 红包id   | string   | 必填     |          |

返回参数

| **参数** | **名字** | **类型** | **说明** |
| -------- | -------- | -------- | -------- |
| mark   | 抢到的红包份额标识  | string   |          |
| amount   | 金额     | float   |          |
| total    | 红包总数 | int      |          |
| remain   | 剩余数量 | int      |          |

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
		"mark" : "2ca28310-94a6-11e8-b452-631c93934368",
	    "amount": 4,
        "total": 10,
        "remain": 9
	}
}
```

### (未登录用户)收红包入账 /red-packet/entry

`post`

请求参数

| **参数**        | **名字** | **类型** | **约束** | **说明** |
| --------------- | -------- | -------- | -------- | -------- |
| mark |  抢到的红包份额标识   | string   | 必填     |          |
| account | 用户手机号  | string   | 必填     |          |

返回参数

| **参数** | **名字** | **类型** | **说明** |
| -------- | -------- | -------- | -------- |
| packet_id   | 红包id     | string   |          |
| amount    | 抢到的红包金额 | float      |          |
| mark   | 抢到的红包份额标识 | string      |   老用户领取领取时为空       |
| coin     | 币种     | int      |    1:BTY 2:YCC      |
| recv_type | 领取类型 | int | 1: 老用户领取  2: 未注册用户领取|
| message | 领取信息 | string | 提示信息|

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
	    "packet_id" : "b84e7d70-94a0-11e8-b452-631c93934368",
		"amount" : 7,
   		"mark" : "72a15da0-94a6-11e8-b452-631c93934368",
		"coin" : 1,
		"recv_type": 2
	}
}
```

### (未登录用户)注册领取红包 /red-packet/register-entry

`post`

请求参数

| **参数**        | **名字** | **类型** | **约束** | **说明** |
| --------------- | -------- | -------- | -------- | -------- |
| mark |  抢到的红包份额标识   | string   | 必填     |          |
| account | 注册手机号   | string   | 必填     |          |
| password | 注册密码   | string   | 必填     |          |
| captcha | 注册验证码   | string   | 必填     |          |

返回参数

| **参数** | **名字** | **类型** | **说明** |
| -------- | -------- | -------- | -------- |
| packet_id   | 红包id     | string   |          |
| amount    | 抢到的红包金额 | float      |          |
| mobile   | 已注册手机号 | string      |          |
| coin     | 币种     | int      |          |
| uid     | 新注册用户uid     | string      |          |

```json
{
	"result": 0,
	"message": "操作成功",
	"data": {
	    "red_packet_id" : "b84e7d70-94a0-11e8-b452-631c93934368",
		"amount" : 7,
		"mobile" : "18668169201",
   		"coin" : 1,
   		"uid" : "200093"
	}
}
```

### 查询用户红包收发记录 /red-packet/record

`post`

请求参数：

| **参数** | **名字**     | **约束** | 类型   | **说明** |
| -------- | ------------ | -------- | ------ | -------- |
| date     | 查询的月份或年份   | 必填     | string |  格式 2018-08 或 2018  |
| coin     | 查询的红包币种 | 预留字段，非必填  | int |  预留字段，非必填，默认查询1:BTY |

Header里带上用户token

返回参数

| **参数** | **名字** | **类型** | **说明** |
| -------- | -------- | -------- | -------- |
| receive_num | 收到红包总数     |  int  |          |
| receive_amount | 收到红包总金额 | float |       |
| send_num | 发送红包总数     |  int  |          |
| send_amount | 发送红包总金额 | float |       |
| receive  | 红包领取记录 |  list[]  |          |
| send     | 红包发送记录 |  list[]  |          |
| coin     | 币种     |  int  |          |
| amount   | 红包发送(领取)金额     |  float  |          |
| type     | 红包类型  | int     |          |
| red_packet_id | 红包id | string |        |
| receive_time | 红包领取时间 | int64 |  ms  |
| username | 红包发放者名称 | string |     |
| rsend_time | 红包发送时间 | int64 |  ms  |

```json
{
    "result": 0,
    "message": "操作成功",
    "data": {
        "receive_num": 1,
        "receive_amount": 10,
        "receive": [
            {
                "red_packet_id": "6f2fs7f9-ak17-11p8-b13c-9eee95375208",
                "type": 1,
                "amount": 10,
                "receive_time": 1535608229000,
                "coin": 3,
                "username": "zhaobiBAazYr48"
            }
        ],
        "send_num": 1,
        "send_amount": 10,
        "send": [
            {
                "red_packet_id": "6f2fs7f9-ak17-11p8-b13c-9eee95375208",
                "type": 1,
                "amount": 10,
                "send_time": 1535608229000,
                "coin": 3
            }
        ]
    }
}
```

### 查询用户红包领取记录 /red-packet/receive-record

`post`

请求参数：

| **参数** | **名字**     | **约束** | 类型   | **说明** |
| -------- | ------------ | -------- | ------ | -------- |
| date     | 查询的月份或年份   | 必填     | string |  格式 2018-08 或 2018  |
| coin     | 查询的红包币种 | 预留字段，非必填  | int |  预留字段，非必填，默认查询1:BTY |

Header里带上用户token

返回参数

| **参数** | **名字** | **类型** | **说明** |
| -------- | -------- | -------- | -------- |
| totalnum | 总条数     |  int  |          |
| total_amount | 红包总金额 | float |       |
| receive  | 红包领取记录 |  list[]  |          |
| coin     | 币种     |  int  |          |
| amount   | 数量     |  float  |          |
| type     | 红包类型  | int     |          |
| red_packet_id | 红包id | string |        |
| receive_time | 红包领取时间 | int64 |  ms  |
| username | 红包发放者名称 | string |     |

```json
{
    "result": 0,
    "message": "操作成功",
    "data": {
        "totalnum": 1,
        "total_amount": 10,
        "receive": [
            {
                "red_packet_id": "6f2fs7f9-ak17-11p8-b13c-9eee95375208",
                "type": 1,
                "amount": 10,
                "receive_time": 1535608229000,
                "coin": 1,
                "username": "zhaobiBAazYr48"
            }
        ]
    }
}
```

### 查询用户发放红包记录 /red-packet/send-record

`post`

请求参数：

| **参数** | **名字**     | **约束** | 类型   | **说明** |
| -------- | ------------ | -------- | ------ | -------- |
| date     | 查询的月份或年份   | 必填     | string |  格式 2018-08 或 2018  |
| coin     | 查询的红包币种 | 预留字段，非必填  | int |  预留字段，非必填，默认查询1:BTY |

Header里带上用户token

返回参数

| **参数** | **名字** | **类型** | **说明** |
| -------- | -------- | -------- | -------- |
| totalnum | 总条数     |  int  |          |
| total_amount | 红包总金额 | float |       |
| send  | 红包发放记录 |  list[]  |          |
| coin     | 币种     |  int  |          |
| amount   | 数量     |  float  |          |
| type     | 红包类型  | int     |          |
| red_packet_id | 红包id | string |        |
| send_time | 红包发放时间 | int64 |  ms  |

```json
{
    "result": 0,
    "message": "操作成功",
    "data": {
        "totalnum": 1,
        "total_amount": 10,
        "send": [
            {
                "red_packet_id": "6f2fs7f9-ak17-11p8-b13c-9eee95375208",
                "type": 1,
                "amount": 10,
                "send_time": 1535608229000,
                "coin": 3
            }
        ]
    }
}
```

## webSocket接口：
### 普通消息

##### Msg 格式说明：

0 系统消息: {"content":"发个公告试试"}

1 文字消息：{content:"heeleee"}

2 语音: {"mediaUrl":"http://zb-chat.oss-cn-shanghai.aliyuncs.com/chatList/voice/20180802/201808021644248_200408.amr","time":6}	

3 图片：{"height":2560,"imageUrl":"https://zb-chat.oss-cn-shanghai.aliyuncs.com/chatList/picture/20180814/20180814203203689_4.jpg","width":1216}	

4 红包: {"coin":3,"packetId":"2168bf60-9f90-11e8-b129-6343f801ef65","packetType":2,"packetUrl":"http://47.74.190.154/packets/ycc/pages/open.html?id=2168bf60-9f90-11e8-b129-6343f801ef65"}

6 通知消息{"content":"消息内容"}

**发送：**

| **参数**    | **名字** | **类型** | **说明**                                               |
| ----------- | -------- | -------- | ------------------------------------------------------ |
| eventType   | 事件类型 | int      | 0: 普通消息                                            |
| msgId       | 消息id   | string   |                                                        |
| channelType | 类型     | int      | 1: 聊天室； 2：群组；3：好友                           |
| targetId    | 接收者id | string   |                                                        |
| msgType     | 消息类型 | int      | 0：系统消息，1:文字，2:音频，3：图片，4：红包，5：视频 |
| msg         | 消息内容 | object   |                                                        |

发送示例：

```json
{
	"eventType": 0,
	"msgId": "123123",
	"channelType": 1,
	"targetId": "9",
	"msgType": 1,
    "msg": {
        "content":"文本消息"
    }
}
```

**接收：**

消息发送成功时返回：

| **参数**    | **名字**   | **类型** | **说明**                                               |
| ----------- | ---------- | -------- | ------------------------------------------------------ |
| eventType   | 事件类型   | int      | 0: 普通消息                                            |
| msgId       | 消息id     | string   | 与发送时的msgId相同                                    |
| channelType | 类型       | int      | 1: 聊天室； 2：群组；3：好友                           |
| fromId      | 发送者id   | string   |                                                        |
| targetId    | 接收者id   | string   |                                                        |
| datetime    | 发送时间   | int64    |                                                        |
| logId       | 消息记录id | int      | 消息记录的id号                                         |
| msgType     | 消息类型   | int      | 0：系统消息，1:文字，2:音频，3：图片，4：红包，5：视频 |
| msg         | 消息内容   | object   |                                                        |
| senderInfo  | 发送者信息 | object   |                                                        |

**senderInfo 格式**

| **参数** | **名字**   | **类型** | **说明** |
| -------- | ---------- | -------- | -------- |
| nickname | 发送者昵称 | string   |          |
| avatar   | 头像       | string   |          |

接收示例：

```json
{
	"eventType": 0,
	"msgId": "123123",
	"channelType": 1,
    "fromId": "1",
	"targetId": "9",
    "datetime": 1568738777731,
    "logId":"123",
	"msgType": 1,
    "msg": {
        "content":"文本消息"
    },
    "senderInfo":{
         "nickname":"昵称",
         "avatar":"http://xxx/xxx/xxx.jpg"
    }
}
```

消息发送失败时返回：

| **参数**  | **名字** | **类型** | **说明**    |
| --------- | -------- | -------- | ----------- |
| eventType | 事件类型 | int      | 0: 普通消息 |
| msgId     | 消息id   | string   |             |
| code      | 错误代码 | int      |             |
| content   | 内容     | string   |             |

具体错误码见[websocket错误代码](#统一错误代码) **注意：**消息发送成功时没有错误码

### 登录聊天室 C->S

| **参数**  | **名字** | **类型** | **说明**     |
| --------- | -------- | -------- | ------------ |
| eventType | 事件类型 | int      | 1:登录聊天室 |
| msgId     | 消息id   | string   |              |
| groupId   | 群id     | string   |              |

返回数据：

| **参数**  | **名字** | **类型** | **说明**   |
| --------- | -------- | -------- | ---------- |
| eventType | 事件类型 | int      | 1          |
| msgId     | 消息id   | string   |            |
| content   | 内容     | string   |            |
| code      | 错误代码 | int      | 0 表示成功 |

### 退出聊天室 C->S

| **参数**  | **名字** | **类型** | **说明**     |
| --------- | -------- | -------- | ------------ |
| eventType | 事件类型 | int      | 2:退出聊天室 |
| msgId     | 消息id   | string   |              |
| groupId   | 群id     | string   |              |

返回数据：

| **参数**  | **名字** | **类型** | **说明**   |
| --------- | -------- | -------- | ---------- |
| eventType | 事件类型 | int      | 1          |
| msgId     | 消息id   | string   |            |
| content   | 内容     | string   |            |
| code      | 错误代码 | int      | 0 表示成功 |

### 用户禁言通知 S->C

| **参数**      | **名字**     | **类型** | **说明**           |
| ------------- | ------------ | -------- | ------------------ |
| eventType     | 事件类型     | int      | 3: 踢出聊天室/禁言 |
| mutedTime     | 禁言开始时间 | datetime |                    |
| mutedLastTime | 禁言持续时间 | datetime |                    |
| datetime      | 发送时间     | datetime |                    |

### 关闭聊天室通知 S->C

| **参数**  | **名字** | **类型** | **说明**      |
| --------- | -------- | -------- | ------------- |
| eventType | 事件类型 | int      | 4: 关闭聊天室 |
| groupId   | 聊天群id | string   |               |
| datetime  | 发送时间 | datetime |               |

### 删除聊天室通知 S->C

| **参数**  | **名字** | **类型** | **说明**      |
| --------- | -------- | -------- | ------------- |
| eventType | 事件类型 | int      | 5: 删除聊天室 |
| groupId   | 聊天群id | string   |               |
| datetime  | 发送时间 | datetime |               |

### 开启聊天室通知 S->C

| **参数**   | **名字** | **类型** | **说明**        |
| ---------- | -------- | -------- | --------------- |
| eventType | 事件类型 | int      | 7: 聊天室被开启 |
| groupId    | 聊天室id       | string   |                  |
| groupName  | 聊天室名称     | string   |                  |
| avatar | 聊天室头像 | string | 头像的地址 |
| description | 聊天室描述     | string   | 预留             |
| createTime | 添加时间       | datetime |                  |
| openTime   | 开放时间       | datetime | 预留             |
| closeTime  | 关闭时间       | datetime | 预留             |
| status      | 聊天室开关状态 | int      | 0：开启，1：关闭 |
| totalNumber | 聊天室内总人数   | string         |          |
| userNumber | 用户数           | string         |          |
| visitorNumber | 游客数           | string         |       |
返回数据：无

### 消息有关通知 S->C（暂不用）

| **参数**   | **名字** | **类型** | **说明**        |
| ---------- | -------- | -------- | --------------- |
| eventType | 事件类型 | int      | 10: 消息操作有关 |
| logType | 所属类型 | int   | 1 :群消息 2：好友消息 |
| logId | 消息id   | string   |                  |
| operateType | 操作类型 | int | 1：撤回  2：删除 |
| userName | 用户名称 | string |  |
返回数据：无

### 入群通知 S->C  //1.创建群 2.被邀请者入群 3.直接入群回复

| **参数**  | **名字**   | **类型** | **说明**    |
| --------- | ---------- | -------- | ----------- |
| eventType | 事件类型   | int      | 20:入群通知 |
| roomId    | 群id       | string   |             |
| userId    | 入群用户id | string   |             |
| datetime  | 发送时间   | datetime |             |

返回数据：无

### 退群通知 S->C  //主动退出群 或被踢出群聊

| **参数**  | **名字**   | **类型** | **说明**                |
| --------- | ---------- | -------- | ----------------------- |
| eventType | 事件类型   | int      | 21:退群通知             |
| roomId    | 群id       | string   |                         |
| userId    | 退群用户id | string   |                         |
| type      | 退群原因   | int      | 1:主动退群 2:被踢出群聊 |
| content   | 提示信息   | string   |                         |

返回数据：无

### 解散群通知 S->C

| **参数**  | **名字** | **类型** | **说明**   |
| --------- | -------- | -------- | ---------- |
| eventType | 事件类型 | int      | 22: 解散群 |
| roomId    | 群id     | string   |            |
| datetime  | 发送时间 | datetime |            |

### 入群请求和回复通知 S->C

| **参数**  | **名字** | **类型** | **说明**         |
| --------- | -------- | -------- | ---------------- |
| eventType | 事件类型 | int      | 23: 入群请求推送 |

```json
{
    "eventType":23,
    senderInfo:{
        "id":"1123",
        "name":"用户1",
        "avatar":"http://...../***.jpg",
        "position":"产品"
    },
    receiveInfo:{
        "id":"1123",
        "name":"用户1",
        "avatar":"http://...../***.jpg",
        "position":"产品"
    },
    "id":123,
    "type":1, //1 群 2 好友
    "applyReason":"申请理由",
    "status":1, //1:等待验证 2:已拒绝 3:已同意 
    "datetime":1676764266167
}
```

### 群在线人数更新通知 S->C

| **参数**  | **名字** | **类型** | **说明**           |
| --------- | -------- | -------- | ------------------ |
| eventType | 事件类型 | int      | 24: 群在线人数更新 |
| roomId    | 群id     | string   |                    |
| number    | 人数     | string   |                    |
| datetime  | 发送时间 | datetime |                    |

### 添加好友申请和回复通知 S->C

| **参数**  | **名字** | **类型** | **说明**             |
| --------- | -------- | -------- | -------------------- |
| eventType | 事件类型 | int      | 31: 添加好友消息通知 |

```json
{
    "eventType":31,
    senderInfo:{
        "id":"1123",
        "name":"用户1",
        "avatar":"http://...../***.jpg",
        "position":"产品"
    },
    receiveInfo:{
        "id":"1123",
        "name":"用户1",
        "avatar":"http://...../***.jpg",
        "position":"产品"
    },
    "id":123,
    "type":1, //1 群 2 好友
    "applyReason":"申请理由",
    "status":1, //1:等待验证 2:已拒绝 3:已同意 
    "datetime":1676764266167
}
```

返回数据：无

