package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"sync"
)

type Driver struct {
	db *sqlx.DB
}

var once sync.Once
var driver *Driver

// NewDriver
// @Description: 初始化驱动
func NewDriver(username, password, host string, port int, dbName string) {
	once.Do(func() {
		// 拼接 dsn, 由调用方项目进行传入数据库相关配置
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, dbName)
		db, err := sqlx.Connect("mysql", dsn)

		if err != nil {
			panic("error connecting to mysql： " + err.Error())
		}
		driver = &Driver{db: db}
	})
}

// GetDriver
// @Description: 获取驱动
// @return *Driver
func GetDriver() *Driver {
	return driver
}

func (d *Driver) Execute(sql string) (sql.Result, error) {
	result, err := d.db.Exec(sql)
	return result, err
}

func (d *Driver) Query(sql string) (*sql.Rows, error) {
	return d.db.Query(sql)
}

func (d *Driver) Select(dest interface{}, sql string) error {
	return d.db.Select(dest, sql)
}

func (d *Driver) Debug() *sqlx.DB {
	return d.db
}

func (d *Driver) Close() error {
	return d.db.Close()
}
