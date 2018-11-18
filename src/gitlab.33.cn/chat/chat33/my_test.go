package main

import (
	"encoding/json"
	"net/http"
	_ "net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"

	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	_ "time"

	"github.com/gorilla/websocket"
	l "github.com/inconshreveable/log15"
)

var testLog = l.New("module", "chat/my_test")

// func TestFruits(t *testing.T) {

// 	// // create http.Handler
// 	// handler := FruitsHandler()

// 	// // run server using httptest
// 	// server := httptest.NewServer(handler)
// 	// defer server.Close()

// 	// create httpexpect instance
// 	e := httpexpect.New(t, "http://localhost:8080")

// 	// is it working?
// 	e.POST("/hello").
// 		Expect().
// 		Status(http.StatusOK)
// }

// func TestMy(t *testing.T){

// 	e := httpexpect.WithConfig(httpexpect.Config{
// 		Reporter: httpexpect.NewAssertReporter(t),
// 		Client: &http.Client{
// 			Jar: httpexpect.NewJar(), // used by default if Client is nil
// 		},
// 	})

// 	e.POST("/hello").
// 		Expect().
// 		Status(http.StatusOK)
// }

var addr = flag.String("addr", "localhost:8080", "http service address")
var reqBaseUrl = "http://localhost:8080"

var User1Account = "my_test_user1"
var User1Password = "123456"

var User2Account = "my_test_user2"
var User2Password = "123456"

var CustomAccount = "my_test_custom1"
var CustomPassword = "123456"

// //用户注册接口测试
// func TestSignUpUser1(t *testing.T){
// 	//注册
// 	e := httpexpect.New(t, reqBaseUrl)

// 	user := map[string]interface{}{
// 		"account": User1Account,
// 		"password": User1Password,
// 	}

// 	response := e.POST("/app/pwdSignUp").WithHeader("COntent-type","application/json").WithJSON(user).Expect()
// 	// is it working?
// 	obj := response.Status(http.StatusOK).JSON().Object()
// 	//返回 result 0 注册成功
// 	obj.ContainsKey("result").ValueEqual("result", 0)
// }

// //用户注册接口测试
// func TestSignUpUser2(t *testing.T){
// 	//注册
// 	e := httpexpect.New(t, reqBaseUrl)

// 	user := map[string]interface{}{
// 		"account": User2Account,
// 		"password": User2Password,
// 	}

// 	response := e.POST("/app/pwdSignUp").WithHeader("COntent-type","application/json").WithJSON(user).Expect()
// 	// is it working?
// 	obj := response.Status(http.StatusOK).JSON().Object()
// 	//返回 result 0 注册成功
// 	obj.ContainsKey("result").ValueEqual("result", 0)
// }

// //用户注册接口测试
// func TestSignUpCustom1(t *testing.T){
// 	//注册
// 	e := httpexpect.New(t, reqBaseUrl)

// 	user := map[string]interface{}{
// 		"account": CustomAccount,
// 		"password": CustomPassword,
// 	}

// 	response := e.POST("/app/pwdSignUp").WithHeader("COntent-type","application/json").WithJSON(user).Expect()
// 	// is it working?
// 	obj := response.Status(http.StatusOK).JSON().Object()
// 	//返回 result 0 注册成功
// 	obj.ContainsKey("result").ValueEqual("result", 0)
// }

var sessionId1 string
var sessionId2 string
var sessionId1Custom1 string

//用户1 登录接口测试
func TestLoginUser1(t *testing.T) {
	//登录
	e := httpexpect.New(t, reqBaseUrl)

	user := map[string]interface{}{
		"account":  User1Account,
		"password": User1Password,
	}

	response := e.POST("/app/pwdLogin").WithJSON(user).Expect()

	obj := response.Status(http.StatusOK).JSON().Object()
	//返回 result 0 登录成功
	obj.ContainsKey("result").ValueEqual("result", 0)

	cookie := response.Status(http.StatusOK).Cookie("session-login")
	sessionId1 = cookie.Raw().Name + "=" + cookie.Raw().Value
}

//
func TestLoginUser2(t *testing.T) {
	//登录
	e := httpexpect.New(t, reqBaseUrl)

	user := map[string]interface{}{
		"account":  User2Account,
		"password": User2Password,
	}

	response := e.POST("/app/pwdLogin").WithJSON(user).Expect()

	obj := response.Status(http.StatusOK).JSON().Object()
	//返回 result 0 登录成功
	obj.ContainsKey("result").ValueEqual("result", 0)

	cookie := response.Status(http.StatusOK).Cookie("session-login")
	sessionId2 = cookie.Raw().Name + "=" + cookie.Raw().Value
}

