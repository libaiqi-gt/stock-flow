package dao

import (
	"fmt"
	"log"
	"os"
	"stock-flow/internal/config"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库连接对象
var DB *gorm.DB

// InitDB 初始化数据库连接
// 建立与 MySQL 的连接，设置连接池参数，并根据配置模式设置日志级别
func InitDB() {
	c := config.AppConfig.Database

	// 优先从环境变量获取配置 (Docker部署)
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = c.Host
	}

	portStr := os.Getenv("DB_PORT")
	port := c.Port
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = c.User
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = c.Password
	}

	name := os.Getenv("DB_NAME")
	if name == "" {
		name = c.Name
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, name)

	var err error

	// 设置 GORM 日志级别
	// release 模式下仅打印 Error，debug 模式下打印 Info (包含 SQL 语句)
	logLevel := logger.Info
	if config.AppConfig.Server.Mode == "release" {
		logLevel = logger.Error
	}

	// 建立数据库连接
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 获取底层 sql.DB 对象以设置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get generic database object: %v", err)
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大存活时间

	log.Println("Database connection established successfully")
}
