package test

import (
	"testing"

	"github.com/simonalong/tdorm/column"
	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	// 普通用法
	assert.Equal(t, "`id`, `name`, `age`", column.Of("id", "name", "age").ToSql())

	// 全部列：通配符
	assert.Equal(t, "*", column.Of("*").ToSql())
	assert.Equal(t, "tableName.*", column.Of("tableName.*").ToSql())

	// 别名
	assert.Equal(t, "id as dataId, `name`, `age`", column.Of("id as dataId", "name", "age").ToSql())

	// 标签列
	assert.Equal(t, "tags id, `name`, `age`", column.Of("tags id", "name", "age").ToSql())
	assert.Equal(t, "tags `id`, `name`, `age`", column.Tags().Of("id", "name", "age").ToSql())

	// 去重
	assert.Equal(t, "distinct `id`, `name`, `age`", column.Distinct().Of("id", "name", "age").ToSql())

	// Hints 提示
	assert.Equal(t, "/*+ BATCH_SCAN() */ `ts`", column.Hints(&column.Hint{BatchScan: true}).Of("ts").ToSql())
	assert.Equal(t, "/*+ NO_BATCH_SCAN() */ `ts`", column.Hints(&column.Hint{NoBatchScan: true}).Of("ts").ToSql())
	assert.Equal(t, "/*+ SORT_FOR_GROUP() */ `ts`", column.Hints(&column.Hint{SortForGroup: true}).Of("ts").ToSql())
	assert.Equal(t, "/*+ PARTITION_FIRST() */ `ts`", column.Hints(&column.Hint{PartitionFirst: true}).Of("ts").ToSql())
	assert.Equal(t, "/*+ PARA_TABLES_SORT() */ `ts`", column.Hints(&column.Hint{ParaTablesSort: true}).Of("ts").ToSql())
	assert.Equal(t, "/*+ SMALLDATA_TS_SORT() */ `ts`", column.Hints(&column.Hint{SmalldataTsSort: true}).Of("ts").ToSql())
	assert.Equal(t, "/*+ SKIP_TSMA() */ `ts`", column.Hints(&column.Hint{SkipTsma: true}).Of("ts").ToSql())

	// 多个提示
	assert.Equal(t, "/*+ BATCH_SCAN() SKIP_TSMA() */ `ts`", column.Hints(&column.Hint{BatchScan: true, SkipTsma: true}).Of("ts").ToSql())
}
