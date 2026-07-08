package postgres

import (
	"os"
)

func Exec(sql string) error {
	_, err := GetDriver().Execute(sql)
	return err
}

func ExecFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return Exec(string(data))
}
