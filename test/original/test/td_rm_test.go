package test

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
)

func TestDelete2(t *testing.T) {
	var taosDSN = "root:taosdata@tcp(localhost:6030)/"
	//var taosDSN = "root:taosdata@tcp(localhost:6030)/log"
	taos, err := sql.Open("taosSql", taosDSN)
	if err != nil {
		log.Fatalln("failed to connect TDengine, err:", err)
		return
	}
	defer taos.Close()

	_, err = taos.Exec("delete from td_orm.td_china where `ts` > ?", "now -2d")
	checkErr(err, "删除异常")
}
func checkErr(err error, prompt string) {
	if err != nil {
		fmt.Printf("%s\n", prompt)
		panic(err)
	}
}
