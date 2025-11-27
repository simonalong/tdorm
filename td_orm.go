package tdorm

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/simonalong/gole/logger"
	"github.com/simonalong/gole/maps"
	"github.com/simonalong/tdorm/column"
	"github.com/simonalong/tdorm/condition"
	"github.com/simonalong/tdorm/constants"
	_ "github.com/taosdata/driver-go/v3/taosWS"
)

type TdClient struct {
	// 连接类型：connect_original，connect_rest，connect_websocket
	connectType byte
	// 命令执行器：根据connectType来设置不同的
	cmdExecutor Cmd
	// 数据库名：我们这里一个库对应一个对象
	DbName string
}

// NewConnectOriginal 创建原生连接
func NewConnectOriginal(host string, port int, user, password, dbName string) *TdClient {
	//_, err := af.Open(host, user, password, dbName, port)
	//if err != nil {
	//	logger.Errorf("tdengine连接异常，请检查配置 host=%v, port=%v, user=%v, password=%v, DbName=%v，异常：%v", host, port, user, password, dbName, err.Error())
	//	return nil
	//}
	return nil

	//return &TdClient{
	//	cmdExecutor: NewOriginalCmd(conn, dbName),
	//	DbName:      dbName,
	//}
}

// NewConnectRest 创建rest连接
func NewConnectRest(host string, port int, user, password, dbName string) *TdClient {
	//"root:taosdata@http(localhost:6041)/log"
	taosDSN := fmt.Sprintf("%s:%s@http(%s:%d)/%s", user, password, host, port, dbName)
	taos, err := sql.Open("taosRestful", taosDSN)
	if err != nil {
		logger.Errorf("tdengine连接异常，请检查配置 taosDSN=%v，异常：%v", taosDSN, err.Error())
		return nil
	}
	return &TdClient{
		cmdExecutor: &RestfulCmd{tdSqlDb: taos, dbName: dbName, hooks: NewTdHooks()},
		connectType: constants.ConnectRestful,
		DbName:      dbName,
	}
}

// NewConnectWebsocket 创建websocket连接
func NewConnectWebsocket(host string, port int, user, password, dbName string) *TdClient {
	//root:taosdata@ws(localhost:6041)/xx
	taosDSN := fmt.Sprintf("%s:%s@ws(%s:%d)/%s", user, password, host, port, dbName)
	taos, err := sql.Open("taosWS", taosDSN)
	if err != nil {
		logger.Errorf("tdengine连接异常，请检查配置 taosDSN=%v，异常：%v", taosDSN, err.Error())
		return nil
	}

	return &TdClient{
		cmdExecutor: &WebsocketCmd{tdSqlDb: taos, dbName: dbName, hooks: NewTdHooks()},
		connectType: constants.ConnectWebsocket,
		DbName:      dbName,
	}
}

func (tdClient *TdClient) getCmdExecutor() Cmd {
	return tdClient.cmdExecutor
}

func (tdClient *TdClient) getHooks() *TdHooks {
	return tdClient.getCmdExecutor().GetHooks()
}

func (tdClient *TdClient) AddHook(hook TdHook) {
	tdClient.getCmdExecutor().AddHook(hook)
}

func (tdClient *TdClient) GetHooks() *TdHooks {
	return tdClient.getCmdExecutor().GetHooks()
}

func (tdClient *TdClient) Exec(sql string, args ...driver.Value) (driver.Result, error) {
	return tdClient.getCmdExecutor().Exec(sql, args...)
}

func (tdClient *TdClient) Query(sql string, args ...driver.Value) ([]*maps.GoleMap, error) {
	return tdClient.getCmdExecutor().Query(sql, args...)
}

// Select
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

func (tdClient *TdClient) Select(selectColumns *column.Columns, fromClause string, whereCondition *condition.Condition) ([]*maps.GoleMap, error) {
	return tdClient.Query(generateSelectSql(selectColumns, fromClause, whereCondition))
}

// Insert 注意：使用该函数要保证表一定存在；否则请使用InsertWithTag函数
func (tdClient *TdClient) Insert(tableName string, toInsertMap *maps.GoleMap) (int, error) {
	if toInsertMap.IsEmpty() {
		logger.Warnf("insert 数据为空")
		return 0, nil
	}
	return tdClient.getCmdExecutor().Insert(tableName, toInsertMap)
}

// InsertWithTag 子表不存在则会自动创建
func (tdClient *TdClient) InsertWithTag(tableName, stableName string, tagsMap *maps.GoleMap, toInsertMap *maps.GoleMap) (int, error) {
	if toInsertMap.IsEmpty() {
		logger.Warnf("insert 数据为空")
		return 0, nil
	}

	return tdClient.getCmdExecutor().InsertWithTag(tableName, stableName, tagsMap, toInsertMap)
}

