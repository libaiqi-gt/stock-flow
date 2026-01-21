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
	fmt.Println("DB_HOST:", host)
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

	// ---------------------------------------------------------
	// 自动创建数据库 (Fix: Error 1049 Unknown database)
	// ---------------------------------------------------------
	createDatabase(user, password, host, port, name)

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

// createDatabase 尝试连接 MySQL 并创建数据库（如果不存在）
func createDatabase(user, password, host string, port int, dbName string) {
	// 连接 DSN 不包含数据库名
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4", user, password, host, port)

	// 使用标准库 sql 或者 gorm 打开连接
	// 这里为了简单直接复用 gorm，虽然有点重，但兼容性好
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("[Warning] Failed to connect to MySQL server to check database existence: %v. Proceeding to connect to DB directly...", err)
		return
	}

	// 获取通用数据库对象以关闭连接
	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	defer sqlDB.Close()

	// 创建数据库 SQL
	createSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", dbName)
	if err := db.Exec(createSQL).Error; err != nil {
		log.Printf("[Warning] Failed to create database '%s': %v", dbName, err)
	} else {
		log.Printf("[Info] Database '%s' checked/created successfully.", dbName)
	}
}
