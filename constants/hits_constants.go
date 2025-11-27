package constants

const BATCH_SCAN = "BATCH_SCAN"               // 采用批量读表的方式	超级表 JOIN 语句
const NO_BATCH_SCAN = "NO_BATCH_SCAN"         // 采用sort方式进行分组, 与PARTITION_FIRST冲突	partition by 列表有普通列时
const SORT_FOR_GROUP = "SORT_FOR_GROUP"       // 在聚合之前使用PARTITION计算分组, 与SORT_FOR_GROUP冲突	partition by 列表有普通列时
const PARA_TABLES_SORT = "PARA_TABLES_SORT"   // 超级表的数据按时间戳排序时, 不使用临时磁盘空间, 只使用内存。当子表数量多, 行长比较大时候, 会使用大量内存, 可能发生OOM	超级表的数据按时间戳排序时
const SMALLDATA_TS_SORT = "SMALLDATA_TS_SORT" // 超级表的数据按时间戳排序时, 查询列长度大于等于256, 但是行数不多, 使用这个提示, 可以提高性能	超级表的数据按时间戳排序时
const SKIP_TSMA = "SKIP_TSMA"                 // 用于显示的禁用TSMA查询优化	带Agg函数的查询语句
