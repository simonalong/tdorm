package tdorm

import (
	"database/sql"
	"database/sql/driver"
	"strings"
	"time"

	"github.com/simonalong/gole/logger"
	"github.com/simonalong/gole/maps"
	"github.com/simonalong/gole/util"
	"github.com/simonalong/tdorm/column"
	"github.com/simonalong/tdorm/condition"
	"github.com/simonalong/tdorm/constants"
)

type WebsocketCmd struct {
	// 数据库名：我们这里一个库对应一个对象
	dbName string
	// 数据库链接：rest连接或者websocket连接
	tdSqlDb *sql.DB
	// 回调执行器
	hooks *TdHooks
}

func (cmd *WebsocketCmd) AddHook(hook TdHook) {
	if cmd.hooks == nil {
		return
	}
	cmd.hooks.add(hook)
}

func (cmd *WebsocketCmd) GetHooks() *TdHooks {
	return cmd.hooks
}

// Exec websocket 不支持占位符，因此args这个参数无法使用
func (cmd *WebsocketCmd) Exec(sql string, args ...driver.Value) (driver.Result, error) {
	// 前置处理
	hookContext, err := cmd.hooks.preProcess(cmd.NewTdHookContextOfWebsocket(sql, constants.EXE))
	if err != nil {
		logger.Group("tdorm").Errorf("执行前置处理异常【executeSql ==> %v】报错：%v", sql, err.Error())
		return nil, err
	}

	var argsAny []any
	for i := range args {
		argsAny = append(argsAny, args[i])
	}
	result, err := cmd.tdSqlDb.Exec(sql, argsAny...)

	if err != nil {
		logger.Group("tdorm").Errorf("执行execute异常，【executeSql ==> %v】【params ==> %v】报错：%v", sql, args, err.Error())
	} else {
		logger.Group("tdorm").Debugf("【executeSql ==> %v】【params ==> %v】", sql, args)
	}
	cmd.hooks.afterProcess(hookContext.EndExe(result, err, args))
	return result, err
}

func (cmd *WebsocketCmd) Query(sql string, args ...driver.Value) ([]*maps.GoleMap, error) {
	// 前置处理
	hookContext, err := cmd.hooks.preProcess(cmd.NewTdHookContextOfWebsocket(sql, constants.QUERY))
	if err != nil {
		logger.Group("tdorm").Errorf("执行前置处理异常【executeSql ==> %v】报错：%v", sql, err.Error())
		return nil, err
	}

	var parameters []interface{}
	for _, arg := range args {
		parameters = append(parameters, arg)
	}
	rows, err := cmd.tdSqlDb.Query(sql, parameters...)
	if err != nil {
		logger.Group("tdorm").Errorf("执行execute异常，【executeSql ==> %v】【params ==> %v】报错：%v", sql, args, err.Error())
	} else {
		logger.Group("tdorm").Debugf("【executeSql ==> %v】【params ==> %v】", sql, args)
	}
	cmd.hooks.afterProcess(hookContext.EndQueryOfSqlRows(rows, err, args))
	result := maps.FromSqlRows(rows)
	if err != nil {
		logger.Group("tdorm").Errorf("执行query异常，【querySql ==> %v】【params ==> %v】报错：%v", sql, args, err.Error())
	} else {
		logger.Group("tdorm").Debugf("【querySql ==> %v】【params ==> %v】【result count ==> %v】", sql, args, len(result))
	}
	return result, err
}

func (cmd *WebsocketCmd) Select(selectColumns *column.Columns, fromClause string, whereCondition *condition.Condition) ([]*maps.GoleMap, error) {
	return cmd.Query(generateSelectSql(selectColumns, fromClause, whereCondition))
}

