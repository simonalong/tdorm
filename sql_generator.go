package tdorm

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/simonalong/gole/maps"
	goleTime "github.com/simonalong/gole/time"
	"github.com/simonalong/gole/util"
	"github.com/simonalong/tdorm/column"
	"github.com/simonalong/tdorm/condition"
	"github.com/taosdata/driver-go/v3/types"
)

func generateInsertSql(dbName, tableName string, dataMap *maps.GoleMap) string {
	sqlFormat := "insert into %s (%s) values (%s)"
	return fmt.Sprintf(sqlFormat, generateTable(dbName, tableName), generateFieldsNameSql(dataMap), generateValueSql(dataMap))
}

func generateInsertWithTagSql(dbName, tableName, stableName string, tagMap *maps.GoleMap, dataMap *maps.GoleMap) string {
	sqlFormat := "insert into %s using %s (%s) tags (%s) (%s) values (%s)"
	return fmt.Sprintf(sqlFormat, generateTable(dbName, tableName), generateTable(dbName, stableName), generateTagsNameSql(tagMap), generateFullValueSql(tagMap), generateFieldsNameSql(dataMap), generateValueSql(dataMap))
}

func generateInsertBatchWithTagSql(dbName, tableName, stableName string, tagMap *maps.GoleMap, dataMaps []*maps.GoleMap) string {
	sqlFormat := "insert into %s using %s (%s) tags (%s) (%s) values %s"
	return fmt.Sprintf(sqlFormat, generateTable(dbName, tableName), generateTable(dbName, stableName), generateTagsNameSql(tagMap), generateFullValueSql(tagMap), generateFieldsNameSql(dataMaps[0]), generateValuesSql(dataMaps))
}

func generateInsertBatchFullSql(dbName, tableName string, dataMaps []*maps.GoleMap) string {
	sqlFormat := "insert into %s (%s) values %s"
	return fmt.Sprintf(sqlFormat, generateTable(dbName, tableName), generateFieldsNameSql(dataMaps[0]), generateBatchFullValueSql(dataMaps))
}

func generateInsertBatchFullWithTagSql(dbName, tableName, stableName string, tagsMap *maps.GoleMap, dataMaps []*maps.GoleMap) string {
	sqlFormat := "insert into %s using %s (%s) tags (%s) (%s) values %s"
	return fmt.Sprintf(sqlFormat, generateTable(dbName, tableName), generateTable(dbName, stableName), generateTagsNameSql(tagsMap), generateFullValueSql(tagsMap), generateFieldsNameSql(dataMaps[0]), generateBatchFullValueSql(dataMaps))
}

func generateDeleteSql(dbName, tableName string, query *condition.Condition) string {
	sqlFormat := "delete from %s %s"
	return fmt.Sprintf(sqlFormat, generateTable(dbName, tableName), query.ToSql())
}

// 非占位符的全sql
// 主要是给websocket用；目前发现 websocket的delete预处理报错
func generateDeleteSqlOfFull(dbName, tableName string, query *condition.Condition) string {
	sqlFormat := "delete from %s %s"
	return fmt.Sprintf(sqlFormat, generateTable(dbName, tableName), query.ToSqlOfFull())
}

func generateQueryOneSql(dbName, tableName string, columns *column.Columns, query *condition.Condition) string {
	sqlFormat := "select %s from %s %s limit 1"
	return fmt.Sprintf(sqlFormat, columns.ToSql(), generateTable(dbName, tableName), query.ToSql())
}

func generateQueryListSql(dbName, tableName string, columns *column.Columns, query *condition.Condition) string {
	sqlFormat := "select %s from %s %s"
	return fmt.Sprintf(sqlFormat, columns.ToSql(), generateTable(dbName, tableName), query.ToSql())
}

func generateQueryListOfDistinctSql(dbName, tableName string, columns *column.Columns, query *condition.Condition) string {
	sqlFormat := "select distinct %s from %s %s"
	return fmt.Sprintf(sqlFormat, columns.ToSql(), generateTable(dbName, tableName), query.ToSql())
}

func generateQueryValueSql(dbName, tableName, columns string, query *condition.Condition) string {
	sqlFormat := "select %s from %s %s limit 1"
	return fmt.Sprintf(sqlFormat, column.ToDbField(columns), generateTable(dbName, tableName), query.ToSql())
}

func generateQueryValuesSql(dbName, tableName, columns string, query *condition.Condition) string {
	sqlFormat := "select %s from %s %s"
	return fmt.Sprintf(sqlFormat, column.ToDbField(columns), generateTable(dbName, tableName), query.ToSql())
}

func generateQueryValuesOfDistinctSql(dbName, tableName, columns string, query *condition.Condition) string {
	sqlFormat := "select distinct %s from %s %s"
	return fmt.Sprintf(sqlFormat, column.ToDbField(columns), generateTable(dbName, tableName), query.ToSql())
}

