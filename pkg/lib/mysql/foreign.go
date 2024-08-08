package mysql

import "fmt"

type foreign struct {
	Reference string
	Table     string
	OnUpdate  string
	OnDelete  string
}

func (f *foreign) generateSql(table string, name string) string {
	s := fmt.Sprintf("CONSTRAINT `%s` FOREIGN KEY (`%s`) REFERENCES `%s`(`%s`)",
		fmt.Sprintf("fk_%s_%s", table, name),
		name,
		f.Table,
		f.Reference,
	)

	if f.OnUpdate != "" {
		s += fmt.Sprintf(" ON UPDATE %s", f.OnUpdate)
	}

	if f.OnDelete != "" {
		s += fmt.Sprintf(" ON DELETE %s", f.OnDelete)
	}

	return s
}