func (cmd *WebsocketCmd) Insert(tableName string, toInsertMap *maps.GoleMap) (int, error) {
	newInsertMap := toInsertMap.SetSort(true).Clone()

	preSql := generateInsertSql(cmd.dbName, tableName, newInsertMap)

	stmt, err := cmd.tdSqlDb.Prepare(preSql)
	if err != nil {
		logger.Errorf("预执行Insert的sql=%v异常：%v", preSql, err)
		return 0, err
	}

	// 获取val列表
	var valList []interface{}
	keys := newInsertMap.Keys()
	for _, key := range keys {
		val, _ := newInsertMap.Get(key)
		valList = append(valList, val)
	}
	// 前置处理
	hookContext, err := cmd.hooks.preProcess(cmd.NewTdHookContextOfWebsocket(preSql, constants.INSERT))
	if err != nil {
		logger.Group("tdorm").Errorf("执行insert，前置处理异常【preSql ==> %v】报错：%v", preSql, err.Error())
		return 0, err
	}

	// 执行
	result, err := stmt.Exec(valList...)
	if err != nil {
		logger.Group("tdorm").Errorf("执行insert，Prepare【preSql ==> %v】【value ==> %v】报错：%v", preSql, valList, err.Error())
		cmd.hooks.afterProcess(hookContext.EndInsert(0, err, newInsertMap))
		return 0, err
	}
	affected, _ := result.RowsAffected()
	cmd.hooks.afterProcess(hookContext.EndInsert(affected, nil, newInsertMap))
	return util.ToInt(affected), nil
}

func (cmd *WebsocketCmd) InsertWithTag(tableName, stableName string, tagsMap *maps.GoleMap, toInsertMap *maps.GoleMap) (int, error) {
	newInsertMap := toInsertMap.Clone()
	newInsertMap.SetSort(true)

	preSql := generateInsertWithTagSql(cmd.dbName, tableName, stableName, tagsMap, newInsertMap)

	stmt, err := cmd.tdSqlDb.Prepare(preSql)
	if err != nil {
		logger.Errorf("预执行Insert的sql=%v异常：%v", preSql, err)
		return 0, err
	}

	// 获取val列表
	valList := newInsertMap.Values()

	logger.Group("tdorm").Debugf("执行insert，【preSql ==> %v】【tags ==> %v】【params ==> %v】", preSql, tagsMap.ToJsonOfSort(), newInsertMap.ToJsonOfSort())

	// 前置处理
	hookContext, err := cmd.hooks.preProcess(cmd.NewTdHookContextOfWebsocket(preSql, constants.SAVE))
	if err != nil {
		logger.Group("tdorm").Errorf("执行insert，前置处理异常【preSql ==> %v】【tags ==> %v】【params ==> %v】报错：%v", preSql, tagsMap.ToJsonOfSort(), newInsertMap.ToJsonOfSort(), err.Error())
		return 0, err
	}

	// 执行
	result, err := stmt.Exec(valList...)
	if err != nil && err.Error() != "[0x914] success" {
		logger.Group("tdorm").Errorf("执行insert，Prepare【preSql ==> %v】【tags ==> %v】【value ==> %v】报错：%v", preSql, tagsMap.ToJsonOfSort(), newInsertMap.ToJsonOfSort(), err.Error())
		cmd.hooks.afterProcess(hookContext.EndInsertFull(0, err, tagsMap, newInsertMap))
		return 0, err
	}
	affected, _ := result.RowsAffected()
	cmd.hooks.afterProcess(hookContext.EndInsertFull(affected, nil, tagsMap, newInsertMap))
	return util.ToInt(affected), nil
}

func (cmd *WebsocketCmd) InsertEntity(tableName string, entity interface{}) (int, error) {
	if entity == nil {
		logger.Warnf("insert entity 数据为空")
		return 0, nil
	}

	return cmd.Insert(tableName, maps.FromEntity(entity).SetSort(true))
}

func (cmd *WebsocketCmd) InsertEntityWithTag(tableName, stableName string, tagEntity interface{}, entity interface{}) (int, error) {
	if entity == nil {
		logger.Warnf("insert entity 数据为空")
		return 0, nil
	}

	return cmd.InsertWithTag(tableName, stableName, maps.FromEntity(tagEntity).SetSort(true), maps.FromEntity(entity).SetSort(true))
}

