package op

import (
	"reflect"

	"github.com/simonalong/gole/maps"
	"github.com/simonalong/tdorm/constants"
)

// And 输出：(`age` = ? and `name` = ?)
func And(parameters ...any) *LogicOperate {
	var operates []*LogicOperate
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

			if haveValue {
				operates = append(operates, NewLogic(constants.AND, false).AddRelationOperate(Eq(key, value)))
			}
			index++
		} else if reflect.TypeOf(data) == reflect.PointerTo(reflect.TypeOf((*LogicOperate)(nil))).Elem() {
			operates = append(operates, data.(*LogicOperate))
		} else if reflect.TypeOf(data) == reflect.PointerTo(reflect.TypeOf((*RelationOperate)(nil))).Elem() {
			operates = append(operates, Append(data.(*RelationOperate)))
		} else if reflect.TypeOf(data) == reflect.PointerTo(reflect.TypeOf((*maps.GoleMap)(nil))).Elem() {
			pMap := data.(*maps.GoleMap)
			for _, key := range pMap.Keys() {
				val, _ := pMap.Get(key)
				operates = append(operates, AndEm(key, val))
			}
		}
	}

	return NewLogic(constants.AND, true).AddChildrenLogicOperate(operates)
}

// AndEm 输出：`age` = ? and `name` = ?
func AndEm(parameters ...any) *LogicOperate {
	var operates []*LogicOperate
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

			if haveValue {
				operates = append(operates, NewLogic(constants.AND, false).AddRelationOperate(Eq(key, value)))
			}
			index++
		} else if reflect.TypeOf(data) == reflect.PointerTo(reflect.TypeOf((*LogicOperate)(nil))).Elem() {
			operates = append(operates, data.(*LogicOperate))
		} else if reflect.TypeOf(data) == reflect.PointerTo(reflect.TypeOf((*RelationOperate)(nil))).Elem() {
			operates = append(operates, Append(data.(*RelationOperate)))
		} else if reflect.TypeOf(data) == reflect.PointerTo(reflect.TypeOf((*maps.GoleMap)(nil))).Elem() {
			pMap := data.(*maps.GoleMap)
			for _, key := range pMap.Keys() {
				val, _ := pMap.Get(key)
				operates = append(operates, AndEm(key, val))
			}
		}
	}
	return NewLogic(constants.AND, false).AddChildrenLogicOperate(operates)
}

// Or 输出：(`age` = ? or `name` = ?)
func Or(parameters ...any) *LogicOperate {
	var operates []*LogicOperate
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

			if haveValue {
				operates = append(operates, NewLogic(constants.OR, false).AddRelationOperate(Eq(key, value)))
			}
			index++
		} else if reflect.TypeOf(data) == reflect.PointerTo(reflect.TypeOf((*LogicOperate)(nil))).Elem() {
			operates = append(operates, data.(*LogicOperate))
		} else if reflect.TypeOf(data) == reflect.PointerTo(reflect.TypeOf((*RelationOperate)(nil))).Elem() {
			operates = append(operates, Append(data.(*RelationOperate)))
		} else if reflect.TypeOf(data) == reflect.PointerTo(reflect.TypeOf((*maps.GoleMap)(nil))).Elem() {
			pMap := data.(*maps.GoleMap)
			for _, key := range pMap.Keys() {
				val, _ := pMap.Get(key)
				operates = append(operates, OrEm(key, val))
			}
		}
	}

	return NewLogic(constants.OR, true).AddChildrenLogicOperate(operates)
}

// OrEm 输出：`age` = ? or `name` = ?
func OrEm(parameters ...any) *LogicOperate {
	var operates []*LogicOperate
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

			if haveValue {
				operates = append(operates, NewLogic(constants.OR, false).AddRelationOperate(Eq(key, value)))
			}
			index++
		} else if reflect.TypeOf(data) == reflect.PointerTo(reflect.TypeOf((*LogicOperate)(nil))).Elem() {
			operates = append(operates, data.(*LogicOperate))
		} else if reflect.TypeOf(data) == reflect.PointerTo(reflect.TypeOf((*RelationOperate)(nil))).Elem() {
			operates = append(operates, Append(data.(*RelationOperate)))
		} else if reflect.TypeOf(data) == reflect.PointerTo(reflect.TypeOf((*maps.GoleMap)(nil))).Elem() {
			pMap := data.(*maps.GoleMap)
			for _, key := range pMap.Keys() {
				val, _ := pMap.Get(key)
				operates = append(operates, OrEm(key, val))
			}
		}
	}
	return NewLogic(constants.OR, false).AddChildrenLogicOperate(operates)
}

func Append(operate *RelationOperate) *LogicOperate {
	return NewLogic(constants.EMPTY, false).AddRelationOperate(operate)
}
