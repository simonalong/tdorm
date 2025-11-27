package tdorm

import (
	"database/sql"
	"database/sql/driver"

	"github.com/simonalong/gole/maps"
	"github.com/simonalong/tdorm/column"
	"github.com/simonalong/tdorm/condition"
)

type RestfulCmd struct {
	// 数据库名：我们这里一个库对应一个对象
	dbName string
	// 数据库链接：rest连接或者websocket连接
	tdSqlDb *sql.DB
	// 回调执行器
	hooks *TdHooks
}

func (cmd *RestfulCmd) Exec(sql string, args ...driver.Value) (driver.Result, error) {
	//TODO implement me
	panic("implement me")
}

func (cmd *RestfulCmd) Query(sql string, args ...driver.Value) ([]*maps.GoleMap, error) {
	//TODO implement me
	panic("implement me")
}

func (cmd *RestfulCmd) Select(selectColumns *column.Columns, fromClause string, whereCondition *condition.Condition) ([]*maps.GoleMap, error) {
	return cmd.Query(generateSelectSql(selectColumns, fromClause, whereCondition))
}

func (cmd *RestfulCmd) Insert(tableName string, toInsertMap *maps.GoleMap) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (cmd *RestfulCmd) InsertWithTag(tableName, stableName string, tagsMap *maps.GoleMap, toInsertMap *maps.GoleMap) (int, error) {
	panic("implement me")
}

func (cmd *RestfulCmd) InsertEntity(tableName string, entity interface{}) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (cmd *RestfulCmd) InsertEntityWithTag(tableName, stableName string, tagEntity interface{}, entity interface{}) (int, error) {
	panic("implement me")
}

func (cmd *RestfulCmd) InsertBatch(tableName string, toInsertMaps []*maps.GoleMap) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (cmd *RestfulCmd) InsertBatchWithTag(tableName, stableName string, tagsMap *maps.GoleMap, toInsertMaps []*maps.GoleMap) (int64, error) {
	panic("implement me")
}

func (cmd *RestfulCmd) InsertEntityBatch(tableName string, entities []interface{}) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (cmd *RestfulCmd) InsertEntityBatchWithTag(tableName, stableName string, tagEntity interface{}, entities []interface{}) (int64, error) {
	panic("implement me")
}

func (cmd *RestfulCmd) Delete(tableName string, query *condition.Condition) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (cmd *RestfulCmd) One(tableName string, columns *column.Columns, query *condition.Condition) (*maps.GoleMap, error) {
	//TODO implement me
	panic("implement me")
}

func (cmd *RestfulCmd) List(tableName string, columns *column.Columns, query *condition.Condition) ([]*maps.GoleMap, error) {
	//TODO implement me
	panic("implement me")
}

func (cmd *RestfulCmd) ListOfDistinct(tableName string, columns *column.Columns, query *condition.Condition) ([]*maps.GoleMap, error) {
	panic("implement me")
}

func (cmd *RestfulCmd) Value(tableName, columnName string, query *condition.Condition) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (cmd *RestfulCmd) Values(tableName, columnName string, query *condition.Condition) ([]interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (cmd *RestfulCmd) ValuesOfDistinct(tableName, columnName string, query *condition.Condition) ([]interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (cmd *RestfulCmd) Count(tableName string, query *condition.Condition) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (cmd *RestfulCmd) AddHook(hook TdHook) {
	if cmd.hooks == nil {
		return
	}
	cmd.hooks.add(hook)
}

func (cmd *RestfulCmd) GetHooks() *TdHooks {
	return cmd.hooks
}
