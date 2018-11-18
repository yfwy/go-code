package db

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	cmn "dev.33.cn/33/common"
	"github.com/astaxie/beego/orm"
	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"gitlab.33.cn/chat/chat33/types"
	"gitlab.33.cn/chat/chat33/utility"
)

var log = log15.New("module", "db/dblogic")

const (
	G_table      = "`group`"
	G_groupid    = "group_id"
	G_groupname  = "group_name"
	G_appid      = "app_id"
	G_csid       = "cs_id"
	G_status     = "`status`"
	G_createtime = "create_time"
)

// 获得所有appid列表
func GetAllAppId() ([]string, error) {
	o := orm.NewOrm()
	var maps []orm.Params
	ret := make([]string, 0)
	_, err := o.Raw("select app_id from app").Values(&maps)
	if err != nil {
		log.Error("db failed", "location", "GetAllAppId")
		return ret, err
	}
	for _, k := range maps {
		ret = append(ret, k["app_id"].(string))
	}
	return ret, nil
}

// 获取app信息
func GetAppInfo() []map[string]string {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("select * from app").Values(&maps)

	var ret = make([]map[string]string, 0)

	if err != nil {
		log.Error("GetAppInfo query db failed", "err_msg", err)
		return ret
	}
	for _, app := range maps {
		temp := make(map[string]string)
		temp["app_id"] = utility.ToString(app["app_id"])
		temp["app_name"] = utility.ToString(app["app_name"])
		temp["icon_uri"] = utility.ToString(app["icon_uri"])
		ret = append(ret, temp)
	}
	return ret
}

// 根据用户id获取所在appID
func GetAppIDByUserID(userID string) (string, error) {
	o := orm.NewOrm()
	var maps []orm.Params

	num, err := o.Raw("select app_id from `user` where user_id = ? limit 1", userID).Values(&maps)
	if err != nil {
		log.Error("GetAppIDByUserID failed", "err_msg", err)
		return "", err
	}
	if num < 1 {
		var app_id string
		if devMap, ok := utility.Usermap[userID]; ok {
			for _, v := range devMap {
				if app_id = v.GetAppID(); app_id != "" {
					return app_id, nil
				}
			}
			return app_id, nil
		} else {
			log.Warn("User not exists", "user_id", userID)
			return "", errors.New("用户不存在")
		}
	}
	return utility.ToString(maps[0]["app_id"]), nil
}

// 获取用户最后一条私聊消息
func GetLastPrivMsg(id string) (map[string]interface{}, error) {
	ret := make(map[string]interface{})

	o := orm.NewOrm()
	var maps []orm.Params

	sqlStr := "SELECT * FROM `private_chat_log` WHERE `status` in(1,2,3) AND (sender_id = ? OR receive_id =? ) order by id desc limit 1"
	num, err := o.Raw(sqlStr, id, id).Values(&maps)

	if err != nil {
		log.Error("GetLastPrivMsg", "err_msg", err)
		return nil, err
	}

	if num < 1 {
		return ret, nil
	}

	ret["id"] = utility.ToInt(maps[0]["id"])
	ret["sender_id"] = utility.ToString(maps[0]["sender_id"])
	ret["receive_id"] = utility.ToString(maps[0]["receive_id"])
	ret["app_id"] = utility.ToString(maps[0]["app_id"])
	ret["msg_type"] = utility.ToInt(maps[0]["msg_type"])
	ret["content"] = utility.ToString(maps[0]["content"])
	ret["status"] = utility.ToInt(maps[0]["status"])
	ret["send_time"] = utility.ToInt64(maps[0]["send_time"])

	return ret, nil
}

// 获取权限信息
func GetPermissionInfo() []map[string]interface{} {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("select * from cs_permission").Values(&maps)

	var ret = make([]map[string]interface{}, 0)

	if err != nil {
		log.Error("GetPermissionInfo query db failed", "err_msg", err)
		return ret
	}
	for _, perm := range maps {
		temp := make(map[string]interface{})
		temp["permission_id"] = utility.ToInt(perm["id"])
		temp["permission_name"] = utility.ToString(perm["permission_name"])
		ret = append(ret, temp)
	}
	return ret
}

// 计算用户的userLevel信息
func GetUserLevelWithAppID(id, appID string) (int, error) {
	o := orm.NewOrm()
	var maps []orm.Params

	num, err := o.Raw("select user_level from `user` where user_id = ? limit 1", id).Values(&maps)

	if err != nil || num < 1 {
		return types.LevelVisitor, err
	}

	userLevel := utility.ToInt(maps[0]["user_level"])
	if userLevel != types.LevelCs {
		return userLevel, nil
	}

	num, err = o.Raw("select permission from custom_service where cs_id = ? and `delete` = 0 and app_id = ? limit 1", id, appID).Values(&maps)

	if err != nil || num < 1 {
		return types.LevelMember, err
	}

	return types.LevelCs, nil
}

