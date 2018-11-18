package db

import (
	"strconv"

	"gitlab.33.cn/chat/chat33/types"
	"gitlab.33.cn/chat/chat33/utility"
)

/*
	获取好友列表
*/
func GetFriendListByTime(id string, tp int, time int) ([]map[string]string, error) {
	var sqlStr string
	if tp == 3 {
		sqlStr = "select * from `friends` as f left join `user` as u on f.friend_id=u.user_id where f.user_id = ? and f.add_time >= ?"
		return conn.Query(sqlStr, id, time)
	} else {
		sqlStr = "select * from `friends` as f left join `user` as u on f.friend_id=u.user_id where f.user_id = ? and f.type = ? and f.add_time >= ?"
		return conn.Query(sqlStr, id, strconv.Itoa(tp), time)
	}

}

func GetFriendList(id string, tp, isDelete int) ([]map[string]string, error) {
	var sqlStr string
	if tp == 3 {
		sqlStr = "select * from `friends` as f left join `user` as u on f.friend_id=u.user_id where f.user_id = ? and f.is_delete = ?"
		return conn.Query(sqlStr, id, isDelete)
	} else {
		sqlStr = "select * from `friends` as f left join `user` as u on f.friend_id=u.user_id where f.user_id = ? and f.type = ? and f.is_delete = ?"
		return conn.Query(sqlStr, id, strconv.Itoa(tp), isDelete)
	}

}

/*
	添加好友申请
*/
func AddFriendRequest(tp int, userID, friendID, reason, remark string) (bool, error) {
	//如果发送过申请 更新请求时间和原因（多次发送或者删除好友之后重新加）
	const sqlStr = "insert into `apply` values(?, ?, ?, ?, ?, ?, ?) on duplicate key update `state` = 0, `apply_reason` = ?, `datetime` = ?,`remark` = ?"
	now := utility.NowMillionSecond()
	_, _, err := conn.Exec(sqlStr, tp, userID, friendID, reason, types.FriendStatusUnHandle, remark, now, reason, now, remark)
	if err != nil {
		return false, err
	}
	return true, nil
}

//添加好友请求
func AddApply(tp int, userID, friendID, reason, remark string, roomId int) (int, error) {
	sqlStr := "insert into `apply` (type,apply_user,target,apply_reason,state,remark,datetime,room_id)values(?, ?, ?, ?, ?, ?, ?, ?)"
	now := utility.NowMillionSecond()
	_, id, err := conn.Exec(sqlStr, tp, userID, friendID, reason, types.FriendStatusUnHandle, remark, now, roomId)
	return int(id), err
}

//更新好友申请数据
func UpdateApply(reason, remark, uid, fid string, state, tp, roomId int) error {
	sql := "update apply set room_id = ?,apply_reason = ?,state = ?,remark = ?,datetime = ? where apply_user = ? and target = ? and type = ?"
	now := utility.NowMillionSecond()
	_, _, err := conn.Exec(sql, roomId, reason, state, remark, now, uid, fid, tp)
	return err
}

/*
	同意好友申请
*/
func AcceptFriend(userID, friendID string) error {
	tx, err := conn.NewTx()
	if err != nil {
		return err
	}

	const updateStatus = "update `apply` set `state` = ?,datetime = ? where apply_user = ? and target = ? and type = ?"
	_, _, err = tx.Exec(updateStatus, types.FriendStatusAccept, utility.NowMillionSecond(), friendID, userID, types.FriendApply)
	if err != nil {
		tx.RollBack()
		return err
	}

	//添加好友
	const addFriend = "insert into `friends` values(?, ?, ?, ?, ?, ?, ?, ?) on duplicate key update add_time = ?, DND = ?, top = ?, type = ?, is_delete = ?"
	// id friend_id remark add_time, DND, top

	now := utility.NowMillionSecond()
	_, _, err = tx.Exec(addFriend, userID, friendID, "", now, types.FriendIsNotDND, types.FriendIsNotTop, types.FriendCommon, types.FriendIsNotDelete, now, types.FriendIsNotDND, types.FriendIsNotTop, types.FriendCommon, types.FriendIsNotDelete)

	if err != nil {
		tx.RollBack()
		return err
	}

	//读取请求的备注信息
	readRemark := "select remark from apply where apply_user = ? and target = ? and type = ?"
	res, err := tx.Query(readRemark, friendID, userID, types.FriendApply)
	if err != nil {
		tx.RollBack()
		return err
	}
	remark := res[0]["remark"]

	//添加好友
	const addFriend1 = "insert into `friends` values(?, ?, ?, ?, ?, ?, ?, ?) on duplicate key update remark = ?, add_time = ?, DND = ?, top = ?, type = ?, is_delete = ?"
	_, _, err = tx.Exec(addFriend1, friendID, userID, remark, now, types.FriendIsNotDND, types.FriendIsNotTop, types.FriendCommon, types.FriendIsNotDelete, remark, now, types.FriendIsNotDND, types.FriendIsNotTop, types.FriendCommon, types.FriendIsNotDelete)
	if err != nil {
		tx.RollBack()
		return err
	}

	return tx.Commit()
}

