package condition

import (
	"reflect"
	"strings"

	"github.com/simonalong/gole/maps"
	"github.com/simonalong/tdorm/constants"
	"github.com/simonalong/tdorm/op"
)

type Condition struct {
	LogicOperates []*op.LogicOperate
	// 其他的条件
	//    [partition_by_clause]
	//    [interp_clause]
	//    [window_clause]
	//    [group_by_clause]
	//    [order_by_clause]
	//    [SLIMIT limit_val [OFFSET offset_val]]
	//    [LIMIT limit_val [OFFSET offset_val]]
	//    [>> export_file]
	OtherClause string
}

func New() *Condition {
	return &Condition{
		LogicOperates: make([]*op.LogicOperate, 0),
	}
}

func (condition *Condition) ToSqlOfFull() string {
	logicOperates := condition.LogicOperates
	var resultSql string
	innerNeedWhere := false
	for _, logicOperate := range logicOperates {
		if logicOperate.NeedWhere() {
			innerNeedWhere = true
		}
		resultSql += logicOperate.GenerateSqlOfFull()
	}

	// 去头
	if len(resultSql) != 0 {
		resultSql = strings.TrimSpace(resultSql)
		if strings.HasPrefix(resultSql, constants.AND) {
			resultSql = resultSql[len(constants.AND):]
		}

		resultSql = strings.TrimSpace(resultSql)
		if strings.HasPrefix(resultSql, constants.OR) {
			resultSql = resultSql[len(constants.OR):]
		}

		resultSql = strings.TrimSpace(resultSql)
		if innerNeedWhere {
			resultSql = constants.WHERE + constants.EMPTY + resultSql
		}
	}

	resultSql += condition.OtherClause
	return resultSql
}

func (condition *Condition) ToSql() string {
	if condition == nil {
		return ""
	}
	logicOperates := condition.LogicOperates
	var resultSql string
	innerNeedWhere := false
	for _, logicOperate := range logicOperates {
		if logicOperate.NeedWhere() {
			innerNeedWhere = true
		}
		resultSql += logicOperate.GenerateSql()
	}

	// 去头
	if len(resultSql) != 0 {
		resultSql = strings.TrimSpace(resultSql)
		if strings.HasPrefix(resultSql, constants.AND) {
			resultSql = resultSql[len(constants.AND):]
		}

		resultSql = strings.TrimSpace(resultSql)
		if strings.HasPrefix(resultSql, constants.OR) {
			resultSql = resultSql[len(constants.OR):]
		}

		resultSql = strings.TrimSpace(resultSql)
		if innerNeedWhere {
			resultSql = constants.WHERE + constants.EMPTY + resultSql
		}
	}

	resultSql += condition.OtherClause
	return strings.TrimSpace(resultSql)
}

func (condition *Condition) Append(otherClause string) *Condition {
	condition.OtherClause = condition.OtherClause + " " + otherClause
	return condition
}

// And 输出：(`age` = ? and `name` = ?)
func (condition *Condition) And(datas ...any) *Condition {
	if datas == nil || len(datas) == 0 {
		return condition
	}

	condition.LogicOperates = append(condition.LogicOperates, op.NewLogic(constants.AND, true).AddChildrenLogicOperate(parseLogicOperator(logicAndEm, datas...)))
	return condition
}

// AndEm 输出：`age` = ? and `name` = ?
func (condition *Condition) AndEm(datas ...any) *Condition {
	if datas == nil || len(datas) == 0 {
		return condition
	}
	condition.LogicOperates = append(condition.LogicOperates, op.NewLogic(constants.AND, false).AddChildrenLogicOperate(parseLogicOperator(logicAndEm, datas...)))
	return condition
}

// Or 输出：(`age` = ? or `name` = ?)
func (condition *Condition) Or(datas ...any) *Condition {
	if datas == nil || len(datas) == 0 {
		return condition
	}
	condition.LogicOperates = append(condition.LogicOperates, op.NewLogic(constants.OR, true).AddChildrenLogicOperate(parseLogicOperator(logicOrEm, datas...)))
	return condition
}