func (tdClient *TdClient) InsertEntity(tableName string, entity interface{}) (int, error) {
	if entity == nil {
		logger.Warnf("insert entity 数据为空")
		return 0, nil
	}
	return tdClient.getCmdExecutor().InsertEntity(tableName, entity)
}

func (tdClient *TdClient) InsertEntityWithTag(tableName, stableName string, tagEntity interface{}, entity interface{}) (int, error) {
	if entity == nil {
		logger.Warnf("insert entity 数据为空")
		return 0, nil
	}

	return tdClient.InsertWithTag(tableName, stableName, maps.FromEntity(tagEntity), maps.FromEntity(entity))
}

func (tdClient *TdClient) InsertBatch(tableName string, toInsertMaps []*maps.GoleMap) (int64, error) {
	if toInsertMaps == nil || len(toInsertMaps) == 0 || maps.AllIsEmpty(toInsertMaps) {
		logger.Warnf("InsertBatch 数据为空")
		return 0, nil
	}

	return tdClient.getCmdExecutor().InsertBatch(tableName, toInsertMaps)
}

func (tdClient *TdClient) InsertBatchWithTag(tableName, stableName string, tagsMap *maps.GoleMap, toInsertMaps []*maps.GoleMap) (int64, error) {
	if toInsertMaps == nil || len(toInsertMaps) == 0 || maps.AllIsEmpty(toInsertMaps) {
		logger.Warnf("InsertBatch 数据为空")
		return 0, nil
	}

	return tdClient.getCmdExecutor().InsertBatchWithTag(tableName, stableName, tagsMap, toInsertMaps)
}

func (tdClient *TdClient) InsertEntityBatch(tableName string, entities []interface{}) (int64, error) {
	if entities == nil || len(entities) == 0 {
		logger.Warnf("InsertEntityBatch 数据为空")
		return 0, nil
	}
	var insertMaps []*maps.GoleMap
	for _, entity := range entities {
		insertMaps = append(insertMaps, maps.FromEntity(entity))
	}
	return tdClient.InsertBatch(tableName, insertMaps)
}

func (tdClient *TdClient) InsertEntityBatchWithTag(tableName, stableName string, tagEntity interface{}, entities []interface{}) (int64, error) {
	if entities == nil {
		logger.Warnf("InsertBatch 数据为空")
		return 0, nil
	}

	return tdClient.getCmdExecutor().InsertEntityBatchWithTag(tableName, stableName, tagEntity, entities)
}

func (tdClient *TdClient) Delete(tableName string, query *condition.Condition) (int64, error) {
	if query == nil || tableName == "" {
		logger.Warnf("delete 参数为空")
		return 0, nil
	}

	return tdClient.getCmdExecutor().Delete(tableName, query)
}

func (tdClient *TdClient) One(tableName string, columns *column.Columns, query *condition.Condition) (*maps.GoleMap, error) {
	if query == nil || tableName == "" {
		logger.Warnf("delete 参数为空")
		return maps.New(), nil
	}

	return tdClient.getCmdExecutor().One(tableName, columns, query)
}

func (tdClient *TdClient) List(tableName string, columns *column.Columns, query *condition.Condition) ([]*maps.GoleMap, error) {
	if query == nil || tableName == "" {
		logger.Warnf("delete 参数为空")
		return []*maps.GoleMap{}, nil
	}
	return tdClient.getCmdExecutor().List(tableName, columns, query)
}

func (tdClient *TdClient) ListOfDistinct(tableName string, columns *column.Columns, query *condition.Condition) ([]*maps.GoleMap, error) {
	if query == nil || tableName == "" {
		logger.Warnf("delete 参数为空")
		return []*maps.GoleMap{}, nil
	}
	return tdClient.getCmdExecutor().ListOfDistinct(tableName, columns, query)
}

func (tdClient *TdClient) Value(tableName, columnName string, query *condition.Condition) (interface{}, error) {
	if query == nil || tableName == "" {
		logger.Warnf("delete 参数为空")
		return maps.New(), nil
	}

	return tdClient.getCmdExecutor().Value(tableName, columnName, query)
}

func (tdClient *TdClient) Values(tableName, columnName string, query *condition.Condition) ([]interface{}, error) {
	if query == nil || tableName == "" {
		logger.Warnf("delete 参数为空")
		return nil, nil
	}
	return tdClient.getCmdExecutor().Values(tableName, columnName, query)
}