func (cmd *WebsocketCmd) InsertBatch(tableName string, toInsertMaps []*maps.GoleMap) (int64, error) {
	var toInsertMapsTem []*maps.GoleMap
	for _, insertMap := range toInsertMaps {
		toInsertMapsTem = append(toInsertMapsTem, insertMap.SetSort(true).Clone())
	}
	// 对insertMaps进行检查和分组
	groupInsertMaps := map[string][]*maps.GoleMap{}
	for _, insertMap := range toInsertMapsTem {
		groupKey := strings.Join(insertMap.Keys(), ",")
		if insertMaps, exist := groupInsertMaps[groupKey]; exist {
			insertMaps = append(insertMaps, insertMap)
			groupInsertMaps[groupKey] = insertMaps
		} else {
			var insertMaps []*maps.GoleMap
			insertMaps = append(insertMaps, insertMap)
			groupInsertMaps[groupKey] = insertMaps
		}
	}

	// 将所有分组进行批量插入
	if len(groupInsertMaps) == 1 {
		return cmd.doBatchInsert(tableName, toInsertMaps)
	} else {
		var finalCount int64
		for _, insertMaps := range groupInsertMaps {
			count, err := cmd.doBatchInsert(tableName, insertMaps)
			if err != nil {
				return count, err
			}
			finalCount += count
		}
		return finalCount, nil
	}
}

func (cmd *WebsocketCmd) InsertEntityBatch(tableName string, entities []interface{}) (int64, error) {
	if entities == nil || len(entities) == 0 {
		logger.Warnf("InsertEntityBatch 数据为空")
		return 0, nil
	}
	var insertMaps []*maps.GoleMap
	for _, entity := range entities {
		insertMaps = append(insertMaps, maps.FromEntity(entity).SetSort(true))
	}
	return cmd.InsertBatch(tableName, insertMaps)
}

func (cmd *WebsocketCmd) InsertBatchWithTag(tableName, stableName string, tagsMap *maps.GoleMap, toInsertMaps []*maps.GoleMap) (int64, error) {
	var toInsertMapsTem []*maps.GoleMap
	for _, insertMap := range toInsertMaps {
		toInsertMapsTem = append(toInsertMapsTem, insertMap.SetSort(true).Clone())
	}
	// 对insertMaps进行检查和分组
	groupInsertMaps := map[string][]*maps.GoleMap{}
	for _, insertMap := range toInsertMapsTem {
		groupKey := strings.Join(insertMap.Keys(), ",")
		if insertMaps, exist := groupInsertMaps[groupKey]; exist {
			insertMaps = append(insertMaps, insertMap)
			groupInsertMaps[groupKey] = insertMaps
		} else {
			var insertMaps []*maps.GoleMap
			insertMaps = append(insertMaps, insertMap)
			groupInsertMaps[groupKey] = insertMaps
		}
	}

	// 将所有分组进行批量插入
	if len(groupInsertMaps) == 1 {
		return cmd.doBatchInsertWithTag(tableName, stableName, tagsMap, toInsertMaps)
	} else {
		var finalCount int64
		for _, insertMaps := range groupInsertMaps {
			count, err := cmd.doBatchInsertWithTag(tableName, stableName, tagsMap, insertMaps)
			if err != nil {
				return count, err
			}
			finalCount += count
		}
		return finalCount, nil
	}
}

func (cmd *WebsocketCmd) InsertEntityBatchWithTag(tableName, stableName string, tagEntity interface{}, entities []interface{}) (int64, error) {
	if entities == nil || len(entities) == 0 {
		logger.Warnf("InsertEntityBatch 数据为空")
		return 0, nil
	}
	var insertMaps []*maps.GoleMap
	for _, entity := range entities {
		insertMaps = append(insertMaps, maps.FromEntity(entity).SetSort(true))
	}
	return cmd.InsertBatchWithTag(tableName, stableName, maps.FromEntity(tagEntity).SetSort(true), insertMaps)
}

