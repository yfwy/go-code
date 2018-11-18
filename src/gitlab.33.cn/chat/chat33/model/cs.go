package model

import (
	"fmt"
	"strconv"

	cmn "dev.33.cn/33/common"
	"github.com/astaxie/beego/orm"
	"github.com/inconshreveable/log15"
	dbtools "gitlab.33.cn/chat/chat33/db"
	"gitlab.33.cn/chat/chat33/result"
	"gitlab.33.cn/chat/chat33/types"
	"gitlab.33.cn/chat/chat33/utility"
)

var logCS = log15.New("module", "model/user")

// ================== 客服接口 ==============

type CSStatistics struct {
	AppId string `json:"app_id"`
	CSNum int    `json:"cs_num"`
}

type CSStatisticsRet struct {
	RltList []CSStatistics `json:"rlt_list"`
}

func CustomServiceStatistics() *result.Error {
	var retList []CSStatistics
	appList, _ := dbtools.GetAllAppId()
	for _, app := range appList {
		num, err := dbtools.CountCSNumByAppId(app)
		if err != nil {
			return &result.Error{ErrorCode: result.ParamsError, Message: ""}
		}

		retList = append(retList, CSStatistics{app, num})
	}

	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

type CSInfo struct {
	CsUid          string `json:"cs_uid"`
	CsId           string `json:"cs_id"`
	CsName         string `json:"cs_name"`
	CreateTime     int64  `json:"create_time"`
	PermissionList []int  `json:"permission_list"`
}

type CSList struct {
	CSList []CSInfo `json:"cs_list"`
}

/*
	根据用户权限值生成用户权限列表
*/
func GetCsPermissionList(permission int) []int {
	var ret = make([]int, 0)
	for i := uint(1); i < MaxPermission; i++ {
		if permission&(1<<i) != 0 {
			ret = append(ret, int(i))
		}
	}
	logCS.Info("Get permission", "permission", permission, "pList", ret)
	return ret
}

/*
	获得客服id列表
*/
func GetCSIdList() (map[string]bool, error) {
	var ret = make(map[string]bool)

	maps, err := dbtools.GetCSIdList()
	if err != nil {
		logCS.Error("GetCSList: Select cs failed", "err_msg", err)
		return ret, err
	}

	for i := 0; i < len(maps); i++ {
		ret[maps[i]["cs_id"]] = true
	}
	return ret, nil
}

/*
	客服信息列表
*/
func CSInfoList(params map[string]interface{}) (interface{}, *result.Error) {
	appIDRaw := params["app_id"]
	appID, ok := appIDRaw.(string)
	if !ok {
		return nil, &result.Error{ErrorCode: result.ParamsError, Message: ""}
	}

	startRaw, sOk := params["start_time"]
	var startTime, endTime int64
	if !sOk || cmn.ToString(startRaw) == "" {
		sOk = false
	}
	if sOk {
		startTime = utility.ToInt64(startRaw)
	}

	endRaw, eOk := params["end_time"]
	if !eOk || cmn.ToString(endRaw) == "" {
		eOk = false
	}
	if eOk {
		endTime = utility.ToInt64(endRaw)
	}

	logUser.Debug("CSInfo List", "sOK", sOk, "startRaw", startRaw)
	logUser.Debug("CSInfo List", "eOK", eOk, "endRaw", endRaw)
	//logUser.Info("CSInfo List", "startTime", startTime, "endTime", endTime)
	o := orm.NewOrm()
	var maps []orm.Params
	var err error
	var qStr = "SELECT uid, cs_id, username, A.create_time, A.permission FROM `custom_service` as A LEFT JOIN `user` on cs_id = user_id WHERE A.app_id = ? and `delete` = 0"

	csName := utility.ToString(params["cs_name"])
	if csName != "" {
		qStr += fmt.Sprintf(` AND username LIKE '%%%v%%'`, csName)
	}

	csUid := utility.ToString(params["cs_uid"])
	if csUid != "" {
		qStr += fmt.Sprintf(` AND uid LIKE '%%%v%%'`, csUid)
	}

	if !sOk && !eOk {
		_, err = o.Raw(qStr, appID).Values(&maps)
	} else if sOk && !eOk {
		_, err = o.Raw(qStr+" and A.create_time > ?", appID, startTime).Values(&maps)
	} else if !sOk && eOk {
		_, err = o.Raw(qStr+" and A.create_time < ?", appID, endTime).Values(&maps)
	} else {
		_, err = o.Raw(qStr+" and A.create_time > ? and A.create_time < ?", appID, startTime, endTime).Values(&maps)
	}
	if err != nil {
		return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}
	csList := make([]CSInfo, 0)
	for _, v := range maps {
		permission, ok := v["permission"]
		if !ok || permission == nil {
			permission = "0"
		}
		i, _ := strconv.Atoi(permission.(string))
		csList = append(csList, CSInfo{v["uid"].(string), v["cs_id"].(string), v["username"].(string), utility.ToInt64(v["create_time"]), GetCsPermissionList(i)})
	}

	return CSList{csList}, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

func AddCS(appID, csUID, csName string, pList []interface{}) *result.Error {
	o := orm.NewOrm()
	var maps []orm.Params
	// TODO
	num, err := o.Raw("select user_id from `user` where app_id = ? and uid = ?", ZhaoBi, csUID).Values(&maps)
	if err != nil {
		logCS.Error("AddCS", "msg", err.Error())
		return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}
	if num < 1 {
		return &result.Error{ErrorCode: result.UserNotExists, Message: ""}
	}

	csID := maps[0]["user_id"]

	num, err = o.Raw("select cs_id, `delete` from custom_service where cs_id = ? and app_id=?", csID.(string), appID).Values(&maps)
	if err != nil {
		logCS.Error("AddCS", "msg", err.Error())
		return &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}
	if num > 0 && utility.ToInt(maps[0]["delete"]) == 0 {
		return &result.Error{ErrorCode: result.UserExists, Message: ""}
	}

	var permission = 0
	for _, i := range pList {
		p := utility.ToInt(i)
		permission = permission | (1 << uint(p))
	}
	dbtools.RWMutex.Lock()
	defer dbtools.RWMutex.Unlock()
	if num < 1 {
		_, err = o.Raw("insert into custom_service (cs_id, app_id, permission, create_time, `delete`) "+
			"values(?, ?, ?, ?, 0)", csID, appID, permission, utility.NowMillionSecond()).Exec()
		if err != nil {
			logCS.Error("AddCS", "msg", err.Error())
			return &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
		}
		_, err = o.Raw("update `user` set username = ? where user_id = ?", csName, csID).Exec()
		if err != nil {
			logCS.Error("AddCS", "msg", err.Error())
			return &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
		}
	} else {
		_, err = o.Raw("update custom_service set create_time = ?, `delete` = ? where cs_id = ? and app_id=?",
			utility.NowMillionSecond(), 0, csID, appID).Exec()
		if err != nil {
			logCS.Error("AddCS", "msg", err.Error())
			return &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
		}
		_, err = o.Raw("update `user` set username = ? where user_id = ?", csName, csID).Exec()
		if err != nil {
			logCS.Error("AddCS", "msg", err.Error())
			return &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
		}
	}
	admin, _ := CheckAdmin(csID.(string))
	if !admin {
		_, err = o.Raw("update `user` set user_level=2 where user_id=?", csID).Exec()
	}
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

/*
	删出客服（仅标记）
*/
func RemoveCustomService(csId, appID string) *result.Error {
	dbtools.RWMutex.Lock()
	defer dbtools.RWMutex.Unlock()

	err := dbtools.RemoveCS(csId, appID)

	if err != nil {
		logCS.Error("Remove CS failed", "err_msg", err)
		return &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
	}

	o := orm.NewOrm()
	var maps []orm.Params

	num, err := o.Raw("select cs_id from custom_service where cs_id = ? and `delete` = 0", csId).Values(&maps)
	if err != nil {
		logCS.Error("RM cs select CS failed", "err_msg", err)
	}
	admin, _ := CheckAdmin(csId)
	if num < 1 && !admin {
		_, err = o.Raw("update `user` set user_level = 1 where user_id = ?", csId).Exec()
		if err != nil {
			logCS.Error("RM cs update user_level failed", "err_msg", err)
		}
	}
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

/*
	编辑客服名称
*/
func EditCSName(csID, appID, csName string) *result.Error {
	err := dbtools.UpdateCSName(csID, csName)
	if err != nil {
		logCS.Error("EditCSName update db", "err_msg", err)
		return &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
	}
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

/*
	编辑客服权限
*/
func EditCSPermission(appID, csID string, pList []interface{}) *result.Error {
	err := dbtools.UpdateCSPermission(appID, csID, pList)
	if err != nil {
		logCS.Error("AddCS", "msg", err)
		return &result.Error{ErrorCode: result.WriteDbFailed, Message: ""}
	}
	return &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

type OperatorLog struct {
	TarId   string `json:"tar_id"`
	TarUid  string `json:"tar_uid"`
	Account string `json:"account"`
	Content string `json:"content"`
	Time    int64  `json:"time"`
}

type OperatorLogList struct {
	Totalnum    int           `json:"totalnum"`
	OperateList []OperatorLog `json:"operate_list"`
}

/*
	客服操作记录
*/
func CSOperateLog(csID, appID string, page, number int) (interface{}, *result.Error) {
	if page < 0 || number < 0 {
		return nil, &result.Error{ErrorCode: result.ParamsError, Message: "参数值为负"}
	}

	num, err := dbtools.CountCSOperateLog(csID, appID)
	if err != nil {
		logCS.Error("CSOperateLog query db failed", "msg", err)
		return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}

	maps, err := dbtools.GetCSOperateLog(csID, appID, page*number, number)
	if err != nil {
		logCS.Error("CSOperateLog query db failed", "msg", err)
		return nil, &result.Error{ErrorCode: result.QueryDbFailed, Message: ""}
	}

	optList := make([]OperatorLog, 0)
	for _, v := range maps {
		optList = append(optList, OperatorLog{v["user_id"], v["uid"], v["account"],
			v["content"], utility.ToInt64(v["operate_time"])})
	}

	return OperatorLogList{num, optList}, &result.Error{ErrorCode: result.CodeOK, Message: ""}
}

// 根据id和app_id查询客服在某一app中是否有某一权限
func HasPermission(id, appID string, permission int) (bool, error) {
	userLevel, err := dbtools.GetUserLevel(id)
	if err != nil || userLevel == 0 {
		return false, err
	}

	/*if userLevel == types.LevelAdmin {
		return true, nil
	}*/
	if userLevel == types.LevelMember {
		return false, nil
	}

	pers, err := dbtools.GetCSPermission(id, appID)
	if err != nil {
		return false, err
	}

	if pers&(1<<uint(permission)) != 0 {
		return true, nil
	}
	return false, nil
}

func CheckAdmin(id string) (bool, error) {
	userLevel, err := dbtools.GetUserLevel(id)
	if err != nil {
		logUser.Error("CheckAdmin", "query userLevel failed", err)
		return false, err
	}

	return userLevel == types.LevelAdmin, nil
}