/*
	拒绝好友申请
*/
func RejectFriend(userID, friendID string) error {
	const sqlStr = "update `apply` set `state` = ?,datetime = ? where apply_user = ? and target = ? and type = ?"
	_, _, err := conn.Exec(sqlStr, types.FriendStatusReject, utility.NowMillionSecond(), friendID, userID, types.FriendApply)
	return err
}

/*
	设置好友备注
*/
func SetFriendRemark(userID, friendID, remark string) error {
	const sqlStr = "update `friends` set `remark` = ? where user_id = ? and friend_id = ?"
	_, _, err := conn.Exec(sqlStr, remark, userID, friendID)
	return err
}

/*
	设置好友免打扰
*/
func SetFriendDND(userID, friendID string, DND int) error {
	const sqlStr = "update `friends` set `DND` = ? where user_id = ? and friend_id = ?"
	_, _, err := conn.Exec(sqlStr, DND, userID, friendID)
	return err
}

/*
	设置好友置顶
*/
func SetFriendTop(userID, friendID string, top int) error {
	const sqlStr = "update `friends` set `top` = ? where user_id = ? and friend_id = ?"
	_, _, err := conn.Exec(sqlStr, top, userID, friendID)
	return err
}

//删除好友
func DeleteFriend(userID, friendID string) error {
	const sqlStr = "update `friends` set is_delete = ?,add_time = ? where user_id = ? and friend_id = ?"
	tx, err := conn.NewTx()
	if err != nil {
		return err
	}
	now := utility.NowMillionSecond()

	_, _, err = tx.Exec(sqlStr, types.FriendIsDelete, now, userID, friendID)
	if err != nil {
		tx.RollBack()
		return err
	}
	_, _, err = tx.Exec(sqlStr, types.FriendIsDelete, now, friendID, userID)
	if err != nil {
		tx.RollBack()
		return err
	}
	return tx.Commit()
}

// 检查是否是好友关系
func CheckFriend(userID, friendID string, isDelete int) (bool, error) {
	sqlStr := "select * from `friends` where user_id = ? and friend_id = ? and is_delete = ?"
	rows, err := conn.Query(sqlStr, userID, friendID, isDelete)
	if err != nil {
		return false, err
	}
	return len(rows) > 0, nil
}

// 用户是否存在
func UserIsExists(userID string) (bool, error) {
	sqlStr := "select * from `user` where user_id = ?"
	rows, err := conn.Query(sqlStr, userID)
	if err != nil {
		return false, err
	}
	return len(rows) > 0, nil
}

//查询好友请求是否存在
func FindFriendRequest(userID, friendID string) ([]map[string]string, error) {
	sqlStr := "SELECT `state` FROM apply WHERE apply_user = ? AND target = ? and type = ?"
	return conn.Query(sqlStr, friendID, userID, types.FriendApply)
}

//查询好友请求信息
func FindFriendRequestInfo(userID, friendID string) ([]map[string]string, error) {
	sqlStr := "SELECT * FROM apply WHERE apply_user = ? AND target = ? and type = ?"
	return conn.Query(sqlStr, friendID, userID, types.FriendApply)
}

//查询好友请求数量
func FindApplyCount(userID, friendID string) (int32, error) {
	sqlStr := "SELECT *  FROM apply WHERE apply_user = ? AND target = ? and type = ?"
	rows, err := conn.Query(sqlStr, userID, friendID, types.FriendApply)
	if err != nil {
		return 0, err
	}
	return int32(len(rows)), nil
}

//查询好友请求数量
func FindApplyId(userID, friendID string) ([]map[string]string, error) {
	sqlStr := "SELECT id  FROM apply WHERE apply_user = ? AND target = ? and type = ?"
	return conn.Query(sqlStr, userID, friendID, types.FriendApply)
}

