package tdorm

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/simonalong/gole/logger"
	"github.com/simonalong/gole/maps"
	"github.com/simonalong/tdorm/constants"
)

type TdHooks struct {
	hooks []TdHook
}

type TdHook interface {
	Before(c *TdHookContext) (*TdHookContext, error)
	After(c *TdHookContext) error
}

type TdHookContext struct {
	context.Context
	// 连接类型：connect_original，connect_rest，connect_websocket
	ConnectType byte
	// 库名
	DbName string
	// 开始时间
	Start time.Time
	// 运行类型：exe、query、insert、batchInsert
	RunType string
	Sql     string

	ArgsOfExe          []driver.Value
	ArgsOfQuery        []driver.Value
	FieldsArgsOfInsert *maps.GoleMap
	TagsArgsOfInsert   *maps.GoleMap
	ArgsOfBatchInsert  []*maps.GoleMap

	ResultOfExe               sql.Result
	ResultOfQueryOfDriverRows driver.Rows
	ResultOfQueryOfSqlRows    *sql.Rows
	// insert和batchInsert
	ResultOfInsert int64

	// 执行耗时
	ExecuteTime time.Duration
	Err         error
}

func (hookOfTd *TdHooks) add(hook TdHook) {
	hookOfTd.hooks = append(hookOfTd.hooks, hook)
}

func NewTdHooks() *TdHooks {
	return &TdHooks{
		hooks: []TdHook{},
	}
}

func (hookOfTd *TdHooks) preProcess(thc *TdHookContext) (*TdHookContext, error) {
	for _, hook := range hookOfTd.hooks {
		thc, err := hook.Before(thc)
		if err != nil {
			return thc, err
		}
	}
	return thc, nil
}

func (hookOfTd *TdHooks) afterProcess(thc *TdHookContext) {
	for _, hook := range hookOfTd.hooks {
		err := hook.After(thc)
		if err != nil && err.Error() != "[0x914] success" {
			logger.Errorf("hook 执行异常：%v", err.Error())
			return
		}
	}
}

func (thc *TdHookContext) GetConnectTypeStr() string {
	switch thc.ConnectType {
	case constants.ConnectOriginal:
		return "original"
	case constants.ConnectRestful:
		return "restful"
	case constants.ConnectWebsocket:
		return "websocket"
	}
	return ""
}

func (thc *TdHookContext) EndExe(result sql.Result, err error, args []driver.Value) *TdHookContext {
	thc.ArgsOfExe = args
	thc.ResultOfExe = result
	thc.Err = err
	thc.ExecuteTime = time.Since(thc.Start)
	return thc
}

func (thc *TdHookContext) EndExeAny(result sql.Result, err error, args []any) *TdHookContext {
	var argsFinal []driver.Value
	for _, arg := range args {
		argsFinal = append(argsFinal, arg)
	}
	thc.ArgsOfExe = argsFinal
	thc.ResultOfExe = result
	thc.Err = err
	thc.ExecuteTime = time.Since(thc.Start)
	return thc
}

func (thc *TdHookContext) EndQuery(rows driver.Rows, err error, args []driver.Value) *TdHookContext {
	thc.ArgsOfQuery = args
	thc.ResultOfQueryOfDriverRows = rows
	thc.Err = err
	thc.ExecuteTime = time.Since(thc.Start)
	return thc
}

func (thc *TdHookContext) EndQueryOfSqlRows(rows *sql.Rows, err error, args []driver.Value) *TdHookContext {
	thc.ArgsOfQuery = args
	thc.ResultOfQueryOfSqlRows = rows
	thc.Err = err
	thc.ExecuteTime = time.Since(thc.Start)
	return thc
}

func (thc *TdHookContext) EndInsert(result int64, err error, args *maps.GoleMap) *TdHookContext {
	thc.FieldsArgsOfInsert = args
	thc.ResultOfInsert = result
	thc.Err = err
	thc.ExecuteTime = time.Since(thc.Start)
	return thc
}

func (thc *TdHookContext) EndInsertFull(result int64, err error, tagsMap *maps.GoleMap, fieldsMap *maps.GoleMap) *TdHookContext {
	thc.TagsArgsOfInsert = tagsMap
	thc.FieldsArgsOfInsert = fieldsMap
	thc.ResultOfInsert = result
	thc.Err = err
	thc.ExecuteTime = time.Since(thc.Start)
	return thc
}

func (thc *TdHookContext) EndBatchInsert(result int64, err error, args []*maps.GoleMap) *TdHookContext {
	thc.ArgsOfBatchInsert = args
	thc.ResultOfInsert = result
	thc.Err = err
	thc.ExecuteTime = time.Since(thc.Start)
	return thc
}
