package tdorm

import (
	"testing"

	"github.com/simonalong/tdorm/column"
	"github.com/simonalong/tdorm/condition"
	"github.com/stretchr/testify/assert"
)

//SELECT {DATABASE() | CLIENT_VERSION() | SERVER_VERSION() | SERVER_STATUS() | NOW() | TODAY() | TIMEZONE() | CURRENT_USER() | USER() }
//
//SELECT [hints] [DISTINCT] [TAGS] select_list
//    from_clause
//    [WHERE condition]
//    [partition_by_clause]
//    [interp_clause]
//    [window_clause]
//    [group_by_clause]
//    [order_by_clasue]
//    [SLIMIT limit_val [SOFFSET offset_val]]
//    [LIMIT limit_val [OFFSET offset_val]]
//    [>> export_file]
//
//hints: /*+ [hint([hint_param_list])] [hint([hint_param_list])] */
//
//hint:
//    BATCH_SCAN | NO_BATCH_SCAN | SORT_FOR_GROUP | PARTITION_FIRST | PARA_TABLES_SORT | SMALLDATA_TS_SORT
//
//select_list:
//    select_expr [, select_expr] ...
//
//select_expr: {
//    *
//  | query_name.*
//  | [schema_name.] {table_name | view_name} .*
//  | t_alias.*
//  | expr [[AS] c_alias]
//}
//
//from_clause: {
//    table_reference [, table_reference] ...
//  | table_reference join_clause [, join_clause] ...
//}
//
//table_reference:
//    table_expr t_alias
//
//table_expr: {
//    table_name
//  | view_name
//  | ( subquery )
//}
//
//join_clause:
//    [INNER|LEFT|RIGHT|FULL] [OUTER|SEMI|ANTI|ASOF|WINDOW] JOIN table_reference [ON condition] [WINDOW_OFFSET(start_offset, end_offset)] [JLIMIT jlimit_num]
//
//window_clause: {
//    SESSION(ts_col, tol_val)
//  | STATE_WINDOW(col)
//  | INTERVAL(interval_val [, interval_offset]) [SLIDING (sliding_val)] [WATERMARK(watermark_val)] [FILL(fill_mod_and_val)]
//  | EVENT_WINDOW START WITH start_trigger_condition END WITH end_trigger_condition
//  | COUNT_WINDOW(count_val[, sliding_val])
//
//interp_clause:
//    RANGE(ts_val [, ts_val]) EVERY(every_val) FILL(fill_mod_and_val)
//
//partition_by_clause:
//    PARTITION BY expr [, expr] ...
//
//group_by_clause:
//    GROUP BY expr [, expr] ... HAVING condition
//
//order_by_clasue:
//    ORDER BY order_expr [, order_expr] ...
//
//order_expr:
//    {expr | position | c_alias} [DESC | ASC] [NULLS FIRST | NULLS LAST]

// 提示测试
func TestSelectSqlHits(t *testing.T) {
	assert.Equal(t, "select /*+ BATCH_SCAN() */ `ts` from td_china", generateSelectSql(column.Hints(&column.Hint{BatchScan: true}).Of("ts"), "td_china", nil))
	assert.Equal(t, "select /*+ NO_BATCH_SCAN() */ `ts` from td_china", generateSelectSql(column.Hints(&column.Hint{NoBatchScan: true}).Of("ts"), "td_china", nil))
	assert.Equal(t, "select /*+ SORT_FOR_GROUP() */ `ts` from td_china", generateSelectSql(column.Hints(&column.Hint{SortForGroup: true}).Of("ts"), "td_china", nil))
	assert.Equal(t, "select /*+ PARTITION_FIRST() */ `ts` from td_china", generateSelectSql(column.Hints(&column.Hint{PartitionFirst: true}).Of("ts"), "td_china", nil))
	assert.Equal(t, "select /*+ PARA_TABLES_SORT() */ `ts` from td_china", generateSelectSql(column.Hints(&column.Hint{ParaTablesSort: true}).Of("ts"), "td_china", nil))
	assert.Equal(t, "select /*+ SMALLDATA_TS_SORT() */ `ts` from td_china", generateSelectSql(column.Hints(&column.Hint{SmalldataTsSort: true}).Of("ts"), "td_china", nil))
	assert.Equal(t, "select /*+ SKIP_TSMA() */ `ts` from td_china", generateSelectSql(column.Hints(&column.Hint{SkipTsma: true}).Of("ts"), "td_china", nil))

	// 多个提示
	assert.Equal(t, "select /*+ BATCH_SCAN() SKIP_TSMA() */ `ts` from td_china", generateSelectSql(column.Hints(&column.Hint{BatchScan: true, SkipTsma: true}).Of("ts"), "td_china", nil))
}

// distinct 测试
func TestDistinct(t *testing.T) {
	assert.Equal(t, "select distinct `name`, `age` from td_china", generateSelectSql(column.Distinct().Of("name", "age"), "td_china", nil))
	assert.Equal(t, "select distinct distinct name, `age` from td_china", generateSelectSql(column.Distinct().Of("distinct name", "age"), "td_china", nil))

	// ！！！！！！！！错误用法！！！！！！！！distinct 作为属性名了，不要这样用
	assert.Equal(t, "select `distinct`, `name`, `age` from td_china", generateSelectSql(column.Of("distinct", "name", "age"), "td_china", nil))
}

