package mysql

import (
	"fmt"
	"github.com/abelpx/go-migrate/pkg/interfaces"
	"github.com/laijunbin/go-solve-kit"
)

type Foundation struct {
	metadata []*metadata
}

func NewFoundation() interfaces.Foundation {
	return &Foundation{}
}

func (bp *Foundation) Id(name string, length int) interfaces.Foundation {
	bp.metadata = append(bp.metadata, &metadata{
		Name:          name,
		Type:          "BIGINT",
		Length:        length,
		AutoIncrement: true,
		Primary:       true,
		unsigned:      true,
		Comment:       "索引ID",
	})
	return bp
}

func (bp *Foundation) String(name string, length int) interfaces.Foundation {
	bp.metadata = append(bp.metadata, &metadata{
		Name:   name,
		Type:   "VARCHAR",
		Length: length,
	})
	return bp
}

func (bp *Foundation) Text(name string) interfaces.Foundation {
	bp.metadata = append(bp.metadata, &metadata{
		Name: name,
		Type: "TEXT",
	})
	return bp
}

func (bp *Foundation) CustomSql(sql string) interfaces.Foundation {
	bp.metadata = append(bp.metadata, &metadata{
		Custom: sql,
	})
	return bp
}

func (bp *Foundation) MediumText(name string) interfaces.Foundation {
	bp.metadata = append(bp.metadata, &metadata{
		Name: name,
		Type: "MEDIUMTEXT",
	})
	return bp
}

func (bp *Foundation) LongText(name string) interfaces.Foundation {
	bp.metadata = append(bp.metadata, &metadata{
		Name: name,
		Type: "LONGTEXT",
	})
	return bp
}

func (bp *Foundation) BigInt(name string, length int) interfaces.Foundation {
	bp.metadata = append(bp.metadata, &metadata{
		Name:   name,
		Type:   "BIGINT",
		Length: length,
	})
	return bp
}

func (bp *Foundation) Collate(collate string) interfaces.Foundation {
	if len(bp.metadata) != 0 {
		bp.metadata[len(bp.metadata)-1].Collate = collate
	}
	return bp
}

func (bp *Foundation) Decimal(name string, length, decimals int) interfaces.Foundation {
	bp.metadata = append(bp.metadata, &metadata{
		Name:     name,
		Type:     "DECIMAL",
		Length:   length,
		Decimals: decimals,
	})
	return bp
}

func (bp *Foundation) Integer(name string, length int) interfaces.Foundation {
	bp.metadata = append(bp.metadata, &metadata{
		Name:   name,
		Type:   "INT",
		Length: length,
	})
	return bp
}

func (bp *Foundation) Date(name string) interfaces.Foundation {
	bp.metadata = append(bp.metadata, &metadata{
		Name: name,
		Type: "DATE",
	})
	return bp
}

func (bp *Foundation) Boolean(name string) interfaces.Foundation {
	bp.metadata = append(bp.metadata, &metadata{
		Name: name,
		Type: "TINYINT",
	})
	return bp
}

func (bp *Foundation) Comment(value string) interfaces.Foundation {
	if len(bp.metadata) != 0 {
		bp.metadata[len(bp.metadata)-1].Comment = value
	}
	return bp
}

func (bp *Foundation) TableComment(value string) interfaces.Foundation {
	bp.metadata = append(bp.metadata, &metadata{
		TableComment: value,
	})
	return bp
}

func (bp *Foundation) DateTime(name string) interfaces.Foundation {
	bp.metadata = append(bp.metadata, &metadata{
		Name: name,
		Type: "DATETIME",
	})
	return bp
}

func (bp *Foundation) Timestamps() {
	bp.metadata = append(bp.metadata, &metadata{
		Name:    "created_at",
		Type:    "DATETIME",
		Default: "CURRENT_TIMESTAMP",
		Comment: "创建时间",
	})

	bp.metadata = append(bp.metadata, &metadata{
		Name:     "updated_at",
		Type:     "DATETIME",
		Nullable: true,
		Default:  "CURRENT_TIMESTAMP",
		Comment:  "更新时间",
	})
}

