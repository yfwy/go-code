package utility

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
	"unsafe"

	cmn "dev.33.cn/33/common"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	l "github.com/inconshreveable/log15"
	"github.com/satori/go.uuid"
)

const (
	SESSION_LOGIN = "session-login"
)

//var SessionStore = sessions.NewCookieStore([]byte("something-very-secret"))

func GetUserIdFromSession(context *gin.Context) (string, error) {
	/*session, err := SessionStore.Get(r, SESSION_LOGIN)
	if err != nil {
		return "", err
	}
	val, ok := session.Values["user_id"]
	if !ok {
		return "", fmt.Errorf("not find user_id key")
	}*/
	session := sessions.Default(context)
	val := session.Get("user_id")
	if val == nil {
		return "", fmt.Errorf("not find user_id key")
	}
	return cmn.ToString(val), nil
}

var Usermap = make(map[string]DeviceClientMap)

//key:userId value:详细对接信息
var DockingMap = make(map[string]DockingInfo)

type DockingInfo interface {
	Init()
	TimerReset()
	GetLastCs() string
}

//key:聊天室id value：用户列表
var GroupList = make(map[string]UserList)

//key:用户Id value:用户连接实例
type UserList map[string]DeviceClientMap

type DeviceClientMap map[string]Client

type Client interface {
	GetDevice() string
	GetIsVisitor() bool
	IsManager() bool
	GetID() string
	ClearGroupId() bool
	WriteMsg([]byte) error
	GetGroupID() string
	GetAppID() string
	DoBroadcast([]byte, ErrorBase) ([]byte, bool)
	GetLoginTime() int64
	GetId() string
}

var ManagerList = make(map[string]Client)

type ErrorBase interface {
	Get(...interface{}) interface{}
	/**
		event,code,msgid,errmsg
	**/
	ForDoBroadcast(event int, code int, msg_id string, errmsg string) []byte
}

var utility_log = l.New("module", "chat/utility/common")

var loc = local()

func local() *time.Location {
	loc, err := time.LoadLocation("Asia/Chongqing")
	if err != nil {
		utility_log.Warn("LoadLocation err", "err", err)
		panic(err)
	}
	return loc
}

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const RoomRandomId = "0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// generate random string; fast
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func RandStringBytesMaskImprSrc(n int, lib string) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(lib) {
			b[i] = lib[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func RandomRoomId() string {
	return RandStringBytesMaskImprSrc(9, RoomRandomId)
}

func RandomUsername() string {
	return "chat" + RandStringBytesMaskImprSrc(10, letterBytes)
}

func NowMillionSecond() int64 {
	return time.Now().UnixNano() / 1e6
}

/**
	随机生成ID
**/
func RandomID() string {
	_uuid, _ := uuid.NewV4()
	rlt := fmt.Sprintf("%v", _uuid)
	return strings.Replace(rlt, "-", "", -1)
}

func RandInt(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

/*
func TimeStampIntToTimeStr(ts int64) string {
	tm := time.Unix(ts/1000, ts%1000*1000000)
	return tm.Format(TimeLayoutMillionSecond)
}

func TimeStampToTimeStr(ts string) string {
	val, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return ts
	}
	return TimeStampIntToTimeStr(val)
}

const (
	TimeLayoutMillionSecond = "2006-01-02 15:04:05.000"
	TimeLayoutSecond        = "2006-01-02 15:04:05"
)

func TimeStrToTimeStamp(tm string) string {
	if tm == "" {
		return "0"
	}
	ts, err := time.ParseInLocation(TimeLayoutMillionSecond, tm, loc)
	if err != nil {
		ts, err = time.ParseInLocation(TimeLayoutSecond, tm, loc)
		fmt.Println("TimeStrToTimeStamp: Warn time layout format")
	}
	return strconv.FormatInt(ts.UnixNano()/1000000, 10)
}*/

func RFC3339ToTimeStampMillionSecond(rfc string) int64 {
	const layout = "2006-01-02T15:04:05.000Z"
	if rfc == "" {
		return 0
	}
	ts, err := time.Parse(layout, rfc)
	if err != nil {
		l.Error("RFC3339ToTimeStampMillionSecond", "err_msg", err, "rfc", rfc)
		return 0
	}
	return ts.UnixNano() / 1000000
}

func ParseString(format string, args ...interface{}) string {
	if len(args) == 0 {
		return format
	}
	return fmt.Sprintf(format, args...)
}

func ToInt(val interface{}) int {
	return int(ToInt32(val))
}

func ToInt32(o interface{}) int32 {
	if o == nil {
		return 0
	}
	switch t := o.(type) {
	case int:
		return int32(t)
	case int32:
		return t
	case int64:
		return int32(t)
	case float64:
		return int32(t)
	case string:
		if o == "" {
			return 0
		}
		temp, err := strconv.ParseInt(o.(string), 10, 32)
		if err != nil {
			panic(err)
		}
		return int32(temp)
	default:
		panic(reflect.TypeOf(t).String())
	}
}

func ToInt64(val interface{}) int64 {
	if val == nil {
		return 0
	}
	switch val.(type) {
	case int:
		return int64(val.(int))
	case string:
		if val.(string) == "" {
			return 0
		}
		ret, err := strconv.ParseInt(val.(string), 10, 64)
		if err != nil {
			utility_log.Error("func ToInt64 error")
			debug.PrintStack()
			return 0
		}
		return ret
	case float64:
		return int64(val.(float64))
	case int64:
		return val.(int64)
	default:
		//utility_log.Error("func ToInt error unknow type")
		return 0
	}
}

func ToString(val interface{}) string {
	if val == nil {
		return ""
	}
	return fmt.Sprintf("%v", val)
}

func StructToString(val interface{}) string {
	if val == nil {
		return ""
	}

	switch val.(type) {
	case interface{}:
		bytes, err := json.Marshal(val)
		if err != nil {
			return ""
		}
		return *(*string)(unsafe.Pointer(&bytes))
	default:
		return ""
	}
}

func StringToJobj(val interface{}) map[string]interface{} {
	var rlt = make(map[string]interface{})
	switch val.(type) {
	case string:
		err := json.Unmarshal([]byte(val.(string)), &rlt)
		if err != nil {
			return nil
		}
		return rlt
	default:
		return nil
	}
}

func VisitorNameSplit(src string) string {
	if len(src) >= 8 {
		return src[0:8]
	} else {
		return src[0:]
	}
}
