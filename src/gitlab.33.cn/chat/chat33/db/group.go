package db

import (
	"gitlab.33.cn/chat/chat33/utility"
)

// 更新聊天室用户头像
func UpdateGroupAvatar(groupId, avatar string) (int64, int64, error) {
	const sqlStr = "update `group` set avatar = ? where group_id = ?"
	return conn.Exec(sqlStr, avatar, groupId)
}

// 获取聊天室信息列表
func GetGroupList(groupStatus, startTime, endTime int, groupName string) ([]map[string]string, error) {
	var queryStr = "select * from `group` WHERE `status` IN (1, 2)"
	if groupStatus != 3 {
		queryStr += " AND `status` = " + utility.ToString(groupStatus)
	}
	if groupName != "" {
		queryStr += " AND `group_name` LIKE '%" + groupName + "%'"
	}
	if startTime > 0 && startTime < endTime {
		queryStr += " AND `create_time` >= " + utility.ToString(startTime) + " AND `create_time` <= " + utility.ToString(endTime)
	}
	return conn.Query(queryStr)
}

// 获取聊天室详情
func GetGroupInfo(groupId string) ([]map[string]string, error) {
	const sqlStr = "SELECT * FROM `group` WHERE `group`.group_id = ?"
	return conn.Query(sqlStr, groupId)
}

// 获取聊天室聊天记录
func GetGroupChatLog(groupId string, startId, number int) ([]map[string]string, error) {
	if startId > 0 {
		const sqlStr = "select * from chat_log left join `user` on CHAR_LENGTH(chat_log.sender_id)=CHAR_LENGTH(`user`.user_id) AND chat_log.sender_id = `user`.user_id where chat_log.log_type = 1 and chat_log.receive_id = ? and id < ? order by `id` DESC LIMIT ?,?"
		return conn.Query(sqlStr, groupId, startId, 0, number+1)
	} else {
		const sqlStr = "select * from chat_log left join `user` on CHAR_LENGTH(chat_log.sender_id)=CHAR_LENGTH(`user`.user_id) AND chat_log.sender_id = `user`.user_id where chat_log.log_type = 1 and chat_log.receive_id = ? order by `id` DESC LIMIT ?,?"
		return conn.Query(sqlStr, groupId, 0, number+1)
	}
}

// 获取开放的聊天室
func GetEnableGroups() ([]map[string]string, error) {
	const sqlStr = "select group_id from `group` where `status`= 1"
	return conn.Query(sqlStr)
}

// 添加聊天室
func AddGroup(groupName, avatar string) (int64, int64, error) {
	const sqlStr = "insert into `group`(group_id,group_name,`status`,create_time,open_time,close_time,avatar) values(?,?,?,?,?,?,?)"
	return conn.Exec(sqlStr, nil, groupName, "1", utility.NowMillionSecond(), utility.NowMillionSecond(), utility.NowMillionSecond(), avatar)
}

func AlterGroupState(groupId string, state int) (int64, int64, error) {
	const sqlStr = "update `group` set `status`=? where group_id=?"
	return conn.Exec(sqlStr, state, groupId)
}

func AlterGroupName(groupId, groupName string) (int64, int64, error) {
	const sqlStr = "update `group` set group_name =? where group_id =?"
	return conn.Exec(sqlStr, groupName, groupId)
}

// 添加聊天日志
func AppendGroupChatLog(senderId, receiveId, msgType, content, logType string) (int64, int64, error) {
	const sqlStr = "insert into chat_log(id,sender_id,receive_id,msg_type,content,log_type,send_time) values(?,?,?,?,?,?,?)"
	return conn.Exec(sqlStr, nil, senderId, receiveId, msgType, content, logType, utility.NowMillionSecond())
}
