package config

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化MySQL连接（GORM）
func InitDB() {
	// 替换为实际的MySQL连接信息：用户名:密码@tcp(IP:端口)/数据库名?charset=utf8mb4&parseTime=True&loc=Local
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

	// 自定义日志配置，便于调试
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 输出到控制台
		logger.Config{
			SlowThreshold: time.Second, // 慢查询阈值
			LogLevel:      logger.Info, // 日志级别
			Colorful:      true,        // 彩色日志
		},
	)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("数据库连接失败：" + err.Error())
	}

	// 设置连接池
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)  // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100) // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	log.Println("数据库连接成功")
}