func (cmd *WebsocketCmd) Delete(tableName string, query *condition.Condition) (int64, error) {
	preSql := generateDeleteSqlOfFull(cmd.dbName, tableName, query)
	logger.Group("tdorm").Debugf("执行 delete【Sql ==> %v】", preSql)

	// 前置处理
	hookContext, err := cmd.hooks.preProcess(cmd.NewTdHookContextOfWebsocket(preSql, constants.DELETE))
	if err != nil {
		logger.Group("tdorm").Errorf("执行 delete，前置处理异常【Sql ==> %v】报错：%v", preSql, err.Error())
		return 0, err
	}
	result, err := cmd.tdSqlDb.Exec(preSql)
	if err != nil {
		logger.Errorf("执行 delete【Sql ==> %v】异常：%v", preSql, err)
		cmd.hooks.afterProcess(hookContext.EndExe(result, err, nil))
		return 0, err
	} else {
		logger.Group("tdorm").Debugf("【executeSql ==> %v】", preSql)
	}
	cmd.hooks.afterProcess(hookContext.EndExe(result, err, nil))
	affected, _ := result.RowsAffected()
	return affected, nil
}

func (cmd *WebsocketCmd) One(tableName string, columns *column.Columns, query *condition.Condition) (*maps.GoleMap, error) {
	preSql := generateQueryOneSql(cmd.dbName, tableName, columns, query)
	args := generateQueryParams(query)
	var parameters []interface{}
	for _, arg := range args {
		parameters = append(parameters, arg)
	}

	// 前置处理
	hookContext, err := cmd.hooks.preProcess(cmd.NewTdHookContextOfWebsocket(preSql, constants.DELETE))
	if err != nil {
		logger.Group("tdorm").Errorf("执行 queryOne，前置处理异常【Sql ==> %v】报错：%v", preSql, err.Error())
		return nil, err
	}

	rows, err := cmd.tdSqlDb.Query(preSql, parameters...)
	if err != nil {
		logger.Errorf("【preSql ==> %v】【params ==> %v】异常：%v", preSql, args, err)
		cmd.hooks.afterProcess(hookContext.EndQueryOfSqlRows(rows, err, args))
		return nil, err
	} else {
		logger.Group("tdorm").Debugf("【executeSql ==> %v】【params ==> %v】", preSql, parameters)
	}
	resultMap := maps.FromSqlRows(rows)
	cmd.hooks.afterProcess(hookContext.EndQueryOfSqlRows(rows, nil, args))
	if len(resultMap) == 1 {
		return resultMap[0], nil
	}
	return maps.New(), nil
}

func (cmd *WebsocketCmd) List(tableName string, columns *column.Columns, query *condition.Condition) ([]*maps.GoleMap, error) {
	preSql := generateQueryListSql(cmd.dbName, tableName, columns, query)
	args := generateQueryParams(query)
	var parameters []interface{}
	for _, arg := range args {
		parameters = append(parameters, arg)
	}

	// 前置处理
	hookContext, err := cmd.hooks.preProcess(cmd.NewTdHookContextOfWebsocket(preSql, constants.DELETE))
	if err != nil {
		logger.Group("tdorm").Errorf("执行 queryList，前置处理异常【Sql ==> %v】报错：%v", preSql, err.Error())
		return nil, err
	}

	rows, err := cmd.tdSqlDb.Query(preSql, parameters...)
	if err != nil {
		logger.Errorf("【preSql ==> %v】【params ==> %v】异常：%v", preSql, args, err)
		cmd.hooks.afterProcess(hookContext.EndQueryOfSqlRows(rows, err, args))
		return nil, err
	} else {
		logger.Group("tdorm").Debugf("【executeSql ==> %v】【params ==> %v】", preSql, parameters)
	}
	resultMap := maps.FromSqlRows(rows)
	cmd.hooks.afterProcess(hookContext.EndQueryOfSqlRows(rows, nil, args))
	return resultMap, nil
}

