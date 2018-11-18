package comet

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	l "github.com/inconshreveable/log15"
	"gitlab.33.cn/chat/chat33/db"
	"gitlab.33.cn/chat/chat33/model"
	"gitlab.33.cn/chat/chat33/proto"
	"gitlab.33.cn/chat/chat33/result"
	logic "gitlab.33.cn/chat/chat33/router"
	"gitlab.33.cn/chat/chat33/types"
	"gitlab.33.cn/chat/chat33/utility"
)

var wsLog = l.New("module", "chat/comet")

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024 * 7
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// allow cross-origin
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	//
	id string
	//
	user *logic.User
	// The websocket connection.
	conn *websocket.Conn

	device string
	// Buffered channel of outbound messages.
	send chan interface{}
	sync.RWMutex
	closed bool

	loginTime int64

	L sync.Mutex
}

func (c *Client) Binding(user *logic.User) bool {
	c.user = user
	return true
}

func (c *Client) GetSender() chan interface{} {
	return c.send
}

func (c *Client) Send(t interface{}) error {
	c.RLock()
	defer c.RUnlock()
	if c.closed {
		return fmt.Errorf("queue is closed")
	}
	c.send <- t
	return nil
}

func (c *Client) CloseSender() {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return
	}
	close(c.send)
	c.closed = true
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) GetDevice() string {
	return c.device
}

func (c *Client) WriteMsg(msg []byte) error {
	c.L.Lock()
	defer c.L.Unlock()
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	err := c.conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		wsLog.Error("websocket writer throw error", "err", err)
		c.conn.Close()
		return err
	}
	return nil
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.user.Disconnect(c)
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				wsLog.Error("READ ERROR", "error", err)
			}
			wsLog.Info("break readPump nomal break", "client id", c.id, "loginTime", c.loginTime)
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		//Prase data
		eventPrase := proto.NewProto()
		eventErr := eventPrase.Prase(c.user, c.device, message)
		if eventErr != nil {
			ret := result.ComposeWsError(eventErr)
			wsLog.Info("Enjoy Ack", "ack", string(ret[:]))
			c.WriteMsg(ret)
			continue
		}

		switch eventPrase.Event {
		// TODO
		case 0:
			msgPrase, msgErr := proto.NewProtoMsg(eventPrase.Data)
			if msgErr != nil {
				ret := result.ComposeWsError(msgErr)
				c.WriteMsg(ret)
				break
			}
			switch msgPrase.MessageType {
			case proto.SYSTEM:
			}
			if msgPrase.Target == proto.TOGROUP {
				//check permisson before send
				msgErr := msgPrase.CheckGroupPush(c.user.Id, c.device, c.user.Level)
				if msgErr != nil {
					ret := result.ComposeWsError(msgErr)
					c.WriteMsg(ret)
					break
				}
				//compose write data
				proto.FormatTargetMsg(c.user.Id, msgPrase)
				msgErr = proto.AppendGroupChatLog(c.user.Id, msgPrase)
				if msgErr != nil {
					ret := result.ComposeWsError(msgErr)
					c.WriteMsg(ret)
					break
				}
				msgPrase.ComposeTargetData()
			}
			if msgPrase.Target == proto.TOROOM {
				//compose write data
				proto.FormatTargetMsg(c.user.Id, msgPrase)
				msgErr = proto.AppendRoomChatLog(c.user.Id, msgPrase)
				if msgErr != nil {
					ret := result.ComposeWsError(msgErr)
					c.WriteMsg(ret)
					break
				}
				msgPrase.ComposeTargetData()
				//TODO 获取未建立连接用户
				unConnectUsers := model.GetNotStayConnectedRoomUsers(msgPrase.TargetId)
				for range unConnectUsers {
					// TODO 添加群推的cid
					// TODO 添加未读聊天日志
					db.AppendRoomMemberReceiveLog(msgPrase.LogId, c.user.Id, "2")
				}
				//把群推push出去
			}
			if msgPrase.Target == proto.TOUSER {
				//compose write data
				proto.FormatTargetMsg(c.user.Id, msgPrase)
				// Append private chat log ,log state is 2(not send success)
				msgErr = proto.AppendFriendChatLog(c.user.Id, msgPrase, types.NotRead)
				if msgErr != nil {
					ret := result.ComposeWsError(msgErr)
					c.WriteMsg(ret)
					break
				}
				msgPrase.ComposeTargetData()
				// TODO
				msgStr := string(msgPrase.TargetMsgData)
				friendId := msgPrase.GetRouter()
				text := "收到一条好友消息"
				model.SendMsgWithGT(friendId, text, msgStr, msgPrase)
				//将消息发给自己
				err = c.WriteMsg(msgPrase.TargetMsgData)
				if err != nil {
					wsLog.Error("WRITE ERROR", "error", err)
					break
				}
				break
			}

			channelRoute := msgPrase.GetRouter()
			channel := c.user.GetChannel(channelRoute)
			if channel != nil {
				channel.Broadcast(msgPrase)
			}
		case 1:
		}
		wsLog.Info("Get Message", "message", message)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.L.Lock()
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				c.L.Unlock()
				wsLog.Info("break writePump", "client id", c.id, "loginTime", c.loginTime)
				return
			}
			//
			switch message.(type) {
			case *proto.ProtoMsg:
				var proMsg *proto.ProtoMsg
				if proMsg, ok = message.(*proto.ProtoMsg); !ok {
					panic(0)
				}

				switch proMsg.Target {
				case proto.TOGROUP:
					//send message
					err := c.WriteMsg(proMsg.TargetMsgData)
					if err != nil {
						wsLog.Error("WRITE ERROR", "error", err)
						return
					}
				case proto.TOROOM:
					//send message
					err := c.WriteMsg(proMsg.TargetMsgData)
					if err != nil {
						//TODO 单推 消息推送:推给单个cid
						//添加接收日志
						db.AppendRoomMemberReceiveLog(proMsg.LogId, c.user.Id, "2")
						wsLog.Error("WRITE ERROR", "error", err)
						return
					}
					//群成员消息接收日志
					db.AppendRoomMemberReceiveLog(proMsg.LogId, c.user.Id, "1")
				case proto.TOUSER:
					if proMsg.TargetId != c.user.Id {
						continue
					}
					//send message
					err := c.WriteMsg(proMsg.TargetMsgData)
					if err != nil {
						wsLog.Error("WRITE ERROR", "error", err)
						return
					}
					//change the log state
					db.ChangePrivateChatLogStstus(utility.ToInt(proMsg.LogId), types.HadRead)
				}
			default:
				msg, ok := message.([]byte)
				if ok {
					err := c.WriteMsg(msg)
					if err != nil {
						wsLog.Error("WRITE ERROR", "error", err)
						return
					}
				} else {
					fmt.Printf("%T", message)
				}
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			c.L.Lock()
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.L.Unlock()
				return
			}
			c.L.Unlock()
		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(context *gin.Context) {
	userid, clientId, device, uuid, appId, level := GetUserInfo(context)
	wsLog.Debug("[Process] Get Info", "user id", userid, "client id", clientId, "device", device, "uuid", uuid, "app id", appId, "level", level)
	switch level {
	case logic.VISITOR:
		if uuid != "" {
			userid = uuid
		} else {
			userid = "yk_" + utility.RandomID()
		}
		wsLog.Debug("this is visitor", "user_id", userid)
	case logic.NOMALUSER:
		wsLog.Debug("this is user", "user_id", userid)
	case logic.MANAGER:
		wsLog.Debug("this is manager", "user_id", userid)
	}

	conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		wsLog.Error("[Process] Ws Upgrade Faild", "err", err)
		return
	}
	wsLog.Debug("[Process] Ws Upgrade Success")

	client := &Client{id: clientId, device: device, conn: conn, send: make(chan interface{})}
	user, ok := logic.AppendUser(userid, device, client, level)
	if ok {
		//private channel
		cl, ok := logic.ChannelMap["default"]
		if ok {
			user.DeviceRegister(device, cl)
		}
		//room channel
		rooms, err := model.GetUserJoinedRooms(userid)
		if err != nil {
			//out put log
		}
		for _, v := range rooms {
			roomChannelId := logic.GetRoomRouteById(v)
			cl, ok := logic.ChannelMap[roomChannelId]
			if ok {
				user.Subscribe(cl)
			}
		}
	}
	// new goroutines.
	wsLog.Debug("[Process]Start Write/Read Pump")
	go client.writePump()
	go client.readPump()
}

