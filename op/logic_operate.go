/*
逻辑运算符
and or
*/

package op

import (
	"strings"

	"github.com/simonalong/tdorm/constants"
)

type LogicOperate struct {
	LogicSymbol      string
	RelationOperates []*RelationOperate
	// 是否有括号
	HaveBracket        bool
	ChildLogicOperates []*LogicOperate
}

func NewLogic(symbol string, haveBracket bool) *LogicOperate {
	return &LogicOperate{
		LogicSymbol:        symbol,
		RelationOperates:   []*RelationOperate{},
		HaveBracket:        haveBracket,
		ChildLogicOperates: []*LogicOperate{},
	}
}

func (receiver *LogicOperate) AddRelationOperate(op *RelationOperate) *LogicOperate {
	if op == nil {
		return receiver
	}
	receiver.RelationOperates = append(receiver.RelationOperates, op)
	return receiver
}

func (receiver *LogicOperate) AddRelationOperates(ops []*RelationOperate) *LogicOperate {
	if ops == nil || len(ops) == 0 {
		return receiver
	}
	receiver.RelationOperates = append(receiver.RelationOperates, ops...)
	return receiver
}

func (receiver *LogicOperate) AddChildLogicOperate(child *LogicOperate) *LogicOperate {
	if child == nil {
		return receiver
	}
	receiver.ChildLogicOperates = append(receiver.ChildLogicOperates, child)
	return receiver
}

func (receiver *LogicOperate) AddChildrenLogicOperate(children []*LogicOperate) *LogicOperate {
	if children == nil || len(children) == 0 {
		return receiver
	}
	receiver.ChildLogicOperates = append(receiver.ChildLogicOperates, children...)
	return receiver
}

func (receiver *LogicOperate) NeedWhere() bool {
	return true
}

func (receiver *LogicOperate) GenerateSqlOfFull() string {
	OperateSymbol := receiver.LogicSymbol
	if len(receiver.ChildLogicOperates) == 0 {
		childOps := receiver.RelationOperates

		var ops []string
		for _, op := range childOps {
			if val := op.GenerateOperateSqlOfFull(); val != "" {
				ops = append(ops, val)
			}
		}
		if receiver.HaveBracket {
			if len(ops) != 0 {
				return constants.EMPTY + OperateSymbol + " (" + strings.Join(ops, OperateSymbol) + ")"
			}
		} else {
			if len(ops) != 0 {
				return constants.EMPTY + OperateSymbol + constants.EMPTY + strings.Join(ops, OperateSymbol)
			}
		}
		return ""
	} else {
		var logicSqlParts []string
		for _, logicOperate := range receiver.ChildLogicOperates {
			if val := trimLogic(logicOperate.GenerateSqlOfFull()); val != "" {
				logicSqlParts = append(logicSqlParts, val)
			}
		}

		if receiver.HaveBracket {
			if len(logicSqlParts) != 0 {
				return constants.EMPTY + OperateSymbol + " (" + strings.TrimSpace(strings.Join(logicSqlParts, constants.EMPTY+OperateSymbol)) + ")"
			}
		} else {
			if len(logicSqlParts) != 0 {
				return constants.EMPTY + OperateSymbol + constants.EMPTY + strings.TrimSpace(strings.Join(logicSqlParts, constants.EMPTY+OperateSymbol))
			}
		}
		return ""
	}
}

func (receiver *LogicOperate) GenerateSql() string {
	OperateSymbol := receiver.LogicSymbol
	if len(receiver.ChildLogicOperates) == 0 {
		childOps := receiver.RelationOperates

		var ops []string
		for _, op := range childOps {
			if val := op.GenerateOperateSql(); val != "" {
				ops = append(ops, val)
			}
		}
		if receiver.HaveBracket {
			if len(ops) != 0 {
				return constants.EMPTY + OperateSymbol + " (" + strings.Join(ops, OperateSymbol) + ")"
			}
		} else {
			if len(ops) != 0 {
				return constants.EMPTY + OperateSymbol + constants.EMPTY + strings.Join(ops, OperateSymbol)
			}
		}
		return ""
	} else {
		var logicSqlParts []string
		for _, logicOperate := range receiver.ChildLogicOperates {
			if val := trimLogic(logicOperate.GenerateSql()); val != "" {
				logicSqlParts = append(logicSqlParts, val)
			}
		}

		if receiver.HaveBracket {
			if len(logicSqlParts) != 0 {
				return constants.EMPTY + OperateSymbol + " (" + strings.TrimSpace(strings.Join(logicSqlParts, constants.EMPTY+OperateSymbol)) + ")"
			}
		} else {
			if len(logicSqlParts) != 0 {
				return constants.EMPTY + OperateSymbol + constants.EMPTY + strings.TrimSpace(strings.Join(logicSqlParts, constants.EMPTY+OperateSymbol))
			}
		}
		return ""
	}
}

func (receiver *LogicOperate) ValueLegal() bool {
	return true
}

func trimLogic(logicSql string) string {
	if len(logicSql) != 0 {
		logicSql = strings.TrimSpace(logicSql)
		if strings.HasPrefix(logicSql, constants.AND) {
			logicSql = logicSql[len(constants.AND):]
		}

		logicSql = strings.TrimSpace(logicSql)
		if strings.HasPrefix(logicSql, constants.OR) {
			logicSql = logicSql[len(constants.OR):]
		}
		logicSql = strings.TrimSpace(logicSql)
	}
	if logicSql != "" {
		return constants.EMPTY + strings.TrimSpace(logicSql)
	} else {
		return logicSql
	}
}
