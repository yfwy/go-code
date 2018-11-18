package db

import (
	"gitlab.33.cn/chat/chat33/utility"
)

// 获取申请列表
func GetApplyList(userId string, id int64, number int, rooms string) ([]map[string]string, error) {
	if id == 0 {
		const sqlStr = "SELECT * FROM `apply` where apply_user = ? or (target = ? and `type` = 2) or (target in(?) and `type` = 1) ORDER BY datetime desc limit 0,?"
		return conn.Query(sqlStr, userId, userId, rooms, number)
	} else {
		const sqlStr = "SELECT * FROM `apply` WHERE apply_user = ? or (target = ? and `type` = 2) or (target in(?) and `type` = 1) and datetime <= (SELECT datetime FROM `apply` WHERE id = ?) ORDER BY datetime desc limit 0,?"
		return conn.Query(sqlStr, userId, userId, rooms, id, number)
	}
}

func GetApplyListNumber() (int, error) {
	const sqlStr = "SELECT COUNT(*) as count FROM `apply`"
	maps, err := conn.Query(sqlStr)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, v := range maps {
		count += utility.ToInt(v["count"])
	}
	return count, err
}
