package db

import (
	"dev.33.cn/33/common/mysql"
	"gitlab.33.cn/chat/chat33/types"
	"gitlab.33.cn/chat/chat33/utility"
)

func CheckRoomMarkIdExist(id string) (bool, error) {
	const sqlStr = "select * from room where mark_id = ?"
	maps, err := conn.Query(sqlStr, id)
	if err != nil {
		return false, err
	}
	if len(maps) > 0 {
		return true, nil
	}
	return false, nil
}

// 获取群中某成员信息
func GetRoomMemberInfo(roomId, userId string) ([]map[string]string, error) {
	const sqlStr = "select * from room_user left join `user` on room_user.user_id = `user`.user_id where room_id = ? and room_user.user_id = ? and is_delete = ?"
	return conn.Query(sqlStr, roomId, userId, types.RoomUserNotDeleted)
}

// 获取群中某成员信息
func GetRoomMemberInfoByName(roomId, name string) ([]map[string]string, error) {
	var sqlStr = "select * from room_user left join `user` on room_user.user_id = `user`.user_id where room_id = ? and is_delete = ? and (user_nickname LIKE '%" + name + "%' or username LIKE '%" + name + "%')"
	return conn.Query(sqlStr, roomId, types.RoomUserNotDeleted)
}

// 获取群成员总数
func GetRoomMemberNumber(roomId string) (int, error) {
	const sqlStr = "select count(*) as count from room_user where room_id = ? and is_delete = ?"
	maps, err := conn.Query(sqlStr, roomId, types.RoomUserNotDeleted)
	if err != nil || len(maps) == 0 {
		return 0, err
	}
	return utility.ToInt(maps[0]["count"]), nil
}

// 获取群成员信息
func GetRoomMembers(roomId string, searchNumber int) ([]map[string]string, error) {
	if searchNumber == -1 {
		const sqlStr = "select * from room_user left join `user` on room_user.user_id = `user`.user_id where room_id = ? and is_delete = ?"
		return conn.Query(sqlStr, roomId, types.RoomUserNotDeleted)
	} else {
		const sqlStr = "select * from room_user left join `user` on room_user.user_id = `user`.user_id where room_id = ? and is_delete = ? limit 0,?"
		return conn.Query(sqlStr, roomId, types.RoomUserNotDeleted, searchNumber)
	}
}

// 获取群中管理员和群主信息
func GetRoomManagerAndMaster(roomId string) ([]map[string]string, error) {
	const sqlStr = "select * from room_user left join `user` on room_user.user_id = `user`.user_id where room_id = ? and `level` > 1 and is_delete = ?"
	return conn.Query(sqlStr, roomId, types.RoomUserNotDeleted)
}

// 获取(群成员对群的配置)详情
func GetRoomsInfoAsUser(roomId, userId string) ([]map[string]string, error) {
	const sqlStr = "select * from room left join room_user on room.id = room_user.room_id where room.id = ? and user_id = ?"
	return conn.Query(sqlStr, roomId, userId)
}

func DeleteRoomById(roomId string) error {
	const sqlStr = "update room set is_delete = ? where id = ?"
	_, _, err := conn.Exec(sqlStr, types.RoomDeleted, roomId)
	return err
}

// 删除群成员
func DeleteRoomMemberById(userId, roomId string) error {
	const sqlStr = "delete from room_user WHERE room_id = ? and user_id = ?"
	_, _, err := conn.Exec(sqlStr, roomId, userId)
	return err
}

//
func GetRoomList(user string, Type int) ([]map[string]string, error) {
	var sql = "select * from room_user left join room on room_user.room_id = room.id where user_id = ? and room_user.is_delete = ? and room.is_delete = ?"
	switch Type {
	case 1:
		sql += " and common_use = " + utility.ToString(types.RoomUncommonUse)
	case 2:
		sql += " and common_use = " + utility.ToString(types.RoomCommonUse)
	case 3:
	}
	sql += " order by name"
	return conn.Query(sql, user, types.RoomUserNotDeleted, types.RoomNotDeleted)
}

// 获取所有群
func GetEnabledRooms() ([]map[string]string, error) {
	const sqlStr = "select * from room where is_delete = ?"
	return conn.Query(sqlStr, types.RoomNotDeleted)
}

// 获取群详情
func GetRoomsInfo(roomId string) ([]map[string]string, error) {
	const sqlStr = "select * from room where room.id = ? and is_delete = ?"
	return conn.Query(sqlStr, roomId, types.RoomNotDeleted)
}

// 获取群详情通过markId
func GetRoomsInfoByMarkId(markId string) ([]map[string]string, error) {
	const sqlStr = "select * from room where room.mark_id = ? and is_delete = ?"
	return conn.Query(sqlStr, markId, types.RoomNotDeleted)
}