// OrEm 输出：`age` = ? or `name` = ?
func (condition *Condition) OrEm(datas ...any) *Condition {
	if datas == nil || len(datas) == 0 {
		return condition
	}
	condition.LogicOperates = append(condition.LogicOperates, op.NewLogic(constants.OR, false).AddChildrenLogicOperate(parseLogicOperator(logicOrEm, datas...)))
	return condition
}

// Em 输出：`age` = ?
func (condition *Condition) Em(datas ...any) *Condition {
	if datas == nil || len(datas) == 0 {
		return condition
	}
	condition.LogicOperates = append(condition.LogicOperates, op.NewLogic(constants.EMPTY, false).AddChildrenLogicOperate(parseLogicOperator(logicEmpty, datas...)))
	return condition
}

// Gt 输出：`age` > ?
func (condition *Condition) Gt(key string, value interface{}) *Condition {
	condition.LogicOperates = append(condition.LogicOperates, op.NewLogic(constants.AND, false).AddRelationOperate(op.Gt(key, value)))
	return condition
}

// Ge 输出：`age` >= ?
func (condition *Condition) Ge(key string, value interface{}) *Condition {
	condition.LogicOperates = append(condition.LogicOperates, op.NewLogic(constants.AND, false).AddRelationOperate(op.Ge(key, value)))
	return condition
}

// Lt 输出：`age` < ?
func (condition *Condition) Lt(key string, value interface{}) *Condition {
	condition.LogicOperates = append(condition.LogicOperates, op.NewLogic(constants.AND, false).AddRelationOperate(op.Lt(key, value)))
	return condition
}

// Le 输出：`age` <= ?
func (condition *Condition) Le(key string, value interface{}) *Condition {
	condition.LogicOperates = append(condition.LogicOperates, op.NewLogic(constants.AND, false).AddRelationOperate(op.Le(key, value)))
	return condition
}

// Eq 输出：`age` = ?
func (condition *Condition) Eq(key string, value interface{}) *Condition {
	condition.LogicOperates = append(condition.LogicOperates, op.NewLogic(constants.AND, false).AddRelationOperate(op.Eq(key, value)))
	return condition
}

// UnEq 输出：`age` != ?
func (condition *Condition) UnEq(key string, value interface{}) *Condition {
	condition.LogicOperates = append(condition.LogicOperates, op.NewLogic(constants.AND, false).AddRelationOperate(op.UnEq(key, value)))
	return condition
}

// IsNull 输出：`name` is null
func (condition *Condition) IsNull(key string) *Condition {
	condition.LogicOperates = append(condition.LogicOperates, op.NewLogic(constants.AND, false).AddRelationOperate(op.IsNull(key)))
	return condition
}

// IsNotNull 输出：`name` is not null
func (condition *Condition) IsNotNull(key string) *Condition {
	condition.LogicOperates = append(condition.LogicOperates, op.NewLogic(constants.AND, false).AddRelationOperate(op.IsNotNull(key)))
	return condition
}

// BetweenAnd 输出：`age` between xx and xx
func (condition *Condition) BetweenAnd(key string, leftValue, rightValue interface{}) *Condition {
	condition.LogicOperates = append(condition.LogicOperates, op.NewLogic(constants.AND, false).AddRelationOperate(op.BetweenAnd(key, leftValue, rightValue)))
	return condition
}

// NotBetweenAnd 输出：`age` not between xx and xx
func (condition *Condition) NotBetweenAnd(key string, leftValue, rightValue interface{}) *Condition {
	condition.LogicOperates = append(condition.LogicOperates, op.NewLogic(constants.AND, false).AddRelationOperate(op.NotBetweenAnd(key, leftValue, rightValue)))
	return condition
}

