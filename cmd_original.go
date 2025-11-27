package tdorm

//
//import (
//	"database/sql/driver"
//	"errors"
//	"fmt"
//	"github.com/taosdata/driver-go/v3/af"
//	"github.com/taosdata/driver-go/v3/common/param"
//	"github.com/taosdata/driver-go/v3/types"
//	"github.com/simonalong/tdorm/column"
//	"github.com/simonalong/tdorm/condition"
//	"github.com/simonalong/tdorm/constants"
//	"github.com/simonalong/gole/logger"
//	"github.com/simonalong/gole/maps"
//	"github.com/simonalong/gole/util"
//	"strings"
//	"sync"
//	"time"
//)
//
//type OriginalCmd struct {
//	// 数据库链接：原生连接
//	tdConnect *af.Connector
//	// 超级表中的每个属性对应的类型关系，其中key：stableName value:key，fieldName/tagName；value：taosType
//	stableFieldTypeMap map[string]map[string]TdFieldMeta
//	// 子表对应的超级表map
//	tableStableMap map[string]string
//	// 数据库客户端
//	tdClient *TdClient
//	// 数据库名：我们这里一个库对应一个对象
//	dbName string
//	// 数据库是否已经加载完成
//	dbLoaded bool
//	// 数据库数据处理器
//	dbLoadLocker sync.Locker
//	// 回调执行器
//	hooks *TdHooks
//}
//
//type TdFieldMeta struct {
//	// taos类型：例如：types.TaosBoolType
//	ColType interface{}
//	// 类型长度
//	ColLen int
//}
//
//func (cmd *OriginalCmd) GetHooks() *TdHooks {
//	return cmd.hooks
//}
//
//func (cmd *OriginalCmd) AddHook(hook TdHook) {
//	if cmd.hooks == nil {
//		return
//	}
//	cmd.hooks.add(hook)
//}
//
//func NewOriginalCmd(conn *af.Connector, dbName string) *OriginalCmd {
//	cmd := &OriginalCmd{
//		tdConnect:          conn,
//		dbName:             dbName,
//		stableFieldTypeMap: make(map[string]map[string]TdFieldMeta),
//		dbLoaded:           false,
//		dbLoadLocker:       &sync.Mutex{},
//		hooks:              NewTdHooks(),
//	}
//	if dbName != "" {
//		cmd.LoadMetaOfDb()
//		cmd.dbLoaded = true
//	}
//	return cmd
//}
//
//// LoadMetaOfDb 加载元数据。在表结构变化或者指定数据库的时候，请调用该函数
//func (cmd *OriginalCmd) LoadMetaOfDb() {
//	cmd.loadStableFieldOfDb()
//	cmd.loadTableStableMapOfDb()
//}
//
//func (cmd *OriginalCmd) loadStableFieldOfDb() {
//	conn := cmd.tdConnect
//	dbName := cmd.dbName
//	cmd.stableFieldTypeMap = make(map[string]map[string]TdFieldMeta)
//	rows, err := conn.Query(fmt.Sprintf("select `table_name`,`col_name`,`col_type`,`col_length` from information_schema.ins_columns where `db_name` = '%s' and (`table_type` = \"SUPER_TABLE\" or `table_type`=\"NORMAL_TABLE\")", dbName))
//	if err != nil {
//		logger.Errorf("获取数据库的元数据异常：%v", err.Error())
//		return
//	}
//	dest := make([]driver.Value, 4)
//	fieldMap := map[string]TdFieldMeta{}
//	var lastStableName string
//	for rows.Next(dest) == nil {
//		currentStableName := util.ToString(dest[0])
//		if lastStableName != "" && lastStableName != currentStableName {
//			cmd.stableFieldTypeMap[lastStableName] = fieldMap
//			fieldMap = make(map[string]TdFieldMeta)
//		}
//		lastStableName = currentStableName
//		fieldMap[util.ToString(dest[1])] = TdFieldMeta{
//			ColType: tdengineColTypeToTaosType(util.ToString(dest[2])),
//			ColLen:  util.ToInt(dest[3]),
//		}
//	}
//	cmd.stableFieldTypeMap[lastStableName] = fieldMap
//	fieldMap = make(map[string]TdFieldMeta)
//}
//
//func (cmd *OriginalCmd) loadTableStableMapOfDb() {
//	conn := cmd.tdConnect
//	dbName := cmd.dbName
//	cmd.tableStableMap = make(map[string]string)
//	rows, err := conn.Query(fmt.Sprintf("select `table_name`, `stable_name` from information_schema.ins_tables where db_name= '%s'", dbName))
//	if err != nil {
//		logger.Errorf("获取数据库的元数据子表和超表关系异常：%v", err.Error())
//		return
//	}
//	dest := make([]driver.Value, 2)
//
//	tableStableMap := map[string]string{}
//	for rows.Next(dest) == nil {
//		tableName := util.ToString(dest[0])
//		stableName := util.ToString(dest[1])
//		tableStableMap[tableName] = stableName
//	}
//	cmd.tableStableMap = tableStableMap
//}
//
//func (cmd *OriginalCmd) checkDb(tableName string) {
//	if cmd.dbLoaded {
//		return
//	}
//	if !strings.Contains(tableName, ".") {
//		return
//	}
//
//	cmd.dbLoadLocker.Lock()
//	defer cmd.dbLoadLocker.Unlock()
//
//	if cmd.dbLoaded {
//		return
//	}
//	names := strings.SplitN(tableName, ".", 2)
//	cmd.dbName = names[0]
//	cmd.LoadMetaOfDb()
//	cmd.dbLoaded = true
//}
//
//func (cmd *OriginalCmd) getTableFiledTaosTypeMap(tableNameOfOriginal string) map[string]TdFieldMeta {
//	tableName := tableNameOfOriginal
//	if strings.Contains(tableNameOfOriginal, ".") {
//		tableName = strings.SplitN(tableNameOfOriginal, ".", 2)[1]
//	}
//	if _, exit := cmd.tableStableMap[tableName]; !exit {
//		cmd.loadTableStableMapOfDb()
//	}
//
//	var finalTableName string
//	stableName, exit := cmd.tableStableMap[tableName]
//	if !exit {
//		logger.Errorf("当前表不存在, %v", tableName)
//		return nil
//	}
//	if stableName == "" {
//		finalTableName = tableName
//	} else {
//		finalTableName = stableName
//	}
//
//	fieldMap := cmd.stableFieldTypeMap[finalTableName]
//	return fieldMap
//}
//
//func (cmd *OriginalCmd) Query(sql string, args ...driver.Value) ([]*maps.GoleMap, error) {
//	// before
//	hookContext, err := cmd.hooks.preProcess(cmd.NewTdHookContextOfOriginal(sql, constants.QUERY))
//
//	// 执行
//	rows, err := cmd.tdConnect.Query(sql, args...)
//
//	// after
//	cmd.hooks.afterProcess(hookContext.EndQuery(rows, err, args))
//	result := maps.FromRows(rows)
//	if err != nil {
//		logger.Group("tdorm").Errorf("执行query异常，【querySql ==> %v】【params ==> %v】报错：%v", sql, args, err.Error())
//	} else {
//		logger.Group("tdorm").Debugf("【querySql ==> %v】【params ==> %v】【result count ==> %v】", sql, args, len(result))
//	}
//	return result, err
//}
//
//func (cmd *OriginalCmd) Insert(tableName string, toInsertMap *maps.GoleMap) (int, error) {
//	newInsertMap := toInsertMap.Clone()
//
//	cmd.checkDb(tableName)
//	newInsertMap.SetSort(true)
//	insertStmt := cmd.tdConnect.InsertStmt()
//	defer func() {
//		err := insertStmt.Close()
//		if err != nil {
//			logger.Warnf("执行insert报错")
//		}
//	}()
//
//	preInsertSql := generateInsertSql(cmd.dbName, tableName, newInsertMap)
//	bindType := generateBindParamsType(cmd, tableName, newInsertMap)
//	params, unFindKeys := generateParams(cmd, tableName, newInsertMap)
//
//	var values []driver.Value
//	for _, p := range params {
//		values = append(values, p.GetValues())
//	}
//	logger.Group("tdorm").Debugf("【preInsertSql ==> %v】【params ==> %v】", preInsertSql, newInsertMap.CloneExceptKeys(unFindKeys).ToJsonOfSort())
//
//	// before
//	hookContext, err := cmd.hooks.preProcess(cmd.NewTdHookContextOfOriginal(preInsertSql, constants.INSERT))
//
//	// 执行
//	err = insertStmt.Prepare(preInsertSql)
//	if err != nil {
//		logger.Errorf("执行insert，Prepare【preInsertSql ==> %v】【params ==> %v】报错：%v", preInsertSql, newInsertMap.CloneExceptKeys(unFindKeys).ToJsonOfSort(), err.Error())
//		// after
//		cmd.hooks.afterProcess(hookContext.EndInsert(0, err, newInsertMap))
//		return 0, err
//	}
//
//	err = insertStmt.BindParam(params, bindType)
//	if err != nil {
//		logger.Errorf("执行insert，BindParam【preInsertSql ==> %v】【params ==> %v】报错：%v", preInsertSql, newInsertMap.CloneExceptKeys(unFindKeys).ToJsonOfSort(), err.Error())
//		// after
//		cmd.hooks.afterProcess(hookContext.EndInsert(0, err, newInsertMap))
//		return 0, err
//	}
//
//	err = insertStmt.AddBatch()
//	if err != nil {
//		logger.Errorf("执行insert，AddBatch【preInsertSql ==> %v】【params ==> %v】报错，%v", preInsertSql, newInsertMap.CloneExceptKeys(unFindKeys).ToJsonOfSort(), err.Error())
//		// after
//		cmd.hooks.afterProcess(hookContext.EndInsert(0, err, newInsertMap))
//		return 0, err
//	}
//
//	err = insertStmt.Execute()
//	if err != nil {
//		logger.Errorf("执行insert，Execute【preInsertSql ==> %v】【params ==> %v】报错，%v", preInsertSql, newInsertMap.CloneExceptKeys(unFindKeys).ToJsonOfSort(), err.Error())
//		// after
//		cmd.hooks.afterProcess(hookContext.EndInsert(0, err, newInsertMap))
//		return 0, err
//	}
//	result := insertStmt.GetAffectedRows()
//	// after
//	cmd.hooks.afterProcess(hookContext.EndInsert(util.ToInt64(result), nil, newInsertMap))
//	return result, nil
//}
//
//func (cmd *OriginalCmd) InsertWithTag(tableName, stableName string, tagsMap *maps.GoleMap, toInsertMap *maps.GoleMap) (int, error) {
//	newInsertMap := toInsertMap.Clone()
//
//	cmd.checkDb(tableName)
//	newInsertMap.SetSort(true)
//
//	insertStmt := cmd.tdConnect.InsertStmt()
//	defer func() {
//		err := insertStmt.Close()
//		if err != nil {
//			logger.Warnf("执行insert报错")
//		}
//	}()
//
//	preInsertSql := generateInsertWithTagSql(cmd.dbName, tableName, stableName, tagsMap, newInsertMap)
//	bindType := generateBindParamsTagsType(cmd, tableName, tagsMap, newInsertMap)
//	params, unFindKeys := generateParamsTags(cmd, tableName, tagsMap, newInsertMap)
//
//	var values []driver.Value
//	for _, p := range params {
//		values = append(values, p.GetValues())
//	}
//	logger.Group("tdorm").Debugf("【preInsertSql ==> %v】【tags ==> %v】【params ==> %v】", preInsertSql, tagsMap.ToJsonOfSort(), newInsertMap.CloneExceptKeys(unFindKeys).ToJsonOfSort())
//
//	// before
//	hookContext, err := cmd.hooks.preProcess(cmd.NewTdHookContextOfOriginal(preInsertSql, constants.INSERT))
//
//	// 执行
//	err = insertStmt.Prepare(preInsertSql)
//	if err != nil {
//		logger.Errorf("执行insert，Prepare【preInsertSql ==> %v】【tags ==> %v】【params ==> %v】报错：%v", preInsertSql, tagsMap.ToJsonOfSort(), newInsertMap.CloneExceptKeys(unFindKeys).ToJsonOfSort(), err.Error())
//		// after
//		cmd.hooks.afterProcess(hookContext.EndInsert(0, err, newInsertMap))
//		return 0, err
//	}
//
//	err = insertStmt.BindParam(params, bindType)
//	if err != nil {
//		logger.Errorf("执行insert，BindParam【preInsertSql ==> %v】【tags ==> %v】【params ==> %v】报错：%v", preInsertSql, tagsMap.ToJsonOfSort(), newInsertMap.CloneExceptKeys(unFindKeys).ToJsonOfSort(), err.Error())
//		// after
//		cmd.hooks.afterProcess(hookContext.EndInsert(0, err, newInsertMap))
//		return 0, err
//	}
//
//	err = insertStmt.AddBatch()
//	if err != nil {
//		logger.Errorf("执行insert，AddBatch【preInsertSql ==> %v】【tags ==> %v】【params ==> %v】报错，%v", preInsertSql, tagsMap.ToJsonOfSort(), newInsertMap.CloneExceptKeys(unFindKeys).ToJsonOfSort(), err.Error())
//		// after
//		cmd.hooks.afterProcess(hookContext.EndInsert(0, err, newInsertMap))
//		return 0, err
//	}
//
//	err = insertStmt.Execute()
//	if err != nil {
//		logger.Errorf("执行insert，Execute【preInsertSql ==> %v】【tags ==> %v】【params ==> %v】报错，%v", preInsertSql, tagsMap.ToJsonOfSort(), newInsertMap.CloneExceptKeys(unFindKeys).ToJsonOfSort(), err.Error())
//		// after
//		cmd.hooks.afterProcess(hookContext.EndInsert(0, err, newInsertMap))
//		return 0, err
//	}
//	result := insertStmt.GetAffectedRows()
//	// after
//	cmd.hooks.afterProcess(hookContext.EndInsert(util.ToInt64(result), err, newInsertMap))
//	return result, nil
//}
//
//func (cmd *OriginalCmd) InsertEntity(tableName string, entity interface{}) (int, error) {
//	return cmd.Insert(tableName, maps.FromEntity(entity))
//}
//
//func (cmd *OriginalCmd) InsertEntityWithTag(tableName, stableName string, tagEntity interface{}, entity interface{}) (int, error) {
//	return cmd.InsertWithTag(tableName, stableName, maps.FromEntity(tagEntity), maps.FromEntity(entity))
//}
//
//func (cmd *OriginalCmd) InsertBatch(tableName string, toInsertMaps []*maps.GoleMap) (int64, error) {
//	// 对insertMaps进行检查和分组
//	groupInsertMaps := map[string][]*maps.GoleMap{}
//	for _, insertMap := range toInsertMaps {
//		groupKey := strings.Join(insertMap.Keys(), ",")
//		if insertMaps, exist := groupInsertMaps[groupKey]; exist {
//			insertMaps = append(insertMaps, insertMap)
//			groupInsertMaps[groupKey] = insertMaps
//		} else {
//			var insertMaps []*maps.GoleMap
//			insertMaps = append(insertMaps, insertMap)
//			groupInsertMaps[groupKey] = insertMaps
//		}
//	}
//
//	// 将所有分组进行批量插入
//	if len(groupInsertMaps) == 1 {
//		return cmd.doBatchInsert(tableName, toInsertMaps)
//	} else {
//		var finalCount int64
//		for _, insertMaps := range groupInsertMaps {
//			count, err := cmd.doBatchInsert(tableName, insertMaps)
//			if err != nil {
//				return count, err
//			}
//			finalCount += count
//		}
//		return finalCount, nil
//	}
//}
//
//func (cmd *OriginalCmd) InsertEntityBatch(tableName string, entities []interface{}) (int64, error) {
//	if entities == nil || len(entities) == 0 {
//		logger.Warnf("InsertEntityBatch 数据为空")
//		return 0, nil
//	}
//	var insertMaps []*maps.GoleMap
//	for _, entity := range entities {
//		insertMaps = append(insertMaps, maps.FromEntity(entity))
//	}
//	return cmd.InsertBatch(tableName, insertMaps)
//}
//
//func (cmd *OriginalCmd) InsertBatchWithTag(tableName, stableName string, tagsMap *maps.GoleMap, toInsertMaps []*maps.GoleMap) (int64, error) {
//	logger.Error("original 暂时不支持 InsertBatchWithTag")
//	return 0, errors.New("original 暂时不支持 InsertBatchWithTag")
//}
//
//func (cmd *OriginalCmd) InsertEntityBatchWithTag(tableName, stableName string, tagEntity interface{}, entities []interface{}) (int64, error) {
//	logger.Error("original 暂时不支持 InsertEntityBatchWithTag")
//	return 0, errors.New("original 暂时不支持 InsertEntityBatchWithTag")
//}
//
//func (cmd *OriginalCmd) Delete(tableName string, query *condition.Condition) (int64, error) {
//	preDeleteSql := generateDeleteSql(cmd.dbName, tableName, query)
//	params := generateQueryParams(query)
//	logger.Group("tdorm").Debugf("【preDeleteSql ==> %v】【params ==> %v】", preDeleteSql, params)
//	rlt, err := cmd.Exec(preDeleteSql, params...)
//	if err != nil {
//		logger.Errorf("数据删除异常：%v", err.Error())
//		return 0, err
//	}
//	return rlt.RowsAffected()
//}
//
//func (cmd *OriginalCmd) One(tableName string, columns *column.Columns, query *condition.Condition) (*maps.GoleMap, error) {
//	preQueryOneSql := generateQueryOneSql(cmd.dbName, tableName, columns, query)
//	params := generateQueryParams(query)
//	result, err := cmd.Query(preQueryOneSql, params...)
//
//	resultOneVal := maps.New()
//	defer func() {
//		logger.Group("tdorm").Debugf("【preQueryOneSql ==> %v】【params ==> %v】【result=%v】", preQueryOneSql, params, resultOneVal.ToJsonOfSort())
//	}()
//	if err != nil {
//		return nil, err
//	}
//	if result != nil && len(result) > 0 {
//		resultOneVal = result[0]
//		return result[0], err
//	}
//	return maps.New(), err
//}
//
//func (cmd *OriginalCmd) List(tableName string, columns *column.Columns, query *condition.Condition) ([]*maps.GoleMap, error) {
//	preQueryListSql := generateQueryListSql(cmd.dbName, tableName, columns, query)
//	params := generateQueryParams(query)
//	logger.Group("tdorm").Debugf("【preQueryListSql ==> %v】【params ==> %v】", preQueryListSql, params)
//	result, err := cmd.Query(preQueryListSql, params...)
//	if err != nil {
//		return nil, err
//	}
//	if result != nil && len(result) > 0 {
//		return result, err
//	}
//	return result, nil
//}
//
//func (cmd *OriginalCmd) ListOfDistinct(tableName string, columns *column.Columns, query *condition.Condition) ([]*maps.GoleMap, error) {
//	preQueryListOfDistinctSql := generateQueryListOfDistinctSql(cmd.dbName, tableName, columns, query)
//	params := generateQueryParams(query)
//	logger.Group("tdorm").Debugf("【preQueryListOfDistinctSql ==> %v】【params ==> %v】", preQueryListOfDistinctSql, params)
//	result, err := cmd.Query(preQueryListOfDistinctSql, params...)
//	if err != nil {
//		return nil, err
//	}
//	if result != nil && len(result) > 0 {
//		return result, err
//	}
//	return result, nil
//}
//
//func (cmd *OriginalCmd) Value(tableName, columnName string, query *condition.Condition) (interface{}, error) {
//	preQueryValueSql := generateQueryValueSql(cmd.dbName, tableName, columnName, query)
//	params := generateQueryParams(query)
//	var resultValue interface{}
//	defer func() {
//		logger.Group("tdorm").Debugf("【preQueryValueSql ==> %v】【params ==> %v】【result=%v】", preQueryValueSql, params, resultValue)
//	}()
//	result, err := cmd.Query(preQueryValueSql, params...)
//	if err != nil {
//		return nil, err
//	}
//	if result != nil && len(result) > 0 {
//		if val, exist := result[0].Get(columnName); exist {
//			resultValue = val
//			return val, nil
//		} else {
//			return nil, nil
//		}
//	}
//	return nil, err
//}
//
//func (cmd *OriginalCmd) Values(tableName, columnName string, query *condition.Condition) ([]interface{}, error) {
//	preQueryValuesSql := generateQueryValuesSql(cmd.dbName, tableName, columnName, query)
//	params := generateQueryParams(query)
//	defer func() {
//		logger.Group("tdorm").Debugf("【preQueryValuesSql ==> %v】【params ==> %v】", preQueryValuesSql, params)
//	}()
//
//	result, err := cmd.Query(preQueryValuesSql, params...)
//	if err != nil {
//		return nil, err
//	}
//	if result != nil && len(result) > 0 {
//		var vals []interface{}
//		for _, ormMap := range result {
//			if val, exist := ormMap.Get(columnName); exist {
//				vals = append(vals, val)
//			}
//		}
//		return vals, nil
//	}
//	return nil, err
//}
//
//func (cmd *OriginalCmd) ValuesOfDistinct(tableName, columnName string, query *condition.Condition) ([]interface{}, error) {
//	preQueryValuesOfDistinctSql := generateQueryValuesOfDistinctSql(cmd.dbName, tableName, columnName, query)
//	params := generateQueryParams(query)
//	logger.Group("tdorm").Debugf("【preQueryValuesOfDistinctSql ==> %v】【params ==> %v】", preQueryValuesOfDistinctSql, params)
//	result, err := cmd.Query(preQueryValuesOfDistinctSql, params...)
//	if err != nil {
//		return nil, err
//	}
//	if result != nil && len(result) > 0 {
//		var vals []interface{}
//		for _, ormMap := range result {
//			if val, exist := ormMap.Get(columnName); exist {
//				vals = append(vals, val)
//			}
//		}
//		return vals, nil
//	}
//	return nil, err
//}
//
//func (cmd *OriginalCmd) Count(tableName string, query *condition.Condition) (int, error) {
//	preQueryCountSql := generateQueryCountSql(cmd.dbName, tableName, query)
//	params := generateQueryParams(query)
//	result, err := cmd.Query(preQueryCountSql, params...)
//	var countVal int
//	defer func() {
//		logger.Group("tdorm").Debugf("【preQueryCountSql ==> %v】【params ==> %v】【result=%v】", preQueryCountSql, params, countVal)
//	}()
//	if err != nil {
//		return 0, nil
//	}
//	if result != nil && len(result) > 0 {
//		if val, exist := result[0].GetInt("cnt"); exist {
//			countVal = val
//			return val, nil
//		} else {
//			return 0, nil
//		}
//	}
//	return 0, nil
//}
//
//func (cmd *OriginalCmd) Exec(sql string, args ...driver.Value) (driver.Result, error) {
//	// before
//	hookContext, err := cmd.hooks.preProcess(cmd.NewTdHookContextOfOriginal(sql, constants.EXE))
//
//	// 执行
//	result, err := cmd.tdConnect.Exec(sql, args...)
//	if err != nil {
//		logger.Group("tdorm").Errorf("执行execute异常，【executeSql ==> %v】【params ==> %v】报错：%v", sql, args, err.Error())
//	} else {
//		logger.Group("tdorm").Debugf("【executeSql ==> %v】【params ==> %v】", sql, args)
//	}
//
//	// after
//	cmd.hooks.afterProcess(hookContext.EndExe(result, err, args))
//	return result, err
//}
//
//func (cmd *OriginalCmd) Select(selectColumns *column.Columns, fromClause string, whereCondition *condition.Condition) ([]*maps.GoleMap, error) {
//	return cmd.Query(generateSelectSql(selectColumns, fromClause, whereCondition))
//}
//
//func generateBindParamsTagsType(cmd *OriginalCmd, tableName string, tagsMap *maps.GoleMap, dataMap *maps.GoleMap) *param.ColumnType {
//	finalColumnType := param.NewColumnType(tagsMap.Size() + dataMap.Size())
//	fieldMap := cmd.getTableFiledTaosTypeMap(tableName)
//
//	for i, key := range tagsMap.Keys() {
//		if i >= tagsMap.Size() {
//			break
//		}
//		fieldMeta := fieldMap[key]
//		switch fieldMeta.ColType {
//		case types.TaosBoolType:
//			finalColumnType.AddBool()
//		case types.TaosTinyintType:
//			finalColumnType.AddTinyint()
//		case types.TaosSmallintType:
//			finalColumnType.AddSmallint()
//		case types.TaosIntType:
//			finalColumnType.AddInt()
//		case types.TaosBigintType:
//			finalColumnType.AddBigint()
//		case types.TaosUTinyintType:
//			finalColumnType.AddUTinyint()
//		case types.TaosUSmallintType:
//			finalColumnType.AddUSmallint()
//		case types.TaosUIntType:
//			finalColumnType.AddUInt()
//		case types.TaosUBigintType:
//			finalColumnType.AddUBigint()
//		case types.TaosFloatType:
//			finalColumnType.AddFloat()
//		case types.TaosDoubleType:
//			finalColumnType.AddDouble()
//		case types.TaosBinaryType:
//			finalColumnType.AddBinary(fieldMeta.ColLen)
//		case types.TaosVarBinaryType:
//			finalColumnType.AddVarBinary(fieldMeta.ColLen)
//		case types.TaosNcharType:
//			finalColumnType.AddNchar(fieldMeta.ColLen)
//		case types.TaosTimestampType:
//			finalColumnType.AddTimestamp()
//		case types.TaosJsonType:
//			finalColumnType.AddJson(fieldMeta.ColLen)
//		case types.TaosGeometryType:
//			finalColumnType.AddGeometry(fieldMeta.ColLen)
//		}
//	}
//
//	for i, key := range dataMap.Keys() {
//		if i >= dataMap.Size() {
//			break
//		}
//		fieldMeta := fieldMap[key]
//		switch fieldMeta.ColType {
//		case types.TaosBoolType:
//			finalColumnType.AddBool()
//		case types.TaosTinyintType:
//			finalColumnType.AddTinyint()
//		case types.TaosSmallintType:
//			finalColumnType.AddSmallint()
//		case types.TaosIntType:
//			finalColumnType.AddInt()
//		case types.TaosBigintType:
//			finalColumnType.AddBigint()
//		case types.TaosUTinyintType:
//			finalColumnType.AddUTinyint()
//		case types.TaosUSmallintType:
//			finalColumnType.AddUSmallint()
//		case types.TaosUIntType:
//			finalColumnType.AddUInt()
//		case types.TaosUBigintType:
//			finalColumnType.AddUBigint()
//		case types.TaosFloatType:
//			finalColumnType.AddFloat()
//		case types.TaosDoubleType:
//			finalColumnType.AddDouble()
//		case types.TaosBinaryType:
//			finalColumnType.AddBinary(fieldMeta.ColLen)
//		case types.TaosVarBinaryType:
//			finalColumnType.AddVarBinary(fieldMeta.ColLen)
//		case types.TaosNcharType:
//			finalColumnType.AddNchar(fieldMeta.ColLen)
//		case types.TaosTimestampType:
//			finalColumnType.AddTimestamp()
//		case types.TaosJsonType:
//			finalColumnType.AddJson(fieldMeta.ColLen)
//		case types.TaosGeometryType:
//			finalColumnType.AddGeometry(fieldMeta.ColLen)
//		}
//	}
//	return finalColumnType
//}
//
//func (cmd *OriginalCmd) NewTdHookContextOfOriginal(sql, runType string) *TdHookContext {
//	return &TdHookContext{
//		ConnectType: constants.ConnectOriginal,
//		DbName:      cmd.dbName,
//		Start:       time.Now(),
//		Sql:         sql,
//		RunType:     runType,
//	}
//}
//
//func (cmd *OriginalCmd) doBatchInsert(tableName string, toInsertMaps []*maps.GoleMap) (int64, error) {
//	if len(toInsertMaps) == 0 {
//		return 0, nil
//	}
//	insertStmt := cmd.tdConnect.InsertStmt()
//	defer func() {
//		err := insertStmt.Close()
//		if err != nil {
//			logger.Warnf("执行BatchInsert报错")
//		}
//	}()
//
//	// 这里默认所有插入的数据的key都完全一样的，这里不再进行检查（为了效率），只取最开始的这个作为结构sql的输入，对于一些数据不均匀的集合来说需要业务保证数据集合的均匀性
//	toInsertMap := toInsertMaps[0]
//	preInsertSql := generateInsertSql(cmd.dbName, tableName, toInsertMap)
//
//	// before
//	hookContext, err := cmd.hooks.preProcess(cmd.NewTdHookContextOfOriginal(preInsertSql, constants.BATCH_INSERT))
//
//	err = insertStmt.Prepare(preInsertSql)
//	if err != nil {
//		logger.Errorf("执行BatchInsert，Prepare 报错，%v", err.Error())
//		// after
//		cmd.hooks.afterProcess(hookContext.EndBatchInsert(0, err, toInsertMaps))
//		return 0, err
//	}
//	logger.Group("tdorm").Debugf("batch【preInsertSql ==> %v】【params ==> %v】", preInsertSql, toInsertMap.ToJsonOfSort())
//
//	// 批量绑定，这样效率高
//	for _, toInsertMap := range toInsertMaps {
//		bindType := generateBindParamsType(cmd, tableName, toInsertMap)
//		params, unFindKeys := generateParams(cmd, tableName, toInsertMap)
//		err = insertStmt.BindParam(params, bindType)
//		if err != nil {
//			logger.Errorf("执行BatchInsert，BindParam 报错【preInsertSql ==> %v】【params ==> %v】，错误：%v", preInsertSql, toInsertMap.CloneExceptKeys(unFindKeys).ToJsonOfSort(), err.Error())
//			// after
//			cmd.hooks.afterProcess(hookContext.EndBatchInsert(0, err, toInsertMaps))
//			return 0, err
//		}
//	}
//
//	err = insertStmt.AddBatch()
//	if err != nil {
//		logger.Errorf("执行BatchInsert，AddBatch 报错，%v", err.Error())
//		// after
//		cmd.hooks.afterProcess(hookContext.EndBatchInsert(0, err, toInsertMaps))
//		return 0, err
//	}
//	err = insertStmt.Execute()
//	if err != nil {
//		logger.Errorf("执行BatchInsert，Execute 报错，%v", err.Error())
//		// after
//		cmd.hooks.afterProcess(hookContext.EndBatchInsert(0, err, toInsertMaps))
//		return 0, err
//	}
//	result := insertStmt.GetAffectedRows()
//	// after
//	cmd.hooks.afterProcess(hookContext.EndBatchInsert(util.ToInt64(result), err, toInsertMaps))
//	return util.ToInt64(result), nil
//}
//
//func generateBindParamsType(cmd *OriginalCmd, tableName string, dataMap *maps.GoleMap) *param.ColumnType {
//	finalColumnType := param.NewColumnType(dataMap.Size())
//	fieldMap := cmd.getTableFiledTaosTypeMap(tableName)
//	for i, key := range dataMap.Keys() {
//		if i >= dataMap.Size() {
//			break
//		}
//		fieldMeta := fieldMap[key]
//		switch fieldMeta.ColType {
//		case types.TaosBoolType:
//			finalColumnType.AddBool()
//		case types.TaosTinyintType:
//			finalColumnType.AddTinyint()
//		case types.TaosSmallintType:
//			finalColumnType.AddSmallint()
//		case types.TaosIntType:
//			finalColumnType.AddInt()
//		case types.TaosBigintType:
//			finalColumnType.AddBigint()
//		case types.TaosUTinyintType:
//			finalColumnType.AddUTinyint()
//		case types.TaosUSmallintType:
//			finalColumnType.AddUSmallint()
//		case types.TaosUIntType:
//			finalColumnType.AddUInt()
//		case types.TaosUBigintType:
//			finalColumnType.AddUBigint()
//		case types.TaosFloatType:
//			finalColumnType.AddFloat()
//		case types.TaosDoubleType:
//			finalColumnType.AddDouble()
//		case types.TaosBinaryType:
//			finalColumnType.AddBinary(fieldMeta.ColLen)
//		case types.TaosVarBinaryType:
//			finalColumnType.AddVarBinary(fieldMeta.ColLen)
//		case types.TaosNcharType:
//			finalColumnType.AddNchar(fieldMeta.ColLen)
//		case types.TaosTimestampType:
//			finalColumnType.AddTimestamp()
//		case types.TaosJsonType:
//			finalColumnType.AddJson(fieldMeta.ColLen)
//		case types.TaosGeometryType:
//			finalColumnType.AddGeometry(fieldMeta.ColLen)
//		}
//	}
//	return finalColumnType
//}
//
//// 生成如下的列表
////
////	params := []*param.Param{
////			param.NewParam(1).AddTimestamp(timeData1, 0),
////			param.NewParam(1).AddNchar("zhou1--1"),
////			param.NewParam(1).AddInt(19),
////			param.NewParam(1).AddNchar("hangzhou1"),
////		}
//func generateParams(cmd *OriginalCmd, tableName string, dataMap *maps.GoleMap) ([]*param.Param, []string) {
//	size := dataMap.Size()
//	var params []*param.Param
//	fieldMap := cmd.getTableFiledTaosTypeMap(tableName)
//	var unFindKeys []string
//	for i, key := range dataMap.Keys() {
//		if i >= size {
//			break
//		}
//		fieldMeta := fieldMap[key]
//		switch fieldMeta.ColType {
//		case types.TaosBoolType:
//			if val, exist := dataMap.GetBool(key); exist {
//				params = append(params, param.NewParam(1).AddBool(val))
//			}
//		case types.TaosTinyintType:
//			if val, exist := dataMap.GetInt(key); exist {
//				params = append(params, param.NewParam(1).AddTinyint(val))
//			}
//		case types.TaosSmallintType:
//			if val, exist := dataMap.GetInt(key); exist {
//				params = append(params, param.NewParam(1).AddSmallint(val))
//			}
//		case types.TaosIntType:
//			if val, exist := dataMap.GetInt(key); exist {
//				params = append(params, param.NewParam(1).AddInt(val))
//			}
//		case types.TaosBigintType:
//			if val, exist := dataMap.GetInt(key); exist {
//				params = append(params, param.NewParam(1).AddBigint(val))
//			}
//		case types.TaosUTinyintType:
//			if val, exist := dataMap.GetUInt(key); exist {
//				params = append(params, param.NewParam(1).AddUTinyint(val))
//			}
//		case types.TaosUSmallintType:
//			if val, exist := dataMap.GetUInt(key); exist {
//				params = append(params, param.NewParam(1).AddUSmallint(val))
//			}
//		case types.TaosUIntType:
//			if val, exist := dataMap.GetUInt(key); exist {
//				params = append(params, param.NewParam(1).AddUInt(val))
//			}
//		case types.TaosUBigintType:
//			if val, exist := dataMap.GetUInt(key); exist {
//				params = append(params, param.NewParam(1).AddUBigint(val))
//			}
//		case types.TaosFloatType:
//			if val, exist := dataMap.GetFloat32(key); exist {
//				params = append(params, param.NewParam(1).AddFloat(val))
//			}
//		case types.TaosDoubleType:
//			if val, exist := dataMap.GetFloat64(key); exist {
//				params = append(params, param.NewParam(1).AddDouble(val))
//			}
//		case types.TaosBinaryType:
//			if val, exist := dataMap.GetBytes(key); exist {
//				params = append(params, param.NewParam(1).AddBinary(val))
//			}
//		case types.TaosVarBinaryType:
//			if val, exist := dataMap.GetBytes(key); exist {
//				params = append(params, param.NewParam(1).AddVarBinary(val))
//			}
//		case types.TaosNcharType:
//			if val, exist := dataMap.GetString(key); exist {
//				params = append(params, param.NewParam(1).AddNchar(val))
//			}
//		case types.TaosTimestampType:
//			if val, exist := dataMap.GetTime(key); exist {
//				params = append(params, param.NewParam(1).AddTimestamp(val, 0))
//			}
//		case types.TaosJsonType:
//			if val, exist := dataMap.GetBytes(key); exist {
//				params = append(params, param.NewParam(1).AddJson(val))
//			}
//		case types.TaosGeometryType:
//			if val, exist := dataMap.GetBytes(key); exist {
//				params = append(params, param.NewParam(1).AddGeometry(val))
//			}
//		default:
//			logger.Warnf("key【%v】在表【%v】没有找到对应的key", key, tableName)
//			unFindKeys = append(unFindKeys, key)
//		}
//	}
//	return params, unFindKeys
//}
//
//func generateParamsTags(cmd *OriginalCmd, tableName string, tagsMap *maps.GoleMap, dataMap *maps.GoleMap) ([]*param.Param, []string) {
//	size := tagsMap.Size() + dataMap.Size()
//	var params []*param.Param
//	fieldMap := cmd.getTableFiledTaosTypeMap(tableName)
//	var unFindKeys []string
//	for i, key := range tagsMap.Keys() {
//		if i >= size {
//			break
//		}
//		fieldMeta := fieldMap[key]
//		switch fieldMeta.ColType {
//		case types.TaosBoolType:
//			if val, exist := tagsMap.GetBool(key); exist {
//				params = append(params, param.NewParam(1).AddBool(val))
//			}
//		case types.TaosTinyintType:
//			if val, exist := tagsMap.GetInt(key); exist {
//				params = append(params, param.NewParam(1).AddTinyint(val))
//			}
//		case types.TaosSmallintType:
//			if val, exist := tagsMap.GetInt(key); exist {
//				params = append(params, param.NewParam(1).AddSmallint(val))
//			}
//		case types.TaosIntType:
//			if val, exist := tagsMap.GetInt(key); exist {
//				params = append(params, param.NewParam(1).AddInt(val))
//			}
//		case types.TaosBigintType:
//			if val, exist := tagsMap.GetInt(key); exist {
//				params = append(params, param.NewParam(1).AddBigint(val))
//			}
//		case types.TaosUTinyintType:
//			if val, exist := tagsMap.GetUInt(key); exist {
//				params = append(params, param.NewParam(1).AddUTinyint(val))
//			}
//		case types.TaosUSmallintType:
//			if val, exist := tagsMap.GetUInt(key); exist {
//				params = append(params, param.NewParam(1).AddUSmallint(val))
//			}
//		case types.TaosUIntType:
//			if val, exist := tagsMap.GetUInt(key); exist {
//				params = append(params, param.NewParam(1).AddUInt(val))
//			}
//		case types.TaosUBigintType:
//			if val, exist := tagsMap.GetUInt(key); exist {
//				params = append(params, param.NewParam(1).AddUBigint(val))
//			}
//		case types.TaosFloatType:
//			if val, exist := tagsMap.GetFloat32(key); exist {
//				params = append(params, param.NewParam(1).AddFloat(val))
//			}
//		case types.TaosDoubleType:
//			if val, exist := tagsMap.GetFloat64(key); exist {
//				params = append(params, param.NewParam(1).AddDouble(val))
//			}
//		case types.TaosBinaryType:
//			if val, exist := tagsMap.GetBytes(key); exist {
//				params = append(params, param.NewParam(1).AddBinary(val))
//			}
//		case types.TaosVarBinaryType:
//			if val, exist := tagsMap.GetBytes(key); exist {
//				params = append(params, param.NewParam(1).AddVarBinary(val))
//			}
//		case types.TaosNcharType:
//			if val, exist := tagsMap.GetString(key); exist {
//				params = append(params, param.NewParam(1).AddNchar(val))
//			}
//		case types.TaosTimestampType:
//			if val, exist := tagsMap.GetTime(key); exist {
//				params = append(params, param.NewParam(1).AddTimestamp(val, 0))
//			}
//		case types.TaosJsonType:
//			if val, exist := tagsMap.GetBytes(key); exist {
//				params = append(params, param.NewParam(1).AddJson(val))
//			}
//		case types.TaosGeometryType:
//			if val, exist := tagsMap.GetBytes(key); exist {
//				params = append(params, param.NewParam(1).AddGeometry(val))
//			}
//		default:
//			logger.Warnf("key【%v】在表【%v】没有找到对应的key", key, tableName)
//			unFindKeys = append(unFindKeys, key)
//		}
//	}
//
//	for i, key := range dataMap.Keys() {
//		if i >= size {
//			break
//		}
//		fieldMeta := fieldMap[key]
//		switch fieldMeta.ColType {
//		case types.TaosBoolType:
//			if val, exist := dataMap.GetBool(key); exist {
//				params = append(params, param.NewParam(1).AddBool(val))
//			}
//		case types.TaosTinyintType:
//			if val, exist := dataMap.GetInt(key); exist {
//				params = append(params, param.NewParam(1).AddTinyint(val))
//			}
//		case types.TaosSmallintType:
//			if val, exist := dataMap.GetInt(key); exist {
//				params = append(params, param.NewParam(1).AddSmallint(val))
//			}
//		case types.TaosIntType:
//			if val, exist := dataMap.GetInt(key); exist {
//				params = append(params, param.NewParam(1).AddInt(val))
//			}
//		case types.TaosBigintType:
//			if val, exist := dataMap.GetInt(key); exist {
//				params = append(params, param.NewParam(1).AddBigint(val))
//			}
//		case types.TaosUTinyintType:
//			if val, exist := dataMap.GetUInt(key); exist {
//				params = append(params, param.NewParam(1).AddUTinyint(val))
//			}
//		case types.TaosUSmallintType:
//			if val, exist := dataMap.GetUInt(key); exist {
//				params = append(params, param.NewParam(1).AddUSmallint(val))
//			}
//		case types.TaosUIntType:
//			if val, exist := dataMap.GetUInt(key); exist {
//				params = append(params, param.NewParam(1).AddUInt(val))
//			}
//		case types.TaosUBigintType:
//			if val, exist := dataMap.GetUInt(key); exist {
//				params = append(params, param.NewParam(1).AddUBigint(val))
//			}
//		case types.TaosFloatType:
//			if val, exist := dataMap.GetFloat32(key); exist {
//				params = append(params, param.NewParam(1).AddFloat(val))
//			}
//		case types.TaosDoubleType:
//			if val, exist := dataMap.GetFloat64(key); exist {
//				params = append(params, param.NewParam(1).AddDouble(val))
//			}
//		case types.TaosBinaryType:
//			if val, exist := dataMap.GetBytes(key); exist {
//				params = append(params, param.NewParam(1).AddBinary(val))
//			}
//		case types.TaosVarBinaryType:
//			if val, exist := dataMap.GetBytes(key); exist {
//				params = append(params, param.NewParam(1).AddVarBinary(val))
//			}
//		case types.TaosNcharType:
//			if val, exist := dataMap.GetString(key); exist {
//				params = append(params, param.NewParam(1).AddNchar(val))
//			}
//		case types.TaosTimestampType:
//			if val, exist := dataMap.GetTime(key); exist {
//				params = append(params, param.NewParam(1).AddTimestamp(val, 0))
//			}
//		case types.TaosJsonType:
//			if val, exist := dataMap.GetBytes(key); exist {
//				params = append(params, param.NewParam(1).AddJson(val))
//			}
//		case types.TaosGeometryType:
//			if val, exist := dataMap.GetBytes(key); exist {
//				params = append(params, param.NewParam(1).AddGeometry(val))
//			}
//		default:
//			logger.Warnf("key【%v】在表【%v】没有找到对应的key", key, tableName)
//			unFindKeys = append(unFindKeys, key)
//		}
//	}
//	return params, unFindKeys
//}
//
//func generateInsertParams(cmd *OriginalCmd, tableName string, dataMap *maps.GoleMap) *param.Param {
//	len := dataMap.Size()
//	finalParam := param.NewParam(dataMap.Size())
//	fieldMap := cmd.getTableFiledTaosTypeMap(tableName)
//	for i, key := range dataMap.Keys() {
//		if i >= len {
//			break
//		}
//		fieldMeta := fieldMap[key]
//		switch fieldMeta.ColType {
//		case types.TaosBoolType:
//			if val, exist := dataMap.GetBool(key); exist {
//				finalParam.AddBool(val)
//			}
//		case types.TaosTinyintType:
//			if val, exist := dataMap.GetInt(key); exist {
//				finalParam.AddTinyint(val)
//			}
//		case types.TaosSmallintType:
//			if val, exist := dataMap.GetInt(key); exist {
//				finalParam.AddSmallint(val)
//			}
//		case types.TaosIntType:
//			if val, exist := dataMap.GetInt(key); exist {
//				finalParam.AddInt(val)
//			}
//		case types.TaosBigintType:
//			if val, exist := dataMap.GetInt(key); exist {
//				finalParam.AddBigint(val)
//			}
//		case types.TaosUTinyintType:
//			if val, exist := dataMap.GetUInt(key); exist {
//				finalParam.AddUTinyint(val)
//			}
//		case types.TaosUSmallintType:
//			if val, exist := dataMap.GetUInt(key); exist {
//				finalParam.AddUSmallint(val)
//			}
//		case types.TaosUIntType:
//			if val, exist := dataMap.GetUInt(key); exist {
//				finalParam.AddUInt(val)
//			}
//		case types.TaosUBigintType:
//			if val, exist := dataMap.GetUInt(key); exist {
//				finalParam.AddUBigint(val)
//			}
//		case types.TaosFloatType:
//			if val, exist := dataMap.GetFloat32(key); exist {
//				finalParam.AddFloat(val)
//			}
//		case types.TaosDoubleType:
//			if val, exist := dataMap.GetFloat64(key); exist {
//				finalParam.AddDouble(val)
//			}
//		case types.TaosBinaryType:
//			if val, exist := dataMap.GetBytes(key); exist {
//				finalParam.AddBinary(val)
//			}
//		case types.TaosVarBinaryType:
//			if val, exist := dataMap.GetBytes(key); exist {
//				finalParam.AddVarBinary(val)
//			}
//		case types.TaosNcharType:
//			if val, exist := dataMap.GetString(key); exist {
//				finalParam.AddNchar(val)
//			}
//		case types.TaosTimestampType:
//			if val, exist := dataMap.GetTime(key); exist {
//				finalParam.AddTimestamp(val, 0)
//			}
//		case types.TaosJsonType:
//			if val, exist := dataMap.GetBytes(key); exist {
//				finalParam.AddJson(val)
//			}
//		case types.TaosGeometryType:
//			if val, exist := dataMap.GetBytes(key); exist {
//				finalParam.AddGeometry(val)
//			}
//		}
//	}
//	return finalParam
//}
