package test

import (
	"testing"

	"github.com/simonalong/gole/maps"

	"github.com/simonalong/tdorm/condition"
	"github.com/simonalong/tdorm/op"
	"github.com/stretchr/testify/assert"
)

func TestQueryAnd(t *testing.T) {
	assert.Equal(t, "where (`age` = ? and `name` = ? and `address` = ?)", condition.New().And("age", 12, "name", "zhou", "address", "杭州").ToSql())
	assert.Equal(t, "where `age` = ? and `name` = ? and `address` = ?", condition.New().AndEm("age", 12, "name", "zhou", "address", "杭州").ToSql())

	assert.Equal(t, "where (`age` = ?) and (`name` = ?) and (`address` = ?)", condition.New().And("age", 12).And("name", "zhou").And("address", "杭州").ToSql())
	assert.Equal(t, "where `age` = ? and `name` = ? and `address` = ?", condition.New().AndEm("age", 12).AndEm("name", "zhou").AndEm("address", "杭州").ToSql())

	assert.Equal(t, "where (`age` = ?)", condition.New().And("age", 12).ToSql())
	assert.Equal(t, "where `age` = ?", condition.New().AndEm("age", 12).ToSql())

	// value数据为nil，则sql不拼接
	assert.Equal(t, "where `age` = ? and `name` = ?", condition.New().AndEm("age", 12, "name", "zhou", "address", nil).ToSql())
	assert.Equal(t, "where `age` = ? and `name` = ? and `address` = ?", condition.New().AndEm("age", 12, "name", "zhou", "address", "").ToSql())
	assert.Equal(t, "", condition.New().AndEm("age", nil).ToSql())

	assert.Equal(t, "where (`age` = ?) and (`name` = ? and `address` = ?)", condition.New().And("age", 12).And("name", "zhou", "address", "杭州").ToSql())

	// OrmMap支持，OrmMap则默认全是 and
	assert.Equal(t, "where (`age` = ? and `name` = ? and `address` = ?)", condition.New().And(maps.OfSort("age", 12, "name", "zhou", "address", "杭州")).ToSql())

	// 全部空
	assert.Equal(t, "", condition.New().AndEm().ToSql())
	assert.Equal(t, "", condition.New().AndEm("name", nil).ToSql())
}

func TestQueryOr(t *testing.T) {
	assert.Equal(t, "where (`age` = ? or `name` = ? or `address` = ?)", condition.New().Or("age", 12, "name", "zhou", "address", "杭州").ToSql())
	assert.Equal(t, "where `age` = ? or `name` = ? or `address` = ?", condition.New().OrEm("age", 12, "name", "zhou", "address", "杭州").ToSql())

	assert.Equal(t, "where (`age` = ?) or (`name` = ?) or (`address` = ?)", condition.New().Or("age", 12).Or("name", "zhou").Or("address", "杭州").ToSql())
	assert.Equal(t, "where `age` = ? or `name` = ? or `address` = ?", condition.New().OrEm("age", 12).OrEm("name", "zhou").OrEm("address", "杭州").ToSql())

	assert.Equal(t, "where (`age` = ?)", condition.New().Or("age", 12).ToSql())
	assert.Equal(t, "where `age` = ?", condition.New().OrEm("age", 12).ToSql())

	// value数据为nil，则sql不拼接
	assert.Equal(t, "where `age` = ? or `name` = ?", condition.New().OrEm("age", 12, "name", "zhou", "address", nil).ToSql())
	assert.Equal(t, "where `age` = ? or `name` = ? or `address` = ?", condition.New().OrEm("age", 12, "name", "zhou", "address", "").ToSql())
	assert.Equal(t, "", condition.New().OrEm("age", nil).ToSql())

	assert.Equal(t, "where (`age` = ?) or (`name` = ? or `address` = ?)", condition.New().Or("age", 12).Or("name", "zhou", "address", "杭州").ToSql())
}