func TestLoginCustom1(t *testing.T) {
	//登录
	e := httpexpect.New(t, reqBaseUrl)

	user := map[string]interface{}{
		"account":  CustomAccount,
		"password": CustomPassword,
	}

	response := e.POST("/app/pwdLogin").WithJSON(user).Expect()

	obj := response.Status(http.StatusOK).JSON().Object()
	//返回 result 0 登录成功
	obj.ContainsKey("result").ValueEqual("result", 0)

	cookie := response.Status(http.StatusOK).Cookie("session-login")
	sessionId1Custom1 = cookie.Raw().Name + "=" + cookie.Raw().Value
}

//1、用户1登录聊天室1
//2、用户2登录聊天室1
//2、游客登录聊天室1
//3、客服登录聊天室1

//4、游客发群聊消息	没有发送权限
//5、用户1发群聊消息 用户1、2 客服 游客 收到消息
//6、游客给用户1发送私聊消息 没有发送权限
//7、游客给客服发送私聊消息 客服收到私聊消息
//6、用户1给客服发送私聊消息 客服收到私聊消息
//7、客服给游客发送私聊消息 游客收到私聊消息
//8、用户2切换到聊天群2
//9、客服给用户2发送私聊消息 用户2收到私聊消息
//10、用户2发送群聊消息到群1 没有发送权限
//11、用户2发送群聊消息到群2 发送成功 游客，用户，客服收不到群聊消息

var wsVisitor *websocket.Conn
var wsUser1 *websocket.Conn

func TestWsInit(t *testing.T) {
	//----------------游客-------------//
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	var err error
	wsVisitor, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
		t.Error("visitor connet failes")
	}

	//---------------用户1-------------//
	var reqHeader = make(http.Header)
	reqHeader.Add("Cookie", sessionId1)

	wsUser1, _, err = websocket.DefaultDialer.Dial(u.String(), reqHeader)
	if err != nil {
		log.Fatal("dial:", err)
		t.Error("visitor connet failes")
	}
}

// type compare struct{
// 	Event_type int `json:"event_type"`
// 	Msg_id string `json:"msg_id"`
// 	Code int `json:"code"`
// 	Content string `json:"content"`
// }

func TestVisitorLoginGroup(t *testing.T) {

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan string)

	go func() {
		defer close(done)
		for {
			_, message, err := wsVisitor.ReadMessage()
			if err != nil {
				testLog.Error("读取消息错误", "err", err)
				t.Error()
				return
			}

			log.Printf("recv: %s", message)
			var data = make(map[string]interface{})
			err = json.Unmarshal(message, &data)
			if err != nil {
				testLog.Error("解析错误", "err", err)
				t.Error("登录失败")
			}

			if data["event_type"] == 1 && data["code"] == 0 {
				testLog.Debug("错误的返回结果")
				t.Error("登录失败")
			}

			done <- "done"
			return
		}
	}()

	//加入聊天室
	//群组号为 9
	err := wsVisitor.WriteMessage(websocket.TextMessage, []byte("{\"event_type\": 1,\"msg_id\": \"123123123\",\"from_id\": \"\",\"group_id\": \"9\"}"))
	if err != nil {
		log.Println("write:", err)
		return
	}

	<-done
}

func TestUser1LoginGroup(t *testing.T) {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	done := make(chan string)

	go func() {
		defer close(done)
		for {
			_, message, err := wsUser1.ReadMessage()
			if err != nil {
				testLog.Error("读取消息错误", "err", err)
				t.Error()
				return
			}
			log.Printf("recv: %s", message)
			var data = make(map[string]interface{})
			err = json.Unmarshal(message, &data)
			if err != nil {
				testLog.Error("消息格式错误", "err", err)
				t.Error("消息格式错误")
			}

			if data["event_type"] == 1 && data["code"] == 0 {
				testLog.Debug("错误的返回结果")
				t.Error("登录失败")
			}
			done <- "done"
			return
		}
	}()

	//加入聊天室
	//群组号为 9
	err := wsUser1.WriteMessage(websocket.TextMessage, []byte("{\"event_type\": 1,\"msg_id\": \"123123123\",\"from_id\": \"\",\"group_id\": \"9\"}"))
	if err != nil {
		log.Println("write:", err)
		return
	}

	<-done
}

