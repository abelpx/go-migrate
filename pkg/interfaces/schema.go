package interfaces

type Schema interface {
	Create(table string, schemaFunc func(Foundation)) Seeds
	Table(table string, schemaFunc func(Foundation)) error
	DropIfExists(table string) error
}
