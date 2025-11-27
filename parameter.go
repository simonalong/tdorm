package tdorm

/* 按照官方的Sql的查询类语法
SELECT [hints] [DISTINCT] [TAGS] select_list
    from_clause
    [WHERE condition]
    [partition_by_clause]
    [interp_clause]
    [window_clause]
    [group_by_clause]
    [order_by_clause]
    [SLIMIT limit_val [OFFSET offset_val]]
    [LIMIT limit_val [OFFSET offset_val]]
    [>> export_file]
*/

// hints: /*+ [hint([hint_param_list])] [hint([hint_param_list])] */
/*
hint:
    BATCH_SCAN | NO_BATCH_SCAN | SORT_FOR_GROUP | PARTITION_FIRST | PARA_TABLES_SORT | SMALLDATA_TS_SORT

select_list:
    select_expr [, select_expr] ...

select_expr: {
    *
  | query_name.*
  | [schema_name.] {table_name | view_name} .*
  | t_alias.*
  | expr [[AS] c_alias]
}

from_clause: {
    table_reference [, table_reference] ...
  | table_reference join_clause [, join_clause] ...
}

table_reference:
    table_expr t_alias

table_expr: {
    table_name
  | view_name
  | ( subquery )
}

join_clause:
    [INNER|LEFT|RIGHT|FULL] [OUTER|SEMI|ANTI|ASOF|WINDOW] JOIN table_reference [ON condition] [WINDOW_OFFSET(start_offset, end_offset)] [LIMIT limit_num]

window_clause: {
    SESSION(ts_col, tol_val)
  | STATE_WINDOW(col)
  | INTERVAL(interval_val [, interval_offset]) [SLIDING (sliding_val)] [WATERMARK(watermark_val)] [FILL(fill_mod_and_val)]
  | EVENT_WINDOW START WITH start_trigger_condition END WITH end_trigger_condition
  | COUNT_WINDOW(count_val[, sliding_val])

interp_clause:
    RANGE(ts_val [, ts_val]) EVERY(every_val) FILL(fill_mod_and_val)

partition_by_clause:
    PARTITION BY expr [, expr] ...

group_by_clause:
    GROUP BY expr [, expr] ... HAVING condition

order_by_clause:
    ORDER BY order_expr [, order_expr] ...

order_expr:
    {expr | position | c_alias} [DESC | ASC] [NULLS FIRST | NULLS LAST]
*/
