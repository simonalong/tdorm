package test

import (
	"fmt"
	"testing"

	"github.com/simonalong/gole/maps"
	goleTime "github.com/simonalong/gole/time"
	"github.com/simonalong/tdorm/condition"
)

func TestValues1(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()

	timeData1, _ := goleTime.ParseTimeYmdHmsS("2024-07-16 11:19:23.211")
	timeData2, _ := goleTime.ParseTimeYmdHmsS("2024-07-16 11:19:23.221")

	// 新增：使用map
	insertMap1 := maps.New()
	insertMap1.SetSort(true)
	insertMap1.Put("ts", timeData1)
	insertMap1.Put("name", "大牛市1")
	insertMap1.Put("age", "18")
	insertMap1.Put("address", "浙江杭州市")
	_, err := tdorm.Insert("td_china", insertMap1)

	insertMap2 := maps.New()
	insertMap2.SetSort(true)
	insertMap2.Put("ts", timeData2)
	insertMap2.Put("name", "大牛市2")
	insertMap2.Put("age", "18")
	insertMap2.Put("address", "浙江杭州市")
	_, err = tdorm.Insert("td_china", insertMap2)
	checkErr(err, "插入异常")

	//data, err := tdorm.Value("td_china", "age", query.New().Eq("ts", timeData))
	//data, err := tdorm.Value("td_china", "age", query.New().Eq("ts", "2024-07-16 11:19:23.291"))
	names, err := tdorm.Values("td_china", "name", condition.New().Eq("age", 18))
	// 新增：使用entity
	//tdChinaDomain := OrmChinaDomain1{}
	//tdorm.OneEntity(&tdChinaDomain, "td_china", query.New().Eq("ts", timeData))
	fmt.Println(names)
}

func TestValues2(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()

	timeData1, _ := goleTime.ParseTimeYmdHmsS("2024-07-16 11:19:23.211")
	timeData2, _ := goleTime.ParseTimeYmdHmsS("2024-07-16 11:19:23.221")

	// 新增：使用map
	insertMap1 := maps.New()
	insertMap1.SetSort(true)
	insertMap1.Put("ts", timeData1)
	insertMap1.Put("name", "大牛市1")
	insertMap1.Put("age", "18")
	insertMap1.Put("address", "浙江杭州市")
	_, err := tdorm.Insert("td_china", insertMap1)

	insertMap2 := maps.New()
	insertMap2.SetSort(true)
	insertMap2.Put("ts", timeData2)
	insertMap2.Put("name", "大牛市2")
	insertMap2.Put("age", "18")
	insertMap2.Put("address", "浙江杭州市")
	_, err = tdorm.Insert("td_china", insertMap2)

	insertMap3 := maps.New()
	insertMap3.SetSort(true)
	insertMap3.Put("ts", timeData2)
	insertMap3.Put("name", "大牛市2")
	insertMap3.Put("age", "18")
	insertMap3.Put("address", "浙江杭州市")
	_, err = tdorm.Insert("td_china", insertMap3)
	checkErr(err, "插入异常")

	//data, err := tdorm.Value("td_china", "age", query.New().Eq("ts", timeData))
	//data, err := tdorm.Value("td_china", "age", query.New().Eq("ts", "2024-07-16 11:19:23.291"))
	names, err := tdorm.ValuesOfDistinct("td_china", "name", condition.New().Eq("age", 18))
	// 新增：使用entity
	//tdChinaDomain := OrmChinaDomain1{}
	//tdorm.OneEntity(&tdChinaDomain, "td_china", query.New().Eq("ts", timeData))
	fmt.Println(names)
}
