package test

import (
	"fmt"
	"testing"

	"github.com/simonalong/gole/maps"
	goleTime "github.com/simonalong/gole/time"
	"github.com/simonalong/tdorm/condition"
)

func TestRmv1(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()

	timeData1, _ := goleTime.ParseTimeYmdHmsS("2024-07-15 14:11:23.391")

	// 新增：使用map
	insertMap := maps.New()
	insertMap.SetSort(true)
	insertMap.Put("ts", timeData1)
	insertMap.Put("name", "大牛市")
	insertMap.Put("age", "18")
	insertMap.Put("address", "浙江杭州市")
	//insertRlt, err := tdorm.Insert("td_china", insertMap)
	_, err := tdorm.Insert("td_china", insertMap)
	checkErr(err, "插入异常")

	deleteRlt, err := tdorm.Delete("td_china", condition.New().Eq("ts", "2024-07-15 14:11:23.391"))
	//cnt, err := tdorm.Delete("td_china", query.New().Eq("ts", "'2024-07-15 14:11:23.391'"))
	//cnt, err := tdorm.Delete("td_china", query.New().Eq("ts", timeData1))
	checkErr(err, "删除异常")
	fmt.Println(deleteRlt)
}

func TestRmv4(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()

	timeData1, _ := goleTime.ParseTimeYmdHmsS("2024-07-15 14:11:23.391")

	// 新增：使用map
	insertMap := maps.New()
	insertMap.SetSort(true)
	insertMap.Put("ts", timeData1)
	insertMap.Put("name", "大牛市")
	insertMap.Put("age", "18")
	insertMap.Put("address", "浙江杭州市")
	//insertRlt, err := tdorm.Insert("td_china", insertMap)
	_, err := tdorm.Insert("td_china", insertMap)
	checkErr(err, "插入异常")
	//fmt.Println(insertRlt)

	deleteRlt, err := tdorm.Delete("td_china", condition.New().Eq("ts", "2024-07-15 14:11:23.391"))
	//_, err = tdorm.Delete("td_china", query.New().Eq("ts", "'2024-07-15 14:11:23.391'"))
	checkErr(err, "删除异常")
	fmt.Println(deleteRlt)
}