func TestQueryOpAndOr(t *testing.T) {
	// 使用op.And
	assert.Equal(t, "where ((`age` = ?) and `name` = ? and `address` = ?)", condition.New().And(op.And("age", 12), "name", "zhou", "address", "杭州").ToSql())
	assert.Equal(t, "where ((`age` = ? and `name` = ?) and `address` = ?)", condition.New().And(op.And("age", 12, "name", "zhou"), "address", "杭州").ToSql())
	assert.Equal(t, "where ((`age` = ? and `name` = ? and `address` = ?))", condition.New().And(op.And("age", 12, "name", "zhou", "address", "杭州")).ToSql())
	assert.Equal(t, "where (`age` = ? and `name` = ? and `address` = ?)", condition.New().And(op.AndEm("age", 12, "name", "zhou", "address", "杭州")).ToSql())

	// 使用op.Or
	assert.Equal(t, "where ((`age` = ?) or `name` = ? or `address` = ?)", condition.New().Or(op.Or("age", 12), "name", "zhou", "address", "杭州").ToSql())
	assert.Equal(t, "where ((`age` = ? or `name` = ?) or `address` = ?)", condition.New().Or(op.Or("age", 12, "name", "zhou"), "address", "杭州").ToSql())
	assert.Equal(t, "where ((`age` = ? or `name` = ? or `address` = ?))", condition.New().Or(op.Or("age", 12, "name", "zhou", "address", "杭州")).ToSql())
	assert.Equal(t, "where (`age` = ? or `name` = ? or `address` = ?)", condition.New().Or(op.OrEm("age", 12, "name", "zhou", "address", "杭州")).ToSql())

	// 使用And 和op.Or混用
	assert.Equal(t, "where ((`age` = ?) and `name` = ? and `address` = ?)", condition.New().And(op.Or("age", 12), "name", "zhou", "address", "杭州").ToSql())
	assert.Equal(t, "where ((`age` = ? or `name` = ?) and `address` = ?)", condition.New().And(op.Or("age", 12, "name", "zhou"), "address", "杭州").ToSql())
	assert.Equal(t, "where ((`age` = ? or `name` = ? or `address` = ?))", condition.New().And(op.Or("age", 12, "name", "zhou", "address", "杭州")).ToSql())

	// 使用And和Or
	assert.Equal(t, "where (`age` = ? and `name` = ?) or (`address` = ?)", condition.New().And("age", 12, "name", "zhou").Or("address", "杭州").ToSql())
	assert.Equal(t, "where `age` = ? and `name` = ? or (`address` = ?)", condition.New().AndEm("age", 12, "name", "zhou").Or("address", "杭州").ToSql())
	assert.Equal(t, "where (`age` = ? and `name` = ?) or `address` = ?", condition.New().And("age", 12, "name", "zhou").OrEm("address", "杭州").ToSql())
}

