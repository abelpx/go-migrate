package mysql

import (
	"fmt"
	"github.com/abelpx/go-migrate/pkg/interfaces"
	"github.com/abelpx/go-migrate/pkg/model"
	"github.com/laijunbin/go-solve-kit"
	"sort"
)

type seeds struct {
	*model.Seeds
	table string
}

func NewSeeder(table string, err error) interfaces.Seeds {
	return &seeds{
		Seeds: model.NewSeed(err),
		table: table,
	}
}

func (s *seeds) Seed(data ...map[string]interface{}) error {
	return runSeed(GetDriver(), s.table, data...)
}

func runSeed(driver interfaces.Driver, table string, data ...map[string]interface{}) error {
	for i := 0; i < len(data); i++ {
		var keys []string
		var values []string
		for k, _ := range data[i] {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for j, k := range keys {
			values = append(values, fmt.Sprintf("'%s'", data[i][k]))
			keys[j] = fmt.Sprintf("`%s`", keys[j])
		}

		sql := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s);",
			table,
			go_solve_kit.FromStringArray(keys).Join(", "),
			go_solve_kit.FromStringArray(values).Join(", "),
		)

		fmt.Println(sql + "\n")

		if _, err := driver.Execute(sql); err != nil {
			return err
		}
	}

	return nil
}
