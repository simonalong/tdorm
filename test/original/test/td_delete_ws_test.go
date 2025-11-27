package test

import (
	"fmt"
	"testing"

	goleTime "github.com/simonalong/gole/time"
)

// 删除这个的搜索条件只支持唯一的时间字段
func TestDeleteOfWs(t *testing.T) {
	newDbWs()

	// 新增数据
	_, err := dbOfWs.Exec("insert into td_orm.td_china(`ts`,`name`,`age`,`address`) values('2024-07-12 12:43:14.321','zhou',13,'hangzhou')")
	checkErr(err, "插入失败")

	// 查询数据
	rowData := dbOfWs.QueryRow("select ts from td_orm.td_china where `ts` = '2024-07-12 12:43:14.321'")
	var timeStr string
	err = rowData.Scan(&timeStr)
	checkErr(err, "查询失败")
	fmt.Println(timeStr)

	// 删除数据
	_, err = dbOfWs.Exec("delete from td_orm.td_china where `ts` = '2024-07-12 12:43:14.321'")
	checkErr(err, "删除异常")
}

func TestDeleteOfWs2(t *testing.T) {
	taosDb := newDbWs()
	defer taosDb.Close()

	// 新增数据
	//_, err := dbOfWs.Exec("")
	//checkErr(err, "插入失败")
	st1, _ := dbOfWs.Prepare("insert into td_orm.td_china(`ts`,`name`,`age`,`address`) values(?,?,?,?)")
	ti, _ := goleTime.ParseTimeYmdHmsS("2024-07-12 12:43:14.321")
	_, err := st1.Exec(ti, "zhou", 13, "hangzhou")
	checkErr(err, "添加异常")

	// 删除数据
	//st, _ := dbOfWs.Prepare("delete from td_orm.td_china where `ts` = ?")
	//ti, _ := goleTime.ParseTimeYmdHmsS("2024-07-12 12:43:14.321")
	_, err = dbOfWs.Exec("delete from td_orm.td_china where `ts` = '2024-07-12 12:43:14.321'")
	checkErr(err, "删除异常")

	// websocket：如下有问题；stmt exec error: wrong number of parameters
	st, _ := dbOfWs.Prepare("delete from td_orm.td_china where `ts` = ?")
	ti, _ = goleTime.ParseTimeYmdHmsS("2024-07-12 12:43:14.321")
	//_, err = st.Exec("2024-07-12 12:43:14.321")
	_, err = st.Exec(ti)
	checkErr(err, "删除异常")
}