// 更新用户信息
func UpdateUserInfo(id, remark, description, picUri, csID string) (int64, error) {
	RWMutex.Lock()
	defer RWMutex.Unlock()

	oldInfo, err := GetUserInfoByID(id)
	if err != nil {
		return 0, err
	}
	oldRemark, ok := oldInfo[0]["remark"]
	if !ok || oldRemark != remark {
		appID, _ := GetAppIDByUserID(id)
		err := AddOperationLog(id, csID, appID, "修改备注为\""+remark+"\"")
		if err != nil {
			log.Error("UpdateUserInfo add opt log failed", "err_msg", err)
		}
		sqlStr := fmt.Sprintf("update user set remark=?, description=?, pic_uri=?, remark_cs=? where user_id=?")
		num, _, err := conn.Exec(sqlStr, remark, description, picUri, csID, id)
		return num, err
	}
	sqlStr := fmt.Sprintf("update user set remark=?, description=?, pic_uri=? where user_id=?")
	num, _, err := conn.Exec(sqlStr, remark, description, picUri, id)

	return num, err
}

func MuteUser(id, csID, appID string, optTime, endTime int64) (int64, error) {
	RWMutex.Lock()
	defer RWMutex.Unlock()

	o := orm.NewOrm()
	sqlStr := fmt.Sprintf("insert into muted (user_id, cs_id, app_id, operate_time, end_time) VALUES (?, ?, ?, ?, ?) on DUPLICATE KEY UPDATE cs_id = ?, operate_time = ?, end_time = ?")

	_, err := o.Raw(sqlStr, id, csID, appID, optTime, endTime, csID, optTime, endTime).Exec()
	return 1, err
}

func KickOutUser(id, csID, appID string, optTime, endTime int64) (int64, error) {
	RWMutex.Lock()
	defer RWMutex.Unlock()

	o := orm.NewOrm()
	sqlStr := fmt.Sprintf("insert into kickout (user_id, cs_id, app_id, operate_time, end_time) VALUES (?, ?, ?, ?, ?) on DUPLICATE KEY UPDATE cs_id = ?, operate_time = ?, end_time = ?")

	_, err := o.Raw(sqlStr, id, csID, appID, optTime, endTime, csID, optTime, endTime).Exec()
	return 1, err
}

/*
	添加客服操作记录
*/
func AddOperationLog(userID, csID, appID, content string) error {
	now := utility.NowMillionSecond()

	o := orm.NewOrm()
	sqlStr := "insert into operation_log (user_id, cs_id, app_id, content, operate_time) values(?, ?, ?, ?, ?)"
	_, err := o.Raw(sqlStr, userID, csID, appID, content, now).Exec()
	return err
}

// 根据用户id获取用户信息
func GetUserInfoByID(id string) ([]map[string]string, error) {
	return conn.Query("select * from user where user_id=?", id)
}

// 获取客服管理的app
func GetManageAppMapByID(id string) (map[string]bool, error) {
	o := orm.NewOrm()
	var maps []orm.Params

	ret := make(map[string]bool)

	_, err := o.Raw("select app_id from custom_service where cs_id = ? and `delete` = 0", id).Values(&maps)
	if err != nil {
		return ret, err
	}
	for i := range maps {
		ret[utility.ToString(maps[i]["app_id"])] = true
	}
	return ret, nil
}

func GetCsNum() (int, error) {
	rows, err := conn.Query("select count(*) as cnt from custom_service WHERE `delete` = 0")
	if err != nil {
		return 0, err
	}
	return cmn.ToInt(rows[0]["cnt"]), nil
}

// GetCSInfoByID 根据客服id查询客服详情
func GetCSInfoByID(id, appId string) ([]map[string]string, error) {
	return conn.Query("select * from custom_service where cs_id=? and app_id = ? and `delete` = 0", id, appId)
}

// GetMutedInfo 根据用户Id和时间筛选获取用户禁言信息
func GetMutedInfo(id string, timestamp int64) ([]map[string]string, error) {

	return conn.Query("SELECT * FROM (SELECT * FROM muted WHERE `user_id` =? "+
		"ORDER BY `operate_time` DESC LIMIT 1) a WHERE `end_time` >? ",
		id, timestamp)
}

// GetKickOutInfo 根据用户Id和时间获取用户移除信息
func GetKickOutInfo(id string, timestamp int64) ([]map[string]string, error) {
	return conn.Query("SELECT * FROM (SELECT * FROM kickout WHERE `user_id` =? "+
		"ORDER BY `operate_time` DESC LIMIT 1) a WHERE `end_time` >? ",
		id, timestamp)
}

