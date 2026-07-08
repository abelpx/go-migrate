package postgres

import (
	"fmt"
	"strings"

	"github.com/abelpx/go-migrate/pkg/interfaces"
	"github.com/abelpx/go-migrate/pkg/model"
)

type migrate struct{}

func InitMigrate() interfaces.Migrate {
	return &migrate{}
}

func (m *migrate) CheckTable() (bool, error) {
	var exists bool
	err := GetDriver().Debug().Get(&exists, `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = 'public'
			  AND table_name = 'migrations'
		)
	`)
	return exists, err
}

func (m *migrate) CreateTable() error {
	_, err := GetDriver().Execute(`
		CREATE TABLE IF NOT EXISTS public.migrations (
			id SERIAL PRIMARY KEY,
			migration VARCHAR(255) NOT NULL UNIQUE,
			batch INTEGER NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	return err
}

func (m *migrate) DropTableIfExists() error {
	_, err := GetDriver().Execute("DROP TABLE IF EXISTS public.migrations;")
	return err
}

func (m *migrate) DropAllTable() error {
	tables := []string{}
	err := GetDriver().Select(&tables, `
		SELECT quote_ident(table_schema) || '.' || quote_ident(table_name)
		FROM information_schema.tables
		WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
		  AND table_type = 'BASE TABLE'
	`)
	if err != nil {
		return err
	}
	if len(tables) == 0 {
		return nil
	}
	_, err = GetDriver().Execute(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", strings.Join(tables, ",")))
	return err
}

func (m *migrate) GetMigrations() ([]model.Migration, error) {
	migrations := []model.Migration{}
	err := GetDriver().Select(&migrations, "SELECT id, migration, batch FROM public.migrations ORDER BY id ASC")
	return migrations, err
}

func (m *migrate) WriteRecord(migration string, batch int) error {
	_, err := GetDriver().Debug().Exec(
		"INSERT INTO public.migrations(migration, batch) VALUES ($1, $2)",
		migration,
		batch,
	)
	return err
}

func (m *migrate) DeleteRecord(id int) error {
	_, err := GetDriver().Debug().Exec("DELETE FROM public.migrations WHERE id = $1", id)
	return err
}
