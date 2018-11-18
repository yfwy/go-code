package mysql_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"dev.33.cn/33/common/mysql"
)

var host = "localhost"
var port = "32768"
var user = "root"
var pwd = "123"

func TestWithDb(t *testing.T) {
	require := require.New(t)
	conn := testConn(t, "test")
	require.NotNil(conn)

	testQuery(t, conn)

	testInsert(t, conn)
	testQuery(t, conn)

	testUpdate(t, conn)
	testQuery(t, conn)

	testDelete(t, conn)
	testQuery(t, conn)

	testTx(t, conn)

	testConnSelectDb(t, conn, "parse")
	res, err := conn.Query("select * from blockchain")
	require.Nil(err)
	t.Log(res)

	testCloseConn(t, conn)
}

func TestConnNoDb(t *testing.T) {
	require := require.New(t)

	conn := testConnNoDb(t)
	require.NotNil(conn)

	testConnSelectDb(t, conn, "test")
	testQuery(t, conn)
	testCloseConn(t, conn)
}

func TestTxNoDb(t *testing.T) {
	require := require.New(t)

	conn := testConnNoDb(t)
	require.NotNil(conn)

	tx, err := conn.NewTx()
	require.Nil(err)

	_, _, err = tx.Exec("use test")
	require.Nil(err)

	testTxSelectDb(t, tx, "test")

	for i := 2000; i < 3000; i++ {
		tx.Exec("insert into user(uid,name,age) values(?,?,?)", i, fmt.Sprintf("foo_%d", i), i%100+1)
	}

	rows, err := tx.Query("select * from user where uid=2588")
	require.Nil(err)
	require.Equal(1, len(rows))

	num, lastId, err := tx.Exec("delete from user where uid>=2000")
	require.Nil(err)
	t.Log(num, lastId)

	_, _, err = tx.Exec("use parse")
	require.Nil(err)

	num, lastId, err = tx.Exec("update blockchain set update_time=?", "2017-08-08 18:18:28")
	require.Nil(err)
	t.Log(num, lastId)

	err = tx.Commit()
	require.Nil(err)

	testCloseConn(t, conn)
}

func testConn(t *testing.T, dbName string) *mysql.MysqlConn {
	require := require.New(t)
	conn, err := mysql.NewMysqlConn(host, port, user, pwd, dbName)
	require.Nil(err)
	return conn
}

func testConnNoDb(t *testing.T) *mysql.MysqlConn {
	require := require.New(t)
	conn, err := mysql.NewMysqlConn(host, port, user, pwd, "")
	require.Equal(err, nil)
	return conn
}

func testConnSelectDb(t *testing.T, conn *mysql.MysqlConn, dbName string) {
	require := require.New(t)
	_, _, err := conn.Exec("use " + dbName)
	require.Nil(err)
}

func testTxSelectDb(t *testing.T, tx *mysql.MysqlTx, dbName string) {
	require := require.New(t)
	_, _, err := tx.Exec("use " + dbName)
	require.Nil(err)
}

func testQuery(t *testing.T, conn *mysql.MysqlConn) {
	require := require.New(t)
	res, err := conn.Query("select * from user")
	require.Nil(err)
	t.Log(res)
}

func testInsert(t *testing.T, conn *mysql.MysqlConn) {
	require := require.New(t)
	num, lastId, err := conn.Exec("insert into user(uid,name,age) values(?,?,?)", 1234, "zhuxueting", 28)
	require.Nil(err)
	t.Log(num, lastId)
}

func testUpdate(t *testing.T, conn *mysql.MysqlConn) {
	require := require.New(t)
	num, lastId, err := conn.Exec("update user set age=? where uid=?", 18, 1234)
	require.Nil(err)
	t.Log(num, lastId)
}

func testDelete(t *testing.T, conn *mysql.MysqlConn) {
	require := require.New(t)
	num, lastId, err := conn.Exec("delete from user where uid=?", 1234)
	require.Nil(err)
	t.Log(num, lastId)
}

func testCloseConn(t *testing.T, conn *mysql.MysqlConn) {
	require := require.New(t)
	err := conn.Close()
	require.Nil(err)
}

func testTx(t *testing.T, conn *mysql.MysqlConn) {
	require := require.New(t)
	tx, err := conn.NewTx()
	require.Nil(err)

	for i := 2000; i < 3000; i++ {
		tx.Exec("insert into user(uid,name,age) values(?,?,?)", i, fmt.Sprintf("foo_%d", i), i%100+1)
	}

	rows, err := tx.Query("select * from user where uid=2588")
	require.Nil(err)
	require.Equal(1, len(rows))

	num, lastId, err := tx.Exec("delete from user where uid>=2000")
	require.Nil(err)
	t.Log(num, lastId)

	err = tx.Commit()
	require.Nil(err)
}