func (tdClient *TdClient) ValuesOfDistinct(tableName, columnName string, query *condition.Condition) ([]interface{}, error) {
	if query == nil || tableName == "" {
		logger.Warnf("delete 参数为空")
		return nil, nil
	}
	return tdClient.getCmdExecutor().ValuesOfDistinct(tableName, columnName, query)
}

func (tdClient *TdClient) Count(tableName string, query *condition.Condition) (int, error) {
	if query == nil || tableName == "" {
		logger.Warnf("delete 参数为空")
		return 0, nil
	}
	return tdClient.getCmdExecutor().Count(tableName, query)
}

//func (tdClient *TdClient) GetAfConnect() (*af.Connector, error) {
//	//if tdClient.connectType == constants.ConnectOriginal {
//	//	return tdClient.getCmdExecutor().(*OriginalCmd).tdConnect, nil
//	//} else {
//	//	return nil, errors.New(fmt.Sprintf("类型[%v]不支持*af.Connect", tdClient.connectType))
//	//}
//	return nil, nil
//}

func (tdClient *TdClient) GetSqlDb() (*sql.DB, error) {
	if tdClient.connectType == constants.ConnectRestful {
		return tdClient.getCmdExecutor().(*RestfulCmd).tdSqlDb, nil
	} else if tdClient.connectType == constants.ConnectWebsocket {
		return tdClient.getCmdExecutor().(*WebsocketCmd).tdSqlDb, nil
	} else {
		return nil, errors.New(fmt.Sprintf("类型[%v]不支持*Sql.Db", tdClient.connectType))
	}
}

func (tdClient *TdClient) SelectDatabase() (string, error) {
	return doQueryOneString(tdClient, "database()")
}

func (tdClient *TdClient) SelectClientVersion() (string, error) {
	return doQueryOneString(tdClient, "client_version()")
}

func (tdClient *TdClient) SelectServerVersion() (string, error) {
	return doQueryOneString(tdClient, "server_version()")
}

// SelectServerStatus 服务器状态检测语句。如果服务器正常，返回一个数字（例如 1）。如果服务器异常，返回 error code。该 SQL 语法能兼容连接池对于 TDengine 状态的检查及第三方工具对于数据库服务器状态的检查。并可以避免出现使用了错误的心跳检测 SQL 语句导致的连接池连接丢失的问题
func (tdClient *TdClient) SelectServerStatus() (string, error) {
	return doQueryOneString(tdClient, "server_status()")
}

func (tdClient *TdClient) SelectNow() (string, error) {
	return doQueryOneString(tdClient, "now()")
}

func (tdClient *TdClient) SelectToday() (string, error) {
	return doQueryOneString(tdClient, "today()")
}

func (tdClient *TdClient) SelectTimeZone() (string, error) {
	return doQueryOneString(tdClient, "timezone()")
}

func (tdClient *TdClient) SelectCurrentUser() (string, error) {
	return doQueryOneString(tdClient, "current_user()")
}

func (tdClient *TdClient) SelectUser() (string, error) {
	return doQueryOneString(tdClient, "user()")
}

// ShowApps 显示接入集群的应用（客户端）信息
func (tdClient *TdClient) ShowApps() ([]*maps.GoleMap, error) {
	return tdClient.Query("show apps")
}

// ShowCluster 显示当前集群的信息
func (tdClient *TdClient) ShowCluster() (*maps.GoleMap, error) {
	dataMaps, err := tdClient.Query("show cluster")
	if err != nil {
		return nil, err
	}
	if len(dataMaps) == 0 {
		return nil, nil
	}
	return dataMaps[0], nil
}

// ShowClusterAlive 查询当前集群的状态是否可用，返回值： 0：不可用 1：完全可用 2：部分可用（集群中部分节点下线，但其它节点仍可以正常使用）
func (tdClient *TdClient) ShowClusterAlive() (int, error) {
	dataMaps, err := tdClient.Query("show cluster alive")
	if err != nil {
		return 0, err
	}
	if len(dataMaps) == 0 {
		return 0, nil
	}

	if result, exist := dataMaps[0].GetInt("status"); exist {
		return result, nil
	} else {
		return 0, nil
	}
}

// ShowConnections 显示当前系统中存在的连接的信息
func (tdClient *TdClient) ShowConnections() ([]*maps.GoleMap, error) {
	return tdClient.Query("show connections")
}

// ShowConsumers 显示当前数据库下所有消费者的信息
func (tdClient *TdClient) ShowConsumers() ([]*maps.GoleMap, error) {
	return tdClient.Query("show consumers")
}