// tags 测试
func TestTags(t *testing.T) {
	assert.Equal(t, "select tags `station` from td_china", generateSelectSql(column.Tags().Of("station"), "td_china", nil))
	assert.Equal(t, "select `station` from td_china", generateSelectSql(column.Of("station"), "td_china", nil))
}

// 其他 测试
//
//	[partition_by_clause]
//	[interp_clause]
//	[window_clause]
//	[group_by_clause]
//	[order_by_clasue]
//	[SLIMIT limit_val [SOFFSET offset_val]]
//	[LIMIT limit_val [OFFSET offset_val]]
func TestCondition(t *testing.T) {
	assert.Equal(t, "select `name` from td_china limit 10", generateSelectSql(column.Of("name"), "td_china", condition.New().Append("limit 10")))
	assert.Equal(t, "select `address`, avg(age) from td_china partition by address", generateSelectSql(column.Of("address", "avg(age)"), "td_china", condition.New().Append("partition by address")))

	// 以下用官网的用例然后用我们的语法表示
	sql := generateSelectSql(column.Of("_wstart", "location", "max(current)"), "meters", condition.New().Append("partition by location interval(10m)"))
	assert.Equal(t, "select _wstart, `location`, max(current) from meters partition by location interval(10m)", sql)

	// SELECT * FROM (SELECT COUNT(*) AS cnt, FIRST(ts) AS fst, status FROM temp_tb_1 STATE_WINDOW(status)) t WHERE status = 2;
	sql = generateSelectSql(column.Of("*"), "(select count(*) as cnt, first(ts) as fst, status from tem_tb_1 stat_window(status)) t", condition.New().Eq("status", 2))
	assert.Equal(t, "select * from (select count(*) as cnt, first(ts) as fst, status from tem_tb_1 stat_window(status)) t where `status` = ?", sql)

	// SELECT tbname, _wstart, CASE WHEN voltage >= 205 and voltage <= 235 THEN 1 ELSE 0 END status FROM meters PARTITION BY tbname STATE_WINDOW(CASE WHEN voltage >= 205 and voltage <= 235 THEN 1 ELSE 0 END);
	sql = generateSelectSql(column.Of("tbname", "_wstart", "case when voltage >= 205 and voltage <= 235 then 1 else 0 end status"), "meters", condition.New().Append("partition by tbname state_window(case when voltage >= 205 and voltage <= 235 then 1 else 0 end)"))
	assert.Equal(t, "select tbname, _wstart, case when voltage >= 205 and voltage <= 235 then 1 else 0 end status from meters partition by tbname state_window(case when voltage >= 205 and voltage <= 235 then 1 else 0 end)", sql)

	//SELECT COUNT(*), FIRST(ts) FROM temp_tb_1 SESSION(ts, tol_val)
	sql = generateSelectSql(column.Of("count(*)", "first(ts)"), "temp_tb_1", condition.New().Append("session(ts, tol_val)"))
	assert.Equal(t, "select count(*), first(ts) from temp_tb_1 session(ts, tol_val)", sql)
}

func TestJoin(t *testing.T) {
	// 以下用例，全部来自于官网

	// SELECT ... FROM table_name1 [INNER] JOIN table_name2 [ON ...] [WHERE ...] [...]
	// 或
	// SELECT ... FROM table_name1, table_name2 WHERE ... [...]
	// 内连接：SELECT a.ts, a.voltage, b.voltage FROM d1001 a JOIN d1002 b ON a.ts = b.ts and a.voltage > 220 and b.voltage > 220
	sql := generateSelectSql(column.Of("a.ts", "a.voltage", "b.voltage"), "d1001 a join d1002 b on a.ts = b.ts and a.voltage > 220 and b.voltage > 220", nil)
	assert.Equal(t, "select a.`ts`, a.`voltage`, b.`voltage` from d1001 a join d1002 b on a.ts = b.ts and a.voltage > 220 and b.voltage > 220", sql)

	sql = generateSelectSql(column.Of("a.ts", "a.voltage", "b.voltage"), "d1001 a join d1002 b on a.ts = b.ts and a.voltage > 220 and b.voltage > 220", condition.New().Gt("a.data", 12))
	assert.Equal(t, "select a.`ts`, a.`voltage`, b.`voltage` from d1001 a join d1002 b on a.ts = b.ts and a.voltage > 220 and b.voltage > 220 where a.`data` > ?", sql)

	// 左（右）连接： SELECT ... FROM table_name1 LEFT|RIGHT [OUTER] JOIN table_name2 ON ... [WHERE ...] [...]
	// SELECT a.ts, a.voltage, b.voltage FROM d1001 a LEFT JOIN d1002 b ON a.ts = b.ts and a.voltage > 220 and b.voltage > 220
}
