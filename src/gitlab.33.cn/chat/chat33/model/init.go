package model

import (
	"fmt"

	cmn "dev.33.cn/33/common"
	"github.com/astaxie/beego/orm"
	"gitlab.33.cn/chat/chat33/db"
	"gitlab.33.cn/chat/chat33/types"
	"gitlab.33.cn/chat/chat33/utility"

	logic "gitlab.33.cn/chat/chat33/router"
)

var cfg *types.Config
var coinList []*types.Coin
var appList []*types.App

func Init(c *types.Config) {
	cfg = c
	// TODO no db operation in this package
	orm.RegisterDriver("mysql", orm.DRMySQL)
	ds := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4", cfg.Mysql.User, cfg.Mysql.Pwd,
		cfg.Mysql.Host, cfg.Mysql.Port, cfg.Mysql.Db)
	orm.RegisterDataBase("default", "mysql", ds, 2, 100)

	loadGroup()
	loadRoom()
	loadCoins()
}

func loadGroup() {
	rlt, err := GetEnableGroups()
	if err != nil {
		panic(err)
	}
	//load group
	for _, v := range rlt {
		channelId := logic.GetGroupRouteById(utility.ToString(v))
		logic.ChannelMap[channelId] = logic.NewChannel(channelId)
	}
}

func loadRoom() {
	rooms, err := GetEnableRoomIds()
	if err != nil {
		panic(err)
	}
	//load room
	for _, v := range rooms {
		channelId := logic.GetRoomRouteById(utility.ToString(v))
		logic.ChannelMap[channelId] = logic.NewChannel(channelId)
	}
}

func loadCoins() {
	coins, err := db.GetAllCoins()
	if err != nil {
		panic(err)
	}

	for _, c := range coins {
		one := &types.Coin{
			CoinId:   cmn.ToInt(c["coin_id"]),
			CoinName: c["coin_name"],
		}
		coinList = append(coinList, one)
	}
}
