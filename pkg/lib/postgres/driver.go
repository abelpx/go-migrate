package postgres

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Driver struct {
	db *sqlx.DB
}

var once sync.Once
var driver *Driver

func NewDriver(username, password, host string, port int, dbName string) {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host,
			port,
			username,
			password,
			dbName,
		)
		db, err := sqlx.Connect("postgres", dsn)
		if err != nil {
			panic("error connecting to postgres: " + err.Error())
		}
		driver = &Driver{db: db}
	})
}

func GetDriver() *Driver {
	return driver
}

func (d *Driver) Execute(sql string) (sql.Result, error) {
	return d.db.Exec(sql)
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
