package test

import "testing"

func TestExeWS(t *testing.T) {
	newDbWs()
	_, err := dbOfWs.Exec("create database if not exists td_orm")
	checkErr(err, "插入失败")
}
