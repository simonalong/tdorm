package test

import (
	"database/sql"
	"fmt"
	//_ "github.com/taosdata/driver-go/v3/taosSql"
	"log"
	"testing"
)

func TestOriginalInsert(t *testing.T) {
	var taosDSN = "root:taosdata@tcp(localhost:6030)/original"
	taos, err := sql.Open("taosSql", taosDSN)
	if err != nil {
		log.Fatalln("failed to connect TDengine, err:", err)
		return
	}
	fmt.Println("Connected")
	defer taos.Close()

	//// 建超级表
	//_, err = taos.Exec("create stable if not exists td_demo1(ts timestamp, name nchar(32), age int, address nchar(128)) tags (station nchar(128))")
	//checkErr(err, "建超级表失败")
	//
	//// 建子表
	//_, err = taos.Exec("create table if not exists td_america using td_demo1(`station`) tags(\"america\")")
	//checkErr(err, "建超级表失败")
	//
	//// 建子表
	//_, err = taos.Exec("create table if not exists td_china using td_demo1(`station`) tags(\"china\")")
	//checkErr(err, "建超级表失败")

	//sql := "insert into original.td_china values(?, ?, ?, ?) (?, ?, ?, ?)"
	////rlt, err := taos.Exec(sql, time.Now(), "zhou", 12, "hangzhou")
	//smt, err := taos.Prepare(sql)
	//checkErr(err, "prepare sql error")

	//rlt, err := smt.Exec(time.Now(), "zhou", 12, "hangzhou", time.Now(), "zhou", 12, "hangzhou1")
	//checkErr(err, "插入失败")
	//
	//fmt.Println(rlt.LastInsertId())
	//fmt.Println(rlt.RowsAffected())

	//
	//now := time.Now()
	//stmt := taos.InsertStmt()
	//err = stmt.Prepare("insert into example_stmt.tb1 values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	//if err != nil {
	//	panic(err)
	//}
	//now := time.Now()
	//params := make([]*param.Param, 14)
	//params[0] = param.NewParam(2).
	//	AddTimestamp(now, common.PrecisionMilliSecond).
	//	AddTimestamp(now.add(time.Second), common.PrecisionMilliSecond)
	//params[1] = param.NewParam(2).AddBool(true).AddNull()
	//params[2] = param.NewParam(2).AddTinyint(2).AddNull()
	//params[3] = param.NewParam(2).AddSmallint(3).AddNull()
	//params[4] = param.NewParam(2).AddInt(4).AddNull()
	//params[5] = param.NewParam(2).AddBigint(5).AddNull()
	//params[6] = param.NewParam(2).AddUTinyint(6).AddNull()
	//params[7] = param.NewParam(2).AddUSmallint(7).AddNull()
	//params[8] = param.NewParam(2).AddUInt(8).AddNull()
	//params[9] = param.NewParam(2).AddUBigint(9).AddNull()
	//params[10] = param.NewParam(2).AddFloat(10).AddNull()
	//params[11] = param.NewParam(2).AddDouble(11).AddNull()
	//params[12] = param.NewParam(2).AddBinary([]byte("binary")).AddNull()
	//params[13] = param.NewParam(2).AddNchar("nchar").AddNull()
	//
	//paramTypes := param.NewColumnType(14).
	//	AddTimestamp().
	//	AddBool().
	//	AddTinyint().
	//	AddSmallint().
	//	AddInt().
	//	AddBigint().
	//	AddUTinyint().
	//	AddUSmallint().
	//	AddUInt().
	//	AddUBigint().
	//	AddFloat().
	//	AddDouble().
	//	AddBinary(6).
	//	AddNchar(5)
	//err = stmt.BindParam(params, paramTypes)
	//if err != nil {
	//	panic(err)
	//}
	//err = stmt.AddBatch()
	//if err != nil {
	//	panic(err)
	//}
	//err = stmt.Execute()
	//if err != nil {
	//	panic(err)
	//}
	//err = stmt.Close()
	//if err != nil {
	//	panic(err)
	//}
}
