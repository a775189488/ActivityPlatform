package db

import (
	"fmt"
	"log"

	"entrytask/internal/conf"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type DbMysql struct {
	Conn *gorm.DB
}

//Connect 初始化数据库配置
func (d *DbMysql) Connect() error {
	var (
		dbType, dbName, user, pwd, host string
	)

	conf := conf.Config.Database
	dbType = conf.Type
	dbName = conf.Name
	user = conf.User
	pwd = conf.Password
	host = conf.Host

	db, err := gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, pwd, host, dbName))
	if err != nil {
		log.Fatal("connecting mysql error: ", err)
		return err
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return defaultTableName
	}

	// todo 以下的配置写到配置文件中
	// 不打印sql语句
	db.LogMode(false)
	db.SingularTable(true)
	// 设置最大最小连接数
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	d.Conn = db

	log.Println("Connect Mysql Success")

	return nil
}

//DB 返回DB
func (d *DbMysql) DB() *gorm.DB {
	return d.Conn
}
