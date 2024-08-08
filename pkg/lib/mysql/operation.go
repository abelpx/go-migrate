package mysql

import (
	"fmt"
	"github.com/laijunbin/go-solve-kit"
)

type operation interface {
	generateSql(table string, metadata []*metadata) []string
}

type createOperation struct{}

func (c createOperation) generateSql(table string, meta []*metadata) []string {
	columns := go_solve_kit.FromInterfaceArray(meta)
	tableComment := ""
	proceedColumns := columns.Map(func(v go_solve_kit.Type, i int) interface{} {
		m := v.ValueOf().(*metadata)

		if m.Type == "DROP" {
			return nil
		}

		if m.Foreign != nil {
			return m.Foreign.generateSql(table, m.Name)
		}

		s := ""
		if m.Type != "" {
			s += fmt.Sprintf("`%s` %s", m.Name, m.Type)
		}
		if m.Custom != "" {
			return m.Custom
		}
		if m.Decimals != 0 {
			s += fmt.Sprintf("(%d,%d)", m.Length, m.Decimals)
		} else if m.Length != 0 {
			s += fmt.Sprintf("(%d)", m.Length)
		}

		if m.Collate != "" {
			s += " COLLATE " + m.Collate
		}

		if m.unsigned {
			s += " UNSIGNED"
		}

		if s != "" && !m.Nullable {
			s += " NOT NULL"
		}

		if m.AutoIncrement {
			s += " AUTO_INCREMENT"
		}

		if m.Comment != "" && m.TableComment == "" {
			s += fmt.Sprintf(` COMMENT "%s"`, m.Comment)
		}

		if m.Default != nil {
			s += fmt.Sprintf(" DEFAULT %v", m.Default)
		}

		if m.Primary {
			if s != "" {
				s += ", "
			}
			s += fmt.Sprintf("PRIMARY KEY (`%s`)", m.Name)
		}

		var indexNameSql string
		if m.Unique.b {
			if len(m.Unique.columns) > 0 {
				return fmt.Sprintf("%s", m.Unique.generateSql(true, m.IndexName, ""))
			} else {
				if m.IndexName != "" {
					if s != "" {
						indexNameSql = ", "
					}
					indexNameSql = fmt.Sprintf("%s%s", indexNameSql, m.Unique.generateSql(true, m.IndexName, m.Name))
				} else {
					if s != "" {
						s += ", "
					}
					s += fmt.Sprintf("UNIQUE (`%s`)", m.Name)
				}
			}
		}

		if m.Index.b {
			if len(m.Index.columns) > 0 {
				return fmt.Sprintf("%s", m.Index.generateSql(false, m.IndexName, ""))
			} else {
				if m.IndexName != "" && indexNameSql == "" {
					if s != "" {
						indexNameSql = ", "
					}
					indexNameSql = fmt.Sprintf("%s%s", indexNameSql, m.Index.generateSql(false, m.IndexName, m.Name))
				} else {
					if s != "" {
						s += ", "
					}
					s += fmt.Sprintf("INDEX (`%s`)", m.Name)
				}
			}
		}

		if m.TableComment != "" {
			tableComment = m.TableComment
			return nil
		}
		return s + indexNameSql
	})

	columnsStr := proceedColumns.Filter(func(s go_solve_kit.Type, i int) bool {
		return s.ValueOf() != nil
	}).ToStringArray().Join(", ").ValueOf()
	sql := fmt.Sprintf("CREATE TABLE `%s` (%s)", table, columnsStr)
	if tableComment != "" {
		sql += fmt.Sprintf(" comment='%s'", tableComment)
	}
	fmt.Println(sql + ";\n")
	return []string{sql + ";"}
}

type alterOperation struct{}

func (a alterOperation) generateSql(table string, meta []*metadata) []string {
	columns := go_solve_kit.FromInterfaceArray(meta)
	tableComment := ""
	sql := fmt.Sprintf("ALTER TABLE `%s` %s", table, columns.Map(func(v go_solve_kit.Type, i int) interface{} {
		m := v.ValueOf().(*metadata)

		if m.Type == "DROP" {
			if m.Primary {
				return "DROP PRIMARY KEY"
			}

			if m.Index.b || m.Unique.b {
				return fmt.Sprintf("DROP INDEX `%s`", m.Name)
			}

			if m.Foreign != nil {
				return fmt.Sprintf("DROP FOREIGN KEY `%[1]s`, DROP INDEX `%[1]s`", fmt.Sprintf("fk_%s_%s", table, m.Name))
			}

			return fmt.Sprintf("DROP `%s`", m.Name)
		}

		if m.Foreign != nil {
			return fmt.Sprintf("ADD %s", m.Foreign.generateSql(table, m.Name))
		}

		s := "ADD "
		if m.Modify {
			s = "MODIFY "
		}

		if m.Type != "" {
			s += fmt.Sprintf("`%s` %s", m.Name, m.Type)
		}
		if m.Custom != "" {
			return m.Custom
		}
		if m.Decimals != 0 {
			s += fmt.Sprintf("(%d,%d)", m.Length, m.Decimals)
		} else if m.Length != 0 {
			s += fmt.Sprintf("(%d)", m.Length)
		}

		if m.Collate != "" {
			s += " COLLATE " + m.Collate
		}

		if m.unsigned {
			s += " UNSIGNED"
		}

		if s != "" && !m.Nullable {
			s += " NOT NULL"
		}

		if m.AutoIncrement {
			s += " AUTO_INCREMENT"
		}

		if m.Default != nil {
			s += fmt.Sprintf(" DEFAULT %v", m.Default)
		}

		if m.Comment != "" && m.TableComment == "" {
			s += fmt.Sprintf(` COMMENT "%s"`, m.Comment)
		}

		if m.Primary {
			if s != "" {
				s += ", "
			}
			fmt.Println(s)
			s += fmt.Sprintf("ADD PRIMARY KEY (`%s`)", m.Name)
		}

		var indexNameSql string
		if m.Unique.b {
			if len(m.Unique.columns) > 0 {
				return fmt.Sprintf("ADD %s", m.Unique.generateSql(true, m.IndexName, ""))
			} else {
				if m.IndexName != "" {
					indexNameSql = fmt.Sprintf(", ADD %s", m.Unique.generateSql(true, m.IndexName, m.Name))
				} else {
					s += fmt.Sprintf("ADD UNIQUE (`%s`)", m.Name)
				}
			}
		}

		if m.Index.b {
			if len(m.Index.columns) > 0 {
				return fmt.Sprintf("ADD %s", m.Index.generateSql(false, m.IndexName, ""))
			} else {
				if m.IndexName != "" && indexNameSql == "" {
					if s != "" {
						indexNameSql += ", "
					}
					indexNameSql = fmt.Sprintf("%sADD %s", indexNameSql, m.Index.generateSql(false, m.IndexName, m.Name))
				} else {
					if s != "" {
						s += ", "
					}
					s += fmt.Sprintf("ADD INDEX (`%s`)", m.Name)
				}
			}
		}
		if m.TableComment != "" {
			tableComment = m.TableComment
			return nil
		}
		return s + indexNameSql
	}).ToStringArray().Join(", ").ValueOf())
	if tableComment != "" {
		sql += fmt.Sprintf(" comment='%s'", tableComment)
	}
	fmt.Println(sql + ";\n")
	return []string{
		sql + ";",
	}
}

var Operation = struct {
	CREATE operation
	ALTER  operation
}{
	CREATE: &createOperation{},
	ALTER:  &alterOperation{},
}
