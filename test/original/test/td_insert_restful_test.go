package test

//
//import (
//	"database/sql"
//	_ "github.com/taosdata/driver-go/v3/taosRestful"
//	"log"
//	"testing"
//	"time"
//)
//
//var dbOfRestful *sql.DB
//
//func newTaosRestful() *sql.DB {
//	if dbOfRestful != nil {
//		return dbOfRestful
//	}
//	taos, err := sql.Open("taosRestful", "root:taosdata@http(localhost:6041)/")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 建库
//	_, err = taos.Exec("create database if not exists td_orm")
//	if err != nil {
//		log.Fatalln("failed to create database, err:", err)
//	}
//
//	// 建超级表
//	_, err = taos.Exec("create stable if not exists td_orm.td_demo1(ts timestamp, name nchar(32), age int, address nchar(128)) tags (station nchar(128))")
//	checkErr(err, "建超级表失败")
//
//	// 创建子表
//	_, err = taos.Exec("create table if not exists td_orm.td_america using td_orm.td_demo1(`station`) tags(\"america\")")
//	checkErr(err, "建超级表失败")
//
//	// 创建子表
//	_, err = taos.Exec("create table if not exists td_orm.td_china using td_orm.td_demo1(`station`) tags(\"china\")")
//	checkErr(err, "建超级表失败")
//	dbOfRestful = taos
//	return taos
//}
//
//func TestInsertRestful(t *testing.T) {
//	newTaosRestful()
//	_, err := dbOfRestful.Exec("insert into td_orm.td_china(`ts`,`name`,`age`,`address`) values('2024-07-12 12:43:14.321','zhou',13,'hangzhou')")
//	checkErr(err, "插入失败")
//}
//
//// 错误
//func TestInsertRestful2(t *testing.T) {
//	newTaosRestful()
//
//	_, err := dbOfRestful.Prepare("insert into td_orm.td_china(`ts`,`name`,`age`,`address`) values(?,?,?,?)")
//	if err != nil {
//		// 执行报错：restful does not support stmt
//		log.Fatalf("执行报错：%v", err)
//	}
//}
//
//func TestInsertRestful3(t *testing.T) {
//	newTaosRestful()
//
//	result, err := dbOfRestful.Exec("insert into td_orm.td_china(`ts`,`name`,`age`,`address`) values(?,?,?,?)", time.Now(), "zhou", 13, "hangzhou")
//	if err != nil {
//		log.Fatalf("执行报错：%v", err)
//	}
//	log.Println("result：", result)
//}