func (bp *Foundation) DeletedAt(b bool) {
	bp.metadata = append(bp.metadata, &metadata{
		Name:      "deleted_at",
		Type:      "DATETIME",
		Nullable:  true,
		Comment:   "删除时间",
		Index:     indexName{b: b},
		IndexName: "idx_deleted_at",
	})
}

func (bp *Foundation) Nullable() interfaces.Foundation {
	if len(bp.metadata) != 0 {
		bp.metadata[len(bp.metadata)-1].Nullable = true
	}
	return bp
}

func (bp *Foundation) Unsigned() interfaces.Foundation {
	if len(bp.metadata) != 0 {
		bp.metadata[len(bp.metadata)-1].unsigned = true
	}
	return bp
}

func (bp *Foundation) Unique(column ...string) interfaces.Foundation {
	if len(column) == 0 && len(bp.metadata) > 0 {
		bp.metadata[len(bp.metadata)-1].Unique = indexName{b: true}
	}
	if len(column) > 0 {
		bp.metadata = append(bp.metadata, &metadata{Unique: indexName{b: true, columns: column}})
	}
	return bp
}

func (bp *Foundation) Index(column ...string) interfaces.Foundation {
	if len(column) == 0 {
		bp.metadata[len(bp.metadata)-1].Index = indexName{b: true}
	} else {
		bp.metadata = append(bp.metadata, &metadata{Index: indexName{b: true, columns: column}})
	}
	return bp
}

func (bp *Foundation) Modify() interfaces.Foundation {
	bp.metadata[len(bp.metadata)-1].Modify = true
	return bp
}

func (bp *Foundation) IndexName(name string) interfaces.Foundation {
	if name != "" && len(bp.metadata) > 0 {
		bp.metadata[len(bp.metadata)-1].IndexName = name
	}
	return bp
}

func (bp *Foundation) Default(value interface{}) interfaces.Foundation {
	if len(bp.metadata) != 0 {
		bp.metadata[len(bp.metadata)-1].Default = fmt.Sprintf("'%v'", value)
	}
	return bp
}

func (bp *Foundation) Foreign(name string) interfaces.ForeignFoundation {
	fb := newForeignFoundation().(*foreignFoundation)
	bp.metadata = append(bp.metadata, &metadata{
		Name:    name,
		Foreign: fb.meta,
	})
	return fb
}

func (bp *Foundation) Primary(name ...string) interfaces.Foundation {
	bp.metadata = append(bp.metadata, &metadata{
		Name:    go_solve_kit.FromStringArray(name).Join("`, `").ValueOf(),
		Primary: true,
	})
	return bp
}

func (bp *Foundation) DropColumn(column string) {
	bp.metadata = append(bp.metadata, &metadata{
		Name: column,
		Type: "DROP",
	})
}

func (bp *Foundation) DropUnique(name string) {
	bp.metadata = append(bp.metadata, &metadata{
		Name:   name,
		Type:   "DROP",
		Unique: indexName{b: true},
	})
}
func (bp *Foundation) DropIndex(name string) {
	bp.metadata = append(bp.metadata, &metadata{
		Name:  name,
		Type:  "DROP",
		Index: indexName{b: true},
	})
}
func (bp *Foundation) DropForeign(name string) {
	bp.metadata = append(bp.metadata, &metadata{
		Name:    name,
		Type:    "DROP",
		Foreign: newForeignFoundation().(*foreignFoundation).meta,
	})
}
func (bp *Foundation) DropPrimary() {
	bp.metadata = append(bp.metadata, &metadata{
		Type:    "DROP",
		Primary: true,
	})
}
func (f *Foundation) GenSql(table string, operation operation) []string {
	return operation.generateSql(table, f.metadata)
}
