package test

import (
	"testing"
	"time"

	"github.com/simonalong/gole/maps"
	goleTime "github.com/simonalong/gole/time"
	"github.com/stretchr/testify/assert"
)

// 测试：
// 连接时候不配置库名，但是所有的表名前面都添加库名，即：表名为：<dbName>.<tableName>
func TestInsertOrmMap1(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()

	// 新增：使用map
	insertMap := maps.New()
	insertMap.SetSort(true)
	insertMap.Put("ts", time.Now())
	insertMap.Put("name", "大牛市1_insert_orm")
	insertMap.Put("age", "28")
	insertMap.Put("address", "浙江杭州市")
	_, err := tdorm.Insert("td_china", insertMap)
	checkErr(err, "插入异常")
}

func TestInsertOrmMap1_1(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()

	timeData, _ := goleTime.ParseTimeYmdHmsS("2024-07-12 12:01:23.391")

	// 新增：使用map
	insertMap := maps.New()
	insertMap.SetSort(true)
	insertMap.Put("ts", timeData)
	insertMap.Put("name", "TestInsertOrmMap1_1")
	insertMap.Put("age", "19")
	insertMap.Put("address", "los angele")
	_, err := tdorm.Insert("td_america", insertMap)
	checkErr(err, "插入异常")
}

func TestInsertOrmMap2(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()

	// 新增：使用map
	insertMap := maps.New()
	insertMap.SetSort(true)
	insertMap.Put("ts", time.Now())
	insertMap.Put("name", "TestInsertOrmMap2")
	insertMap.Put("age", "182")
	insertMap.Put("address", "浙江杭州市")
	_, err := tdorm.Insert("td_orm.td_china", insertMap)
	checkErr(err, "插入异常")
}

// 使用原生标签的实体
func TestInsertEntity1(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()

	// 新增：使用entity
	tdChinaDomain := OrmChinaDomain1{
		Ts:      time.Now(),
		Name:    "TestInsertEntity1",
		Age:     19,
		Address: "浙江温州市",
	}
	_, err := tdorm.InsertEntity("td_china", tdChinaDomain)
	checkErr(err, "插入数据")
}

func TestInsertEntityTag1(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()

	// 新增：使用entity
	tdChinaDomain := OrmChinaDomain1{
		Ts:      time.Now(),
		Name:    "TestInsertEntityTag1",
		Age:     19,
		Address: "浙江温州市",
	}

	tag := OrmChinaDomainTag{Station: "测试"}
	_, err := tdorm.InsertEntityWithTag("td_china", "td_demo1", tag, tdChinaDomain)
	checkErr(err, "插入数据")
}

func TestInsertOrmMapTag1(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()

	// 新增：使用entity
	tdChinaDomain := maps.NewSort().
		Put("ts", time.Now()).
		Put("name", "TestInsertOrmMapTag1").
		Put("age", 19).
		Put("address", "杭州临平")

	tag := maps.NewSort().Put("station", "浙江")
	_, err := tdorm.InsertWithTag("td_china_new", "td_demo1", tag, tdChinaDomain)
	checkErr(err, "插入数据")
}

// 使用带json标签的转换
func TestInsertEntity2(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()

	// 新增：使用entity
	tdChinaDomain := OrmChinaDomain2{
		TsJson:      time.Now(),
		NameJson:    "entity_json",
		AgeJson:     19,
		AddressJson: "浙江温州市",
	}
	_, err := tdorm.InsertEntity("td_china", tdChinaDomain)
	checkErr(err, "插入数据")
}

// 使用带自定义标签column的转换
func TestInsertEntity3(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()

	// 新增：使用entity
	tdChinaDomain := OrmChinaDomain3{
		Timestamp: time.Now(),
		Na:        "entity_column",
		Ag:        19,
		Add:       "浙江温州市",
	}
	_, err := tdorm.InsertEntity("td_china", tdChinaDomain)
	checkErr(err, "插入数据")
}

func TestBatchInsert(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()

	timeData1, _ := goleTime.ParseTimeYmdHmsS("2024-01-15 14:01:23.391")
	timeData2, _ := goleTime.ParseTimeYmdHmsS("2024-02-15 15:01:23.391")
	timeData3, _ := goleTime.ParseTimeYmdHmsS("2024-03-15 16:01:23.391")

	var insertMaps []*maps.GoleMap
	insertMaps = append(insertMaps, maps.New().SetSort(true).Put("ts", timeData1).Put("name", "大牛市1").Put("age", "18").Put("address", "浙江杭州市1"))
	insertMaps = append(insertMaps, maps.New().SetSort(true).Put("ts", timeData2).Put("name", "大牛市2").Put("age", "28").Put("address", "浙江杭州市2"))
	insertMaps = append(insertMaps, maps.New().SetSort(true).Put("ts", timeData3).Put("name", "大牛市3").Put("age", "38").Put("address", "浙江杭州市3"))

	num, err := tdorm.InsertBatch("td_china", insertMaps)
	assert.Equal(t, int64(3), num)
	checkErr(err, "批量插入异常")
}