func (condition *Condition) In(key string, values interface{}) *Condition {
	condition.LogicOperates = append(condition.LogicOperates, op.NewLogic(constants.AND, false).AddRelationOperate(op.In(key, values)))
	return condition
}

func (condition *Condition) NotIn(key string, values interface{}) *Condition {
	condition.LogicOperates = append(condition.LogicOperates, op.NewLogic(constants.AND, false).AddRelationOperate(op.NotIn(key, values)))
	return condition
}

func (condition *Condition) Like(key string, value interface{}) *Condition {
	condition.LogicOperates = append(condition.LogicOperates, op.NewLogic(constants.AND, false).AddRelationOperate(op.Like(key, value)))
	return condition
}

func (condition *Condition) NotLike(key string, value interface{}) *Condition {
	condition.LogicOperates = append(condition.LogicOperates, op.NewLogic(constants.AND, false).AddRelationOperate(op.NotLike(key, value)))
	return condition
}

func (condition *Condition) Match(key string, value interface{}) *Condition {
	condition.LogicOperates = append(condition.LogicOperates, op.NewLogic(constants.AND, false).AddRelationOperate(op.Match(key, value)))
	return condition
}

func (condition *Condition) NotMatch(key string, value interface{}) *Condition {
	condition.LogicOperates = append(condition.LogicOperates, op.NewLogic(constants.AND, false).AddRelationOperate(op.NotMatch(key, value)))
	return condition
}

func (condition *Condition) NMatch(key string, value interface{}) *Condition {
	condition.LogicOperates = append(condition.LogicOperates, op.NewLogic(constants.AND, false).AddRelationOperate(op.NMatch(key, value)))
	return condition
}

func (condition *Condition) Contains(key string, value interface{}) *Condition {
	condition.LogicOperates = append(condition.LogicOperates, op.NewLogic(constants.AND, false).AddRelationOperate(op.Contains(key, value)))
	return condition
}

func parseLogicOperator(logicCmd string, parameters ...any) []*op.LogicOperate {
	if parameters == nil || len(parameters) == 0 {
		return []*op.LogicOperate{}
	}
	var operates []*op.LogicOperate
	// key-value处理：key必须为String类型，key后面必须为对应的value，kv形式默认转为无括号的and
	for index := 0; index < len(parameters); index++ {
		data := parameters[index]
		if reflect.TypeOf(data).Kind() == reflect.String {
			key := data.(string)
			var value interface{}
			haveValue := false
			if (index + 1) < len(parameters) {
				value = parameters[index+1]
				haveValue = true
			}

			switch logicCmd {
			case logicAnd:
				operates = append(operates, op.NewLogic(constants.AND, true).AddRelationOperate(op.NewRelation(key, constants.EQUAL, value)))
				break
			case logicAndEm:
				operates = append(operates, op.NewLogic(constants.AND, false).AddRelationOperate(op.NewRelation(key, constants.EQUAL, value)))
				break
			case logicOr:
				operates = append(operates, op.NewLogic(constants.OR, true).AddRelationOperate(op.NewRelation(key, constants.EQUAL, value)))
				break
			case logicOrEm:
				operates = append(operates, op.NewLogic(constants.OR, false).AddRelationOperate(op.NewRelation(key, constants.EQUAL, value)))
				break
			case logicEmpty:
				if haveValue {
					operates = append(operates, op.NewLogic(constants.EMPTY, false).AddRelationOperate(op.NewRelation(key, constants.EQUAL, value)))
				} else {
					operates = append(operates, op.NewLogic(constants.EMPTY, false).AddRelationOperate(op.NewRelationNoneSymbolAndValue(key)))
				}
			default:
				operates = append(operates, op.NewLogic(constants.EMPTY, false).AddRelationOperate(op.NewRelation(key, constants.EQUAL, value)))
			}
			index++
		} else if reflect.TypeOf(data) == reflect.PointerTo(reflect.TypeOf((*op.LogicOperate)(nil))).Elem() {
			operates = append(operates, data.(*op.LogicOperate))
		} else if reflect.TypeOf(data) == reflect.PointerTo(reflect.TypeOf(([]*op.LogicOperate)(nil))).Elem() {
			operates = append(operates, data.([]*op.LogicOperate)...)
		} else if reflect.TypeOf(data) == reflect.PointerTo(reflect.TypeOf((*op.RelationOperate)(nil))).Elem() {
			operates = append(operates, op.Append(data.(*op.RelationOperate)))
		} else if reflect.TypeOf(data) == reflect.PointerTo(reflect.TypeOf((*maps.GoleMap)(nil))).Elem() {
			pMap := data.(*maps.GoleMap)
			for _, key := range pMap.Keys() {
				val, _ := pMap.Get(key)
				operates = append(operates, op.AndEm(key, val))
			}
		}
	}
	return operates
}

