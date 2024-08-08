package template

var NewTemplate = `package migrations

import (
	"github.com/abelpx/go-migrate/pkg/interfaces"
	"github.com/abelpx/go-migrate/pkg/lib/mysql"
)

type %[2]s struct{}

func Create%[2]s() interfaces.Migration {
	return &%[2]s{}
}

func (t *%[2]s) Up() error {
	
}

func (t *%[2]s) Down() error {

}
`

var CreateTemplate = `package migrations

import (
	"github.com/abelpx/go-migrate/pkg/interfaces"
	"github.com/abelpx/go-migrate/pkg/lib/mysql"
)

type %[2]s struct{}

func Create%[2]s() interfaces.Migration {
	return &%[2]s{}
}

func (t *%[2]s) Up() error {
	return mysql.NewSchema().Create("%[3]s", func(table interfaces.Foundation) {
		table.Id("id", 22)
		table.Timestamps()
	})
}

func (t *%[2]s) Down() error {
	return mysql.NewSchema().DropIfExists("%[3]s")
}
`
var AlterTemplate = `package migrations

import (
	"github.com/abelpx/go-migrate/pkg/interfaces"
	"github.com/abelpx/go-migrate/pkg/lib/mysql"
)

type %[2]s struct{}

func Create%[2]s() interfaces.Migration {
	return &%[2]s{}
}

func (t *%[2]s) Up() error {
	return mysql.NewSchema().Table("%[3]s", func(table interfaces.Foundation) {

	})
}

func (t *%[2]s) Down() error {
	return mysql.NewSchema().Table("%[3]s", func(table interfaces.Foundation) {
		
	})
}
`
