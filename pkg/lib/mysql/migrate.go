package mysql

import (
	"fmt"
	"github.com/abelpx/go-migrate/pkg/interfaces"
	"github.com/abelpx/go-migrate/pkg/model"
	"strings"
)

type migrate struct {
}

func InitMigrate() interfaces.Migrate {
	return &migrate{}
}

func (m *migrate) CheckTable() (bool, error) {
	driver := GetDriver()

	sql := "SHOW TABLES LIKE 'migrations'"
	rows, err := driver.Query(sql)
	if err != nil {
		return false, err
	}

	return rows.Next(), nil
}

func (m *migrate) CreateTable() error {
	driver := GetDriver()

	sqls := []string{
		"CREATE TABLE `migrations` (`id` int(10) UNSIGNED NOT NULL, `migration` varchar(255) NOT NULL, `batch` int(11) NOT NULL);",
		"ALTER TABLE `migrations` ADD PRIMARY KEY (`id`);",
		"ALTER TABLE `migrations` MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;",
	}

	for _, sql := range sqls {
		if _, err := driver.Execute(sql); err != nil {
			return err
		}
	}

	return nil
}

func (m *migrate) DropTableIfExists() error {
	driver := GetDriver()

	sql := "DROP TABLE IF EXISTS migrations;"
	_, err := driver.Execute(sql)
	return err
}

func (m *migrate) DropAllTable() error {
	driver := GetDriver()

	tables := []string{}
	sql := "SHOW TABLES"
	driver.Select(&tables, sql)

	sql = fmt.Sprintf("DROP TABLE IF EXISTS %s;", strings.Join(tables, ","))
	_, err := driver.Execute(sql)
	return err
}

func (m *migrate) GetMigrations() ([]model.Migration, error) {
	driver := GetDriver()

	migrations := []model.Migration{}
	sql := "SELECT id, migration, batch FROM `migrations`"
	err := driver.Select(&migrations, sql)
	return migrations, err
}

func (m *migrate) WriteRecord(migration string, batch int) error {
	driver := GetDriver()

	sql := fmt.Sprintf("INSERT INTO `migrations`(`migration`, `batch`) VALUES ('%s','%d')", migration, batch)
	_, err := driver.Execute(sql)
	return err
}

func (m *migrate) DeleteRecord(id int) error {
	driver := GetDriver()

	sql := fmt.Sprintf("DELETE FROM `migrations` WHERE id = %d", id)
	_, err := driver.Execute(sql)
	return err
}
