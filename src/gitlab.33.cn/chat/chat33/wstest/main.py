import websocket
import json
import requests
import _thread as thread

_url_lists = (
    "127.0.0.1:8090",
    "172.16.103.31:8090"
)
_token_list = (
    "eac5985cd6facc4a697fa4e5df8835ec5b7754a1",
    "0d087c8bf0a35efc8b81869913e78a2f7ae72f2b"
)


class User:

    def __init__(self):
        self.s = requests.session()
        self.cookie = None
        self.token = None
        self.header = None

        print("please input ip addr(default localhost)")
        for i, t in enumerate(_url_lists):
            print(i, t)
        choice = input()
        if choice == "":
            self.url = _url_lists[0]
        else:
            self.url = _url_lists[int(choice)]

        self.http = "http://" + self.url
        self.ws_url = "ws://" + self.url + "/ws"

        print("please choose token you want to use(default visitor)")
        for i, t in enumerate(_token_list):
            print(i, t)
        choice = input()
        if choice != "":
            self.token = _token_list[int(choice)]
            self.header = {
                "FZM-AUTH-TOKEN": self.token,
                "FZM-DEVICE": "py"
            }
            self._login()
        self._ws_connect()

    def _login(self):
        res = self.s.post(self.http + "/user/tokenLogin", headers=self.header)
        js = json.loads(res.content)
        print(js)
        if js['result'] != 0:
            print("登录失败")
            exit(-1)
        self.cookie = "session-login=" + self.s.cookies.get("session-login")
        self.header.update({"Cookie": self.cookie})

    def _ws_connect(self):
        ws = websocket.WebSocketApp(self.ws_url,
                                    on_message=User.on_message,
                                    on_error=User.on_error,
                                    on_close=User.on_close,
                                    header=self.header
                                    )
        ws.on_open = self._on_open()
        ws.run_forever()

    @staticmethod
    def on_message(ws, message):
        print("<<==  " + message)
        pass

    @staticmethod
    def on_error(ws, error):
        print(error)

    @staticmethod
    def on_close(ws):
        print("### closed ###")

    def get_friend_list(self):
        resp = self.s.post(self.http+"/friend/list", json={"type": 0})
        print(json.loads(resp.content))

    def delete_friend(self, _id):
        resp = self.s.post(self.http+"/friend/delete", json={"id": _id})
        print(json.loads(resp.content))

    def add_friend(self, _id, reason):
        resp = self.s.post(self.http+"/friend/add", json={"id": _id, "reason": reason})
        print(json.loads(resp.content))

    def friend_resp(self, _id, agree):
        resp = self.s.post(self.http+"/friend/response", json={"id": _id, "agree": agree})
        print(json.loads(resp.content))

    def get_request_list(self):
        resp = self.s.post(self.http+"/friend/requestList")
        print(json.loads(resp.content))

    def _on_open(self):
        def inner(ws):
            log_in = {
                "eventType": 1,
                "msgId": "123123123",
                "fromId": "",
                "groupId": "15"
            }
            msg = {
                "msgId": "213",
                "channelType": "2",
                "targetId": "16",
                "name": "",
                "userLevel": 1,
                "msgType": 1,
                "msg": {"content": "heeleee"},
                "datetime": ""
            }

            def run(*args):
                print('\n'.join("""enter {id} 进入聊天室
                        switch {id} 切换群
                        list 好友列表
                        add {id} {reason}
                        delete {id}
                        agree {id}
                        disagree {id}
                        a_list 申请列表""".split("\n                        ")))
                while True:
                    m = input()
                    if m.startswith("enter"):
                        room = m.split()[1]
                        log_in.update({"groupId": room})
                        msg.update({"targetId": room, "channelType": "1"})
                        ws.send(json.dumps(log_in))
                    elif m.startswith("switch"):
                        room = m.split()[1]
                        print("switch to " + room)
                        msg.update({"targetId": room, "channelType": "2"})
                    elif m.startswith("list"):
                        self.get_friend_list()
                    elif m.startswith("delete"):
                        p = m.split()
                        self.delete_friend(p[1])
                    elif m.startswith("add"):
                        p = m.split()
                        self.add_friend(p[1], p[2] if len(p) > 2 else "")
                    elif m.startswith("agree"):
                        p = m.split()
                        self.friend_resp(p[1], 1)
                    elif m.startswith("disagree"):
                        p = m.split()
                        self.friend_resp(p[1], 0)
                    elif m.startswith("a_list"):
                        self.get_request_list()
                    else:
                        msg.update({"msg": {"content": "{}".format(m)}})
                        ws.send(json.dumps(msg))
            thread.start_new_thread(run, ())
        return inner


if __name__ == "__main__":
    u = User()
