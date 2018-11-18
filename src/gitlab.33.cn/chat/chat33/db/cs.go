package db

import (
	"gitlab.33.cn/chat/chat33/utility"
)

/*
	编辑客服名称
*/
func UpdateCSName(csID, csName string) error {
	const sqlStr = "update `user` set `username`=? where user_id=?"
	_, _, err := conn.Exec(sqlStr, csName, csID)
	return err
}

/*
	获取客服Id列表
*/
func GetCSIdList() ([]map[string]string, error) {
	const sqlStr = "SELECT DISTINCT cs_id FROM `custom_service` WHERE `delete` = 0"
	return conn.Query(sqlStr)
}

/*
	获取app下的客服列表
*/
func GetCSListWithAppID(appID string) ([]map[string]string, error) {
	const sqlStr = "select cs_id from custom_service where app_id = ? and `delete`=0"
	return conn.Query(sqlStr, appID)
}

/*
	删除客服(仅标记)
*/
func RemoveCS(csID, appID string) error {
	const sqlStr = "update custom_service set `delete`=1 where cs_id=? and app_id=?"
	_, _, err := conn.Exec(sqlStr, csID, appID)
	return err
}

// ======================= CS Permission ============================

func GetPermission(userId int) ([]map[string]string, error) {
	const sqlStr = "select * from custom_service where cs_id=? AND `delete` = 0"
	return conn.Query(sqlStr, userId)
}

func GetAdminPermission(userId int) ([]map[string]string, error) {
	const sqlStr = "SELECT app.app_id, permission FROM app LEFT JOIN custom_service on app.app_id = " +
		" custom_service.app_id and cs_id = ? and `delete` = 0"
	return conn.Query(sqlStr, userId)
}

/*
	编辑客服权限
*/
func UpdateCSPermission(appID, csID string, pList []interface{}) error {
	var permission = 0
	for _, i := range pList {
		p := utility.ToInt(i)
		permission = permission | (1 << uint(p))
	}
	const sqlStr = "update custom_service set permission=? where cs_id = ? and app_id = ?"
	_, _, err := conn.Exec(sqlStr, permission, csID, appID)
	return err
}

/*
	获取客服在某一app下的权限
*/
func GetCSPermission(csID, appID string) (int, error) {
	const sqlStr = "select permission from custom_service where cs_id = ? and `delete` = 0 and app_id = ? limit 1"

	ret, err := conn.Query(sqlStr, csID, appID)

	if err != nil || len(ret) < 1 {
		return 0, err
	}

	return utility.ToInt(ret[0]["permission"]), nil
}

// ======================== CS Operation ================================

/*
	获取客服操作记录数
*/
func CountCSOperateLog(csID, appID string) (int, error) {
	const sqlStr = "select count(cs_id) from operation_log where cs_id = ? and app_id=?"
	ret, err := conn.Query(sqlStr, csID, appID)
	if err != nil {
		return 0, err
	}

	num := utility.ToInt(ret[0]["count(cs_id)"])
	return num, nil
}

/*
	获取客服操作记录
*/
func GetCSOperateLog(csID, appID string, start, end int) ([]map[string]string, error) {
	const sqlStr = "SELECT A.user_id, B.uid, B.account, A.content, A.operate_time from operation_log as A LEFT JOIN `user` " +
		"as B on A.user_id=B.user_id WHERE A.cs_id=? and A.app_id=? limit ?, ?"
	return conn.Query(sqlStr, csID, appID, start, end)
}
