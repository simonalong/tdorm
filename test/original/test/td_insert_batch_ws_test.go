package test

import (
	"fmt"
	"testing"

	goleTime "github.com/simonalong/gole/time"
	_ "github.com/taosdata/driver-go/v3/taosWS"
)

func TestBatchWsInsert1(t *testing.T) {
	newDbWs()
	cnt, err := dbOfWs.Exec("insert into td_orm.td_china(`ts`,`name`,`age`,`address`) values('2024-01-12 12:43:14.321','zhou',13,'hangzhou')('2024-01-22 12:43:14.321','zhou',14,'hangzhou1')")
	checkErr(err, "插入失败")
	fmt.Println(cnt)
}

// 异常报错
func TestBatchWsInsert2(t *testing.T) {
	newDbWs()
	// websocket批量插入不支持预处理，wrong number of parameters
	stmt, err := dbOfWs.Prepare("insert into td_orm.td_china(`ts`,`name`,`age`,`address`) values(?, ?, ?, ?) (?, ?, ?, ?)")
	checkErr(err, "预处理失败")

	t1, _ := goleTime.ParseTimeYmdHmsS("2023-02-12 10:13:14.321")
	t2, _ := goleTime.ParseTimeYmdHmsS("2024-01-12 12:23:14.321")

	result, err := stmt.Exec(t1, "'zhou'", 113, "'hangzhou'", t2, "'zhou12'", 111, "'hangzhou1'")
	checkErr(err, "插入失败")
	fmt.Println(result.RowsAffected())
}
