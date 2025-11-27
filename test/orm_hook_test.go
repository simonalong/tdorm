package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/simonalong/gole/maps"
	goleTime "github.com/simonalong/gole/time"
	"github.com/simonalong/tdorm"
	"github.com/simonalong/tdorm/column"
	"github.com/simonalong/tdorm/condition"
	"github.com/simonalong/tdorm/constants"
)

func TestHook(t *testing.T) {
	//tdClient := newDbOriginal()
	tdClient := newDbWs()

	tdHook := DemoTdHook{}
	tdClient.AddHook(&tdHook)

	_, err := tdClient.Exec("create database if not exists td_orm")
	checkErr(err, "建库失败")
	//
	//// 建超级表：请先创建库 td_orm
	//_, err = tdClient.Exec("create stable if not exists td_demo1(ts timestamp, name nchar(32), age int, address nchar(128)) tags (station nchar(128))")
	//checkErr(err, "建超级表失败")
	//
	//// 建子表
	//_, err = tdClient.Exec("create table if not exists td_china using td_demo1(`station`) tags(\"china\")")
	//checkErr(err, "建子表失败")

	// 新增：使用map，对应SQL：【insert into td_orm.td_china (`ts`,`name`,`age`,`address`) values (?,?,?,?)】
	insertMap := maps.NewSort().Put("ts", time.Now()).Put("name", "大牛市").Put("age", "18").Put("address", "浙江杭州市")
	_, err = tdClient.Insert("td_china", insertMap)
	checkErr(err, "插入数据失败")

	// 新增：使用entity，对应SQL：【insert into td_orm.td_china (`ts`,`name`,`age`,`address`) values (?,?,?,?)】
	type OrmChinaDomain struct {
		Timestamp time.Time `column:"ts"`
		Na        string    `column:"name"`
		Ag        int       `column:"age"`
		add       string    `column:"address"`
	}
	tdChinaDomain := OrmChinaDomain{Timestamp: time.Now(), Na: "大牛市2", Ag: 19, add: "浙江温州市"}
	_, err = tdClient.InsertEntity("td_china", tdChinaDomain)
	checkErr(err, "插入数据失败")

	//// 删除，对应SQL：【delete from td_orm.td_china where `ts` > ?】
	//_, err = tdClient.Delete("td_china", query.New().Gt("ts", "2024-07-12 12:00:00.000"))
	////_, err = tdClient.Delete("td_china", query.New().Gt("ts", "now-2d")) // websocket 不支持
	//checkErr(err, "删除数据失败")

	// 查询：一行，对应SQL：【select `name`, `age` from td_orm.td_china where `ts` = ? limit 1】
	timeData, _ := goleTime.ParseTimeYmdHmsS("2024-07-16 11:19:23.291")
	_, err = tdClient.One("td_china", column.Of("name", "age"), condition.New().Eq("ts", timeData))
	checkErr(err, "查询one数据失败")

	// 查询：多行，对应SQL：【select `name`, `age` from td_orm.td_china where `age` = ?】
	_, err = tdClient.List("td_china", column.Of("name", "age"), condition.New().Eq("age", 18))
	checkErr(err, "查询list数据失败")

	// 查询：一个，对应SQL：【select `name` from td_orm.td_china where `age` > ? and `ts` = ? limit 1】
	_, err = tdClient.Value("td_china", "name", condition.New().Gt("age", 12).Eq("ts", "2024-07-16 11:19:23.291"))
	checkErr(err, "查询value数据失败")

	// 查询：一列，对应SQL：【select `name` from td_orm.td_china where `age` = ?】
	_, err = tdClient.Values("td_china", "name", condition.New().Eq("age", 18))

	// 查询：个数，对应SQL：【select count(*) as cnt from td_orm.td_china where `age` = ?】
	_, _ = tdClient.Count("td_china", condition.New().Eq("age", 18))

	// 新增：批量新增，对应SQL：【insert into td_orm.td_china (`ts`,`name`,`age`,`address`) values (?,?,?,?)】
	timeData1, _ := goleTime.ParseTimeYmdHmsS("2024-07-15 14:01:23.391")
	timeData2, _ := goleTime.ParseTimeYmdHmsS("2024-07-15 15:01:23.391")
	timeData3, _ := goleTime.ParseTimeYmdHmsS("2024-07-15 16:01:23.391")
	var insertMaps []*maps.GoleMap
	insertMaps = append(insertMaps, maps.NewSort().Put("ts", timeData1).Put("name", "大牛市1").Put("age", "18").Put("address", "浙江杭州市1"))
	insertMaps = append(insertMaps, maps.NewSort().Put("ts", timeData2).Put("name", "大牛市2").Put("age", "28").Put("address", "浙江杭州市2"))
	insertMaps = append(insertMaps, maps.NewSort().Put("ts", timeData3).Put("name", "大牛市3").Put("age", "38").Put("address", "浙江杭州市3"))

	cnt, err := tdClient.InsertBatch("td_china", insertMaps)
	checkErr(err, "批量插入数据失败")
	fmt.Println("cnt: ", cnt)
}

type DemoTdHook struct {
}

func (tdHook *DemoTdHook) Before(thc *tdorm.TdHookContext) (*tdorm.TdHookContext, error) {
	fmt.Println("before: ", thc.Sql)
	return thc, nil
}

func (tdHook *DemoTdHook) After(thc *tdorm.TdHookContext) error {
	err := thc.Err
	if err != nil {
		fmt.Println("err: ", err)
		return err
	}
	fmt.Println("after: ")
	switch thc.RunType {
	case constants.EXE:
		fmt.Println(thc.ResultOfExe.RowsAffected())
	case constants.QUERY:
		if thc.ResultOfQueryOfDriverRows != nil {
			fmt.Println(thc.ResultOfQueryOfDriverRows)
		} else if thc.ResultOfQueryOfSqlRows != nil {
			fmt.Println(thc.ResultOfQueryOfDriverRows)
		}
	case constants.INSERT, constants.SAVE:
		fmt.Println(thc.ResultOfInsert)
	case constants.BATCH_INSERT:
		fmt.Println(thc.ResultOfInsert)
	}
	return nil
}