func TestBatchInsertEntity(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()

	timeData1, _ := goleTime.ParseTimeYmdHmsS("2024-07-21 11:11:23.391")
	timeData2, _ := goleTime.ParseTimeYmdHmsS("2024-07-22 12:21:23.391")
	timeData3, _ := goleTime.ParseTimeYmdHmsS("2024-07-23 13:31:23.391")

	var insertEntities []interface{}
	insertEntities = append(insertEntities, OrmChinaDomain1{Ts: timeData1, Name: "TestBatchInsertEntity", Age: 19, Address: "浙江温州市"})
	insertEntities = append(insertEntities, OrmChinaDomain1{Ts: timeData2, Name: "TestBatchInsertEntity", Age: 19, Address: "浙江温州市"})
	insertEntities = append(insertEntities, OrmChinaDomain1{Ts: timeData3, Name: "TestBatchInsertEntity", Age: 19, Address: "浙江温州市"})

	num, err := tdorm.InsertEntityBatch("td_china", insertEntities)
	assert.Equal(t, int64(3), num)
	checkErr(err, "批量插入异常")
}

func TestBatchInsertWithTag(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()

	timeData1, _ := goleTime.ParseTimeYmdHmsS("2024-01-15 14:01:23.391")
	timeData2, _ := goleTime.ParseTimeYmdHmsS("2024-02-15 15:01:23.391")
	timeData3, _ := goleTime.ParseTimeYmdHmsS("2024-03-15 16:01:23.391")

	var insertMaps []*maps.GoleMap
	insertMaps = append(insertMaps, maps.New().SetSort(true).Put("ts", timeData1).Put("name", "TestBatchInsertWithTag").Put("age", "18").Put("address", "浙江杭州市1"))
	insertMaps = append(insertMaps, maps.New().SetSort(true).Put("ts", timeData2).Put("name", "TestBatchInsertWithTag").Put("age", "28").Put("address", "浙江杭州市2"))
	insertMaps = append(insertMaps, maps.New().SetSort(true).Put("ts", timeData3).Put("name", "TestBatchInsertWithTag").Put("age", "38").Put("address", "浙江杭州市3"))

	tagMap := maps.NewSort().Put("station", "浙江")
	num, err := tdorm.InsertBatchWithTag("td_china", "td_demo1", tagMap, insertMaps)
	assert.Equal(t, int64(3), num)
	checkErr(err, "批量插入异常")
}

func TestBatchInsertEntityWithTag(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()

	timeData1, _ := goleTime.ParseTimeYmdHmsS("2024-07-21 11:11:23.391")
	timeData2, _ := goleTime.ParseTimeYmdHmsS("2024-07-22 12:21:23.391")
	timeData3, _ := goleTime.ParseTimeYmdHmsS("2024-07-23 13:31:23.391")

	var insertEntities []interface{}
	insertEntities = append(insertEntities, OrmChinaDomain1{Ts: timeData1, Name: "TestBatchInsertEntityWithTag", Age: 19, Address: "浙江温州市"})
	insertEntities = append(insertEntities, OrmChinaDomain1{Ts: timeData2, Name: "TestBatchInsertEntityWithTag", Age: 19, Address: "浙江温州市"})
	insertEntities = append(insertEntities, OrmChinaDomain1{Ts: timeData3, Name: "TestBatchInsertEntityWithTag", Age: 19, Address: "浙江温州市"})

	tagEntity := OrmChinaDomainTag{Station: "测试"}
	num, err := tdorm.InsertEntityBatchWithTag("td_china_entity_new", "td_demo1", tagEntity, insertEntities)
	assert.Equal(t, int64(3), num)
	checkErr(err, "批量插入异常")
}

// 表不存在则自动创建
func TestInsertExeOrmMap1(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()

	// 新增：使用map
	insertMap := maps.New()
	insertMap.SetSort(true)
	insertMap.Put("ts", time.Now())
	insertMap.Put("name", "大牛市1")
	insertMap.Put("age", 28)
	insertMap.Put("address", "浙江杭州市")

	tagsMap := maps.NewSort()
	tagsMap.Put("station", "hangzhou1")

	_, err := tdorm.InsertWithTag("td_china2_new", "td_demo1", tagsMap, insertMap)
	checkErr(err, "插入异常")
}
