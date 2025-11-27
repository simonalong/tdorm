package test

import (
	"fmt"
	"os"
	"time"

	tmqcommon "github.com/taosdata/driver-go/v3/common/tmq"
	"github.com/taosdata/driver-go/v3/ws/tmq"
)

// func TestShow(t *testing.T) {
func Show() {
	tdClient := newDbWs()

	fmt.Println("================== ShowApps ==================")
	apps, _ := tdClient.ShowApps()
	for _, app := range apps {
		// 示例：{"app_id":3630794682715340801,"current_req":0,"fetch_bytes":11413,"insert_bytes":23370,"insert_req":144,"insert_row":182,"insert_time":105411,"ip":"111.15.0.2","last_access":"2025-05-16T14:55:53.526+08:00","name":"taosadapter","pid":35,"query_time":70486203,"slow_query":12,"start_time":"2025-05-13T14:57:34.893+08:00","total_req":1131}
		fmt.Println(app.ToJson())
	}

	fmt.Println("================== ShowCluster ==================")
	cluster, _ := tdClient.ShowCluster()
	// {"create_time":"2025-05-13T13:51:27.575+08:00","id":142421125691033562,"name":"6c238b14-156c-4a96-99f5-b8f8bdc48803","uptime":250200,"version":"community"}
	fmt.Println(cluster.ToJson())

	fmt.Println("================== ShowClusterAlive ==================")
	alive, _ := tdClient.ShowClusterAlive()
	// 1
	fmt.Println(alive)

	fmt.Println("================== ShowConnections ==================")
	connects, _ := tdClient.ShowConnections()
	for _, connect := range connects {
		// 示例：
		// {"app":"taosadapter","conn_id":824010437,"end_point":"111.15.0.2:50660","last_access":"2025-05-16T15:04:29.406+08:00","login_time":"2025-05-16T15:04:29.406+08:00","pid":35,"user":"root"}
		// {"app":"taos","conn_id":4132423451,"end_point":"111.15.0.2:60512","last_access":"2025-05-16T15:04:29.166+08:00","login_time":"2025-05-14T13:33:31.975+08:00","pid":178,"user":"root"}
		fmt.Println(connect.ToJson())
	}

	_, err := tdClient.Exec("CREATE TOPIC IF NOT EXISTS topic_meters AS SELECT * FROM td_orm.td_demo1")
	if err != nil {
		panic(err)
	}

	go consume()

	time.Sleep(1 * time.Second)

	fmt.Println("================== ShowConsumers ==================")
	consumers, _ := tdClient.ShowConsumers()
	for _, consumer := range consumers {
		// 示例：
		// {"client_id":"test_tmq_client","consumer_group":"test","consumer_id":"0x32632d7f23504c9e","parameters":"tbname:1,commit:0,interval:5000ms,reset:latest","status":"rebalancing","subscribe_time":"2025-05-16T15:14:56.965+08:00","topics":"topic_meters","up_time":"2025-05-16T15:14:56.965+08:00"}
		fmt.Println(consumer.ToJson())
	}

	fmt.Println("================== ShowCreateDatabase ==================")
	// CREATE DATABASE `td_orm` BUFFER 256 CACHESIZE 1 CACHEMODEL 'none' COMP 2 DURATION 14400m WAL_FSYNC_PERIOD 3000 MAXROWS 4096 MINROWS 100 STT_TRIGGER 1 KEEP 5256000m,5256000m,5256000m PAGES 256 PAGESIZE 4 PRECISION 'ms' REPLICA 1 WAL_LEVEL 1 VGROUPS 2 SINGLE_STABLE 0 TABLE_PREFIX 0 TABLE_SUFFIX 0 TSDB_PAGESIZE 4 WAL_RETENTION_PERIOD 3600 WAL_RETENTION_SIZE 0 KEEP_TIME_OFFSET 0
	dbCreateSql, _ := tdClient.ShowCreateDatabase()
	fmt.Println(dbCreateSql)

	fmt.Println("================== ShowCreateStable ==================")
	// CREATE STABLE `td_demo1` (`ts` TIMESTAMP, `name` NCHAR(32), `age` INT, `address` NCHAR(128)) TAGS (`station` NCHAR(128))
	createStableSql, _ := tdClient.ShowCreateStable("td_demo1")
	fmt.Println(createStableSql)

	fmt.Println("================== ShowCreateTable ==================")
	// CREATE TABLE `td_china` USING `td_demo1` (`station`) TAGS ("china")
	createTableSql, _ := tdClient.ShowCreateTable("td_china")
	fmt.Println(createTableSql)

	fmt.Println("================== ShowDatabases ==================")
	// [td_orm]
	databasesOfUser, _ := tdClient.ShowDatabases("user")
	fmt.Println(databasesOfUser)

	fmt.Println("================== ShowDatabases ==================")
	// [information_schema performance_schema td_orm]
	databases, _ := tdClient.ShowDatabases("")
	fmt.Println(databases)

	fmt.Println("================== ShowDNodes ==================")
	dnodes, _ := tdClient.ShowDNodes()
	for _, dnode := range dnodes {
		// {"create_time":"2025-05-13T13:51:27.568+08:00","endpoint":"838c875052ca:6030","id":1,"note":"","reboot_time":"2025-05-13T14:56:48.748+08:00","status":"ready","support_vnodes":16,"vnodes":2}
		fmt.Println(dnode.ToJson())
	}

	fmt.Println("================== ShowLicences ==================")
	licences, _ := tdClient.ShowLicences()
	for _, licence := range licences {
		// {"cpu_cores":"unlimited","dnodes":"unlimited","expire_time":"unlimited","expired":"false","service_time":"limited","state":"ungranted","timeseries":"unlimited","version":"community"}
		fmt.Println(licence.ToJson())
	}

	fmt.Println("================== ShowGrants ==================")
	grants, _ := tdClient.ShowGrants()
	for _, grant := range grants {
		// {"cpu_cores":"unlimited","dnodes":"unlimited","expire_time":"unlimited","expired":"false","service_time":"limited","state":"ungranted","timeseries":"unlimited","version":"community"}
		fmt.Println(grant.ToJson())
	}

	fmt.Println("================== ShowIndexes ==================")
	indexes, _ := tdClient.ShowIndexes("td_demo1")
	for _, index := range indexes {
		// {"column_name":"station","create_time":"2025-05-13T13:56:34.193+08:00","db_name":"td_orm","index_name":"station_td_demo1","index_type":"tag_index","table_name":"td_demo1"}
		fmt.Println(index.ToJson())
	}

	fmt.Println("================== ShowLocalVariables ==================")
	localVariables, _ := tdClient.ShowLocalVariables()
	for _, localVariable := range localVariables {
		// {"name":"firstEp","scope":"both","value":"838c875052ca:6030"}
		//{"name":"secondEp","scope":"both","value":"838c875052ca:6030"}
		//{"name":"fqdn","scope":"server","value":"838c875052ca"}
		//{"name":"serverPort","scope":"server","value":"6030"}
		//{"name":"tempDir","scope":"both","value":"/tmp/"}
		//{"name":"minimalTmpDirGB","scope":"both","value":"1.000000"}
		//{"name":"shellActivityTimer","scope":"both","value":"3"}
		//{"name":"compressMsgSize","scope":"both","value":"-1"}
		//{"name":"queryPolicy","scope":"client","value":"1"}
		//{"name":"enableQueryHb","scope":"client","value":"1"}
		//{"name":"enableScience","scope":"client","value":"0"}
		//{"name":"querySmaOptimize","scope":"client","value":"0"}
		//{"name":"queryPlannerTrace","scope":"client","value":"0"}
		//{"name":"queryNodeChunkSize","scope":"client","value":"32768"}
		//{"name":"queryUseNodeAllocator","scope":"client","value":"1"}
		//{"name":"keepColumnName","scope":"client","value":"0"}
		//{"name":"smlChildTableName","scope":"client","value":""}
		//{"name":"smlAutoChildTableNameDelimiter","scope":"client","value":""}
		//{"name":"smlTagName","scope":"client","value":"_tag_null"}
		//{"name":"smlTsDefaultName","scope":"client","value":"_ts"}
		//{"name":"smlDot2Underline","scope":"client","value":"1"}
		//{"name":"maxShellConns","scope":"client","value":"50000"}
		//{"name":"maxInsertBatchRows","scope":"client","value":"1000000"}
		//{"name":"maxRetryWaitTime","scope":"both","value":"10000"}
		//{"name":"useAdapter","scope":"client","value":"1"}
		//{"name":"crashReporting","scope":"client","value":"1"}
		//{"name":"queryMaxConcurrentTables","scope":"client","value":"200"}
		//{"name":"metaCacheMaxSize","scope":"client","value":"-1"}
		//{"name":"slowLogThreshold","scope":"client","value":"3"}
		//{"name":"slowLogScope","scope":"client","value":""}
		//{"name":"numOfRpcThreads","scope":"both","value":"4"}
		//{"name":"numOfRpcSessions","scope":"both","value":"30000"}
		//{"name":"timeToGetAvailableConn","scope":"both","value":"500000"}
		//{"name":"keepAliveIdle","scope":"both","value":"60"}
		//{"name":"numOfTaskQueueThreads","scope":"client","value":"4"}
		//{"name":"experimental","scope":"both","value":"1"}
		//{"name":"monitor","scope":"server","value":"1"}
		//{"name":"monitorInterval","scope":"server","value":"30"}
		//{"name":"disableCount","scope":"client","value":"1"}
		//{"name":"configDir","scope":"both","value":"/etc/taos/"}
		//{"name":"scriptDir","scope":"both","value":"/etc/taos/"}
		//{"name":"logDir","scope":"both","value":"/var/log/taos/"}
		//{"name":"minimalLogDirGB","scope":"both","value":"1.000000"}
		//{"name":"numOfLogLines","scope":"both","value":"10000000"}
		//{"name":"asyncLog","scope":"both","value":"1"}
		//{"name":"logKeepDays","scope":"both","value":"0"}
		//{"name":"debugFlag","scope":"both","value":"0"}
		//{"name":"simDebugFlag","scope":"both","value":"131"}
		//{"name":"tmrDebugFlag","scope":"both","value":"131"}
		//{"name":"uDebugFlag","scope":"both","value":"131"}
		//{"name":"rpcDebugFlag","scope":"both","value":"131"}
		//{"name":"jniDebugFlag","scope":"client","value":"131"}
		//{"name":"qDebugFlag","scope":"both","value":"131"}
		//{"name":"cDebugFlag","scope":"client","value":"131"}
		//{"name":"timezone","scope":"both","value":"Asia/Shanghai (CST, +0800)"}
		//{"name":"locale","scope":"both","value":"en_US.UTF-8"}
		//{"name":"charset","scope":"both","value":"UTF-8"}
		//{"name":"assert","scope":"both","value":"1"}
		//{"name":"enableCoreFile","scope":"both","value":"1"}
		//{"name":"numOfCores","scope":"both","value":"8.000000"}
		//{"name":"ssd42","scope":"both","value":"0"}
		//{"name":"avx","scope":"both","value":"0"}
		//{"name":"avx2","scope":"both","value":"0"}
		//{"name":"fma","scope":"both","value":"0"}
		//{"name":"avx512","scope":"both","value":"0"}
		//{"name":"simdEnable","scope":"both","value":"0"}
		//{"name":"tagFilterCache","scope":"both","value":"0"}
		//{"name":"openMax","scope":"both","value":"1048576"}
		//{"name":"streamMax","scope":"both","value":"16"}
		//{"name":"pageSizeKB","scope":"both","value":"4"}
		//{"name":"totalMemoryKB","scope":"both","value":"8130824"}
		//{"name":"os sysname","scope":"both","value":"Linux"}
		//{"name":"os nodename","scope":"both","value":"838c875052ca"}
		//{"name":"os release","scope":"both","value":"6.6.22-linuxkit"}
		//{"name":"os version","scope":"both","value":"#1 SMP PREEMPT_DYNAMIC Fri Mar 29 12:23:08 UTC 2024"}
		//{"name":"os machine","scope":"both","value":"x86_64"}
		//{"name":"version","scope":"both","value":"3.2.3.0"}
		//{"name":"compatible_version","scope":"both","value":"3.0.0.0"}
		//{"name":"gitinfo","scope":"both","value":"e27fdcff254b7bd0e0ad2f825e0414da4c0f37dc"}
		//{"name":"buildinfo","scope":"both","value":"Built Linux-x64 at 2024-02-29 18:01:09 +0800"}
		fmt.Println(localVariable.ToJson())
	}

	fmt.Println("================== ShowMNodes ==================")
	mnodes, _ := tdClient.ShowMNodes()
	for _, mnode := range mnodes {
		// {"create_time":"2025-05-13T13:51:27.572+08:00","endpoint":"838c875052ca:6030","id":1,"role":"leader","role_time":"2025-05-13T14:56:49.024+08:00","status":"ready"}
		fmt.Println(mnode.ToJson())
	}

	fmt.Println("================== ShowStables ==================")
	// 模糊搜索
	stables, _ := tdClient.ShowStables("demo1")
	for _, stable := range stables {
		// td_demo1
		fmt.Println(stable)
	}

	fmt.Println("================== ShowTables ==================")
	tables, _ := tdClient.ShowTables("", "china")
	for _, table := range tables {
		//td_china_entity
		//td_china_batch
		//td_china
		//td_china_new
		//td_china2_new
		//td_china_entity_new
		fmt.Println(table)
	}

	fmt.Println("================== ShowTags ==================")
	tags, _ := tdClient.ShowTags("td_china")
	for _, tag := range tags {
		// {"db_name":"td_orm","stable_name":"td_demo1","table_name":"td_china","tag_name":"station","tag_type":"NCHAR(128)","tag_value":"china"}
		fmt.Println(tag.ToJson())
	}

	fmt.Println("================== ShowTopics ==================")
	topics, _ := tdClient.ShowTopics()
	for _, topic := range topics {
		// {"db_name":"td_orm","stable_name":"td_demo1","table_name":"td_china","tag_name":"station","tag_type":"NCHAR(128)","tag_value":"china"}
		fmt.Println(topic)
	}

}

func consume() {
	consumer, err := tmq.NewConsumer(&tmqcommon.ConfigMap{
		"group.id":            "test",
		"auto.offset.reset":   "latest",
		"td.connect.ip":       "127.0.0.1",
		"td.connect.user":     "root",
		"td.connect.pass":     "taosdata",
		"td.connect.port":     "6041",
		"client.id":           "test_tmq_client",
		"enable.auto.commit":  "false",
		"msg.with.table.name": "true",
		"ws.url":              "ws://127.0.0.1:6041",
	})
	if err != nil {
		panic(err)
	}
	err = consumer.Subscribe("topic_meters", nil)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 5; i++ {
		ev := consumer.Poll(500)
		if ev != nil {
			switch e := ev.(type) {
			case *tmqcommon.DataMessage:
				fmt.Printf("get message:%v\n", e)
			case tmqcommon.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				panic(e)
			}
			consumer.Commit()
		}
	}
}
