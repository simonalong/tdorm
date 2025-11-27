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
)

func TestBaseTdOrm(t *testing.T) {
	// 连接：当前暂时只支持原生连接和和websocket，暂时不支持restful；另外建议使用websocket
	tdClient := tdorm.NewConnectWebsocket("localhost", 6041, "root", "taosdata", "td_orm")

	// 建超级表：请先创建库 td_orm
	_, err := tdClient.Exec("create stable if not exists td_demo1(ts timestamp, name nchar(32), age int, address nchar(128)) tags (station nchar(128))")
	checkErr(err, "建超级表失败")

	// 建子表
	_, err = tdClient.Exec("create table if not exists td_china using td_demo1(`station`) tags(\"china\")")
	checkErr(err, "建子表失败")

	// 新增：使用map，对应SQL：【insert into td_orm.td_china (`ts`,`name`,`age`,`address`) values (?,?,?,?)】
	insertMap := maps.NewSort().Put("ts", time.Now()).Put("name", "大牛市").Put("age", "18").Put("address", "浙江杭州市")
	_, err = tdClient.Insert("td_china", insertMap)
	checkErr(err, "插入数据失败")

	// 新增：使用entity，对应SQL：【insert into td_orm.td_china (`ts`,`name`,`age`,`address`) values (?,?,?,?)】
	type OrmChinaDomain struct {
		Timestamp time.Time `column:"ts"`
		Na        string    `column:"name"`
		Ag        int       `column:"age"`
		Add       string    `column:"address"`
	}
	tdChinaDomain := OrmChinaDomain{Timestamp: time.Now(), Na: "大牛市2", Ag: 19, Add: "浙江温州市"}
	_, err = tdClient.InsertEntity("td_china", tdChinaDomain)
	checkErr(err, "插入数据失败")

	// 新增：使用标签，则如果表不存在则会自动创建，对应SQL：【insert into td_orm.td_china2_new using td_orm.td_demo1 (`station`) tags ('hangzhou1') (`ts`,`name`,`age`,`address`) values (?,?,?,?)】
	insertMap1 := maps.NewSort().Put("ts", time.Now()).Put("name", "大牛市1").Put("age", 28).Put("address", "浙江杭州市")
	tagsMap := maps.NewSort().Put("station", "hangzhou1")
	_, err = tdClient.InsertWithTag("td_china2_new", "td_demo1", tagsMap, insertMap1)
	checkErr(err, "插入异常")

	//删除，对应SQL：【delete from td_orm.td_china where `ts` > ?】
	_, err = tdClient.Delete("td_china", condition.New().Gt("ts", "2024-07-12 12:00:00.000"))
	_, err = tdClient.Delete("td_china", condition.New().Gt("ts", "now-2d"))
	checkErr(err, "删除数据失败")

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

	_, err = tdClient.InsertBatch("td_china", insertMaps)
	// 也支持批量插入实体
	//_, err = tdClient.InsertEntityBatch("td_china", insertEntities)

	// 新增：批量新增（待标签），子表可以自动创建新子表
	// 新增：批量新增，对应SQL：【insert into td_orm.td_china (`ts`,`name`,`age`,`address`) values (?,?,?,?)】
	timeData1OfBatch, _ := goleTime.ParseTimeYmdHmsS("2024-07-15 14:01:23.391")
	timeData2OfBatch, _ := goleTime.ParseTimeYmdHmsS("2024-07-15 15:01:23.391")
	timeData3OfBatch, _ := goleTime.ParseTimeYmdHmsS("2024-07-15 16:01:23.391")
	var insertMapsOfBatch []*maps.GoleMap
	insertMapsOfBatch = append(insertMapsOfBatch, maps.NewSort().Put("ts", timeData1OfBatch).Put("name", "大牛市1_batch").Put("age", "18").Put("address", "浙江杭州市1"))
	insertMapsOfBatch = append(insertMapsOfBatch, maps.NewSort().Put("ts", timeData2OfBatch).Put("name", "大牛市2_batch").Put("age", "28").Put("address", "浙江杭州市2"))
	insertMapsOfBatch = append(insertMapsOfBatch, maps.NewSort().Put("ts", timeData3OfBatch).Put("name", "大牛市3_batch").Put("age", "38").Put("address", "浙江杭州市3"))

	tagsMapOfBatch := maps.NewSort().Put("station", "batch")
	_, err = tdClient.InsertBatchWithTag("td_china_batch", "td_demo1", tagsMapOfBatch, insertMapsOfBatch)
}

func checkErr(err error, prompt string) {
	if err != nil {
		panic(fmt.Sprintf("%v：错误：%v", prompt, err.Error()))
	}
}
