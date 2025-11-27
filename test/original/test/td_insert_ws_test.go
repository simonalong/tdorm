package test

import (
	"database/sql"
	"log"
	"testing"
	"time"

	_ "github.com/taosdata/driver-go/v3/taosWS"
)

var dbOfWs *sql.DB

func newDbWs() *sql.DB {
	if dbOfWs != nil {
		return dbOfWs
	}
	taos, err := sql.Open("taosWS", "root:taosdata@ws(localhost:6041)/td_orm")
	if err != nil {
		log.Fatal(err)
	}

	// 建库
	_, err = taos.Exec("create database if not exists td_orm")
	if err != nil {
		log.Fatalln("创建库失败：", err)
	}

	// 建超级表
	_, err = taos.Exec("create stable if not exists td_orm.td_demo1(ts timestamp, name nchar(32), age int, address nchar(128)) tags (station nchar(128))")
	checkErr(err, "建超级表失败")

	// 创建子表
	_, err = taos.Exec("create table if not exists td_orm.td_america using td_orm.td_demo1(`station`) tags(\"america\")")
	checkErr(err, "建超级表失败")

	// 创建子表
	_, err = taos.Exec("create table if not exists td_orm.td_china using td_orm.td_demo1(`station`) tags(\"china\")")
	checkErr(err, "建超级表失败")
	dbOfWs = taos
	return taos
}

func TestInsertWS(t *testing.T) {
	newDbWs()
	_, err := dbOfWs.Exec("insert into td_orm.td_china(`ts`,`name`,`age`,`address`) values('2024-07-12 12:43:14.321','zhou',13,'hangzhou')")
	checkErr(err, "插入失败")
}

func TestInsertWS2(t *testing.T) {
	newDbWs()

	stmt, err := dbOfWs.Prepare("insert into td_orm.td_china(`ts`,`name`,`age`,`address`) values(?,?,?,?)")
	if err != nil {
		log.Fatalf("预处理报错：%v", err)
	}
	// 正确
	//result, err := stmt.Exec("2021-07-11T12:43:14.321Z", "zhou", 13, "hangzhou")

	// 错误
	//result, err := stmt.Exec("2025-09-11 12:43:14.321", "zhou", 13, "hangzhou")

	// 正确
	result, err := stmt.Exec(time.Now(), "zhou", 13, "hangzhou")
	if err != nil {
		log.Fatalf("执行报错：%v", err)
	}
	log.Println("result：", result)
}

// 错误例子
func TestInsertWS3(t *testing.T) {
	newDbWs()

	// 不支持
	// syntax error near 'zhou,13,hangzhou)' (invalid data or symbol)
	result, err := dbOfWs.Exec("insert into td_orm.td_china(`ts`,`name`,`age`,`address`) values(?,?,?,?)", time.Now(), "zhou", 13, "hangzhou")
	if err != nil {
		log.Fatalf("执行报错：%v", err)
	}
	log.Println("result：", result)
}

func TestInsertWS4(t *testing.T) {
	newDbWs()

	stmt, err := dbOfWs.Prepare("insert into td_orm.td_china using td_orm.td_demo1 (`station`) tags ('测试') (`ts`,`name`,`age`,`address`) values (?,?,?,?)")
	if err != nil {
		log.Fatalf("预处理报错：%v", err)
	}

	// 正确
	result, err := stmt.Exec("2025-09-11T12:43:14.321Z", "'zhou'", 13, "'hangzhou'")
	if err != nil {
		log.Fatalf("执行报错：%v", err)
	}
	log.Println("result：", result)
}

// 错误
func TestInsertWS4_err(t *testing.T) {
	newDbWs()

	stmt, err := dbOfWs.Prepare("insert into td_orm.td_china using td_orm.td_demo1 (`station`) tags (?) (`ts`,`name`,`age`,`address`) values (?,?,?,?)")
	if err != nil {
		log.Fatalf("预处理报错：%v", err)
	}

	// 错误，标签不支持占位符
	result, err := stmt.Exec("'测试'", "2025-09-11T12:43:14.321Z", "'zhou'", 13, "'hangzhou'")
	if err != nil {
		log.Fatalf("执行报错：%v", err)
	}
	log.Println("result：", result)
}

// insert into td_orm.td_china using td_orm.td_demo1 (`station`) tags (?) (`ts`,`name`,`age`,`address`) values (?,?,?,?)
func TestInsertFullSqlWS(t *testing.T) {
	newDbWs()

	result, err := dbOfWs.Exec("insert into td_orm.td_china using td_orm.td_demo1 (`station`) tags ('浙江') (`ts`,`name`,`age`,`address`) values ('2024-01-15 14:01:23.391','大牛市1','18','浙江杭州市1') ('2024-02-15 15:01:23.391','大牛市2','28','浙江杭州市2') ('2024-03-15 16:01:23.391','大牛市3','38','浙江杭州市3')")
	if err != nil {
		log.Fatalf("执行报错：%v", err)
	}
	log.Println("result：", result)
}