// 原生简化写法
// 在一些条件比较少的情况下的简化写法
func TestQueryOpOriginal(t *testing.T) {
	// 比较操作符：两种写法，最后生成的sql是一样的
	assert.Equal(t, "where `age` > ?", condition.New().Gt("age", 12).ToSql())
	assert.Equal(t, "where `age` > ?", condition.New().AndEm(op.Gt("age", 12)).ToSql())

	assert.Equal(t, "where `age` > ? and `name` = ?", condition.New().Gt("age", 12).AndEm("name", "zhou").ToSql())
	assert.Equal(t, "where `age` > ? and `name` = ?", condition.New().AndEm(op.Gt("age", 12), "name", "zhou").ToSql())

	// 其他比较操作符
	assert.Equal(t, "where `age` >= ?", condition.New().Ge("age", 12).ToSql())
	assert.Equal(t, "where `age` < ?", condition.New().Lt("age", 12).ToSql())
	assert.Equal(t, "where `age` <= ?", condition.New().Le("age", 12).ToSql())
	assert.Equal(t, "where `age` = ?", condition.New().Eq("age", 12).ToSql())
	assert.Equal(t, "where `age` != ?", condition.New().UnEq("age", 12).ToSql())

	// is null；is not null
	assert.Equal(t, "where `name` is null", condition.New().IsNull("name").ToSql())
	assert.Equal(t, "where `name` is not null", condition.New().IsNotNull("name").ToSql())

	// between and
	assert.Equal(t, "where `age` between ? and ?", condition.New().BetweenAnd("age", 12, 20).ToSql())
	assert.Equal(t, "where `age` not between ? and ?", condition.New().NotBetweenAnd("age", 12, 20).ToSql())

	// in；not in
	assert.Equal(t, "where `age` in (?, ?)", condition.New().In("age", []int{12, 18}).ToSql())
	assert.Equal(t, "where `age` in (?, ?)", condition.New().In("age", []int{12, 18}).ToSql())
	assert.Equal(t, "where `age` in (?, ?, ?, ?, ?, ?, ?, ?, ?)", condition.New().In("age", []int{12, 18, 43, 32, 43, 54, 65, 12, 64}).ToSql())
	assert.Equal(t, "where `age` not in (?, ?, ?)", condition.New().NotIn("age", []int{12, 18, 119}).ToSql())

	// like; not like
	assert.Equal(t, "where `name` like ?", condition.New().Like("name", "%牛_").ToSql())
	assert.Equal(t, "where `name` like ?", condition.New().Like("name", "%牛__").ToSql())
	assert.Equal(t, "where `name` like ?", condition.New().Like("name", "%牛%").ToSql())
	assert.Equal(t, "where `name` not like ?", condition.New().NotLike("name", "%牛%").ToSql())

	// match; nmatch
	assert.Equal(t, "where `name` match ?", condition.New().Match("name", "%牛_").ToSql())
	assert.Equal(t, "where `name` match ?", condition.New().Match("name", "%牛__").ToSql())
	assert.Equal(t, "where `name` match ?", condition.New().Match("name", "%牛%").ToSql())
	assert.Equal(t, "where `name` nmatch ?", condition.New().NotMatch("name", "%牛%").ToSql())
	assert.Equal(t, "where `name` nmatch ?", condition.New().NMatch("name", "%牛%").ToSql())

	// contains
	assert.Equal(t, "where `info` contains ?", condition.New().Contains("info", "k1").ToSql())
}

