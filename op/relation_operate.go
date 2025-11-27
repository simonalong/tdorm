/*
关系运算符
=, !=, <>, >, >=, <, <=
IS [NOT] NULL
[NOT] BETWEEN AND
[NOT] IN
[NOT] LIKE
MATCH, NMATCH
CONTAINS
*/

package op

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/simonalong/gole/logger"
	"github.com/simonalong/tdorm/column"
	"github.com/simonalong/tdorm/constants"
)

type RelationOperate struct {
	ColumnName     string
	RelationSymbol string
	Value          interface{}

	// leftValue和rightValue主要用于"between and" 和 "not between and"
	LeftValue  interface{}
	RightValue interface{}

	// values 主要是给in使用的，注意这个interface{}这个是数组
	Values interface{}
}

func NewRelation(key, relationSymbol string, value interface{}) *RelationOperate {
	return &RelationOperate{
		ColumnName:     key,
		RelationSymbol: relationSymbol,
		Value:          value,
	}
}

func NewRelationMultiValue(key, relationSymbol string, leftValue, rightValue interface{}) *RelationOperate {
	return &RelationOperate{
		ColumnName:     key,
		RelationSymbol: relationSymbol,
		LeftValue:      leftValue,
		RightValue:     rightValue,
	}
}

func NewRelationValues(key, relationSymbol string, values interface{}) *RelationOperate {
	if reflect.TypeOf(values).Kind() != reflect.Slice && reflect.TypeOf(values).Kind() != reflect.Array {
		logger.Warnf("搜索条件使用有误，values应该为数组或者切片类型")
		return nil
	}
	return &RelationOperate{
		ColumnName:     key,
		RelationSymbol: relationSymbol,
		Values:         values,
	}
}

func NewRelationNoneValue(key, relationSymbol string) *RelationOperate {
	return &RelationOperate{
		ColumnName:     key,
		RelationSymbol: relationSymbol,
	}
}

func NewRelationNoneSymbolAndValue(key string) *RelationOperate {
	return &RelationOperate{
		ColumnName: key,
	}
}

func Eq(key string, value interface{}) *RelationOperate {
	return NewRelation(key, constants.EQUAL, value)
}

func UnEq(key string, value interface{}) *RelationOperate {
	return NewRelation(key, constants.NOT_EQUAL, value)
}

// Gt 大于等于(greater than)：>=
func Gt(key string, value interface{}) *RelationOperate {
	return NewRelation(key, constants.GREATHER_THAN, value)
}

// Ge 大于等于(greater equal)：>=
func Ge(key string, value interface{}) *RelationOperate {
	return NewRelation(key, constants.GREATHER_EQUAL, value)
}

// Lt 大于等于(less than)：>=
func Lt(key string, value interface{}) *RelationOperate {
	return NewRelation(key, constants.LESS_THAN, value)
}

// Le 大于等于(less equal)：>=
func Le(key string, value interface{}) *RelationOperate {
	return NewRelation(key, constants.LESS_EQUAL, value)
}

func IsNull(key string) *RelationOperate {
	return NewRelationNoneValue(key, constants.IS_NULL)
}

func IsNotNull(key string) *RelationOperate {
	return NewRelationNoneValue(key, constants.IS_NOT_NULL)
}

func BetweenAnd(key string, leftValue, rightValue interface{}) *RelationOperate {
	return NewRelationMultiValue(key, constants.BETWEEN, leftValue, rightValue)
}

func NotBetweenAnd(key string, leftValue, rightValue interface{}) *RelationOperate {
	return NewRelationMultiValue(key, constants.NOT_BETWEEN, leftValue, rightValue)
}

// In in (?, ?, ?)
// 注意values 这个为数组或者切片类型
func In(key string, values interface{}) *RelationOperate {
	return NewRelationValues(key, constants.IN, values)
}

// NotIn not in (?, ?, ?)
// 注意values 这个为数组或者切片类型
func NotIn(key string, values interface{}) *RelationOperate {
	return NewRelationValues(key, constants.NOT_IN, values)
}

func Like(key string, value interface{}) *RelationOperate {
	return NewRelation(key, constants.LIKE, value)
}

func NotLike(key string, value interface{}) *RelationOperate {
	return NewRelation(key, constants.NOT_LIKE, value)
}

// Match 正则表达式匹配
func Match(key string, value interface{}) *RelationOperate {
	return NewRelation(key, constants.MATCH, value)
}

