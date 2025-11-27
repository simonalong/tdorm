package column

import (
	"strings"

	"github.com/simonalong/gole/util"
	"github.com/simonalong/tdorm/constants"
)

const dom = "`"

/* 按照官方的Sql的查询类语法
SELECT [hints] [DISTINCT] [TAGS] select_list
    from_clause
    [WHERE condition]
    [partition_by_clause]
    [interp_clause]
    [window_clause]
    [group_by_clause]
    [order_by_clasue]
    [SLIMIT limit_val [SOFFSET offset_val]]
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
*/

// 我们这里的columns 提供的是 [hints] [DISTINCT] [TAGS] select_list 部分的sql拼接功能

type Columns struct {
	hints    *Hint
	distinct bool
	tags     bool
	fields   []string
}

func Of(datas ...string) *Columns {
	if len(datas) == 0 || (len(datas) == 1 && datas[0] == "*") {
		return &Columns{}
	}
	return &Columns{
		fields: datas,
	}
}

func Hints(hints *Hint) *Columns {
	return &Columns{
		hints: hints,
	}
}

func Tags() *Columns {
	return &Columns{
		tags: true,
	}
}

func Distinct() *Columns {
	return &Columns{
		distinct: true,
	}
}

func (columns *Columns) Tags() *Columns {
	columns.tags = true
	return columns
}

func (columns *Columns) Distinct() *Columns {
	columns.distinct = true
	return columns
}

func (columns *Columns) Of(datas ...string) *Columns {
	if len(datas) == 0 || (len(datas) == 1 && datas[0] == "*") {
		return &Columns{}
	}
	columns.fields = datas
	return columns
}

func (columns *Columns) Add(datas ...string) *Columns {
	columns.fields = append(columns.fields, datas...)
	return columns
}

func (columns *Columns) ToSql() string {
	if columns == nil {
		return ""
	}
	sqlStr := columns.hints.ToSql()
	if columns.distinct {
		sqlStr += "distinct "
	}

	if columns.tags {
		sqlStr += "tags "
	}

	if len(columns.fields) == 0 {
		return "*"
	}

	var fields []string
	for _, field := range columns.fields {
		fields = append(fields, ToDbField(field))
	}
	return sqlStr + strings.Join(fields, ", ")
}

func ToDbField(key string) string {
	if len(key) == 0 {
		return key
	}

	if key == constants.ALL_FIELD {
		return key
	}

	// 包含空格
	if strings.Contains(key, " ") {
		return key
	}
	// 包含函数
	if strings.Contains(key, "(") || strings.Contains(key, ")") {
		return key
	}
	// 伪列
	if isFakeColumn(key) {
		return key
	}
	// 包含属性列特殊字符
	if strings.Contains(key, dom) {
		return key
	}
	// 包含 别名
	if strings.Contains(key, "as") || strings.Contains(key, "AS") {
		return key
	}
	// 包含*
	if strings.Contains(key, "*") {
		return key
	}

	point := "."
	if strings.Contains(key, ".") {
		datas := strings.SplitN(key, point, 2)
		tableName := datas[0]
		columnName := datas[1]
		return tableName + "." + dom + columnName + dom
	} else {
		return dom + key + dom
	}
}

var defaultFakeColumns = []string{
	"tbname",
	"_qstart",
	"_qend",
	"_wstart",
	"_wend",
	"_wduration",
	"_c0",
	"_rowts",
	"_irowts",
}

// 是否是伪列
func isFakeColumn(columnName string) bool {
	columnNameTem := strings.ToLower(columnName)
	return util.ListContains(defaultFakeColumns, columnNameTem)
}

type Hint struct {
	BatchScan       bool // 说明：采用批量读表的方式；									适用：超级表 JOIN 语句
	NoBatchScan     bool // 说明：采用顺序读表的方式；									适用：超级表 JOIN 语句
	SortForGroup    bool // 说明：采用sort方式进行分组, 与PARTITION_FIRST冲突；			适用：partition by 列表有普通列时
	PartitionFirst  bool // 说明：在聚合之前使用PARTITION计算分组, 与SORT_FOR_GROUP冲突；	适用：partition by 列表有普通列时
	ParaTablesSort  bool // 说明：超级表的数据按时间戳排序时, 不使用临时磁盘空间, 只使用内存。当子表数量多, 行长比较大时候, 会使用大量内存, 可能发生OOM；	适用：超级表的数据按时间戳排序时
	SmalldataTsSort bool // 说明：超级表的数据按时间戳排序时, 查询列长度大于等于256, 但是行数不多, 使用这个提示, 可以提高性能；						适用：超级表的数据按时间戳排序时
	SkipTsma        bool // 说明：用于显示的禁用TSMA查询优化；																				适用：带Agg函数的查询语句
}

func (hint *Hint) ToSql() string {
	if hint == nil {
		return ""
	}
	var hintStrs []string
	if hint.BatchScan {
		hintStrs = append(hintStrs, "BATCH_SCAN()")
	}
	if hint.NoBatchScan {
		hintStrs = append(hintStrs, "NO_BATCH_SCAN()")
	}
	if hint.SortForGroup {
		hintStrs = append(hintStrs, "SORT_FOR_GROUP()")
	}
	if hint.PartitionFirst {
		hintStrs = append(hintStrs, "PARTITION_FIRST()")
	}
	if hint.ParaTablesSort {
		hintStrs = append(hintStrs, "PARA_TABLES_SORT()")
	}
	if hint.SmalldataTsSort {
		hintStrs = append(hintStrs, "SMALLDATA_TS_SORT()")
	}
	if hint.SkipTsma {
		hintStrs = append(hintStrs, "SKIP_TSMA()")
	}
	return "/*+ " + strings.Join(hintStrs, " ") + " */ "
}