// ShowCreateDatabase 显示 db_name 指定的数据库的创建语句
func (tdClient *TdClient) ShowCreateDatabase() (string, error) {
	dbName := tdClient.DbName
	dataMaps, err := tdClient.Query("show create database " + dbName)
	if err != nil {
		return "", err
	}
	if len(dataMaps) == 0 {
		return "", nil
	}

	if result, exist := dataMaps[0].GetString("Create Database"); exist {
		return result, nil
	} else {
		return "", nil
	}
}

// ShowCreateStable 显示 tb_name 指定的超级表的创建语句
func (tdClient *TdClient) ShowCreateStable(stableName string) (string, error) {
	dbName := tdClient.DbName
	dataMaps, err := tdClient.Query("show create stable " + dbName + "." + stableName)
	if err != nil {
		return "", err
	}
	if len(dataMaps) == 0 {
		return "", nil
	}

	if result, exist := dataMaps[0].GetString("Create Table"); exist {
		return result, nil
	} else {
		return "", nil
	}
}

// ShowCreateTable 显示 tb_name 指定的表的创建语句。支持普通表、超级表和子表
func (tdClient *TdClient) ShowCreateTable(tableName string) (string, error) {
	dbName := tdClient.DbName
	dataMaps, err := tdClient.Query("show create table " + dbName + "." + tableName)
	if err != nil {
		return "", err
	}
	if len(dataMaps) == 0 {
		return "", nil
	}

	if result, exist := dataMaps[0].GetString("Create Table"); exist {
		return result, nil
	} else {
		return "", nil
	}
}

// ShowDatabases dbType：system：指定只显示系统数据库; user：指定只显示用户创建的数据库
func (tdClient *TdClient) ShowDatabases(dbType string) ([]string, error) {
	dataMaps, err := tdClient.Query(fmt.Sprintf("show %s databases", dbType))
	if err != nil {
		return nil, err
	}
	if len(dataMaps) == 0 {
		return nil, nil
	}

	var names []string
	for _, dataMap := range dataMaps {
		name, _ := dataMap.GetString("name")
		names = append(names, name)
	}
	return names, nil
}

// ShowDNodes 显示当前系统中 DNODE 的信息
func (tdClient *TdClient) ShowDNodes() ([]*maps.GoleMap, error) {
	return tdClient.Query("show dnodes")
}

// ShowFunctions 显示用户定义的自定义函数。
func (tdClient *TdClient) ShowFunctions() ([]*maps.GoleMap, error) {
	return tdClient.Query("show functions")
}

// ShowLicences 显示企业版许可授权的信息；注：企业版独有
func (tdClient *TdClient) ShowLicences() ([]*maps.GoleMap, error) {
	return tdClient.Query("show licences")
}

// ShowGrants 显示企业版许可授权的信息；注：企业版独有
func (tdClient *TdClient) ShowGrants() ([]*maps.GoleMap, error) {
	return tdClient.Query("show grants")
}

// ShowIndexes 显示已创建的索引
func (tdClient *TdClient) ShowIndexes(tableName string) ([]*maps.GoleMap, error) {
	dbName := tdClient.DbName
	return tdClient.Query("show indexes from " + dbName + "." + tableName)
}

// ShowLocalVariables 显示当前客户端配置参数的运行值。
func (tdClient *TdClient) ShowLocalVariables() ([]*maps.GoleMap, error) {
	return tdClient.Query("show local variables")
}

// ShowMNodes 显示当前系统中 MNODE 的信息
func (tdClient *TdClient) ShowMNodes() ([]*maps.GoleMap, error) {
	return tdClient.Query("show mnodes")
}

// ShowQNodes 显示当前系统中 QNODE （查询节点）的信息
func (tdClient *TdClient) ShowQNodes() ([]*maps.GoleMap, error) {
	return tdClient.Query("show qnodes")
}

// ShowQueries 显示当前系统中正在进行的查询
func (tdClient *TdClient) ShowQueries() ([]*maps.GoleMap, error) {
	return tdClient.Query("show queries")
}

// ShowScores 显示系统被许可授权的容量的信息 注：企业版独有
func (tdClient *TdClient) ShowScores() ([]*maps.GoleMap, error) {
	return tdClient.Query("show queries")
}

// ShowStables 显示当前数据库下的所有超级表的信息。会自动使用 LIKE 对表名进行模糊匹配
func (tdClient *TdClient) ShowStables(stableName string) ([]string, error) {
	dbName := tdClient.DbName
	var queryStr string
	if dbName != "" {
		queryStr = "show " + dbName + ".stables like '%" + stableName + "%'"
	} else {
		queryStr = "show stables like '%" + stableName + "%'"
	}
	dataMaps, err := tdClient.Query(queryStr)
	if err != nil {
		return nil, err
	}
	if len(dataMaps) == 0 {
		return nil, nil
	}

	var names []string
	for _, dataMap := range dataMaps {
		stableNameFind, _ := dataMap.GetString("stable_name")
		names = append(names, stableNameFind)
	}
	return names, nil
}

