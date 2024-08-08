package mysql

import (
	"fmt"
	"strings"
)

type indexName struct {
	b       bool
	columns []string
}

// generateSql
// @Description: 生成索引
func (indexN *indexName) generateSql(isUnique bool, name string, column string) string {
	index := ""
	if isUnique {
		index = "UNIQUE "
	}
	if len(indexN.columns) > 0 {
		for i := range indexN.columns {
			indexN.columns[i] = fmt.Sprintf("`%s`", indexN.columns[i])
		}
		return fmt.Sprintf("%sINDEX %s (%s)", index, name, strings.Join(indexN.columns, ","))
	}
	return fmt.Sprintf("%sINDEX %s (`%s`)", index, name, column)
}