// 根据userId获取加入的所有群
func GetRoomsById(id string) ([]map[string]string, error) {
	const sqlStr = "select * from room_user where user_id = ? and is_delete = ?"
	return conn.Query(sqlStr, id, types.RoomUserNotDeleted)
}

// 根据userId获取管理的所有群
func GetManageRoomsById(id string) ([]map[string]string, error) {
	const sqlStr = "select * from room_user where user_id = ? and `level` > 1 and is_delete = ?"
	return conn.Query(sqlStr, id, types.RoomUserNotDeleted)
}

//添加入群申请记录
func AppendJoinRoomApplyLog(roomId, userId, applyReason string, state int) (int64, int64, error) {
	createTime := utility.NowMillionSecond()
	const sqlStr = "insert into `apply` values(?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE state = ?, datetime = ?"
	return conn.Exec(sqlStr, nil, types.RoomType, userId, roomId, applyReason, state, "", createTime, state, createTime)
}

// 添加群聊聊天日志
func AppendRoomChatLog(sender_id, receive_id, msg_type, content string) (int64, int64, error) {
	const sqlStr = "insert into room_msg_content(room_id,sender_id,msg_type,content,datetime) values(?,?,?,?,?)"
	return conn.Exec(sqlStr, receive_id, sender_id, msg_type, content, utility.NowMillionSecond())
}

// 添加群聊接收日志
func AppendRoomMemberReceiveLog(roomLogId, receive_id, state string) (string, bool) {
	const sqlStr = "insert into room_msg_receive(id, room_msg_id, receive_id, state) values(?,?,?,?)"
	_, _, err := conn.Exec(sqlStr, nil, roomLogId, receive_id, state)
	if err != nil {
		return "-1", false
	}

	maps, err := conn.Query("select Max(id) from room_msg_receive")
	if err != nil || len(maps) < 1 {
		return "-1", false
	}
	return utility.ToString(maps[0]["id"]), true
}

// 根据userId将群消息接收状态更新为已读
func UpdateReceiveStateReaded(roomId, userId string) (int64, int64, error) {
	const sqlStr = "update room_msg_receive set state = ? where room_msg_id = ? and receive_id = ?"
	return conn.Exec(sqlStr, types.HadRead, roomId, userId)
}

//获取聊天消息 startLogId 0:从最新一条消息开始 大于0:从startLogId开始
func GetChatlog(roomId string, startLogId int64, number int) ([]map[string]string, string, error) {
	var sqlStr string
	var maps []map[string]string
	var err error
	if startLogId == 0 {
		sqlStr = "select * from room_msg_content left join `user` on sender_id = `user`.user_id where room_id = ? order by id desc limit 0,?"
		maps, err = conn.Query(sqlStr, roomId, number+1)
	} else {
		sqlStr = "select * from room_msg_content left join `user` on sender_id = `user`.user_id where room_id = ? and id <= ? order by id desc limit 0,?"
		maps, err = conn.Query(sqlStr, roomId, startLogId, number+1)
	}

	var nextLogId string
	if len(maps) > number {
		nextLogId = utility.ToString(maps[len(maps)-1]["id"])
	}

	return maps, nextLogId, err
}

//获取未读消息统计
func GetRoomsUnreadNumber(userId string) ([]map[string]string, error) {
	sqlStr := "SELECT room_id,COUNT(*) as count FROM `room_msg_receive` left join `room_msg_content` on room_msg_id = room_msg_content.id where receive_id = ? and state = 2 group by room_id"
	return conn.Query(sqlStr, userId)
}

// 创建房间 返回 roomId
func CreateRoom(creater, roomName, roomAvatar string, canAddFriend, joinPermission, adminMuted, masterMuted int, members []string, randomRoomId string, createTime int64) (int, error) {
	tx, err := conn.NewTx()
	if err != nil {
		return 0, err
	}
	const insertRoomSql = "insert into room values(?,?,?,?,?,?,?,?,?,?,?)"
	_, _, err = tx.Exec(insertRoomSql, nil, randomRoomId, roomName, roomAvatar, creater, createTime, canAddFriend, joinPermission, adminMuted, masterMuted, types.RoomNotDeleted)
	if err != nil {
		tx.RollBack()
		return 0, err
	}
	maps, err := tx.Query("select LAST_INSERT_ID() as id")

	if err != nil || len(maps) < 1 {
		tx.RollBack()
		return 0, err
	}
	roomId := utility.ToInt(maps[0]["id"])

	//
	const insertMemberSql = "insert into room_user values(?,?,?,?,?,?,?,?,?,?)"
	_, _, err = tx.Exec(insertMemberSql, nil, roomId, creater, "", types.RoomLevelMaster, types.RoomNodisturbingOff, types.RoomUncommonUse, types.RoomNotOnTop, createTime, types.RoomUserNotDeleted)
	if err != nil {
		tx.RollBack()
		return 0, err
	}
	for _, memberId := range members {
		_, _, err = tx.Exec(insertMemberSql, nil, roomId, memberId, "", types.RoomLevelNomal, types.RoomNodisturbingOff, types.RoomUncommonUse, types.RoomNotOnTop, createTime, types.RoomUserNotDeleted)
	}

	return roomId, tx.Commit()
}

