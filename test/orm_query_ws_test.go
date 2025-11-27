package test

import (
	"testing"

	goleTime "github.com/simonalong/gole/time"
)

func TestQuery(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()

	tim, _ := goleTime.ParseTimeYmdHmsS("2025-01-14 14:13:28.699")

	_, err := tdorm.Query(`select * from td_china where ts >= ?`, tim)
	if err != nil {
		t.Error(err)
	}

	// 报错异常
	//_, err = tdorm.QueryForOriginal(`select * from td_china where ts >= ?`, tim)
	//if err != nil {
	//	t.Error(err)
	//}
}
