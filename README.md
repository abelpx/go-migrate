# Go 项目用来管理数据迁移的工具

## 项目描述
go-migrate 是一个用来管理数据迁移的工具，支持 mysql、sqlite3、postgres、es 等数据库的迁移。

## 目录

- [安装](#安装)
- [快速开始](#快速开始)
- [使用方法](#使用方法)
- [版本规划及清单](#版本规划及清单)
- [下版本规划](#下版本规划)
- [贡献指南](#贡献指南)
- [许可证](#许可证)

## 安装

go get -u github.com/abelpx/go-migrate

## 快速开始

一. 配置项目执行命令

```go
package main

import (
	"fmt"
	"github.com/abelpx/go-migrate"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	_ "upgrade_platform/migrations"
)

var dbConfig go_migrate.DbConfig

func main() {
	// 初始化根命令
	var rootCmd = &cobra.Command{
		Use:   "db-migrate",
		Short: "数据库迁移工具",
		Long:  `一个用于数据库迁移的 CLI 工具，支持迁移、回滚和创建迁移文件操作`,
	}

	migrationPath, _ := filepath.Abs("../migrations")

	rootCmd.PersistentFlags().StringVar(&dbConfig.MigratePath, "migrate-path", migrationPath, "迁移文件夹路径")
	rootCmd.PersistentFlags().StringVar(&dbConfig.Username, "username", "root", "数据库用户名")
	rootCmd.PersistentFlags().StringVar(&dbConfig.Password, "password", "", "数据库密码")
	rootCmd.PersistentFlags().StringVar(&dbConfig.Host, "host", "127.0.0.1", "数据库主机")
	rootCmd.PersistentFlags().IntVar(&dbConfig.Port, "port", 3306, "数据库端口")
	rootCmd.PersistentFlags().StringVar(&dbConfig.DbName, "dbname", "upgrade_platform", "数据库名称")
	rootCmd.PersistentFlags().StringVar(&dbConfig.DbType, "dbtype", "mysql", "数据库类型 (mysql, sqlite, postgres, ...)")

	go_migrate.StartCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		go_migrate.Config = dbConfig

		fmt.Println("Migrations:", len(go_migrate.Migrations))
	}

	rootCmd.AddCommand(go_migrate.StartCmd)

	// 执行根命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

```

二. 创建迁移文件

```shell
go run main.go migrate run create_users_tables
```

三. 修改迁移文件
进入你的 `migrate-path` 文件夹下，找到对应的表文件，他会生成类似如下的文件
```go
package migrations

import (
	"github.com/abelpx/go-migrate/pkg/interfaces"
	"github.com/abelpx/go-migrate/pkg/lib/mysql"
)

type UsersTable20240808153702 struct{}

func CreateUsersTable20240808153702() interfaces.Migration {
	return &UsersTable20240808153702{}
}

func (t *UsersTable20240808153702) Up() error {
	return mysql.NewSchema().Create("users", func(table interfaces.Foundation) {
		table.Id("id", 22).Index()
		table.Timestamps()
	})
}

func (t *UsersTable20240808153702) Down() error {
	return mysql.NewSchema().DropIfExists("users")
}

```

然后在你的 Up() 函数中添加你的表结构需要的字段
```go
func (t *UsersTable20240808153702) Up() error {
    return mysql.NewSchema().Create("users", func(table interfaces.Foundation) {
        table.Id("id", 22).Index()
        table.String("name", 255).Unique().Comment("用户名")
        table.String("email", 255).Unique().Comment("邮箱")
        table.String("password", 255).Comment("密码")
        table.Integer("age", 11).Comment("年龄")
        table.Timestamps()
    })
}
```
他所支持的字段类型如下：
```go
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
```
TODO 先让我犯个懒，后面我在整理成具体的表格。

--- 
## 版本规划及清单

- [x] 完成 mysql 迁移逻辑
- [ ] 完成 sqlite3 迁移逻辑
- [ ] 完成 postgres 迁移逻辑
- [ ] 完成 es 迁移逻辑
- [ ] 更多数据库支持 ... 更多敬请期待

## 版本清单

| 版本     | 描述                    | 状态  | 更新时间       | 更新人  | 备注                         |
|--------|-----------------------|-----|------------|------|----------------------------|
| v0.0.1 | 初始化mysql 数据库与相关迁移文件管理 | 已完成 | 2021-07-01 | abel | 初步版本，后期计划会进一步优化 mysql 迁移逻辑 |

