---
title: go-sql-driver 源码分析
layout: post
category: golang
author: 夏泽民
---
1、建立连接
首先是Open，
db, err := sql.Open(“mysql”, “user:password@/dbname”)
db 是一个*sql.DB类型的指针，在后面的操作中，都要用到db
open之后，并没有与数据库建立实际的连接，与数据库建立实际的连接是通过Ping方法完成。此外，db应该在整个程序的生命周期中存在，也就是说，程序一启动，就通过Open获得db，直到程序结束，再Close db，而不是经常Open/Close。
err = db.Ping()

https://studygolang.com/articles/7372
https://github.com/go-sql-driver/mysql/wiki/Examples

2、基本用法
DB的主要方法有：
Query 执行数据库的Query操作，例如一个Select语句，返回*Rows

QueryRow 执行数据库至多返回1行的Query操作，返回*Row

PrePare 准备一个数据库query操作，返回一个*Stmt，用于后续query或执行。这个Stmt可以被多次执行，或者并发执行

Exec 执行数不返回任何rows的据库语句，例如delete操作

Stmt的主要方法：
Exec
Query
QueryRow
Close
用法与DB类似
Rows的主要方法：
Cloumns： 返回[]string，column names
Scan：
Next：
Close：

https://studygolang.com/articles/7372


import (
"database/sql"
_ "github.com/go-sql-driver/mysql"
)


https://github.com/go-sql-driver/mysql/wiki/Examples


driver.go
func init() {
	sql.Register("mysql", &MySQLDriver{})
}


func (d MySQLDriver) Open(dsn string) (driver.Conn, error) {
	cfg, err := ParseDSN(dsn)
 return c.Connect(context.Background())
}


type DialContextFunc func(ctx context.Context, addr string) (net.Conn, error)
var dials     map[string]DialContextFunc

func RegisterDialContext(net string, dial DialContextFunc) {
}



connector.go

type connector struct {
	cfg *Config // immutable private copy.
}

func (c *connector) Connect(ctx context.Context) (driver.Conn, error){
   mc := &mysqlConn{
		maxAllowedPacket: maxPacketSize,
		maxWriteSize:     maxPacketSize - 1,
		closech:          make(chan struct{}),
		cfg:              c.cfg,
    }
   nd := net.Dialer{Timeout: mc.cfg.Timeout}
   mc.netConn, err = dial(dctx, mc.cfg.Addr)
   authResp, err := mc.auth(authData, plugin)
   
}

func (c *connector) Driver() driver.Driver {
	return &MySQLDriver{}
}


connection.go
type mysqlConn struct {
}

func (mc *mysqlConn) Prepare(query string) (driver.Stmt, error) {
   err := mc.writeCommandPacketStr(comStmtPrepare, query)
   stmt := &mysqlStmt{
		mc: mc,
   }
   columnCount, err := stmt.readPrepareResultPacket()
}

func (mc *mysqlConn) Begin() (driver.Tx, error) {
}
func (mc *mysqlConn) Close() (err error){
}


statement.go
type mysqlStmt struct {
	mc         *mysqlConn
	id         uint32
	paramCount int
}
func (stmt *mysqlStmt) Close() error 
func (stmt *mysqlStmt) Exec(args []driver.Value) (driver.Result, error) {
}
func (stmt *mysqlStmt) Query(args []driver.Value) (driver.Rows, error) {
}
func (stmt *mysqlStmt) NumInput() int 


var (
	_ driver.ConnBeginTx        = &mysqlConn{}
	_ driver.ConnPrepareContext = &mysqlConn{}
	_ driver.ExecerContext      = &mysqlConn{}
	_ driver.Pinger             = &mysqlConn{}
	_ driver.QueryerContext     = &mysqlConn{}
)


查询
func (mc *mysqlConn) query(query string, args []driver.Value) (*textRows, error) {
   // try client-side prepare to reduce roundtrip
    prepared, err := mc.interpolateParams(query, args)
    err := mc.writeCommandPacketStr(comQuery, query)
    resLen, err = mc.readResultSetHeaderPacket()
    switch err := rows.NextResultSet();
    rows.rs.columns, err = mc.readColumns(resLen)
}


src/database/sql/sql.go

var 	drivers   = make(map[string]driver.Driver)
func Register(name string, driver driver.Driver) {
    drivers[name] = driver
}

src/database/driver/driver.go
type Driver interface {
    Open(name string) (Conn, error)
}

type Conn interface {
   Prepare(query string) (Stmt, error)
   Close() error
   Begin() (Tx, error)
}


type Stmt interface {
   Close() error
   NumInput() int
   Exec(args []Value) (Result, error)
   Query(args []Value) (Rows, error)
}

type Result interface {
  LastInsertId() (int64, error)
  RowsAffected() (int64, error)
}

type Rows interface {
    Columns() []string
    Close() error
    Next(dest []Value) error
}



初始化：
func Open(driverName, dataSourceName string) (*DB, error) {
 //获取驱动
 driveri, ok := drivers[driverName]
connector, err := driverCtx.OpenConnector(dataSourceName)
//连接数据库
return OpenDB(dsnConnector{dsn: dataSourceName, driver: driveri}), nil
}

