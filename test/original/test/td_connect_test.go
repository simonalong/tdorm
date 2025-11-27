package test

//
//import (
//	"database/sql"
//	"fmt"
//	"github.com/taosdata/driver-go/v3/af"
//	_ "github.com/taosdata/driver-go/v3/taosRestful"
//	_ "github.com/taosdata/driver-go/v3/taosSql"
//	_ "github.com/taosdata/driver-go/v3/taosWS"
//	"log"
//	"testing"
//)
//
//// 原生连接
//func TestOriginal(t *testing.T) {
//	var taosDSN = "root:taosdata@tcp(localhost:6030)/"
//	//var taosDSN = "root:taosdata@tcp(localhost:6030)/log"
//	taos, err := sql.Open("taosSql", taosDSN)
//	if err != nil {
//		log.Fatalln("failed to connect TDengine, err:", err)
//		return
//	}
//	fmt.Println("Connected")
//	defer taos.Close()
//}
//
//// rest连接
//func TestOriginalRest(t *testing.T) {
//	var taosDSN = "root:taosdata@http(localhost:6041)/"
//	//var taosDSN = "root:taosdata@http(localhost:6041)/log"
//	taos, err := sql.Open("taosRestful", taosDSN)
//	if err != nil {
//		log.Fatalln("failed to connect TDengine, err:", err)
//		return
//	}
//	fmt.Println("Connected")
//	defer taos.Close()
//}
//
//// rest连接
//func TestWebsocket(t *testing.T) {
//	var taosDSN = "root:taosdata@ws(localhost:6041)/"
//	//var taosDSN = "root:taosdata@http(localhost:6041)/log"
//	taos, err := sql.Open("taosWS", taosDSN)
//	if err != nil {
//		log.Fatalln("failed to connect TDengine, err:", err)
//		return
//	}
//	fmt.Println("Connected")
//	defer taos.Close()
//}
//
//// 高级原生连接
//func TestOriginalAdvance(t *testing.T) {
//	//conn, err := af.Open("localhost", "root", "taosdata", "log", 6030)
//	conn, err := af.Open("localhost", "root", "taosdata", "", 6030)
//	defer conn.Close()
//	if err != nil {
//		log.Fatalln("failed to connect, err:", err)
//	} else {
//		fmt.Println("connected")
//	}
//
//}
