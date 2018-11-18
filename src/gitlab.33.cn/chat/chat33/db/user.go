package db

import (
	"gitlab.33.cn/chat/chat33/utility"
)

/*
	获取用户UserLevel
*/
func GetUserLevel(id string) (int, error) {
	ret, err := conn.Query("select user_level from `user` where user_id = ? limit 1", id)
	if err != nil || len(ret) < 1 {
		return 0, err
	}
	return utility.ToInt(ret[0]["user_level"]), nil
}

/*
	获取昨日用户数
*/
func GetUserLastDayNumWithAppID(appID string) (int, error) {
	const OneDay = 86400000
	const OffSet = OneDay / 3
	now := utility.NowMillionSecond()
	eTime := now - now%OneDay - OffSet
	stTime := eTime - OneDay

	const sqlStr = "SELECT count(DISTINCT user_id) as `count` from login_log where app_id = ? and login_time >= ? and login_time < ?"
	ret, err := conn.Query(sqlStr, appID, stTime, eTime)
	if err != nil {
		return 0, err
	}
	return utility.ToInt(ret[0]["count"]), nil
}

/*
	获取禁言用户数
*/
func CountMutedNumWithAppID(appID string) (int, error) {
	const sqlStr = "select count(user_id) as `count` from muted where app_id = ? and end_time > ?"
	ret, err := conn.Query(sqlStr, appID, utility.NowMillionSecond())
	if err != nil {
		return 0, err
	}
	return utility.ToInt(ret[0]["count"]), nil
}

/*
	获取移出用户数
*/
func CountKickOutNumWithAppID(appID string) (int, error) {
	const sqlStr = "select count(user_id) as `count` from kickout where app_id = ? and end_time > ?"
	ret, err := conn.Query(sqlStr, appID, utility.NowMillionSecond())
	if err != nil {
		return 0, err
	}
	return utility.ToInt(ret[0]["count"]), nil
}

/*
	获取用户信息
*/
func GetUserInfoWithID(id string) ([]map[string]string, error) {
	const sqlStr = "SELECT user_id, uid, username, account, user_level, avatar from `user` where user_id= ?"
	return conn.Query(sqlStr, id)
}

// 根据被举报用户id获取举报详情
func GetReportInfo(id string) ([]map[string]string, error) {
	sqlStr := "select report.msg_id, datetime, user_id, uid, username, account, msg_type, chat_log.content from report" +
		" left join `user` on `user`.user_id = o_id left join chat_log on report.msg_id = chat_log.id where s_id=?"
	return conn.Query(sqlStr, id)
}

/*
	写入举报信息
*/
func InsertReportInfo(oID, sID, msgID string) error {
	const sqlStr = "INSERT INTO report (o_id, s_id, msg_id, datetime) VALUES (?, ?, ?, ?)"
	_, _, err := conn.Exec(sqlStr, oID, sID, msgID, utility.NowMillionSecond())
	return err
}

/*
	获取用户被举报次数
*/
func CountReportByID(id string) (int, error) {
	const sqlStr = "select count(s_id) from report where s_id=?"
	ret, err := conn.Query(sqlStr, id)
	if err != nil {
		return 0, err
	}
	return utility.ToInt(ret[0]["count(s_id)"]), nil
}

/*
	更新用户头像
*/
func UpdateUserAvatar(id, avatar string) error {
	const sqlStr = "update `user` set avatar = ? where user_id = ?"
	_, _, err := conn.Exec(sqlStr, avatar, id)
	return err
}

// GetUserInfo 获取用户信息 根据uid
func GetUserInfoByUID(uid string) ([]map[string]string, error) {
	const sqlStr = `select * from user where uid=? limit 1`
	return conn.Query(sqlStr, uid)
}

// insert user info
func InsertUser(uid, userName, account, email, phone, userLevel, verified string) (num int64, userId int64, err error) {
	const sqlStr = `INSERT IGNORE INTO user(uid,username,account,email,phone,user_level,verified) VALUES(?,?,?,?,?,?,?)`
	num, userId, err = conn.Exec(sqlStr, uid, userName, account, email, phone, userLevel, verified)
	return
}

/*
	添加用户登录记录
*/
func AddUserLoginLog(userID, deviceType string) error {
	if deviceType == "" {
		deviceType = "Unknown"
	}
	const sqlStr = "insert into login_log (user_id, device, login_time) values(?, ?, ?)"
	_, _, err := conn.Exec(sqlStr, userID, deviceType, utility.NowMillionSecond())
	return err
}

// 更新用户昵称
func UpdateUserName(id, username string) error {
	const sqlStr = "update `user` set username=? where user_id=?"
	_, _, err := conn.Exec(sqlStr, username, id)
	return err
}

/*
	禁言用户-全局
*/
func MuteUserAll(id, operator string, optTime, endTime int64) error {
	const sqlStr = "insert into muted (user_id, operator, operate_time, end_time) VALUES (?, ?, ?, ?) on DUPLICATE KEY UPDATE operator = ?, operate_time = ?, end_time = ?"
	_, _, err := conn.Exec(sqlStr, id, operator, optTime, endTime, operator, optTime, endTime)
	return err
}
