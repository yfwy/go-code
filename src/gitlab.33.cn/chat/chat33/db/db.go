package db

import (
	"fmt"

	"gitlab.33.cn/chat/chat33/types"

	"dev.33.cn/33/common/mysql"
)

var conn *mysql.MysqlConn

func InitDB(cfg *types.Config) {
	c, err := mysql.NewMysqlConn(cfg.Mysql.Host, fmt.Sprintf("%v", cfg.Mysql.Port),
		cfg.Mysql.User, cfg.Mysql.Pwd, cfg.Mysql.Db)
	if err != nil {
		panic(err)
	}
	conn = c
}

func GetNewTx() (*mysql.MysqlTx, error) {
	return conn.NewTx()
}
