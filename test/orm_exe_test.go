package test

import (
	"database/sql"
	"testing"

	"github.com/simonalong/gole/logger"

	"github.com/simonalong/tdorm"
)

func TestCreate(t *testing.T) {
	host := "localhost"
	port := 6041
	user := "root"
	password := "taosdata"
	db := ""

	tdorm := tdorm.NewConnectWebsocket(host, port, user, password, db)

	// 建库
	_, err := tdorm.Exec("create database if not exists td_orm")
	checkErr(err, "建库失败")

	_, err = tdorm.Exec("create stable if not exists td_orm1.td_demo1(ts timestamp, name nchar(32), age int, address nchar(128)) tags (station nchar(128))")
	checkErr(err, "建超级表失败")
	// 建超级表
	_, err = tdorm.Exec("create stable if not exists td_orm.td_demo1(ts timestamp, name nchar(32), age int, address nchar(128)) tags (station nchar(128))")
	checkErr(err, "建超级表失败")

	// 建子表
	_, err = tdorm.Exec("create table if not exists td_orm.td_china using td_orm.td_demo1(`station`) tags(\"china\")")
	checkErr(err, "建超级表失败")
}

func TestCreate2(t *testing.T) {
	//host := "localhost"
	//port := 6041
	//user := "root"
	//password := "taosdata"
	//db := "td_orm1"
	//pOrm := tdorm.NewConnectWebsocket(host, port, user, password, db)
	//_, err := pOrm.Exec("create stable if not exists td_orm1.td_demo1(ts timestamp, name nchar(32), age int, address nchar(128)) tags (station nchar(128))")
	//checkErr(err, "建超级表失败")

	taos, err := sql.Open("taosWS", "root:taosdata@ws(localhost:6041)/td_orm1")
	if err != nil {
		logger.Errorf("tdengine连接异常，请检查配置 taosDSN，异常：%v", err.Error())
		return
	}
	_, err = taos.Exec("create stable if not exists td_orm1.td_demo1(ts timestamp, name nchar(32), age int, address nchar(128)) tags (station nchar(128))")
	if err != nil {
		logger.Group("tdorm").Errorf("执行execute异常，【executeSql ==> %v】【params ==> %v】报错：%v", "", nil, err.Error())
	} else {
		logger.Group("tdorm").Debugf("【executeSql ==> %v】【params ==> %v】", "", nil)
	}
}

func TestCreateWs(t *testing.T) {
	taos, err := sql.Open("taosWS", "root:taosdata@ws(localhost:6041)/td_orm1")
	if err != nil {
		logger.Errorf("tdengine连接异常，请检查配置 taosDSN=%v，异常：%v", "root:taosdata@ws(localhost:6041)/td_orm1", err.Error())
		return
	}
	_, err = taos.Exec("create stable if not exists td_orm1.td_demo1(ts timestamp, name nchar(32), age int, address nchar(128)) tags (station nchar(128))")
	if err != nil {
		t.Fatal(err)
	}
}
