package test

import (
	"fmt"
	"testing"

	"github.com/simonalong/gole/maps"
	goleTime "github.com/simonalong/gole/time"
	"github.com/simonalong/tdorm/condition"
)

func TestValue(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()

	start, _ := goleTime.ParseTimeYmdHmsS("2024-09-10 13:19:23.291")
	end, _ := goleTime.ParseTimeYmdHmsS("2024-09-10 15:19:23.291")

	//var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	//start.In(cstSh)
	//end.In(cstSh)

	//time.ParseInLocation(goleTime.FmtYMdHmsS, "2024-09-10 14:19:23.291", time.Local)

	// 新增：使用map
	insertMap := maps.New()
	insertMap.SetSort(true)
	insertMap.Put("ts", start)
	insertMap.Put("name", "大牛市")
	insertMap.Put("age", "18")
	insertMap.Put("address", "浙江杭州市")
	tdorm.Insert("td_china", insertMap)
	_, err := tdorm.Insert("td_china", insertMap)
	checkErr(err, "插入异常")
	//fmt.Println(insertRlt)

	//data, err := tdorm.Value("td_china", "age", query.New().Eq("ts", timeData))
	//data, err := tdorm.Value("td_china", "age", query.New().Eq("ts", "2024-07-16 11:19:23.291"))
	//_, err = tdorm.Value("td_china", "age", query.New().Eq("ts", "2024-07-16 11:19:23.291"))
	// 新增：使用entity
	//tdChinaDomain := OrmChinaDomain1{}
	//tdorm.OneEntity(&tdChinaDomain, "td_china", query.New().Eq("ts", timeData))
	//fmt.Println(data)

	//ts, _ := tdorm.Value("td_china", "ts", query.New().Eq("age", 12))
	//fmt.Println(ts)

	//logger.Group(orm.DEFAULT_LOG_GROUP).SetLevel(logrus.DebugLevel)

	ts, _ := tdorm.Value("td_china", "ts", condition.New().Ge("ts", start).Lt("ts", end))
	fmt.Println(ts)
}
