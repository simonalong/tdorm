package test

import (
	"fmt"
	"testing"

	"github.com/simonalong/tdorm"
)

func newDbOriginal() *tdorm.TdClient {
	host := "localhost"
	port := 6030
	user := "root"
	password := "taosdata"
	db := "td_orm"

	tdClient := tdorm.NewConnectOriginal(host, port, user, password, db)

	_, err := tdClient.Exec("create database if not exists td_orm")
	checkErr(err, "建库失败")

	// 建超级表
	_, err = tdClient.Exec("create stable if not exists td_demo1(ts timestamp, name nchar(32), age int, address nchar(128)) tags (station nchar(128))")
	checkErr(err, "建超级表失败")

	// 建子表
	_, err = tdClient.Exec("create table if not exists td_america using td_demo1(`station`) tags(\"america\")")
	checkErr(err, "建超级表失败")

	// 建子表
	_, err = tdClient.Exec("create table if not exists td_china using td_demo1(`station`) tags(\"china\")")
	checkErr(err, "建超级表失败")

	return tdClient
}

func newDbWs() *tdorm.TdClient {
	host := "localhost"
	port := 6041
	user := "root"
	password := "taosdata"
	db := "td_orm"

	tdClient := tdorm.NewConnectWebsocket(host, port, user, password, db)

	// 建超级表
	_, err := tdClient.Exec("create stable if not exists td_demo1(ts timestamp, name nchar(32), age int, address nchar(128)) tags (station nchar(128))")
	checkErr(err, "建超级表失败")

	// 建子表
	_, err = tdClient.Exec("create table if not exists td_america using td_demo1(`station`) tags(\"america\")")
	checkErr(err, "建超级表失败")

	// 建子表
	_, err = tdClient.Exec("create table if not exists td_china using td_demo1(`station`) tags(\"china\")")
	checkErr(err, "建超级表失败")

	return tdClient
}

// 测试高级连接
func TestConnectOriginal(t *testing.T) {
	host := "localhost"
	port := 6030
	user := "root"
	password := "taosdata"
	db := ""

	tdorm := tdorm.NewConnectOriginal(host, port, user, password, db)
	if tdorm != nil {
		fmt.Println("连接成功")
	}
}

func TestConnectRest(t *testing.T) {
	host := "localhost"
	port := 6030
	user := "root"
	password := "taosdata"
	db := ""

	ormDb := tdorm.NewConnectRest(host, port, user, password, db)
	if ormDb != nil {
		fmt.Println("连接成功")
	}
}

func TestConnectWebsocket(t *testing.T) {
	host := "localhost"
	port := 6030
	user := "root"
	password := "taosdata"
	db := ""

	ormDb := tdorm.NewConnectWebsocket(host, port, user, password, db)
	if ormDb != nil {
		fmt.Println("连接成功")
	}
}