// NMatch 正则表达式匹配的反向，同 NotMatch
func NMatch(key string, value interface{}) *RelationOperate {
	return NewRelation(key, constants.NOT_MATCH, value)
}

// NotMatch 正则表达式匹配的反向，同 NMatch
func NotMatch(key string, value interface{}) *RelationOperate {
	return NewRelation(key, constants.NOT_MATCH, value)
}

// Contains 用于json字段中的包含
func Contains(key string, value interface{}) *RelationOperate {
	return NewRelation(key, constants.CONTAINS, value)
}

// GenerateOperateSqlOfFull value值不是占位符，而是完整的值
func (receiver *RelationOperate) GenerateOperateSqlOfFull() string {
	key := receiver.ColumnName
	value := receiver.Value
	symbol := receiver.RelationSymbol
	switch symbol {
	case constants.EQUAL, constants.NOT_EQUAL, constants.GREATHER_THAN, constants.GREATHER_EQUAL, constants.LESS_THAN, constants.LESS_EQUAL, constants.LIKE, constants.NOT_LIKE, constants.MATCH, constants.NOT_MATCH, constants.CONTAINS:
		if value != nil {
			return fmt.Sprintf("%v %v %v %v %v", column.ToDbField(key), constants.EMPTY, symbol, constants.EMPTY, ValueToStr(value))
		}
	case constants.IS_NULL, constants.IS_NOT_NULL:
		return column.ToDbField(key) + constants.EMPTY + symbol
	case constants.BETWEEN:
		return fmt.Sprintf("%v %v %v %v %v and %v", column.ToDbField(key), constants.EMPTY, symbol, constants.EMPTY, ValueToStr(receiver.LeftValue), ValueToStr(receiver.RightValue))
	case constants.NOT_BETWEEN:
		return fmt.Sprintf("%v %v %v %v %v and %v", column.ToDbField(key), constants.EMPTY, symbol, constants.EMPTY, ValueToStr(receiver.LeftValue), ValueToStr(receiver.RightValue))
	case constants.IN, constants.NOT_IN:
		values := receiver.Values
		if values == nil {
			return ""
		}
		collectionValues := reflect.ValueOf(values)
		arrayLen := collectionValues.Len()
		if arrayLen == 0 {
			return ""
		}
		var metaValues []string
		for arrayIndex := 0; arrayIndex < arrayLen; arrayIndex++ {
			fieldValueItem := collectionValues.Index(arrayIndex)
			metaValues = append(metaValues, ValueToStr(fieldValueItem.Interface()))
		}

		return column.ToDbField(key) + constants.EMPTY + symbol + constants.EMPTY + "(" + strings.Join(metaValues, ", ") + ")"
	}
	return ""
}

func (receiver *RelationOperate) GenerateOperateSql() string {
	key := receiver.ColumnName
	value := receiver.Value
	symbol := receiver.RelationSymbol
	switch symbol {
	case constants.EQUAL, constants.NOT_EQUAL, constants.GREATHER_THAN, constants.GREATHER_EQUAL, constants.LESS_THAN, constants.LESS_EQUAL, constants.LIKE, constants.NOT_LIKE, constants.MATCH, constants.NOT_MATCH, constants.CONTAINS:
		if value != nil {
			return column.ToDbField(key) + constants.EMPTY + symbol + constants.EMPTY + "?"
		}
	case constants.IS_NULL, constants.IS_NOT_NULL:
		return column.ToDbField(key) + constants.EMPTY + symbol
	case constants.BETWEEN:
		return column.ToDbField(key) + constants.EMPTY + symbol + constants.EMPTY + "? and ?"
	case constants.NOT_BETWEEN:
		return column.ToDbField(key) + constants.EMPTY + symbol + constants.EMPTY + "? and ?"
	case constants.IN, constants.NOT_IN:
		values := receiver.Values
		if values == nil {
			return ""
		}
		collectionValues := reflect.ValueOf(values)
		arrayLen := collectionValues.Len()
		if arrayLen == 0 {
			return ""
		}
		var seizeStr []string
		for arrayIndex := 0; arrayIndex < arrayLen; arrayIndex++ {
			seizeStr = append(seizeStr, "?")
		}
		return column.ToDbField(key) + constants.EMPTY + symbol + constants.EMPTY + "(" + strings.Join(seizeStr, ", ") + ")"
	}
	return ""
}
