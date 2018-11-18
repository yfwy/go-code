package utility

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

//个推工具包

const (
	appkey       = "OwJM9uqPpt5PPFPS4SnGk7"
	appId        = "le1QnKJrAcAEwvinL2awk3"
	mastersecret = "mz0S3rutJ17gamib9Sib12"
)

var authtoken = ""

//创建参数
//如果是群推  cid就填""
func newPara(text, title, transmissionContent, cid string, isOffline, transmissionType bool, offlineExpireTime int) ([]byte, error) {
	var paras = make(map[string]interface{})
	var message = make(map[string]interface{})
	var notification = make(map[string]interface{})
	var style = make(map[string]interface{})

	message["appkey"] = appkey
	message["msgtype"] = "notification"
	message["is_offline"] = isOffline
	message["offline_expire_time"] = offlineExpireTime

	style["type"] = 0
	style["text01"] = text
	style["title"] = title

	notification["style"] = style
	notification["transmission_type"] = transmissionType
	notification["transmission_content"] = transmissionContent

	paras["message"] = message
	paras["notification"] = notification

	if cid != "" {
		paras["cid"] = cid
		uuid := TimeUUID()
		requestid := ""
		for i, u := range uuid {
			if i < 30 {
				requestid += string(u)
			}
		}
		paras["requestid"] = requestid
	}

	return json.Marshal(paras)
}

//获取AuthToken
func getAuthToken() (string, error) {

	timestamp := NowMillionSecond()
	shaP := appkey + strconv.Itoa(int(timestamp)) + mastersecret
	sha256 := sha256.New()
	_, err := sha256.Write([]byte(shaP))
	if err != nil {
		return "", err
	}
	sign := sha256.Sum(nil)
	var p = make(map[string]interface{})
	p["timestamp"] = timestamp
	p["appkey"] = appkey
	p["sign"] = hex.EncodeToString(sign)
	param, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	url := `https://restapi.getui.com/v1/` + appId + `/auth_sign`
	req, err := http.NewRequest("POST", url, strings.NewReader(string(param)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result = make(map[string]interface{})
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	return result["auth_token"].(string), nil

}

//对使用App的某个用户，单独推送消息

//is_offline 是否离线推送
//offline_expire_time 消息离线存储有效期，单位：ms
//transmissionContent  传透内容
func GTPushSingle(text, title, transmissionContent, cid string, isOffline, transmissionType bool, offlineExpireTime int) error {
	client := &http.Client{}
	paramBytes, err := newPara(text, title, transmissionContent, cid, isOffline, transmissionType, offlineExpireTime)
	if err != nil {
		return err
	}
	parans := string(paramBytes)
	url := `https://restapi.getui.com/v1/` + appId + `/push_single`
	req, err := http.NewRequest("POST", url, strings.NewReader(parans))
	if err != nil {
		return err
	}
	req.Header["authtoken"] = []string{authtoken}
	//req.Header.Add("authtoken", "866adb68ce519b9eeb3e80601941cbe3a2017c020fea26b528f097c20863fd14")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result = make(map[string]interface{})
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	if result["result"] == "ok" {
		return nil
	}
	if result["result"] == "not_auth" {
		//计数，限制更换token次数，避免进入死循环
		i := 0
		authtoken, _ = getAuthToken()
		fmt.Println(authtoken)
		i++
		if i < 5 {
			return GTPushSingle(text, title, transmissionContent, cid, isOffline, transmissionType, offlineExpireTime)
		}
		return errors.New(result["result"].(string))
	} else {
		return errors.New(result["result"].(string))
	}
}

//在执行群推任务的时候，需首先执行save_list_body接口，
// 将推送消息保存在服务器上，后面可以重复调用tolist接口将保存的消息发送给不同的目标用户。

//群推 先储存消息 返回taskid
func GTPushGroupSaveMsg(text, title, transmissionContent string, isOffline, transmissionType bool, offlineExpireTime int) (string, error) {
	client := &http.Client{}
	paramBytes, err := newPara(text, title, transmissionContent, "", isOffline, transmissionType, offlineExpireTime)
	if err != nil {
		return "", err
	}
	parans := string(paramBytes)
	url := `https://restapi.getui.com/v1/` + appId + `/save_list_body`
	req, err := http.NewRequest("POST", url, strings.NewReader(parans))
	if err != nil {
		return "", err
	}
	req.Header["authtoken"] = []string{authtoken}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result = make(map[string]interface{})
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	if result["result"] == "ok" {
		return result["taskid"].(string), nil
	}
	if result["result"] == "not_auth" {
		//计数，限制更换token次数，避免进入死循环
		i := 0
		authtoken, _ = getAuthToken()
		fmt.Println(authtoken)
		i++
		if i < 5 {
			return GTPushGroupSaveMsg(text, title, transmissionContent, isOffline, transmissionType, offlineExpireTime)
		}
		return "", errors.New(result["result"].(string))
	} else {
		return "", errors.New(result["result"].(string))
	}
}

// 群推  真正发送消息
func GTPushGroup(taskId string, cid []string) error {
	client := &http.Client{}
	var paraMap = make(map[string]interface{})
	paraMap["cid"] = cid
	paraMap["taskid"] = taskId
	paraMap["need_detail"] = true

	paraByte, err := json.Marshal(paraMap)
	if err != nil {
		return err
	}
	para := string(paraByte)

	url := `https://restapi.getui.com/v1/` + appId + `/push_list`
	req, err := http.NewRequest("POST", url, strings.NewReader(para))
	if err != nil {
		return err
	}
	req.Header["authtoken"] = []string{authtoken}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result = make(map[string]interface{})
	err = json.Unmarshal(body, &result)
	fmt.Println(result)
	if err != nil {
		return err
	}
	if result["result"] == "ok" {
		return nil
	}
	if result["result"] == "not_auth" {
		//计数，限制更换token次数，避免进入死循环
		i := 0
		authtoken, _ = getAuthToken()
		fmt.Println(authtoken)
		i++
		if i < 5 {
			return GTPushGroup(taskId, cid)
		}
		return errors.New(result["result"].(string))
	} else {
		return errors.New(result["result"].(string))
	}

}