//查看好友详情
func FindUserInfo(friendID string) ([]map[string]string, error) {
	sqlStr := "SELECT * FROM `user` WHERE user_id = ?"
	return conn.Query(sqlStr, friendID)
}

//查看好友关系详情 备注等信息
func FindFriend(userID, friendID string) ([]map[string]string, error) {
	sqlStr := "SELECT * FROM `friends` WHERE user_id = ? AND friend_id = ?"
	return conn.Query(sqlStr, userID, friendID)
}

//获取最新的消息id
func FindLastCatLogId(userID, friendID string) ([]map[string]string, error) {
	sqlStr := "SELECT MAX(`id`) FROM private_chat_log WHERE (sender_id = ? AND receive_id = ?) OR (sender_id = ? AND receive_id = ?)"
	return conn.Query(sqlStr, userID, friendID, friendID, userID)
}

//查找消息记录
func FindCatLog(userID, friendID string, start, number int) ([]map[string]string, error) {
	sqlStr := "SELECT * FROM private_chat_log WHERE ((sender_id = ? AND receive_id = ?) OR (sender_id = ? AND receive_id = ?)) AND id < ? ORDER BY id DESC LIMIT ?"
	return conn.Query(sqlStr, userID, friendID, friendID, userID, start, number)
}

//查找user的名字头像
func SenderInfo(userID string) ([]map[string]string, error) {
	sqlStr := "SELECT username as name,avatar FROM user WHERE user_id = ?"
	return conn.Query(sqlStr, userID)
}

//查询该条好友聊天记录是否属于user
func CheckCatLogIsUser(userId, id string) (bool, error) {
	sqlStr := "SELECT * from private_chat_log WHERE id = ? and sender_id = ?"
	rows, err := conn.Query(sqlStr, id, userId)
	if err != nil {
		return false, err
	}
	return len(rows) > 0, nil
}

//删除好友聊天记录
func DeleteCatLog(id string) (int, error) {
	sqlStr := "DELETE FROM private_chat_log WHERE id = ?"
	num, _, err := conn.Exec(sqlStr, id)
	return int(num), err
}

//查询该条群聊天记录是否属于user
func CheckRoomMsgContentIsUser(userId, id string) (bool, error) {
	sqlStr := "SELECT * from room_msg_content WHERE id = ? and sender_id = ?"
	rows, err := conn.Query(sqlStr, id, userId)
	if err != nil {
		return false, err
	}
	return len(rows) > 0, nil
}

//删除群聊天记录
func DeleteRoomMsgContent(id string) (int, error) {
	sqlStr := "DELETE FROM room_msg_content WHERE id = ?"
	num, _, err := conn.Exec(sqlStr, id)
	return int(num), err
}

//根据uid查找用户id
func FindUserByMarkId(uid string) ([]map[string]string, error) {
	sqlStr := "SELECT user_id FROM user WHERE account = ? or uid = ?"
	return conn.Query(sqlStr, uid, uid)
}

//获取所有好友未读消息数
func GetAllFriendUnreadMsgCountByUserId(uid string, status int) (int32, error) {
	sqlStr := "SELECT * FROM friends f RIGHT JOIN private_chat_log c ON f.friend_id = c.sender_id  WHERE f.user_id = ? AND c.status = ?"
	rows, err := conn.Query(sqlStr, uid, status)
	if err != nil {
		return 0, err
	}
	return int32(len(rows)), nil
}

//查询所有好友id
func FindFriendIdByUserId(uid string) ([]map[string]string, error) {
	sql := "SELECT friend_id  FROM friends  WHERE user_id = ?"
	return conn.Query(sql, uid)
}

//查询好友id头像备注（昵称）
func FindFriendInfoByUserId(uid string, fid string) ([]map[string]string, error) {
	sql := "SELECT f.friend_id,u.avatar,  IF(ISNULL(f.remark)||LENGTH(f.remark)<1,u.username,f.remark) AS username  FROM friends f  left join user u on f.friend_id = u.user_id WHERE f.user_id = ? and f.friend_id = ?"
	return conn.Query(sql, uid, fid)
}

//查询未读消息数
func FindUnReadNum(uid, fid string, status int) (int32, error) {
	sql := "SELECT *  FROM private_chat_log WHERE sender_id = ? AND receive_id = ? AND status = ?"
	rows, err := conn.Query(sql, fid, uid, status)
	if err != nil {
		return 0, err
	}
	return int32(len(rows)), nil
}

