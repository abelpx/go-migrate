package go_migrate

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/abelpx/go-migrate/pkg/interfaces"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/abelpx/go-migrate/pkg/lib/mysql"
	"github.com/abelpx/go-migrate/pkg/template"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type DbConfig struct {
	MigratePath string // 迁移文件路径，绝对路径，用 filepath.Abs() 来获取
	Username    string
	Password    string
	Host        string
	Port        int
	DbName      string
	DbType      string // 数据库类型 (例如: mysql, postgres)
}

var (
	Config     DbConfig
	Migrations []interfaces.Migration
	Migrate    interfaces.Migrate
)

var StartCmd = &cobra.Command{
	Use:          "migrate",
	Short:        "该命令用于管理数据库迁移相关操作。",
	SilenceUsage: true,
}

func init() {
	StartCmd.AddCommand(&cobra.Command{
		Use:   "new",
		Short: "创建新的迁移文件",
		PreRun: func(cmd *cobra.Command, args []string) {
			initializeDriver()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return createMigrateFile(args)
		},
	})

	StartCmd.AddCommand(&cobra.Command{
		Use:   "rollback",
		Short: "该命令将迁移到数据库中的表回滚回来",
		PreRun: func(cmd *cobra.Command, args []string) {
			initializeDriver()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return rollback()
		},
	})

	StartCmd.AddCommand(&cobra.Command{
		Use:   "run",
		Short: "迁移进入数据库。",
		PreRun: func(cmd *cobra.Command, args []string) {
			checkDatabase()
			initializeDriver()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	})
}

// 初始化数据库驱动
func initializeDriver() {
	switch Config.DbType {
	case "mysql":
		// 初始化 MySQL 驱动
		mysql.NewDriver(Config.Username, Config.Password, Config.Host, Config.Port, Config.DbName)
		Migrate = mysql.InitMigrate()
	case "sqlite3":
		// TODO: 后续添加 sqlite3 支持
		log.Fatalf("sqlite3 support not implemented")
	case "postgres":
		// TODO: 后续添加 postgres 支持
		log.Fatalf("postgres support not implemented")
	default:
		log.Fatalf("unsupported database type")
	}
}

// 创建迁移文件
func createMigrateFile(args []string) error {
	if len(args) < 1 {
		return errors.New("请输入表名")
	}
	filename := args[0]
	if err := newMigration(filename); err != nil {
		return err
	}
	fmt.Printf("新建迁移文件 %s 成功。\n", filename)
	return nil
}

// 将字符串转换为驼峰命名法
func toCamelCase(s string) string {
	words := strings.Split(s, "_")
	var result string
	for _, word := range words {
		result += cases.Title(language.Und, cases.NoLower).String(word)
	}
	return result
}

// 初始化文件路径
func initFile() string {
	filePath := filepath.Join(Config.MigratePath, "init.go")

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 创建文件并写入初始化内容
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("创建文件失败: %v\n", err)
			return filePath
		}
		defer file.Close()

		// 写入初始内容，比如 package 声明和基础结构
		initialContent := `package migrations

import "github.com/abelpx/go-migrate"

func init() {
}
`
		if _, err := file.WriteString(initialContent); err != nil {
			fmt.Printf("写入初始内容失败: %v\n", err)
		}
	}

	return filePath
}

// 临时文件路径
func tempFile() string {
	return filepath.Join(Config.MigratePath, ".tmp_init.go")
}

// 创建新的迁移文件
func newMigration(filename string) error {
	data, initStr, err := generateMigrationData(filename)
	if err != nil {
		return err
	}

	filename = fmt.Sprintf("%s_%s.go", filename, time.Now().Format("20060102150405"))
	filePath := filepath.Join(Config.MigratePath, filename)
	if err := os.WriteFile(filePath, data, 0o600); err != nil {
		return fmt.Errorf("无法创建迁移模板文件: \n%v", err)
	}

	return updateInitFile(initStr)
}

