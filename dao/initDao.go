package dao

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 访问数据库用的
var DB *gorm.DB

// 初始化连接数据库
func Init() error {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,  // 慢 SQL 阈值
			LogLevel:      logger.Error, // Log level
			Colorful:      true,         // 彩色打印
		},
	)
	var err error
	dsn := "root:openGauss@123@tcp(123.249.111.225:3306)/gorm_class?charset=utf8mb4&parseTime=True&loc=Local"
	//想要正确的处理time.Time,需要带上 parseTime 参数，
	//要支持完整的UTF-8编码，需要将 charset=utf8 更改为 charset=utf8mb4
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 newLogger,
		SkipDefaultTransaction: false,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "gov_", // table name prefix, table for `User` would be `t_users`
			SingularTable: true,   // use singular table name, table for `User` would be `user` with this option enable
		},
		DisableForeignKeyConstraintWhenMigrating: true, //逻辑外键
	})
	if err != nil {
		return err
	}
	return err
	sqlDB, err := DB.DB()
	if err != nil {
		log.Println("dao.Init Second", err)
		return err
	}
	sqlDB.SetMaxIdleConns(10)           //最大的空闲连接数
	sqlDB.SetMaxOpenConns(100)          //连接池最大连接数量
	sqlDB.SetConnMaxLifetime(time.Hour) //连接池中链接的最大可复用时间
	return nil
}
