package mysql

import (
	"fmt"
	"github.com/abelpx/go-migrate/pkg/interfaces"
)

type schema struct {
	driver interfaces.Driver
}

func NewSchema() interfaces.Schema {
	return &schema{}
}

// Create
// @Description: 创建表
func (s *schema) Create(table string, schemaFunc func(foundation interfaces.Foundation)) interfaces.Seeds {
	return NewSeeder(table, createWithDriver(GetDriver(), table, schemaFunc))
}

func createWithDriver(driver interfaces.Driver, table string, schemaFunc func(foundation interfaces.Foundation)) error {
	foundation := NewFoundation().(*Foundation)
	schemaFunc(foundation)
	sqls := foundation.GenSql(table, Operation.CREATE)
	for _, sql := range sqls {
		if _, err := driver.Execute(sql); err != nil {
			return err
		}
	}

	return nil
}

func tableWithDriver(driver interfaces.Driver, table string, schemaFunc func(foundation interfaces.Foundation)) error {

	foundation := NewFoundation().(*Foundation)
	schemaFunc(foundation)

	sqls := foundation.GenSql(table, Operation.ALTER)
	for _, sql := range sqls {
		if _, err := driver.Execute(sql); err != nil {
			return err
		}
	}

	return nil
}

func dropIfExistsWithDriver(driver interfaces.Driver, table string) error {

	sql := fmt.Sprintf("DROP TABLE IF EXISTS %s;", table)
	_, err := driver.Execute(sql)
	return err
}

func (s *schema) Table(table string, schemaFunc func(interfaces.Foundation)) error {
	return tableWithDriver(GetDriver(), table, schemaFunc)
}

func (s *schema) DropIfExists(table string) error {
	return dropIfExistsWithDriver(GetDriver(), table)
}
