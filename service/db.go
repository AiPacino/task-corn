package service

import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/Masterminds/squirrel"
  "fmt"
  "time"
  "os"
)

var db *sql.DB
var dbWriter *sql.DB
var dbCache squirrel.DBProxyBeginner
var dbCacheWriter squirrel.DBProxyBeginner    // 写操作数据库连接
var dbLog = AppLog.With("file", "db")

func OpenDB() squirrel.DBProxyBeginner {
  log := dbLog.With("func", "OpenDB")
  var err error
  DBUrl := os.Getenv("DB_URL")
  db, err = sql.Open("mysql", DBUrl)
  db.SetMaxOpenConns(300)
  db.SetMaxIdleConns(0)
  db.SetConnMaxLifetime(10 * time.Second)
  log.Infof("Connecting to %s ... \n", DBUrl)
  if err != nil {
    log.Errorf("DB conection set up failed, %s\n", err.Error())
    panic(err)
  }
  err = db.Ping()
  if err != nil {
    log.Errorf("DB conection set up failed, %s\n", err.Error())
    panic(err)
  }
  log.Infof("DB conection set up successfully \n")
  dbCache = squirrel.NewStmtCacheProxy(db)
  return dbCache
}

func CloseDB() {
  db.Close()
}

// 写操作数据库连接
func OpenWriterDB() squirrel.DBProxyBeginner {
  log := dbLog.With("func", "OpenWriterDB")
  var err error
  writerDBUrl := os.Getenv("WRITER_DB_URL")
  dbWriter, err = sql.Open("mysql", writerDBUrl)
  dbWriter.SetMaxOpenConns(300)
  dbWriter.SetMaxIdleConns(0)
  dbWriter.SetConnMaxLifetime(10 * time.Second)
  log.Infof("Connecting to %s ... \n", writerDBUrl)
  if err != nil {
    log.Errorf("DB Writer conection set up failed, %s\n", err.Error())
    panic(err)
  }
  err = dbWriter.Ping()
  if err != nil {
    log.Errorf("DB Writer conection set up failed, %s\n", err.Error())
    panic(err)
  }
  log.Infof("DB Writer conection set up successfully \n")
  dbCacheWriter = squirrel.NewStmtCacheProxy(dbWriter)
  return dbCacheWriter
}

func CloseWriterDB() {
  dbWriter.Close()
}

type txDBProxyBeginnerWrapper struct {
  tx *sql.Tx
}

func (w txDBProxyBeginnerWrapper) Begin() (*sql.Tx, error) {
  fmt.Println("return existed tx from txDBProxyBeginnerWrapper")
  return w.tx, nil
}
func (w txDBProxyBeginnerWrapper) Exec(query string, args ...interface{}) (sql.Result, error) {
  return w.tx.Exec(query, args...)
}

func (w txDBProxyBeginnerWrapper) Query(query string, args ...interface{}) (*sql.Rows, error) {
  return w.tx.Query(query, args...)
}
func (w txDBProxyBeginnerWrapper) QueryRow(query string, args ...interface{}) squirrel.RowScanner {
  return w.tx.QueryRow(query, args...)
}
func (w txDBProxyBeginnerWrapper) Prepare(query string) (*sql.Stmt, error) {
  return w.tx.Prepare(query)
}
func txToBeginner(tx *sql.Tx) squirrel.DBProxyBeginner {
  return txDBProxyBeginnerWrapper{tx: tx}
}

func inTx(f func(squirrel.DBProxyBeginner) error) error {
  tx, err := dbCacheWriter.Begin()

  if err != nil {
    fmt.Printf(err.Error())
    return err
  }
  defer tx.Rollback()

  dbProxyBeginner := txToBeginner(tx)
  if err := f(dbProxyBeginner); err != nil {
    return err
  }
  tx.Commit()
  return nil
}