func (cmd *WebsocketCmd) ListOfDistinct(tableName string, columns *column.Columns, query *condition.Condition) ([]*maps.GoleMap, error) {
	preSql := generateQueryListOfDistinctSql(cmd.dbName, tableName, columns, query)
	args := generateQueryParams(query)
	var parameters []interface{}
	for _, arg := range args {
		parameters = append(parameters, arg)
	}
	var resultMap []*maps.GoleMap
	defer func() {
		logger.Group("tdorm").Debugf("【preSql ==> %v】【params ==> %v】【result=%v】", preSql, parameters, resultMap)
	}()

	// 前置处理
	hookContext, err := cmd.hooks.preProcess(cmd.NewTdHookContextOfWebsocket(preSql, constants.DELETE))
	if err != nil {
		logger.Group("tdorm").Errorf("执行 queryList，前置处理异常【Sql ==> %v】报错：%v", preSql, err.Error())
		return nil, err
	}

	rows, err := cmd.tdSqlDb.Query(preSql, parameters...)
	if err != nil {
		logger.Errorf("【preSql ==> %v】【params ==> %v】异常：%v", preSql, args, err)
		cmd.hooks.afterProcess(hookContext.EndQueryOfSqlRows(rows, err, args))
		return nil, err
	} else {
		logger.Group("tdorm").Debugf("【executeSql ==> %v】【params ==> %v】", preSql, parameters)
	}
	resultMap = maps.FromSqlRows(rows)
	cmd.hooks.afterProcess(hookContext.EndQueryOfSqlRows(rows, nil, args))

	return resultMap, nil
}

func (cmd *WebsocketCmd) Value(tableName, columnName string, query *condition.Condition) (interface{}, error) {
	preQueryValueSql := generateQueryValueSql(cmd.dbName, tableName, columnName, query)
	params := generateQueryParams(query)
	var resultValue interface{}
	defer func() {
		logger.Group("tdorm").Debugf("【preQueryValueSql ==> %v】【params ==> %v】【result=%v】", preQueryValueSql, params, resultValue)
	}()
	result, err := cmd.Query(preQueryValueSql, params...)
	if err != nil {
		logger.Group("tdorm").Errorf("【executeSql ==> %v】【params ==> %v】", preQueryValueSql, params)
		return nil, err
	} else {
		logger.Group("tdorm").Debugf("【executeSql ==> %v】【params ==> %v】", preQueryValueSql, params)
	}
	if result != nil && len(result) > 0 {
		if val, exist := result[0].Get(columnName); exist {
			resultValue = val
			return val, nil
		} else {
			return nil, nil
		}
	}
	return nil, err
}

func (cmd *WebsocketCmd) Values(tableName, columnName string, query *condition.Condition) ([]interface{}, error) {
	preQueryValuesSql := generateQueryValuesSql(cmd.dbName, tableName, columnName, query)
	params := generateQueryParams(query)
	defer func() {
		logger.Group("tdorm").Debugf("【preQueryValuesSql ==> %v】【params ==> %v】", preQueryValuesSql, params)
	}()

	result, err := cmd.Query(preQueryValuesSql, params...)
	if err != nil {
		logger.Group("tdorm").Errorf("【executeSql ==> %v】【params ==> %v】", preQueryValuesSql, params)
		return nil, err
	} else {
		logger.Group("tdorm").Debugf("【executeSql ==> %v】【params ==> %v】", preQueryValuesSql, params)
	}
	if result != nil && len(result) > 0 {
		var vals []interface{}
		for _, ormMap := range result {
			if val, exist := ormMap.Get(columnName); exist {
				vals = append(vals, val)
			}
		}
		return vals, nil
	}
	return nil, err
}

func (cmd *WebsocketCmd) ValuesOfDistinct(tableName, columnName string, query *condition.Condition) ([]interface{}, error) {
	preQueryValuesOfDistinctSql := generateQueryValuesOfDistinctSql(cmd.dbName, tableName, columnName, query)
	params := generateQueryParams(query)
	defer func() {
		logger.Group("tdorm").Debugf("【preQueryValuesOfDistinctSql ==> %v】【params ==> %v】", preQueryValuesOfDistinctSql, params)
	}()

	result, err := cmd.Query(preQueryValuesOfDistinctSql, params...)
	if err != nil {
		logger.Group("tdorm").Errorf("【executeSql ==> %v】【params ==> %v】", preQueryValuesOfDistinctSql, params)
		return nil, err
	} else {
		logger.Group("tdorm").Debugf("【executeSql ==> %v】【params ==> %v】", preQueryValuesOfDistinctSql, params)
	}
	if result != nil && len(result) > 0 {
		var vals []interface{}
		for _, ormMap := range result {
			if val, exist := ormMap.Get(columnName); exist {
				vals = append(vals, val)
			}
		}
		return vals, nil
	}
	return nil, err
}

