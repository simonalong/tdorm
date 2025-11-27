
# 查询条件
查询条件主要是query支持功能，也支持GoleMap作为条件，这里主要讲query.Query的功能；支持如下函数
```go
And(datas ...any) *Query
AndEm(datas ...any) *Query
Or(datas ...any) *Query
OrEm(datas ...any) *Query
Em(datas ...any) *Query
Gt(key string, value interface{}) *Query {
Ge(key string, value interface{}) *Query {
Lt(key string, value interface{}) *Query {
Le(key string, value interface{}) *Query {
Eq(key string, value interface{}) *Query {
UnEq(key string, value interface{}) *Query {
IsNull(key string) *Query {
IsNotNull(key string) *Query {
BetweenAnd(key string, leftValue, rightValue interface{}) *Query {
NotBetweenAnd(key string, leftValue, rightValue interface{}) *Query {
In(key string, values interface{}) *Query {
NotIn(key string, values interface{}) *Query {
Like(key string, value interface{}) *Query {
NotLike(key string, value interface{}) *Query {
Match(key string, value interface{}) *Query {
NMatch(key string, value interface{}) *Query {
NotMatch(key string, value interface{}) *Query {
Contains(key string, value interface{}) *Query {
```
query是查询条件
```go
assert.Equal(t, "where (`age` = ? and `name` = ? and `address` = ?)", query.New().And("age", 12, "name", "zhou", "address", "杭州").ToSql())
assert.Equal(t, "where `age` = ? and `name` = ? and `address` = ?", query.New().AndEm("age", 12, "name", "zhou", "address", "杭州").ToSql())

assert.Equal(t, "where (`age` = ?) and (`name` = ?) and (`address` = ?)", query.New().And("age", 12).And("name", "zhou").And("address", "杭州").ToSql())
assert.Equal(t, "where `age` = ? and `name` = ? and `address` = ?", query.New().AndEm("age", 12).AndEm("name", "zhou").AndEm("address", "杭州").ToSql())

assert.Equal(t, "where (`age` = ?)", query.New().And("age", 12).ToSql())
assert.Equal(t, "where `age` = ?", query.New().AndEm("age", 12).ToSql())

// value数据为nil，则sql不拼接
assert.Equal(t, "where `age` = ? and `name` = ?", query.New().AndEm("age", 12, "name", "zhou", "address", nil).ToSql())
assert.Equal(t, "where `age` = ? and `name` = ? and `address` = ?", query.New().AndEm("age", 12, "name", "zhou", "address", "").ToSql())
assert.Equal(t, "", query.New().AndEm("age", nil).ToSql())
```
一个条件简单查询，可以使用简化如下
```go
// 比较操作符：两种写法，最后生成的sql是一样的
assert.Equal(t, "where `age` > ?", query.New().Gt("age", 12).ToSql())
assert.Equal(t, "where `age` > ?", query.New().AndEm(op.Gt("age", 12)).ToSql())

assert.Equal(t, "where `age` > ? and `name` = ?", query.New().Gt("age", 12).AndEm("name", "zhou").ToSql())
assert.Equal(t, "where `age` > ? and `name` = ?", query.New().AndEm(op.Gt("age", 12), "name", "zhou").ToSql())

// 其他比较操作符
assert.Equal(t, "where `age` >= ?", query.New().Ge("age", 12).ToSql())
assert.Equal(t, "where `age` < ?", query.New().Lt("age", 12).ToSql())
assert.Equal(t, "where `age` <= ?", query.New().Le("age", 12).ToSql())
assert.Equal(t, "where `age` = ?", query.New().Eq("age", 12).ToSql())
assert.Equal(t, "where `age` != ?", query.New().UnEq("age", 12).ToSql())

// is null；is not null
assert.Equal(t, "where `name` is null", query.New().IsNull("name").ToSql())
assert.Equal(t, "where `name` is not null", query.New().IsNotNull("name").ToSql())

// between and
assert.Equal(t, "where `age` between ? and ?", query.New().BetweenAnd("age", 12, 20).ToSql())
assert.Equal(t, "where `age` not between ? and ?", query.New().NotBetweenAnd("age", 12, 20).ToSql())

// in；not in
assert.Equal(t, "where `age` in (?, ?)", query.New().In("age", []int{12, 18}).ToSql())
assert.Equal(t, "where `age` in (?, ?)", query.New().In("age", []int{12, 18}).ToSql())
assert.Equal(t, "where `age` in (?, ?, ?, ?, ?, ?, ?, ?, ?)", query.New().In("age", []int{12, 18, 43, 32, 43, 54, 65, 12, 64}).ToSql())
assert.Equal(t, "where `age` not in (?, ?, ?)", query.New().NotIn("age", []int{12, 18, 119}).ToSql())

// like; not like
assert.Equal(t, "where `name` like ?", query.New().Like("name", "%牛_").ToSql())
assert.Equal(t, "where `name` like ?", query.New().Like("name", "%牛__").ToSql())
assert.Equal(t, "where `name` like ?", query.New().Like("name", "%牛%").ToSql())
assert.Equal(t, "where `name` not like ?", query.New().NotLike("name", "%牛%").ToSql())

// match; nmatch
assert.Equal(t, "where `name` match ?", query.New().Match("name", "%牛_").ToSql())
assert.Equal(t, "where `name` match ?", query.New().Match("name", "%牛__").ToSql())
assert.Equal(t, "where `name` match ?", query.New().Match("name", "%牛%").ToSql())
assert.Equal(t, "where `name` nmatch ?", query.New().NotMatch("name", "%牛%").ToSql())
assert.Equal(t, "where `name` nmatch ?", query.New().NMatch("name", "%牛%").ToSql())

// contains
assert.Equal(t, "where `info` contains ?", query.New().Contains("info", "k1").ToSql())
```
也支持嵌套`op.xxx`；op操作可以作为Query.And()或者Or()的参数，如下示例
```go
// 使用And 和op.Or混用
assert.Equal(t, "where ((`age` = ?) and `name` = ? and `address` = ?)", query.New().And(op.Or("age", 12), "name", "zhou", "address", "杭州").ToSql())
assert.Equal(t, "where ((`age` = ? or `name` = ?) and `address` = ?)", query.New().And(op.Or("age", 12, "name", "zhou"), "address", "杭州").ToSql())
assert.Equal(t, "where ((`age` = ? or `name` = ? or `address` = ?))", query.New().And(op.Or("age", 12, "name", "zhou", "address", "杭州")).ToSql())
```
op的函数有如下，与Query.And或者Or函数可以组成更加复杂的条件
```go
// 比较操作符
assert.Equal(t, "where (`age` > ? and `name` = ? and `address` = ?)", query.New().And(op.Gt("age", 12), "name", "zhou", "address", "杭州").ToSql())
assert.Equal(t, "where (`age` >= ? and `name` = ? and `address` = ?)", query.New().And(op.Ge("age", 12), "name", "zhou", "address", "杭州").ToSql())
assert.Equal(t, "where (`age` < ? and `name` = ? and `address` = ?)", query.New().And(op.Lt("age", 12), "name", "zhou", "address", "杭州").ToSql())
assert.Equal(t, "where (`age` <= ? and `name` = ? and `address` = ?)", query.New().And(op.Le("age", 12), "name", "zhou", "address", "杭州").ToSql())
//assert.Equal(t, "where (`age` <> ? and `name` = ? and `address` = ?)", query.New().And(op.UnEq("age", 12), "name", "zhou", "address", "杭州").ToSql())
assert.Equal(t, "where (`age` != ? and `name` = ? and `address` = ?)", query.New().And(op.UnEq("age", 12), "name", "zhou", "address", "杭州").ToSql())

// is null；is not null
assert.Equal(t, "where `name` is null and `address` = ?", query.New().AndEm(op.IsNull("name"), "address", "杭州").ToSql())
assert.Equal(t, "where `name` is not null and `address` = ?", query.New().AndEm(op.IsNotNull("name"), "address", "杭州").ToSql())

// between and
assert.Equal(t, "where `age` between ? and ? and `address` = ?", query.New().AndEm(op.BetweenAnd("age", 12, 20), "address", "杭州").ToSql())
assert.Equal(t, "where `age` not between ? and ? and `address` = ?", query.New().AndEm(op.NotBetweenAnd("age", 12, 20), "address", "杭州").ToSql())

// in；not in
assert.Equal(t, "where `age` in (?, ?) and `address` = ?", query.New().AndEm(op.In("age", []int{12, 18}), "address", "杭州").ToSql())
assert.Equal(t, "where `age` in (?, ?) and `address` = ?", query.New().AndEm(op.In("age", []int{12, 18}), "address", "杭州").ToSql())
assert.Equal(t, "where `age` in (?, ?, ?, ?, ?, ?, ?, ?, ?) and `address` = ?", query.New().AndEm(op.In("age", []int{12, 18, 43, 32, 43, 54, 65, 12, 64}), "address", "杭州").ToSql())
assert.Equal(t, "where `age` not in (?, ?, ?) and `address` = ?", query.New().AndEm(op.NotIn("age", []int{12, 18, 119}), "address", "杭州").ToSql())

// like; not like
assert.Equal(t, "where `name` like ?", query.New().Em(op.Like("name", "%牛_")).ToSql())
assert.Equal(t, "where `name` like ?", query.New().Em(op.Like("name", "%牛__")).ToSql())
assert.Equal(t, "where `name` like ?", query.New().Em(op.Like("name", "%牛%")).ToSql())
assert.Equal(t, "where `name` not like ?", query.New().Em(op.NotLike("name", "%牛%")).ToSql())

// match; nmatch
assert.Equal(t, "where `name` match ?", query.New().Em(op.Match("name", "%牛_")).ToSql())
assert.Equal(t, "where `name` match ?", query.New().Em(op.Match("name", "%牛__")).ToSql())
assert.Equal(t, "where `name` match ?", query.New().Em(op.Match("name", "%牛%")).ToSql())
assert.Equal(t, "where `name` nmatch ?", query.New().Em(op.NotMatch("name", "%牛%")).ToSql())
assert.Equal(t, "where `name` nmatch ?", query.New().Em(op.NMatch("name", "%牛%")).ToSql())

// contains
assert.Equal(t, "where `info` contains ?", query.New().Em(op.Contains("info", "k1")).ToSql())
```