func TestQueryOp(t *testing.T) {
	// 比较操作符
	assert.Equal(t, "where (`age` > ? and `name` = ? and `address` = ?)", condition.New().And(op.Gt("age", 12), "name", "zhou", "address", "杭州").ToSql())
	assert.Equal(t, "where (`age` >= ? and `name` = ? and `address` = ?)", condition.New().And(op.Ge("age", 12), "name", "zhou", "address", "杭州").ToSql())
	assert.Equal(t, "where (`age` < ? and `name` = ? and `address` = ?)", condition.New().And(op.Lt("age", 12), "name", "zhou", "address", "杭州").ToSql())
	assert.Equal(t, "where (`age` <= ? and `name` = ? and `address` = ?)", condition.New().And(op.Le("age", 12), "name", "zhou", "address", "杭州").ToSql())
	//assert.Equal(t, "where (`age` <> ? and `name` = ? and `address` = ?)", query.New().And(op.UnEq("age", 12), "name", "zhou", "address", "杭州").ToSql())
	assert.Equal(t, "where (`age` != ? and `name` = ? and `address` = ?)", condition.New().And(op.UnEq("age", 12), "name", "zhou", "address", "杭州").ToSql())

	// is null；is not null
	assert.Equal(t, "where `name` is null and `address` = ?", condition.New().AndEm(op.IsNull("name"), "address", "杭州").ToSql())
	assert.Equal(t, "where `name` is not null and `address` = ?", condition.New().AndEm(op.IsNotNull("name"), "address", "杭州").ToSql())

	// between and
	assert.Equal(t, "where `age` between ? and ? and `address` = ?", condition.New().AndEm(op.BetweenAnd("age", 12, 20), "address", "杭州").ToSql())
	assert.Equal(t, "where `age` not between ? and ? and `address` = ?", condition.New().AndEm(op.NotBetweenAnd("age", 12, 20), "address", "杭州").ToSql())

	// in；not in
	assert.Equal(t, "where `age` in (?, ?) and `address` = ?", condition.New().AndEm(op.In("age", []int{12, 18}), "address", "杭州").ToSql())
	assert.Equal(t, "where `age` in (?, ?) and `address` = ?", condition.New().AndEm(op.In("age", []int{12, 18}), "address", "杭州").ToSql())
	assert.Equal(t, "where `age` in (?, ?, ?, ?, ?, ?, ?, ?, ?) and `address` = ?", condition.New().AndEm(op.In("age", []int{12, 18, 43, 32, 43, 54, 65, 12, 64}), "address", "杭州").ToSql())
	assert.Equal(t, "where `age` not in (?, ?, ?) and `address` = ?", condition.New().AndEm(op.NotIn("age", []int{12, 18, 119}), "address", "杭州").ToSql())

	// like; not like
	assert.Equal(t, "where `name` like ?", condition.New().Em(op.Like("name", "%牛_")).ToSql())
	assert.Equal(t, "where `name` like ?", condition.New().Em(op.Like("name", "%牛__")).ToSql())
	assert.Equal(t, "where `name` like ?", condition.New().Em(op.Like("name", "%牛%")).ToSql())
	assert.Equal(t, "where `name` not like ?", condition.New().Em(op.NotLike("name", "%牛%")).ToSql())

	// match; nmatch
	assert.Equal(t, "where `name` match ?", condition.New().Em(op.Match("name", "%牛_")).ToSql())
	assert.Equal(t, "where `name` match ?", condition.New().Em(op.Match("name", "%牛__")).ToSql())
	assert.Equal(t, "where `name` match ?", condition.New().Em(op.Match("name", "%牛%")).ToSql())
	assert.Equal(t, "where `name` nmatch ?", condition.New().Em(op.NotMatch("name", "%牛%")).ToSql())
	assert.Equal(t, "where `name` nmatch ?", condition.New().Em(op.NMatch("name", "%牛%")).ToSql())

	// contains
	assert.Equal(t, "where `info` contains ?", condition.New().Em(op.Contains("info", "k1")).ToSql())
}

func TestAppend(t *testing.T) {
	assert.Equal(t, "where (`age` > ? and `name` = ? and `address` = ?) order by `ts` desc", condition.New().And(op.Gt("age", 12), "name", "zhou", "address", "杭州").Append("order by `ts` desc").ToSql())
}

// 使用表名
func TestTableName(t *testing.T) {
	assert.Equal(t, "where (td_table1.`age` > ? and td_table1.`name` = ? and td_table2.`address` = ?)", condition.New().And(op.Gt("td_table1.age", 12), "td_table1.name", "zhou", "td_table2.address", "杭州").ToSql())
}

// tdengine的高级功能
func TestAdvance(t *testing.T) {
	assert.Equal(t, "where `status` = ?", condition.New().Eq("status", 2).ToSql())
	assert.Equal(t, "where `status` = ? limit 10", condition.New().Eq("status", 2).Append("limit 10").ToSql())
	assert.Equal(t, "where `status` = ? partition by location interval(10m)", condition.New().Eq("status", 2).Append("partition by location interval(10m)").ToSql())
	assert.Equal(t, "where `status` = ? partition by tbname state_window(case when voltage >= 205 and voltage <= 235 then 1 else 0 end)", condition.New().Eq("status", 2).Append("partition by tbname state_window(case when voltage >= 205 and voltage <= 235 then 1 else 0 end)").ToSql())
	assert.Equal(t, "session(ts, tol_val)", condition.New().Append("session(ts, tol_val)").ToSql())
}