//查询好友第一条聊天记录
func FindFirstMsg(uid, fid string) ([]map[string]string, error) {
	sql := "SELECT * FROM private_chat_log c WHERE c.sender_id = ? AND c.status = ? AND c.receive_id = ? ORDER BY id DESC LIMIT 0,1"
	return conn.Query(sql, fid, types.NotRead, uid)
}

//FindCid 查找cid
func FindCid(uid string) ([]map[string]string, error) {
	sql := "select getui_cid as cid from user where user_id = ?"
	return conn.Query(sql, uid)
}

//添加私聊聊天记录
func AddPrivateChatLog(senderId, receiveId string, msgType int, content string, status int) (int64, int64, error) {
	sql := "INSERT INTO private_chat_log (sender_id,receive_id,msg_type,content,status,send_time) VALUES (?,?,?,?,?,?)"
	return conn.Exec(sql, senderId, receiveId, msgType, content, status, utility.NowMillionSecond())
}

//修改聊天记录状态
func ChangePrivateChatLogStstus(id, status int) (int64, int64, error) {
	sql := "UPDATE private_chat_log SET status = ? WHERE id = ?"
	return conn.Exec(sql, status, id)
}

//查找聊天记录
func FindPrivateChatLog(senderId, receiveId string) ([]map[string]string, error) {
	sql := "SELECT * FROM private_chat_log WHERE sender_id = ? AND receive_id = ?"
	return conn.Query(sql, senderId, receiveId)
}

//通过状态查找聊天记录
func FindPrivateChatLogByStatus(senderId, receiveId string, status int) ([]map[string]string, error) {
	sql := "SELECT * FROM private_chat_log WHERE sender_id = ? AND receive_id = ? and status = ?"
	return conn.Query(sql, senderId, receiveId, status)
}

//修改聊天记录状态
func ChangePrivateChatLogStstusByUserAndFriendId(uid, fid string) (int64, int64, error) {
	sql := "UPDATE private_chat_log SET status = ? WHERE sender_id = ? and receive_id = ?"
	return conn.Exec(sql, types.HadRead, fid, uid)
}

//判断群是否存在
func CheckRoomIsExist(roomId int) (bool, error) {
	sql := `SELECT is_delete FROM room WHERE id = ?`
	rows, err := conn.Query(sql, roomId)
	if err != nil {
		return false, err
	}
	if len(rows) <= 0 {
		return false, nil
	}

	isDelete, err := strconv.Atoi(rows[0]["is_delete"])
	if err != nil {
		return false, err
	}
	return isDelete == 1, nil

}

//判断用户是否在群里
func CheckUserInRoom(userId string, roomId int) (bool, error) {
	sql := `SELECT is_delete FROM room_user WHERE user_id = ? AND room_id = ?`
	rows, err := conn.Query(sql, userId, roomId)
	if err != nil {
		return false, err
	}
	if len(rows) <= 0 {
		return false, nil
	}

	isDelete, err := strconv.Atoi(rows[0]["is_delete"])
	if err != nil {
		return false, err
	}
	return isDelete == 1, nil
}

//判断该群是否允许添加好友
func CheckRoomIsCanAddFriend(roomId int) (bool, error) {
	sql := `SELECT can_add_friend FROM room WHERE id = ?`
	rows, err := conn.Query(sql, roomId)
	if err != nil {
		return false, err
	}
	if len(rows) <= 0 {
		return false, nil
	}

	canAdd, err := strconv.Atoi(rows[0]["can_add_friend"])
	if err != nil {
		return false, err
	}
	return canAdd == 1, nil
}

//查询群成员id
func FindRoomMemberIds(roomId int) ([]map[string]string, error) {
	sql := `SELECT user_id FROM room_user WHERE room_id = ? AND is_delete = ?`
	return conn.Query(sql, roomId, types.IsNotDelete)
}

//查询群昵称 没有的话 返回用户名称
func FindRoomMemberName(roomId int, userId string) ([]map[string]string, error) {
	sql := `SELECT IF(ISNULL(r.user_nickname)||LENGTH(r.user_nickname)<1,u.username,r.user_nickname) AS name FROM room_user r RIGHT JOIN user u ON r.user_id = u.user_id  WHERE r.room_id =? AND r.user_id = ?`
	return conn.Query(sql, roomId, userId)
}
