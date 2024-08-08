package mysql

import (
	"github.com/abelpx/go-migrate/pkg/interfaces"
	"strings"
)

type foreignFoundation struct {
	meta *foreign
}

func newForeignFoundation() interfaces.ForeignFoundation {
	return &foreignFoundation{
		meta: &foreign{},
	}
}

func (f *foreignFoundation) Reference(name string) interfaces.ForeignFoundation {
	f.meta.Reference = name
	return f
}

func (f *foreignFoundation) On(table string) interfaces.ForeignFoundation {
	f.meta.Table = table
	return f
}

func (f *foreignFoundation) OnUpdate(action string) interfaces.ForeignFoundation {
	f.meta.OnUpdate = strings.ToUpper(action)
	return f
}

func (f *foreignFoundation) OnDelete(action string) interfaces.ForeignFoundation {
	f.meta.OnDelete = strings.ToUpper(action)
	return f
}