// 入群申请，步骤1 添加user
func JoinRoomApproveStepInsert(tx *mysql.MysqlTx, roomId, userId string) error {
	createTime := utility.NowMillionSecond()
	const sqlStr = "insert into room_user values(?,?,?,?,?,?,?,?,?,?)"
	_, _, err := tx.Exec(sqlStr, nil, roomId, userId, "", types.RoomLevelNomal, types.RoomNodisturbingOff, types.RoomUncommonUse, types.RoomNotOnTop, createTime, types.RoomUserNotDeleted)
	if err != nil {
		tx.RollBack()
	}
	return err
}

// 入群申请，步骤2 更改状态
func JoinRoomApproveStepChangeState(tx *mysql.MysqlTx, roomId, userId string, status int) (int64, error) {
	const sqlFindApply = "select id from `apply` where type = ? and apply_user = ? and target = ?"
	maps, err := tx.Query(sqlFindApply, types.RoomType, userId, roomId)
	if err != nil || len(maps) < 1 {
		tx.RollBack()
		return 0, err
	}
	logId := utility.ToInt64(maps[0]["id"])

	const sqlStr = "update `apply` set `state` = ? where id = ?"
	_, _, err = tx.Exec(sqlStr, status, logId)
	if err != nil {
		tx.RollBack()
		return 0, err
	}
	tx.Commit()
	return logId, nil
}

// 群中添加成员
func RoomAddMember(userId, roomId string, createTime int64) (int64, int64, error) {
	const sqlStr = "insert into room_user values(?,?,?,?,?,?,?,?,?,?)"
	return conn.Exec(sqlStr, nil, roomId, userId, "", types.RoomLevelNomal, types.RoomNodisturbingOff, types.RoomUncommonUse, types.RoomNotOnTop, createTime, types.RoomUserNotDeleted)
}

// 修改是否可添加好友
func AlterRoomCanAddFriendPermission(roomId string, permisson int) error {
	if permisson != 0 {
		const sqlStr = "update room set can_add_friend = ? where id = ?"
		_, _, err := conn.Exec(sqlStr, permisson, roomId)
		if err != nil {
			return err
		}
	}
	return nil
}

// 修改加入群权限
func AlterRoomJoinPermission(roomId string, permisson int) error {
	if permisson != 0 {
		const sqlStr = "update room set join_permission = ? where id = ?"
		_, _, err := conn.Exec(sqlStr, permisson, roomId)
		if err != nil {
			return err
		}
	}
	return nil
}

// 设置群成员等级
func SetRoomMemberLevel(userId, roomId string, level int) (int64, int64, error) {
	const sqlStr = "update room_user set level = ? where user_id = ? and room_id = ?"
	return conn.Exec(sqlStr, level, userId, roomId)
}

// 设置群免打扰
func SetRoomNoDisturbing(userId, roomId string, noDisturbing int) (int64, int64, error) {
	const sqlStr = "update room_user set no_disturbing = ? where user_id = ? and room_id = ?"
	return conn.Exec(sqlStr, noDisturbing, userId, roomId)
}

// 设置群置顶
func SetRoomOnTop(userId, roomId string, onTop int) (int64, int64, error) {
	const sqlStr = "update room_user set room_top = ? where user_id = ? and room_id = ?"
	return conn.Exec(sqlStr, onTop, userId, roomId)
}

// 群成员设置昵称
func SetMemberNickname(userId, roomId string, nickname string) (int64, int64, error) {
	const sqlStr = "update room_user set user_nickname = ? where user_id = ? and room_id = ?"
	return conn.Exec(sqlStr, nickname, userId, roomId)
}

// 转让群主
func SetNewMaster(master, userId, roomId string, level int) error {
	tx, err := conn.NewTx()
	if err != nil {
		return err
	}
	const setNewMasterSql = "update room_user set level = ? where user_id = ? and room_id = ?"
	_, _, err = tx.Exec(setNewMasterSql, level, userId, roomId)
	if err != nil {
		tx.RollBack()
		return err
	}
	const removeOldMasterSql = "update room_user set level = ? where user_id = ? and room_id = ?"
	_, _, err = tx.Exec(removeOldMasterSql, types.RoomLevelNomal, master, roomId)
	if err != nil {
		tx.RollBack()
		return err
	}
	return tx.Commit()
}
