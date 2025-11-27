package op

import "fmt"

type Clause struct {
	// 	  [interp_clause]
	//    [window_clause]
	//    [group_by_clause]
	//    [order_by_clasue]
	ClauseSymbol string
	ColumnName   string
	AppendTail   string
}

func (receiver *Clause) GenerateSql() string {
	return fmt.Sprintf("%v %v %v", receiver.ClauseSymbol, receiver.ColumnName, receiver.AppendTail)
}