func (condition *Condition) ToValues() []interface{} {
	var parameters []interface{}
	for _, logicOperate := range condition.LogicOperates {
		parameters = append(parameters, getValuesFromLogicOperate(logicOperate)...)
	}
	return parameters
}

func getValuesFromLogicOperate(logicOperate *op.LogicOperate) []interface{} {
	var parameters []interface{}
	if logicOperate.ChildLogicOperates == nil || len(logicOperate.ChildLogicOperates) == 0 {
		// 那么这个就是叶节点，搜集values
		for _, relationOperate := range logicOperate.RelationOperates {
			relSymbol := relationOperate.RelationSymbol
			switch relSymbol {
			case constants.EQUAL, constants.NOT_EQUAL, constants.GREATHER_THAN, constants.GREATHER_EQUAL, constants.LESS_THAN, constants.LESS_EQUAL,
				constants.LIKE, constants.NOT_LIKE, constants.MATCH, constants.NOT_MATCH, constants.CONTAINS:

				if reflect.TypeOf(relationOperate.Value).Kind() == reflect.String {
					parameters = append(parameters, strChange(relationOperate.Value.(string)))
				} else {
					parameters = append(parameters, relationOperate.Value)
				}
			case constants.IS_NULL, constants.IS_NOT_NULL:
				break
			case constants.BETWEEN:
				parameters = append(parameters, relationOperate.LeftValue)
				parameters = append(parameters, relationOperate.RightValue)
			case constants.NOT_BETWEEN:
				parameters = append(parameters, relationOperate.LeftValue)
				parameters = append(parameters, relationOperate.RightValue)
			case constants.IN, constants.NOT_IN:
				values := relationOperate.Values
				if values == nil {
					break
				}
				collectionValues := reflect.ValueOf(values)
				arrayLen := collectionValues.Len()
				if arrayLen == 0 {
					break
				}
				for arrayIndex := 0; arrayIndex < arrayLen; arrayIndex++ {
					if collectionValues.Index(arrayIndex).Kind() == reflect.String {
						parameters = append(parameters, strChange(collectionValues.Index(arrayIndex).Interface().(string)))
					} else {
						parameters = append(parameters, collectionValues.Index(arrayIndex))
					}
				}
			}
		}
	} else {
		// 非叶节点
		for _, childLogicOperate := range logicOperate.ChildLogicOperates {
			parameters = append(parameters, getValuesFromLogicOperate(childLogicOperate)...)
		}
	}
	return parameters
}

// tdengine的sql的执行时候发现如果是string类型则需要'xxx'这样才能识别
func strChange(val string) string {
	if strings.Contains(val, "now") {
		return val
	}
	if !strings.HasPrefix(val, "'") {
		val = "'" + val
	}
	if !strings.HasSuffix(val, "'") {
		val = val + "'"
	}
	return val
}

const (
	logicAnd   = "and"
	logicAndEm = "andEm"
	logicOr    = "or"
	logicOrEm  = "orEm"
	logicEmpty = "empty"
)
