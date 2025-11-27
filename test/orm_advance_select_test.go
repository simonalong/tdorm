package test

import (
	"fmt"
	"testing"
)

func TestAdvanceSelect(t *testing.T) {
	//tdorm := newDbOriginal()
	tdorm := newDbWs()
	fmt.Println(tdorm.SelectDatabase())
	fmt.Println(tdorm.SelectClientVersion())
	fmt.Println(tdorm.SelectServerVersion())
	fmt.Println(tdorm.SelectServerStatus())
	fmt.Println(tdorm.SelectNow())
	fmt.Println(tdorm.SelectToday())
	fmt.Println(tdorm.SelectTimeZone())
	fmt.Println(tdorm.SelectCurrentUser())
	fmt.Println(tdorm.SelectUser())
}