// ShowStreams 显示当前系统内所有流计算的信息
func (tdClient *TdClient) ShowStreams() ([]*maps.GoleMap, error) {
	return tdClient.Query("show streams")
}

// ShowSubscriptions 显示当前系统内所有的订阅关系
func (tdClient *TdClient) ShowSubscriptions() ([]*maps.GoleMap, error) {
	return tdClient.Query("show subscriptions")
}

// ShowTables 显示当前数据库下的所有普通表和子表的信息。会自动使用 LIKE 对表名进行模糊匹配
// 参数：
// tableRole：可为空
// - normal：指定只显示普通表信息
// - child：指定只显示子表信息
func (tdClient *TdClient) ShowTables(tableRole, tableName string) ([]string, error) {
	dbName := tdClient.DbName
	queryStr := "show"
	if tableRole != "" {
		queryStr += " " + tableRole
	}
	if dbName != "" {
		queryStr += " " + dbName + ".tables like '%" + tableName + "%'"
	} else {
		queryStr += " tables like '%" + tableName + "%'"
	}
	dataMaps, err := tdClient.Query(queryStr)
	if err != nil {
		return nil, err
	}
	if len(dataMaps) == 0 {
		return nil, nil
	}

	var names []string
	for _, dataMap := range dataMaps {
		tableNameFind, _ := dataMap.GetString("table_name")
		names = append(names, tableNameFind)
	}
	return names, nil
}

// ShowTags 显示子表的标签信息
func (tdClient *TdClient) ShowTags(tableName string) ([]*maps.GoleMap, error) {
	dbName := tdClient.DbName
	var queryStr string
	if dbName != "" {
		queryStr = "show tags from " + dbName + "." + tableName
	} else {
		queryStr = "show tags from" + tableName
	}
	return tdClient.Query(queryStr)
}

// ShowTopics 显示当前数据库下的所有主题的信息
func (tdClient *TdClient) ShowTopics() ([]string, error) {
	dataMaps, err := tdClient.Query("show topics")
	if err != nil {
		return nil, err
	}
	if len(dataMaps) == 0 {
		return nil, nil
	}

	var topics []string
	for _, dataMap := range dataMaps {
		topic, _ := dataMap.GetString("topic_name")
		topics = append(topics, topic)
	}
	return topics, nil
}

// ShowTransactions 显示当前系统中正在执行的事务的信息(该事务仅针对除普通表以外的元数据级别)
func (tdClient *TdClient) ShowTransactions() ([]*maps.GoleMap, error) {
	return tdClient.Query("show transactions")
}

// ShowUsers 显示当前系统中所有用户的信息。包括用户自定义的用户和系统默认用户
func (tdClient *TdClient) ShowUsers() ([]*maps.GoleMap, error) {
	return tdClient.Query("show users")
}

// ShowClusterVariables 显示当前系统中各节点需要相同的配置参数的运行值，也可以指定 DNODE 来查看其的配置参数
func (tdClient *TdClient) ShowClusterVariables() ([]*maps.GoleMap, error) {
	return tdClient.Query("show cluster variables")
}

// ShowVGroups 显示当前数据库中所有 VGROUP 的信息
func (tdClient *TdClient) ShowVGroups() ([]*maps.GoleMap, error) {
	dbName := tdClient.DbName
	queryStr := "show vgroups"
	if dbName != "" {
		queryStr = "show " + dbName + ".vgroups"
	}
	return tdClient.Query(queryStr)
}

// ShowVNodes 显示当前系统中所有 VNODE 或某个 DNODE 的 VNODE 的信息
func (tdClient *TdClient) ShowVNodes(dNodeId string) ([]*maps.GoleMap, error) {
	queryStr := "show vnodes"
	if dNodeId != "" {
		queryStr = "show vnodes on dnode " + dNodeId
	}
	return tdClient.Query(queryStr)
}

func doQueryOneString(cmd *TdClient, queryData string) (string, error) {
	dataMaps, err := cmd.Query("select " + queryData)
	if err != nil {
		return "", err
	}
	if len(dataMaps) == 0 {
		return "", nil
	}

	if result, exist := dataMaps[0].GetString(queryData); exist {
		return result, nil
	} else {
		return "", nil
	}
}