//游客发群聊消息
func TestVisitorSendGroupMessage(t *testing.T) {

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan string)

	go func() {
		defer close(done)
		for {
			_, message, err := wsVisitor.ReadMessage()
			if err != nil {
				testLog.Error("读取消息错误", "err", err)
				t.Error()
				return
			}

			log.Printf("recv: %s", message)
			var data = make(map[string]interface{})
			err = json.Unmarshal(message, &data)
			if err != nil {
				testLog.Error("解析错误", "err", err)
				t.Error("登录失败")
			}

			if data["event_type"] == 0 && data["code"] == -2008 {
				testLog.Debug("错误的返回结果")
				t.Error("登录失败")
			}

			done <- "done"
			return
		}
	}()

	//向群组9发送数据
	err := wsVisitor.WriteMessage(websocket.TextMessage, []byte("{\"event_type\":0,\"msg_id\":\"123123\",\"from_id\":\"\",\"from_gid\":\"9\",\"to_id\":\"\",\"name\":\"\",\"user_level\":\"\",\"msg_type\":1,\"msg\":{\"content\":\"youke xiaoxi 1\"},\"datetime\":\"\"}"))
	if err != nil {
		log.Println("write:", err)
		return
	}

	<-done
}

//用户1 发群聊消息
func TestUser1SendGroupMessage(t *testing.T) {

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan string)

	go func() {
		defer close(done)
		for {
			_, message, err := wsVisitor.ReadMessage()
			if err != nil {
				testLog.Error("读取消息错误", "err", err)
				t.Error()
				return
			}

			log.Printf("recv: %s", message)
			var data = make(map[string]interface{})
			err = json.Unmarshal(message, &data)
			if err != nil {
				testLog.Error("解析错误", "err", err)
				t.Error("登录失败")
			}

			if data["event_type"] == 0 && data["code"] == -2008 {
				testLog.Debug("错误的返回结果")
				t.Error("登录失败")
			}

			done <- "done"
			return
		}
	}()

	//向群组9发送数据
	err := wsUser1.WriteMessage(websocket.TextMessage, []byte("{\"event_type\":0,\"msg_id\":\"123123\",\"from_id\":\"\",\"from_gid\":\"9\",\"to_id\":\"\",\"name\":\"\",\"user_level\":\"\",\"msg_type\":1,\"msg\":{\"content\":\"youke xiaoxi 1\"},\"datetime\":\"\"}"))
	if err != nil {
		log.Println("write:", err)
		return
	}

	<-done
}

// //游客登录测试
// func TestWebsocket(t *testing.T){

// 	flag.Parse()
// 	log.SetFlags(0)

// 	interrupt := make(chan os.Signal, 1)
// 	signal.Notify(interrupt, os.Interrupt)

// 	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
// 	log.Printf("connecting to %s", u.String())

// 	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
// 	if err != nil {
// 		log.Fatal("dial:", err)
// 	}
// 	defer c.Close()

// 	done := make(chan string)

// 	type compare struct{
// 		Event_type int `json:"event_type"`
// 		Msg_id string `json:"msg_id"`
// 		Code int `json:"code"`
// 		Content string `json:"content"`
// 	}

// 	var _revData = new(compare)
// 	go func() {
// 		defer close(done)
// 		for {
// 			_, message, err := c.ReadMessage()
// 			if err != nil {
// 				testLog.Error("读取消息错误","err",err)
// 				t.Error()
// 				return
// 			}
// 			log.Printf("recv: %s", message)
// 			err = json.Unmarshal(message,_revData)
// 			if err!=nil{
// 				testLog.Error("解析错误","err",err)
// 				t.Error("登录失败")
// 			}
// 			if _revData.Event_type != 1 || _revData.Msg_id != "123123123" || _revData.Content != "操作成功" || _revData.Code != 0{
// 				testLog.Debug("错误的返回结果")
// 				t.Error("登录失败")
// 			}
// 			done <- "done"
// 			return
// 		}
// 	}()

// 	//加入聊天室
// 	//群组号为 9
// 	err = c.WriteMessage(websocket.TextMessage, []byte("{\"event_type\": 1,\"msg_id\": \"123123123\",\"from_id\": \"\",\"group_id\": \"9\"}"))
// 	if err != nil {
// 		log.Println("write:", err)
// 		return
// 	}

// 	for {
// 		select {
// 		case <-done:
// 			return
// 		}
// 	}
// }
