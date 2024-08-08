package interfaces

type Foundation interface {
	Id(name string, length int) Foundation
	String(name string, length int) Foundation
	Text(name string) Foundation
	LongText(name string) Foundation
	MediumText(name string) Foundation
	CustomSql(sql string) Foundation
	BigInt(name string, length int) Foundation
	Integer(name string, length int) Foundation
	Decimal(name string, length, decimals int) Foundation
	Date(name string) Foundation
	Comment(value string) Foundation
	Collate(collate string) Foundation
	TableComment(value string) Foundation
	Boolean(name string) Foundation
	DateTime(name string) Foundation
	Nullable() Foundation
	Unsigned() Foundation
	Modify() Foundation
	Unique(column ...string) Foundation
	Index(column ...string) Foundation
	IndexName(name string) Foundation
	Default(value interface{}) Foundation
	Foreign(name string) ForeignFoundation
	Primary(name ...string) Foundation
	DropColumn(column string)
	DropUnique(name string)
	DropIndex(name string)
	DropForeign(name string)
	DropPrimary()
	Timestamps()
	DeletedAt(index bool)
}
