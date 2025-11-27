package test

//
//import (
//	"fmt"
//	"github.com/taosdata/driver-go/v3/af"
//	"github.com/taosdata/driver-go/v3/common/param"
//	goleTime "github.com/simonalong/gole/time"
//	"log"
//	"testing"
//	"time"
//)
//
//func TestInsert(t *testing.T) {
//	conn, err := af.Open("localhost", "root", "taosdata", "", 6030)
//	if err != nil {
//		log.Fatalln("failed to connect TDengine, err:", err)
//		return
//	}
//	fmt.Println("Connected")
//	defer conn.Close()
//
//	p := param.NewParam(4).AddTimestamp(time.Now(), 0).AddNchar("zhou").AddInt(19).AddNchar("hangzhou")
//	_, err = conn.StmtExecute("insert into td_orm.td_china(`ts`,`name`,`age`,`address`) values(?,?,?,?)", p)
//	checkErr(err, "插入失败")
//}
//
//func TestBatchInsert1(t *testing.T) {
//	conn, err := af.Open("localhost", "root", "taosdata", "", 6030)
//	if err != nil {
//		log.Fatalln("failed to connect TDengine, err:", err)
//		return
//	}
//	fmt.Println("Connected")
//	defer conn.Close()
//
//	_, err = conn.Exec("insert into td_orm.td_china(`ts`,`name`,`age`,`address`) values('2024-07-12 12:43:14.321','zhou',13,'hangzhou')('2024-07-12 12:43:12.321','zhou',122,'hangzhou12')")
//	checkErr(err, "插入失败")
//}
//
//func TestBatchInsert2(t *testing.T) {
//	conn, err := af.Open("localhost", "root", "taosdata", "", 6030)
//	if err != nil {
//		log.Fatalln("failed to connect TDengine, err:", err)
//		return
//	}
//	fmt.Println("Connected")
//	defer conn.Close()
//
//	_, err = conn.Exec("insert into td_orm.td_china(`ts`,`name`,`age`,`address`) values" +
//		"('2024-07-12 12:43:14.321','zhou',13,'hangzhou')" +
//		"('2024-07-12 12:43:12.322','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.323','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.324','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.325','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.326','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.327','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.328','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.329','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.310','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.321','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.331','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.341','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.351','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.361','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.371','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.381','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.391','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.401','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.411','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.421','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.431','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.441','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.451','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.461','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.471','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.481','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.491','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.501','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.511','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.521','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:12.522','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:22.421','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:32.421','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:42.421','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:43:52.421','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:44:52.421','zhou',122,'hangzhou12')" +
//		"('2024-07-12 12:45:52.421','zhou',122,'hangzhou12')" +
//		"")
//	checkErr(err, "插入失败")
//}
//
//func TestBatchInsert3(t *testing.T) {
//	conn, err := af.Open("localhost", "root", "taosdata", "", 6030)
//	if err != nil {
//		log.Fatalln("failed to connect TDengine, err:", err)
//		return
//	}
//	fmt.Println("Connected")
//	defer conn.Close()
//
//	insertStmt := conn.InsertStmt()
//	err = insertStmt.Prepare("insert into td_orm.td_china(`ts`,`name`,`age`,`address`) values(?,?,?,?)")
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	bindType := param.NewColumnType(4).AddTimestamp().AddNchar(32).AddInt().AddNchar(128)
//
//	timeData1, _ := goleTime.ParseTimeYmdHmsS("2024-07-15 11:01:23.391")
//	params := []*param.Param{
//		param.NewParam(1).AddTimestamp(timeData1, 0),
//		param.NewParam(1).AddNchar("zhou1--1"),
//		param.NewParam(1).AddInt(19),
//		param.NewParam(1).AddNchar("hangzhou1"),
//	}
//	err = insertStmt.BindParam(params, bindType)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	timeData2, _ := goleTime.ParseTimeYmdHmsS("2024-07-12 13:01:23.391")
//	params2 := []*param.Param{
//		param.NewParam(1).AddTimestamp(timeData2, 0),
//		param.NewParam(1).AddNchar("zhou1--3"),
//		param.NewParam(1).AddInt(19),
//		param.NewParam(1).AddNchar("hangzhou1"),
//	}
//	err = insertStmt.BindParam(params2, bindType)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	timeData3, _ := goleTime.ParseTimeYmdHmsS("2024-07-12 13:01:23.391")
//	params3 := []*param.Param{
//		param.NewParam(1).AddTimestamp(timeData3, 0),
//		param.NewParam(1).AddNchar("zhou1--3"),
//		param.NewParam(1).AddInt(19),
//		param.NewParam(1).AddNchar("hangzhou1"),
//	}
//	err = insertStmt.BindParam(params3, bindType)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	//err = insertStmt.SetSubTableName("td_demo1")
//	//if err != nil {
//	//	t.Error(err)
//	//	return
//	//}
//
//	//err = insertStmt.SetTableName("td_china")
//	//if err != nil {
//	//	t.Error(err)
//	//	return
//	//}
//
//	err = insertStmt.AddBatch()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//	err = insertStmt.Execute()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//	err = insertStmt.Close()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//}