func generateQueryCountSql(dbName, tableName string, query *condition.Condition) string {
	sqlFormat := "select count(*) as cnt from %s %s"
	return fmt.Sprintf(sqlFormat, generateTable(dbName, tableName), query.ToSql())
}

func generateQueryParams(query *condition.Condition) []driver.Value {
	var values []driver.Value
	for _, val := range query.ToValues() {
		values = append(values, val)
	}
	return values
}

// 返回：<DbName>.tableName
func generateTable(dbName, tableName string) string {
	if strings.Contains(tableName, ".") {
		return tableName
	} else {
		if dbName != "" {
			return dbName + "." + tableName
		}
		return tableName
	}
}

// 返回：`ts`, `name`, `age`
func generateFieldsNameSql(dataMap *maps.GoleMap) string {
	var keys []string
	for _, key := range dataMap.Keys() {
		if !strings.HasPrefix(key, "`") && !strings.HasSuffix(key, "`") {
			keys = append(keys, "`"+key+"`")
		}
	}
	return strings.Join(keys, ",")
}

// 返回：(23, '宋江', '杭州') (29, '霍元甲', '郴州') ...
func generateBatchFullValueSql(dataMaps []*maps.GoleMap) string {
	if len(dataMaps) == 0 {
		return ""
	}

	var strList []string
	for _, dataMap := range dataMaps {
		strList = append(strList, "("+generateFullValueSql(dataMap)+")")
	}
	return strings.Join(strList, " ")
}

// 返回：(?, ?, ?) (?, ?, ?) (?, ?, ?)...
func generateValuesSql(dataMaps []*maps.GoleMap) string {
	var strList []string
	for _, dataMap := range dataMaps {
		strList = append(strList, generateValueSql(dataMap))
	}
	return strings.Join(strList, " ")
}

// 返回：?, ?, ?
func generateValueSql(dataMap *maps.GoleMap) string {
	var seizes []string
	for range dataMap.Keys() {
		seizes = append(seizes, "?")
	}
	return strings.Join(seizes, ",")
}

func generateFullValueSql(dataMap *maps.GoleMap) string {
	var seizes []string
	for _, key := range dataMap.Keys() {
		val, _ := dataMap.Get(key)
		seizes = append(seizes, typeChangeToString(val))
	}
	return strings.Join(seizes, ",")
}

// 返回：`ts`, `name`, `age`
func generateTagsNameSql(tagsMap *maps.GoleMap) string {
	var keys []string
	for _, key := range tagsMap.Keys() {
		if !strings.HasPrefix(key, "`") && !strings.HasSuffix(key, "`") {
			keys = append(keys, "`"+key+"`")
		}
	}
	return strings.Join(keys, ",")
}

func generateSelectSql(selectColumns *column.Columns, fromClause string, whereCondition *condition.Condition) string {
	sqlFormat := "select %s from %s %s"
	return strings.TrimSpace(fmt.Sprintf(sqlFormat, selectColumns.ToSql(), fromClause, whereCondition.ToSql()))
}

func typeChangeToString(val any) string {
	valType := reflect.TypeOf(val)
	if util.IsNumberType(valType) {
		return util.ToString(val)
	} else if util.IsStringType(valType) {
		return "'" + util.ToString(val) + "'"
	} else if util.IsTimeType(valType) {
		return "'" + goleTime.TimeToStringYmdHmsS(val.(time.Time)) + "'"
	} else if util.IsBoolType(valType) {
		return util.ToString(val)
	}
	return util.ToString(val)
}

// tdengine的数据库类型向taos类型转换
func tdengineColTypeToTaosType(colTypeStr string) interface{} {
	switch colTypeStr {
	case "TIMESTAMP":
		return types.TaosTimestampType
	case "BOOL":
		return types.TaosBoolType
	case "TINYINT":
		return types.TaosTinyintType
	case "SMALLINT":
		return types.TaosSmallintType
	case "INT":
		return types.TaosIntType
	case "BIGINT":
		return types.TaosBigintType
	case "TINYINT UNSIGNED":
		return types.TaosUTinyintType
	case "SMALLINT UNSIGNED":
		return types.TaosUSmallintType
	case "INT UNSIGNED":
		return types.TaosUIntType
	case "BIGINT UNSIGNED":
		return types.TaosUBigintType
	case "FLOAT":
		return types.TaosFloatType
	case "DOUBLE":
		return types.TaosDoubleType
	case "VARBINARY":
		return types.TaosVarBinaryType
	case "GEOMETRY":
		return types.TaosGeometryType
	}

	if strings.HasPrefix(colTypeStr, "VARCHAR") {
		return types.TaosBinaryType
	} else if strings.HasPrefix(colTypeStr, "NCHAR") {
		return types.TaosNcharType
	}
	return nil
}
