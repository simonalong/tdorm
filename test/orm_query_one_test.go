package test

import (
	"fmt"
	"testing"

	"github.com/simonalong/gole/maps"
	goleTime "github.com/simonalong/gole/time"
	"github.com/simonalong/tdorm/column"
	"github.com/simonalong/tdorm/condition"
)

func TestOne1(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()

	timeData, _ := goleTime.ParseTimeYmdHmsS("2024-07-16 11:19:23.291")

	// 新增：使用map
	insertMap := maps.New()
	insertMap.SetSort(true)
	insertMap.Put("ts", timeData)
	insertMap.Put("name", "大牛市")
	insertMap.Put("age", "18")
	insertMap.Put("address", "浙江杭州市")
	insertRlt, err := tdorm.Insert("td_china", insertMap)
	//_, err := tdorm.Insert("td_china", insertMap)
	checkErr(err, "插入异常")
	fmt.Println(insertRlt)

	dataMap, err := tdorm.One("td_china", column.Of("*"), condition.New().Eq("ts", timeData))
	// 新增：使用entity
	//tdChinaDomain := OrmChinaDomain1{}
	//tdorm.OneEntity(&tdChinaDomain, "td_china", query.New().Eq("ts", timeData))
	if dataMap != nil {
		fmt.Println(dataMap.ToJson())
	}
}

func TestOne2(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()

	timeData, _ := goleTime.ParseTimeYmdHmsS("2024-07-16 11:19:23.291")

	// 新增：使用map
	insertMap := maps.New()
	insertMap.SetSort(true)
	insertMap.Put("ts", timeData)
	insertMap.Put("name", "大牛市")
	insertMap.Put("age", "18")
	insertMap.Put("address", "浙江杭州市")
	insertRlt, err := tdorm.Insert("td_china", insertMap)
	//_, err := tdorm.Insert("td_china", insertMap)
	checkErr(err, "插入异常")
	fmt.Println(insertRlt)

	dataMap, err := tdorm.One("td_china", column.Of("name", "age"), condition.New().Eq("ts", timeData))
	// 新增：使用entity
	//tdChinaDomain := OrmChinaDomain1{}
	//tdorm.OneEntity(&tdChinaDomain, "td_china", query.New().Eq("ts", timeData))
	if dataMap != nil {
		fmt.Println(dataMap.ToJson())
	}
}
