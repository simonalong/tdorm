package test

import (
	"fmt"
	"testing"

	"github.com/simonalong/gole/maps"
	goleTime "github.com/simonalong/gole/time"
	"github.com/simonalong/tdorm/column"
	"github.com/simonalong/tdorm/condition"
)

func TestList1(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()

	timeData1, _ := goleTime.ParseTimeYmdHmsS("2024-07-16 11:19:23.291")
	timeData2, _ := goleTime.ParseTimeYmdHmsS("2024-07-16 11:12:23.291")

	// 新增：使用map
	insertMap1 := maps.New()
	insertMap1.SetSort(true)
	insertMap1.Put("ts", timeData1)
	insertMap1.Put("name", "大牛市")
	insertMap1.Put("age", "18")
	insertMap1.Put("address", "浙江杭州市")
	_, err := tdorm.Insert("td_china", insertMap1)

	insertMap2 := maps.New()
	insertMap2.SetSort(true)
	insertMap2.Put("ts", timeData2)
	insertMap2.Put("name", "大牛市")
	insertMap2.Put("age", "18")
	insertMap2.Put("address", "浙江杭州市")
	_, err = tdorm.Insert("td_china", insertMap2)
	checkErr(err, "插入异常")

	dataMaps, err := tdorm.List("td_china", column.Of("name", "age"), condition.New().Eq("age", 18))
	if dataMaps != nil {
		for _, dataMap := range dataMaps {
			fmt.Println(dataMap.ToJson())
		}
	}
}

func TestList2(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()

	timeData1, _ := goleTime.ParseTimeYmdHmsS("2024-07-16 11:19:23.191")
	timeData2, _ := goleTime.ParseTimeYmdHmsS("2024-07-16 11:12:23.391")
	timeData3, _ := goleTime.ParseTimeYmdHmsS("2024-07-16 11:12:23.491")

	// 新增：使用map
	insertMap1 := maps.New()
	insertMap1.SetSort(true)
	insertMap1.Put("ts", timeData1)
	insertMap1.Put("name", "大牛市")
	insertMap1.Put("age", "18")
	insertMap1.Put("address", "浙江杭州市")
	_, err := tdorm.Insert("td_china", insertMap1)

	insertMap2 := maps.New()
	insertMap2.SetSort(true)
	insertMap2.Put("ts", timeData2)
	insertMap2.Put("name", "大牛市")
	insertMap2.Put("age", "18")
	insertMap2.Put("address", "浙江杭州市")
	_, err = tdorm.Insert("td_china", insertMap2)

	insertMap3 := maps.New()
	insertMap3.SetSort(true)
	insertMap3.Put("ts", timeData3)
	insertMap3.Put("name", "大牛市")
	insertMap3.Put("age", "18")
	insertMap3.Put("address", "浙江杭州市")
	_, err = tdorm.Insert("td_china", insertMap3)
	checkErr(err, "插入异常")

	dataMaps, err := tdorm.ListOfDistinct("td_china", column.Of("name", "age"), condition.New().Eq("age", 18))
	if dataMaps != nil {
		for _, dataMap := range dataMaps {
			fmt.Println(dataMap.ToJson())
		}
	}
}
