package tdorm

import (
	"database/sql/driver"

	"github.com/simonalong/gole/maps"
	"github.com/simonalong/tdorm/column"
	"github.com/simonalong/tdorm/condition"
)

type Cmd interface {
	/* --------------  基础api -------------- */

	Exec(sql string, args ...driver.Value) (driver.Result, error)
	Query(sql string, args ...driver.Value) ([]*maps.GoleMap, error)

	/* --------------  高级简化api -------------- */

	// Select 高级查询：支持所有的查询，包括不限于各种特色查询
	Select(hintsDistinctTagsAndColumns *column.Columns, fromClause string, whereConditionAndClause *condition.Condition) ([]*maps.GoleMap, error)

	/* --------------  基础封装 -------------- */

	// Insert 新增
	Insert(tableName string, toInsertMap *maps.GoleMap) (int, error)
	InsertWithTag(tableName, stableName string, tagsMap *maps.GoleMap, toInsertMap *maps.GoleMap) (int, error) // 子表不存在则会自动创建（不支持：original）
	InsertEntity(tableName string, entity interface{}) (int, error)
	InsertEntityWithTag(tableName, stableName string, tagEntity interface{}, entity interface{}) (int, error) // 子表不存在则会自动创建（不支持：original）
	InsertBatch(tableName string, toInsertMaps []*maps.GoleMap) (int64, error)
	InsertBatchWithTag(tableName, stableName string, tagsMap *maps.GoleMap, toInsertMaps []*maps.GoleMap) (int64, error) // 子表不存在则会自动创建（不支持：original）
	InsertEntityBatch(tableName string, entities []interface{}) (int64, error)
	InsertEntityBatchWithTag(tableName, stableName string, tagEntity interface{}, entities []interface{}) (int64, error) // 子表不存在则会自动创建（不支持：original）
	// Delete 删除
	Delete(tableName string, query *condition.Condition) (int64, error)
	// One 查询：多列一行
	One(tableName string, columns *column.Columns, query *condition.Condition) (*maps.GoleMap, error)
	// List 查询：多列多行
	List(tableName string, columns *column.Columns, query *condition.Condition) ([]*maps.GoleMap, error)
	ListOfDistinct(tableName string, columns *column.Columns, query *condition.Condition) ([]*maps.GoleMap, error)
	// Value 查询：一列一值
	Value(tableName, columnName string, query *condition.Condition) (interface{}, error)
	// Values 查询：一列多值
	Values(tableName, columnName string, query *condition.Condition) ([]interface{}, error)
	// ValuesOfDistinct 查询：一列多值，去重
	ValuesOfDistinct(tableName, columnName string, query *condition.Condition) ([]interface{}, error)
	// Count 查询：个数
	Count(tableName string, query *condition.Condition) (int, error)

	/* --------------  其他功能 -------------- */

	AddHook(hook TdHook)
	GetHooks() *TdHooks
}

type InfoCmd interface {
	SelectDatabase() (string, error)
	SelectClientVersion() (string, error)
	SelectServerVersion() (string, error)
	SelectServerStatus() (string, error) // 服务器状态检测语句。如果服务器正常，返回一个数字（例如 1）。如果服务器异常，返回 error code。该 SQL 语法能兼容连接池对于 TDengine 状态的检查及第三方工具对于数据库服务器状态的检查。并可以避免出现使用了错误的心跳检测 SQL 语句导致的连接池连接丢失的问题
	SelectNow() (string, error)
	SelectToday() (string, error)
	SelectTimeZone() (string, error)
	SelectCurrentUser() (string, error)
	SelectUser() (string, error)
}

type ShowCmd interface {
	ShowApps() ([]*maps.GoleMap, error)                           // 显示接入集群的应用（客户端）信息
	ShowCluster() (*maps.GoleMap, error)                          // 显示当前集群的信息
	ShowClusterAlive() (int, error)                               // 查询当前集群的状态是否可用，返回值： 0：不可用 1：完全可用 2：部分可用（集群中部分节点下线，但其它节点仍可以正常使用）
	ShowConnections() ([]*maps.GoleMap, error)                    // 显示当前系统中存在的连接的信息
	ShowConsumers() ([]*maps.GoleMap, error)                      // 显示当前数据库下所有消费者的信息
	ShowCreateDatabase() (string, error)                          // 显示 db_name 指定的数据库的创建语句
	ShowCreateStable(stableName string) (string, error)           // 显示 tb_name 指定的超级表的创建语句
	ShowCreateTable(tableName string) (string, error)             // 显示 tb_name 指定的表的创建语句。支持普通表、超级表和子表
	ShowDatabases(dbType string) ([]string, error)                // dbType：system：指定只显示系统数据库; user：指定只显示用户创建的数据库
	ShowDNodes() ([]*maps.GoleMap, error)                         // 显示当前系统中 DNODE 的信息
	ShowFunctions() ([]*maps.GoleMap, error)                      // 显示用户定义的自定义函数。
	ShowLicences() ([]*maps.GoleMap, error)                       // 显示企业版许可授权的信息；注：企业版独有
	ShowGrants() ([]*maps.GoleMap, error)                         // 显示企业版许可授权的信息；注：企业版独有
	ShowIndexes(tableName string) ([]*maps.GoleMap, error)        // 显示已创建的索引
	ShowLocalVariables() ([]*maps.GoleMap, error)                 // 显示当前客户端配置参数的运行值。
	ShowMNodes() ([]*maps.GoleMap, error)                         // 显示当前系统中 MNODE 的信息
	ShowQNodes() ([]*maps.GoleMap, error)                         // 显示当前系统中 QNODE （查询节点）的信息
	ShowQueries() ([]*maps.GoleMap, error)                        // 显示当前系统中正在进行的查询
	ShowScores() ([]*maps.GoleMap, error)                         // 显示系统被许可授权的容量的信息 注：企业版独有
	ShowStables(stableNamePart string) ([]string, error)          // 显示当前数据库下的所有超级表的信息。会自动使用 LIKE 对表名进行模糊匹配
	ShowStreams() ([]*maps.GoleMap, error)                        // 显示当前系统内所有流计算的信息
	ShowSubscriptions() ([]*maps.GoleMap, error)                  // 显示当前系统内所有的订阅关系
	ShowTables(tableRole, tableNamePart string) ([]string, error) // 显示当前数据库下的所有普通表和子表的信息。会自动使用 LIKE 对表名进行模糊匹配。tableRole：NORMAL 指定只显示普通表信息， CHILD 指定只显示子表信息
	ShowTags(tableName string) ([]*maps.GoleMap, error)           // 显示子表的标签信息
	ShowTopics() ([]string, error)                                // 显示当前数据库下的所有主题的信息
	ShowTransactions() ([]*maps.GoleMap, error)                   // 显示当前系统中正在执行的事务的信息(该事务仅针对除普通表以外的元数据级别)
	ShowUsers() ([]*maps.GoleMap, error)                          // 显示当前系统中所有用户的信息。包括用户自定义的用户和系统默认用户
	ShowClusterVariables() ([]*maps.GoleMap, error)               // 显示当前系统中各节点需要相同的配置参数的运行值，也可以指定 DNODE 来查看其的配置参数
	ShowVGroups() ([]*maps.GoleMap, error)                        // 显示当前数据库中所有 VGROUP 的信息
	ShowVNodes(dNodeId string) ([]*maps.GoleMap, error)           // 显示当前系统中所有 VNODE 或某个 DNODE 的 VNODE 的信息
}
