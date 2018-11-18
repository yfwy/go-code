package utility

import (
	"fmt"
	"testing"
)

func TestGetAuthToken(t *testing.T) {
	fmt.Println(getAuthToken())
}

func TestUUID(t *testing.T) {
	fmt.Println(TimeUUID())
}

//个推
func TestGTPushSingle(t *testing.T) {
	cid := "155a15456252aed09e0036f15ba0d847"
	err := GTPushSingle("111请求加好友", "系统消息", "有一个好友请求", cid, true, true, 9999999)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("ok")
}

//群推先储存消息
func TestGTPushGroupSaveMsg(t *testing.T) {
	taskId, err := GTPushGroupSaveMsg("这是群推消息", "系统消息", "有一个好友请求", true, true, 9999999)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(taskId) //RASL_1022_255f3b226de847ed9e9f994bc3c60de1
}

//群推
func TestGTPushGroup(t *testing.T) {
	cid := []string{"155a15456252aed09e0036f15ba0d847"}
	taskId, err := GTPushGroupSaveMsg("这是群推消息", "系统消息", "有一个好友请求", true, true, 9999999)
	if err != nil {
		fmt.Println(err)
	}
	//taskId := "RASL_1022_0604cdda25e34878b81445463ce94d50"
	err = GTPushGroup(taskId, cid)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("ok")
	}
}
