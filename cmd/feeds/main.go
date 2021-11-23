package main

import (
	"bytes"
	"github.com/lycblank/feeds/internal/feeds"
	"github.com/lycblank/feeds/internal/infrastructure/conf"
	"github.com/lycblank/feeds/internal/infrastructure/persistence"
	"github.com/lycblank/feeds/internal/users"
	"go.uber.org/dig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func main() {
	container := dig.New()
	Provide(container, conf.GetConfig)
	Provide(container, NewMysqlDB)
	Provide(container, persistence.NewMysqlFeedRepository)
	Provide(container, persistence.NewMysqlUserRepository)
	Provide(container, feeds.NewFeedTaskService)
	Provide(container, users.NewUserService)

	Invoke(container, Run)
}

func Run(fts *feeds.FeedTaskService) {
	go fts.Run()
	select{}
}


func NewMysqlDB(config *conf.Config) *gorm.DB {
	var buf bytes.Buffer
	buf.WriteString(config.Mysql.UserName)
	buf.WriteString(":")
	buf.WriteString(config.Mysql.Password)
	buf.WriteString("@tcp(")
	buf.WriteString(config.Mysql.Addr)
	buf.WriteString(`)/`)
	buf.WriteString(config.Mysql.Dbname)
	buf.WriteString("?charset=utf8mb4&parseTime=True&loc=Local")
	gdb, err := gorm.Open(mysql.Open(buf.String()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(config.Mysql.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.Mysql.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.Mysql.ConnMaxLifetime) * time.Second)
	if err = sqlDB.Ping(); err != nil {
		panic(err)
	}
	return gdb.Debug()
}

func Provide(container *dig.Container, constructor interface{}, opts ...dig.ProvideOption) {
	if err := container.Provide(constructor, opts...); err != nil {
		panic(err)
	}
}

func Invoke(container *dig.Container, function interface{}, opts ...dig.InvokeOption) {
	if err := container.Invoke(function, opts...); err != nil {
		panic(err)
	}
}


