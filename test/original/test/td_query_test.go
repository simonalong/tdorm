package test

//
//import (
//	"database/sql"
//	"fmt"
//	"github.com/taosdata/driver-go/v3/af"
//	"log"
//	"testing"
//)
//
//func TestQueryOriginal(t *testing.T) {
//	var taosDSN = "root:taosdata@tcp(localhost:6030)/"
//	//var taosDSN = "root:taosdata@tcp(localhost:6030)/log"
//	taos, err := sql.Open("taosSql", taosDSN)
//	if err != nil {
//		log.Fatalln("failed to connect TDengine, err:", err)
//		return
//	}
//	fmt.Println("Connected")
//	defer taos.Close()
//
//	rows, err := taos.Query("select * from td_orm.td_china where ts > ?", "2024-04-02")
//	checkErr(err, "查询失败")
//	var a interface{}
//	var b interface{}
//	var c interface{}
//	var d interface{}
//	for rows.Next() {
//		err = rows.Scan(&a, &b, &c, &d)
//		checkErr(err, "转换失败")
//		fmt.Println(a, b, c, d)
//	}
//}
//
//func TestQueryAf(t *testing.T) {
//	conn, err := af.Open("localhost", "root", "taosdata", "", 6030)
//	if err != nil {
//		log.Fatalln("failed to connect TDengine, err:", err)
//		return
//	}
//	fmt.Println("Connected")
//	defer conn.Close()
//
//	rows, err := conn.Query("select ts, name, age a, address from td_orm.td_china where ts > ? and age > ?", "2024-04-02", 10)
//	checkErr(err, "插入失败")
//	rowMapList := maps.FromRows(rows)
//
//	for _, ormMap := range rowMapList {
//		fmt.Println(ormMap.ToString())
//	}
//}
//
//func TestQueryAf2(t *testing.T) {
//	conn, err := af.Open("localhost", "root", "taosdata", "", 6030)
//	if err != nil {
//		log.Fatalln("failed to connect TDengine, err:", err)
//		return
//	}
//	fmt.Println("Connected")
//	defer conn.Close()
//
//	rows, err := conn.Query("select * from td_orm.td_china, td_orm.td_america where td_china.ts = td_america.ts  and td_china.age > ?", 12)
//	checkErr(err, "插入失败")
//	rowMapList := maps.FromRows(rows)
//
//	for _, ormMap := range rowMapList {
//		fmt.Println(ormMap.ToString())
//	}
//}
//
//// 查询一行数据
//func TestQueryOne(t *testing.T) {
//	conn, err := af.Open("localhost", "root", "taosdata", "", 6030)
//	if err != nil {
//		log.Fatalln("failed to connect TDengine, err:", err)
//		return
//	}
//	fmt.Println("Connected")
//	defer conn.Close()
//
//	rows, err := conn.Query("select * from td_orm.td_china where td_china.age > ? limit 0,12", 12)
//	checkErr(err, "插入失败")
//	rowMapList := maps.FromRows(rows)
//
//	for _, ormMap := range rowMapList {
//		fmt.Println(ormMap.ToString())
//	}
//}
//
//// 查询多行数据
//func TestQueryList(t *testing.T) {
//	conn, err := af.Open("localhost", "root", "taosdata", "", 6030)
//	if err != nil {
//		log.Fatalln("failed to connect TDengine, err:", err)
//		return
//	}
//	fmt.Println("Connected")
//	defer conn.Close()
//
//	rows, err := conn.Query("select * from td_orm.td_china where td_china.age > ? limit 0,12", 12)
//	checkErr(err, "插入失败")
//	rowMapList := maps.FromRows(rows)
//
//	for _, ormMap := range rowMapList {
//		fmt.Println(ormMap.ToString())
//	}
//}
//
//// 查询某个值
//func TestQueryValue(t *testing.T) {
//	conn, err := af.Open("localhost", "root", "taosdata", "", 6030)
//	if err != nil {
//		log.Fatalln("failed to connect TDengine, err:", err)
//		return
//	}
//	fmt.Println("Connected")
//	defer conn.Close()
//
//	rows, err := conn.Query("select `name`, `age` from td_orm.td_china where `ts` = ? limit 1", 12)
//	checkErr(err, "插入失败")
//	rowMapList := maps.FromRows(rows)
//
//	for _, ormMap := range rowMapList {
//		fmt.Println(ormMap.ToString())
//	}
//}
//
//// 查询某列值
//func TestQueryValues(t *testing.T) {
//	conn, err := af.Open("localhost", "root", "taosdata", "", 6030)
//	if err != nil {
//		log.Fatalln("failed to connect TDengine, err:", err)
//		return
//	}
//	fmt.Println("Connected")
//	defer conn.Close()
//
//	rows, err := conn.Query("select `name` from td_orm.td_china where td_china.age > ?", 12)
//	checkErr(err, "插入失败")
//	rowMapList := maps.FromRows(rows)
//
//	for _, ormMap := range rowMapList {
//		fmt.Println(ormMap.ToString())
//	}
//}
//
//// 查询某列值
//func TestQueryValues2(t *testing.T) {
//	conn, err := af.Open("localhost", "root", "taosdata", "", 6030)
//	if err != nil {
//		log.Fatalln("failed to connect TDengine, err:", err)
//		return
//	}
//	fmt.Println("Connected")
//	defer conn.Close()
//
//	rows, err := conn.Query("select distinct `name` from td_orm.td_china where td_china.age > ?", 12)
//	checkErr(err, "插入失败")
//	rowMapList := maps.FromRows(rows)
//
//	for _, ormMap := range rowMapList {
//		fmt.Println(ormMap.ToString())
//	}
//}
//
//func TestQueryIn(t *testing.T) {
//	conn, err := af.Open("localhost", "root", "taosdata", "", 6030)
//	if err != nil {
//		log.Fatalln("failed to connect TDengine, err:", err)
//		return
//	}
//	fmt.Println("Connected")
//	defer conn.Close()
//
//	rows, err := conn.Query("select * from td_orm.td_china where td_china.name in (?, ?, ?)", "'大牛市'", "'大牛市1'", "'大牛市2'")
//	//rows, err := conn.QueryForOriginal("show create stable sb_fhsj")
//	checkErr(err, "插入失败")
//	rowMapList := maps.FromRows(rows)
//
//	for _, ormMap := range rowMapList {
//		fmt.Println(ormMap.ToString())
//	}
//}
//
//func TestQueryLike(t *testing.T) {
//	conn, err := af.Open("localhost", "root", "taosdata", "", 6030)
//	if err != nil {
//		log.Fatalln("failed to connect TDengine, err:", err)
//		return
//	}
//	fmt.Println("Connected")
//	defer conn.Close()
//
//	rows, err := conn.Query("select * from td_orm.td_china where td_china.name like ?", "'%牛_'")
//	//rows, err := conn.QueryForOriginal("show create stable sb_fhsj")
//	checkErr(err, "插入失败")
//	rowMapList := maps.FromRows(rows)
//
//	for _, ormMap := range rowMapList {
//		fmt.Println(ormMap.ToString())
//	}
//}
//
//func TestName(t *testing.T) {
//	conn, err := af.Open("localhost", "root", "taosdata", "td_orm", 6030)
//	if err != nil {
//		log.Fatalln("failed to connect TDengine, err:", err)
//		return
//	}
//	fmt.Println("Connected")
//	defer conn.Close()
//
//	//times := param.NewParam(1).AddTimestamp(time.Now(), 0)
//	rlt, err := conn.Exec("delete from td_china where ts > ?", "now-5d")
//	checkErr(err, "执行失败")
//	fmt.Println(rlt)
//}
