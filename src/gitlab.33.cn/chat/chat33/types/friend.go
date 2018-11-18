package types

//好友请求状态
const (
	FriendStatusUnHandle = 1
	FriendStatusAccept   = 3
	FriendStatusReject   = 2

	FriendRequestAccept = 1
	FriendRequestReject = 2
)

//好友表里的字段
const (
	//是否删除 1 未删除 2 已删除
	FriendIsNotDelete = 1
	FriendIsDelete    = 2

	//是否消息免打扰 1免打扰 2关闭
	FriendIsDND    = 1
	FriendIsNotDND = 2

	//好友置顶 1置顶 2不置顶
	FriendIsTop    = 1
	FriendIsNotTop = 2

	//好友类型 1 普通 2 常用
	FriendCommon     = 1
	FriendFrequently = 2
)

//群消息
const RoomMsg = 1

//好友消息
const FriendMsg = 2

//私聊消息状态
const HadRead = 1 //已读
const NotRead = 2 //未读

//apply表
//加群申请
const GroupApply = 1

//加好友申请
const FriendApply = 2

//是否删除 1 未删除 2 已删除
const IsNotDelete = 1
const IsDelete = 2