func OpenDB(c driver.Connector) *DB {
   go db.connectionOpener(ctx)
}


// Runs in a separate goroutine, opens new connections when requested.
func (db *DB) connectionOpener(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-db.openerCh:
			db.openNewConnection(ctx)
		}
	}
}

func (db *DB) openNewConnection(ctx context.Context) {
	ci, err := db.connector.Connect(ctx)
	db.maybeOpenNewConnections()
}

func (db *DB) maybeOpenNewConnections() {
     db.openerCh <- struct{}{}
}


func (dc *driverConn) finalClose() error {
dc.db.maybeOpenNewConnections()
}

func (db *DB) conn(ctx context.Context, strategy connReuseStrategy) (*driverConn, error) {
    db.maybeOpenNewConnections()
}

func (db *DB) putConn(dc *driverConn, err error, resetSession bool) {
   db.maybeOpenNewConnections()
}



查询
func (db *DB) query(ctx context.Context, query string, args []interface{}, strategy connReuseStrategy) (*Rows, error) {
dc, err := db.conn(ctx, strategy)
return db.queryDC(ctx, nil, dc, dc.releaseConn, query, args)
}

func (db *DB) queryDC(ctx, txctx context.Context, dc *driverConn, releaseConn func(error), query string, args []interface{}) (*Rows, error) {
  nvdargs, err = driverArgsConnLocked(dc.ci, nil, args)
  rowsi, err = ctxDriverQuery(ctx, queryerCtx, queryer, query, nvdargs)
}

database/sql/ctxutil.go

func ctxDriverQuery(ctx context.Context, queryerCtx driver.QueryerContext, queryer driver.Queryer, query string, nvdargs []driver.NamedValue) (driver.Rows, error) {
   return queryerCtx.QueryContext(ctx, query, nvdargs)
   return queryer.Query(query, dargs)
}

执行
func (db *DB) exec(ctx context.Context, query string, args []interface{}, strategy connReuseStrategy) (Result, error) {
dc, err := db.conn(ctx, strategy)
return db.execDC(ctx, dc, dc.releaseConn, query, args)
}



https://www.cnblogs.com/simpman/p/6510604.html
https://www.mysqlzh.com/api/89.html
https://studygolang.com/articles/1795
https://learnku.com/go/t/49692
SQL 预编译技术
https://blog.csdn.net/weixin_29586681/article/details/113327370
https://www.cnblogs.com/sflik/p/4587368.html
https://blog.csdn.net/qq_37102984/article/details/108988837
https://blog.csdn.net/weixin_33728708/article/details/90620309
https://www.zhihu.com/question/375120061


Buffer.go
buffer 是一个用于给 数据库连接 (net.Conn) 进行缓冲的一个数据结构，其结构为：

type buffer struct {
    buf     []byte     // 缓冲池中的数据
    nc      net.Conn   // 负责缓冲的数据库连接对象
    idx     int        // 已读数据索引
    length  int        // 缓冲池中未读数据的长度
    timeout time.Duration // 数据库连接的超时设置
}
可以看到，因为 数据库连接 (net.Conn) 在通信的时候是 同步 的。而为了让其能够 同时 读/写 ，所以实现了 buffer 这个数据结构，通过该 buffer 进行数据缓冲还能实现 零拷贝 ( zero-copy-ish ) 

Collations.go
collations 包含了 MySQL 所有支持的 字符集 格式，并支持通过 COLLATION_NAME 返回其字符集 ID


Dsn.go
DSN 即 数据源名称 （Data Source Name） ，是 驱动程序连接数据库的变量信息 ，简而言之就是根据你连接的不同数据库使用对应的连接信息。

通常，数据库的连接配置就是在这里定义的：


// Config 基本的数据库连接信息
type Config struct {
    User         string            // Username
    Passwd       string            // Password (requires User)
    Net          string            // Network type
    Addr         string            // Network address (requires Net)
    DBName       string            // Database name
    Params       map[string]string // Connection parameters
    Collation    string            // Connection collation
    Loc          *time.Location    // Location for time.Time values
    TLSConfig    string            // TLS configuration name
    tls          *tls.Config       // TLS configuration
    Timeout      time.Duration     // Dial timeout
    ReadTimeout  time.Duration     // I/O read timeout
    WriteTimeout time.Duration     // I/O write timeout

    AllowAllFiles           bool // 允许文件使用 LOAD DATA LOCAL INFILE 导入数据库
    AllowCleartextPasswords bool // 支持明文密码客户端
    AllowOldPasswords       bool // 允许使用不可靠的旧密码
    ClientFoundRows         bool // 返回匹配的行数而不是受影响的行数
    ColumnsWithAlias        bool // 将表名前置在列名
    InterpolateParams       bool // 将占位符插入查询的SQL字符串
    MultiStatements         bool // 允许一条语句多次查询
    ParseTime               bool // 格式化时间值为 time.Time 变量
    Strict                  bool // 将 warnings 返回 errors
}

Packets.go
接下来就要深入到 MySQL 的通信协议

https://www.jianshu.com/p/8e0bfed0bb90


<!-- more -->