// 生成迁移数据
func generateMigrationData(filename string) ([]byte, string, error) {
	var table string
	var initStr string
	nowTimeStr := time.Now().Format("20060102150405")
	var data []byte

	if strings.HasSuffix(filename, "_table") {
		if strings.HasPrefix(filename, "create_") {
			table = filename[7 : len(filename)-6]
			data = []byte(fmt.Sprintf(template.CreateTemplate, Config.DbType, toCamelCase(table+"_Table"+nowTimeStr), table))
			initStr = fmt.Sprintf(`	go_migrate.Migrations = append(go_migrate.Migrations, %s())`, "Create"+toCamelCase(table+"_Table"+nowTimeStr))
		} else if index := strings.Index(filename, "to_"); index != -1 {
			table = filename[index+3 : len(filename)-6]
			data = []byte(fmt.Sprintf(template.AlterTemplate, Config.DbType, toCamelCase(filename+nowTimeStr), table))
			initStr = fmt.Sprintf(`	go_migrate.Migrations = append(go_migrate.Migrations, %s())`, "Create"+toCamelCase(filename+nowTimeStr))
		}
	}

	// 未匹配到 create 或 alter
	if len(data) == 0 {
		data = []byte(fmt.Sprintf(template.NewTemplate, Config.DbType, toCamelCase(filename+nowTimeStr)))
		initStr = fmt.Sprintf(`	go_migrate.Migrations = append(go_migrate.Migrations, %s())`, "Create"+toCamelCase(filename+nowTimeStr))
	}

	return data, initStr, nil
}

// 更新初始化文件
func updateInitFile(initStr string) error {
	filePath := initFile()
	tmpFilePath := tempFile()

	// 打开初始化文件和临时文件
	file, err := os.OpenFile(filePath, os.O_RDWR, 0o644)
	if err != nil {
		return fmt.Errorf("打开初始化文件失败: %v", err)
	}
	defer file.Close()

	tmpFile, err := os.Create(tmpFilePath)
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %v", err)
	}
	defer tmpFile.Close()

	// 读取初始化文件内容并修改
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "}") {
			tmpFile.WriteString(initStr + "\n")
		}
		tmpFile.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取初始化文件内容失败: %v", err)
	}

	// 移动临时文件覆盖初始化文件
	if err := os.Rename(tmpFilePath, filePath); err != nil {
		return fmt.Errorf("移动临时文件失败: %v", err)
	}

	return nil
}

// 回滚迁移
func rollback() error {
	migrate := Migrate
	migrations, err := migrate.GetMigrations()

	if err != nil {
		return err
	}

	sort.SliceStable(migrations, func(a, b int) bool { return migrations[b].Id < migrations[a].Id }) //nolint
	if len(migrations) > 0 {
		batch := migrations[0].Batch
		for _, m := range migrations {
			if batch != m.Batch {
				break
			}
			for _, v := range Migrations {
				migration := reflect.TypeOf(v).String()
				if migration == m.Migration {
					batch = m.Batch
					if err := v.Down(); err != nil {
						return err
					}
					if err := migrate.DeleteRecord(m.Id); err != nil {
						return err
					}
					_, _ = fmt.Fprintf(os.Stdout, "rollback %s success.\n", migration)
				}
			}
		}
		if batch == 1 {
			return migrate.DropTableIfExists()
		}
	}
	return nil
}

// 往数据库里执行创建表的操作
func run() error {
	return runMysql()
}

// checkDatabase
// @Description: 检查数据库是否存在，不存在则创建 TODO 后续需要优化，兼容其他数据的检查逻辑
func checkDatabase() {
	// 初始化数据库
	sql := fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;`, Config.DbName)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", Config.Username, Config.Password, Config.Host, Config.Port, Config.DbName)
	db, err := sqlx.Connect("mysql", dsn)

	if err != nil && strings.Contains(err.Error(), "Unknown database") {
		// 创建指定数据库
		db, err = sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql",
			Config.Username, Config.Password, Config.Host, Config.Port))
		if _, err = db.Exec(sql); err != nil {
			fmt.Printf("创建数据库失败: %v", err)
			panic(err)
		}
	}
}

func runMysql() error {
	// 迁移Mysql
	migrate := Migrate
	exists, err := migrate.CheckTable()
	if err != nil {
		return err
	}

	if !exists {
		if cErr := migrate.CreateTable(); cErr != nil {
			return cErr
		}
	}

	migrations, err := migrate.GetMigrations()
	if err != nil {
		return err
	}
	batch := 0

	fmt.Println("----------Migrations")
nextMigrate:
	for _, v := range Migrations {
		migration := reflect.TypeOf(v).String()
		for _, m := range migrations {
			batch = m.Batch
			if migration == m.Migration {
				continue nextMigrate
			}
		}
		if err = v.Up(); err != nil {
			var e interfaces.Seeds
			ok := errors.As(err, &e)
			if ok {
				if e.Error() != "" {
					return e.(error)
				}
			} else {
				return err
			}
		}
		if err := migrate.WriteRecord(migration, batch+1); err != nil {
			return err
		}
		_, _ = fmt.Fprintf(os.Stdout, "migrate %s success.\n", migration)
	}

	return nil
}