// 柑橘appid统计该app内的客服数量
func CountCSNumByAppId(appId string) (int, error) {
	ret, err := conn.Query("select count(*) from custom_service where app_id=? and `delete`=0", appId)

	if err != nil {
		return 0, err
	}
	if len(ret) == 0 {
		return 0, errors.New("no such appId: " + appId)
	}

	num := ret[0]["count(*)"]
	i, _ := strconv.Atoi(num)
	return i, nil
}

// 根据客服id查询客服权限
func GetCSPermissionByIDAndAppID(id, appID string) (int, error) {
	ret, err := conn.Query("select permission from custom_service where cs_id=? and app_id=? and `delete` = 0", id, appID)
	if err != nil {
		return 0, err
	}

	if len(ret) == 0 {
		return 0, errors.New("no such cs: " + id)
	}

	permission, ok := ret[0]["permission"]
	if !ok || permission == "" {
		return 0, errors.New("permission is null")
	}
	i, _ := strconv.Atoi(permission)
	return i, nil
}

//-----------------------------------group-----------------------------------//
/**
*	功能：根据群组id查询群组信息
*	参数：id 群组id permission 权限类型 0全部
* 	返回：bool true 是客服 false 不是客服
**/
func GetGroupInfoByid(id string) ([]map[string]string, error) {
	return conn.Query("select * from `group` where group_id=? ", id)
}

var RWMutex sync.RWMutex

/**
*	功能：插入新的群组信息
*	参数：
* 	返回：
**/
func InsertNewGroup(appId, groupName, status, cs_id, createtime int64) (int64, error) {
	RWMutex.Lock()
	sqlstr := fmt.Sprintf("insert into %s(%s,%s,%s,%s,%s,%s) values(?,?,?,?,?,?)", G_table, G_groupid, G_groupname, G_appid, G_status, G_csid, G_createtime)
	num, _, err := conn.Exec(sqlstr, nil, groupName, appId, status, cs_id, createtime)
	RWMutex.Unlock()
	return num, err
}

/**
*	功能：编辑群组信息
*	参数：
* 	返回：
**/
func SetGroupInfo(group_id, field, value string) (int64, error) {
	RWMutex.Lock()
	sqlstr := fmt.Sprintf("update %s set %s=? where %s=?", G_table, field, G_groupid)
	num, _, err := conn.Exec(sqlstr, value, group_id)
	RWMutex.Unlock()
	return num, err
}

func InsertPacket(packetID, userID, toID, tType, size, amount, remark string, cType, coin int, time int64) error {
	const sqlStr = "insert into red_packet_log(packet_id,ctype,user_id,to_id,coin,size,amount,remark,created_at,type) values(?,?,?,?,?,?,?,?,?,?)"
	_, _, err := conn.Exec(sqlStr, packetID, cType, userID, toID, coin, size, amount, remark, time, tType)
	return err
}

func GetRedPacket(packetId string) ([]map[string]string, error) {
	return conn.Query(`SELECT r.packet_id, r.uid as user_id, r.type, r.coin, r.size, r.amount, r.remark, r.created_at, u.username, 
		u.avatar, u.uid FROM red_packet_log as r, user as u WHERE r.packet_id=? AND r.uid=u.user_id`, packetId)
}

func GetAppsPackets(appId string) ([]map[string]string, error) {
	return conn.Query(`SELECT * from red_packet_log WHERE app_id=?`, appId)
}

func GetTodayPackets() ([]map[string]string, error) {
	now := time.Now()
	totayTimeStart := cmn.MillionSecond(cmn.BeginSecOfDay(now))
	totayTimeEnd := cmn.MillionSecond(cmn.EndSecOfDay(now))
	return conn.Query(`SELECT * from red_packet_log WHERE created_at >= ? AND created_at <= ?`, totayTimeStart, totayTimeEnd)
}