func GetUserInfo(c *gin.Context) (userid, clientId, device, uuid, appId string, level int) {
	device = c.GetHeader("FZM-DEVICE")
	uuid = c.GetHeader("FZM-UUID")
	appId = c.GetHeader("FZM-APP-ID")
	var isManager bool

	//获取当前用户id
	session := sessions.Default(c)

	// TODO device
	wsLog.Info("session values:",
		"user_id", session.Get("user_id"), "ismanager", session.Get("ismanager"), "app_id", session.Get("app_id"), "devtype", session.Get("devtype"),
		"id", session.Get("id"), "app_list", session.Get("app_list"))

	_userId := session.Get("user_id")
	if _userId == nil {
		wsLog.Debug("seesion中user_id为空")
	} else {
		userid = utility.ToString(_userId)
	}

	_manager := session.Get("ismanager")
	if _manager == nil {
		isManager = false
	} else {
		isManager = _manager.(bool)
	}

	_id := session.Get("id")
	clientId = utility.ToString(_id)

	//获取类型
	if device == "" {
		_device := session.Get("devtype")
		device = utility.ToString(_device)
	}
	switch device {
	case "Web":
	case "Android":
	case "iOS":
	default:
		device = ""
	}

	if userid == "" {
		level = logic.VISITOR
	} else if isManager {
		level = logic.MANAGER
	} else {
		level = logic.NOMALUSER
	}
	return
}