func (cmd *WebsocketCmd) Count(tableName string, query *condition.Condition) (int, error) {
	preQueryCountSql := generateQueryCountSql(cmd.dbName, tableName, query)
	params := generateQueryParams(query)
	result, err := cmd.Query(preQueryCountSql, params...)
	var countVal int
	defer func() {
		logger.Group("tdorm").Debugf("【preQueryCountSql ==> %v】【params ==> %v】【result=%v】", preQueryCountSql, params, countVal)
	}()
	if err != nil {
		logger.Group("tdorm").Errorf("【executeSql ==> %v】【params ==> %v】", preQueryCountSql, params)
		return 0, nil
	} else {
		logger.Group("tdorm").Debugf("【executeSql ==> %v】【params ==> %v】", preQueryCountSql, params)
	}
	if result != nil && len(result) > 0 {
		if val, exist := result[0].GetInt("cnt"); exist {
			countVal = val
			return val, nil
		} else {
			return 0, nil
		}
	}
	return 0, nil
}

func (cmd *WebsocketCmd) doBatchInsertWithTag(tableName, stableName string, tagsMap *maps.GoleMap, toInsertMaps []*maps.GoleMap) (int64, error) {
	if len(toInsertMaps) == 0 {
		return 0, nil
	}
	// 这里默认所有插入的数据的key都完全一样的，这里不再进行检查（为了效率），只取最开始的这个作为结构sql的输入，对于一些数据不均匀的集合来说需要业务保证数据集合的均匀性
	toInsertMap := toInsertMaps[0]
	fullSql := generateInsertBatchFullWithTagSql(cmd.dbName, tableName, stableName, tagsMap, toInsertMaps)

	defer func() {
		logger.Group("tdorm").Debugf("batch【InsertSql ==> %v】【params[0] ==> %v】", fullSql, toInsertMap.ToJsonOfSort())
	}()

	// before
	hookContext, err := cmd.hooks.preProcess(cmd.NewTdHookContextOfWebsocket(fullSql, constants.BATCH_INSERT))

	result, err := cmd.tdSqlDb.Exec(fullSql)
	if err != nil {
		logger.Errorf("执行BatchInsert，exe 报错，%v", err.Error())
		// after
		cmd.hooks.afterProcess(hookContext.EndBatchInsert(0, err, toInsertMaps))
		return 0, err
	}
	affected, err := result.RowsAffected()
	cmd.hooks.afterProcess(hookContext.EndBatchInsert(affected, err, toInsertMaps))
	return affected, err
}

func (cmd *WebsocketCmd) doBatchInsert(tableName string, toInsertMaps []*maps.GoleMap) (int64, error) {
	if len(toInsertMaps) == 0 {
		return 0, nil
	}
	// 这里默认所有插入的数据的key都完全一样的，这里不再进行检查（为了效率），只取最开始的这个作为结构sql的输入，对于一些数据不均匀的集合来说需要业务保证数据集合的均匀性
	toInsertMap := toInsertMaps[0]
	fullSql := generateInsertBatchFullSql(cmd.dbName, tableName, toInsertMaps)

	defer func() {
		logger.Group("tdorm").Debugf("batch【InsertSql ==> %v】【params[0] ==> %v】", fullSql, toInsertMap.ToJsonOfSort())
	}()

	// before
	hookContext, err := cmd.hooks.preProcess(cmd.NewTdHookContextOfWebsocket(fullSql, constants.BATCH_INSERT))

	result, err := cmd.tdSqlDb.Exec(fullSql)
	if err != nil {
		logger.Errorf("执行BatchInsert，exe 报错，%v", err.Error())
		// after
		cmd.hooks.afterProcess(hookContext.EndBatchInsert(0, err, toInsertMaps))
		return 0, err
	}
	affected, err := result.RowsAffected()
	cmd.hooks.afterProcess(hookContext.EndBatchInsert(affected, err, toInsertMaps))
	return affected, err
}

func (cmd *WebsocketCmd) NewTdHookContextOfWebsocket(sql, runType string) *TdHookContext {
	return &TdHookContext{
		ConnectType: constants.ConnectWebsocket,
		DbName:      cmd.dbName,
		Start:       time.Now(),
		Sql:         sql,
		RunType:     runType,
	}
}