func GetAppsPacketsWithFilter(filter *types.PacketQueryParam) (*types.RedPacketInfoList, error) {
	log.Debug("GetAppsPacketsWithFilter", "filter", filter)

	cntStr := fmt.Sprintf(`SELECT count(*) as cnt `)
	queryStr := fmt.Sprintf(`SELECT * `)
	sqlStr := fmt.Sprintf(" FROM red_packet_log as r left join `user` as u on r.uid = u.user_id WHERE r.app_id=? ")

	args := []interface{}{}
	args = append(args, filter.AppId)

	if filter.PacketId != "" {
		sqlStr += fmt.Sprintf(` AND packet_id LIKE '%%%v%%'`, filter.PacketId)
	} else {
		if filter.PacketType != 0 {
			sqlStr += fmt.Sprintf(` AND type=?`)
			args = append(args, filter.PacketType)
		}

		if filter.Coin != 0 {
			sqlStr += fmt.Sprintf(` AND coin=?`)
			args = append(args, filter.Coin)
		}

		if filter.Uid != "" {
			sqlStr += fmt.Sprintf(` AND u.uid like '%%%v%%'`, filter.Uid)
			//args = append(args, filter.Uid)
		}

		if filter.StartTime != 0 {
			sqlStr += fmt.Sprintf(` AND created_at>=? AND created_at<=?`)
			args = append(args, filter.StartTime, filter.EndTime)
		}
	}

	log.Debug("query", "str", cntStr+sqlStr, "args", args)
	rows, err := conn.Query(cntStr+sqlStr, args...)
	if err != nil {
		return nil, err
	}

	items := []*types.RedPacketItem{}

	cnt := cmn.ToInt(rows[0]["cnt"])
	if cnt == 0 {
		return &types.RedPacketInfoList{
			TotalNum: 0,
			Packets:  items,
		}, nil
	}
	log.Debug("query", "rows", cnt)

	lmtStr := fmt.Sprintf(` ORDER BY created_at DESC LIMIT ?,?`)
	args = append(args, (filter.Page-1)*filter.Number)
	args = append(args, filter.Number)

	log.Debug("query2", "str", queryStr+sqlStr+lmtStr, "args", args)
	rows, err = conn.Query(queryStr+sqlStr+lmtStr, args...)
	if err != nil {
		return nil, err
	}
	log.Debug("query2", "rows", len(rows))

	for _, row := range rows {
		one := &types.RedPacketItem{}
		one.PacketId = cmn.ToString(row["packet_id"])
		one.PacketType = cmn.ToInt(row["type"])
		one.SendUid = cmn.ToString(row["user_id"])
		one.AppUid = cmn.ToString(row["uid"])
		one.Coin = cmn.ToInt(row["coin"])
		one.Amount = cmn.ToInt(row["amount"])
		one.Size = cmn.ToInt(row["size"])
		one.Time = cmn.ToInt64(row["created_at"])
		items = append(items, one)
	}

	ret := &types.RedPacketInfoList{
		TotalNum: cnt,
		Packets:  items,
	}
	return ret, nil
}

func GetAllCoins() ([]map[string]string, error) {
	return conn.Query("select * from coin")
}

func GetAllApps() ([]map[string]string, error) {
	return conn.Query("select * from app")
}

// 获取用户禁言和移除详情
func GetMuteAndRemoveInfo(id string) (map[string]int64, error) {
	ret := make(map[string]int64)
	o := orm.NewOrm()
	var maps []orm.Params
	num, err := o.Raw("select user_id, end_time, operate_time from muted where user_id=?", id).Values(&maps)
	if err != nil {
		log.Error("UserInfo select from muted failed", "err_msg", err)
		return ret, err
	}
	now := utility.NowMillionSecond()

	var mutedTime, mutedLastTime int64
	if num < 1 {
		mutedTime = 0
		mutedLastTime = 0
	} else {
		optTime := maps[0]["operate_time"]
		endTime := maps[0]["end_time"]
		optTimeStamp := utility.ToInt64(optTime)
		endTimeStamp := utility.ToInt64(endTime)
		if endTimeStamp <= now {
			mutedLastTime = 0
			mutedTime = 0
		} else {
			mutedTime = optTimeStamp
			mutedLastTime = endTimeStamp - optTimeStamp
		}
	}

	num, err = o.Raw("select user_id, end_time, operate_time from kickout where user_id=?", id).Values(&maps)
	if err != nil {
		log.Error("UserInfo  select from kickout failed", "err_msg", err)
		return ret, err
	}
	var kickoutTime, kickoutLastTime int64
	if num < 1 {
		kickoutTime = 0
		kickoutLastTime = 0
	} else {
		optTime := maps[0]["operate_time"]
		endTime := maps[0]["end_time"]
		optTimeStamp := utility.ToInt64(optTime)
		endTimeStamp := utility.ToInt64(endTime)
		if endTimeStamp <= now {
			kickoutTime = 0
			kickoutLastTime = 0
		} else {
			kickoutTime = optTimeStamp
			kickoutLastTime = endTimeStamp - optTimeStamp
		}
	}

	ret["muted_time"] = mutedTime
	ret["muted_last_time"] = mutedLastTime
	ret["remove_time"] = kickoutTime
	ret["remove_last_time"] = kickoutLastTime
	return ret, nil
}
