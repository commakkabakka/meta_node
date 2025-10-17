package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBManager struct {
	DB    *gorm.DB
	SqlDB *sql.DB
}

// 单例模式
var Ins *DBManager
var once sync.Once

func InitGlobalDB() error {
	var err error
	once.Do(func() {
		Ins = &DBManager{}
		err = Ins.Init()
	})
	return err
}

func GetDBManager() *DBManager {
	if Ins == nil {
		log.Fatal("DBManager is not initialized. Call InitGlobalDB first.")
	}
	return Ins
}

func (dm *DBManager) Init() error {
	username := "hello"
	password := "Hello12345."
	ip := "localhost"
	port := 3306
	dbname := "hellodb"

	// 连接 Mysql ：获取 *gorm.DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, ip, port, dbname)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
	if err != nil {
		return errors.New(err.Error())
	}
	dm.DB = db

	// 配置数据库连接池 ：获取 *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return errors.New(err.Error())
	}
	dm.SqlDB = sqlDB

	sqlDB.SetMaxIdleConns(10)                  // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)                 // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour)        // 连接最大存活时间
	sqlDB.SetConnMaxIdleTime(30 * time.Minute) // 连接最大空闲时间

	return nil
}

func (dm *DBManager) CreateTable() error {
	err := dm.DB.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
