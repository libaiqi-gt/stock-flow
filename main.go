package main

import (
	"fmt"
	"stock-flow/internal/config"
	"stock-flow/internal/dao"
	"stock-flow/internal/models"
	"stock-flow/internal/routers"

	_ "stock-flow/docs" // for swagger
)

// @title 耗材管理系统 API
// @version 1.0
// @description LIMS-Consumable Backend API
// @host localhost:8080
// @BasePath /
//
// main 程序入口函数
// 初始化配置、数据库、路由并启动 HTTP 服务
func main() {
	// 1. 初始化配置
	// 从 config.yaml 加载应用配置
	if err := config.InitConfig(); err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// 2. 初始化数据库
	// 建立 MySQL 连接并设置连接池
	dao.InitDB()

	// 3. 自动迁移 (可选，仅开发环境)
	// 自动创建或更新数据库表结构
	if config.AppConfig.Database.AutoMigrate {
		dao.DB.AutoMigrate(&models.User{}, &models.Material{}, &models.Inventory{}, &models.Outbound{})
	}

	// 4. 初始化路由
	// 注册 Gin 路由和中间件
	r := routers.InitRouter()

	// 5. 启动服务
	// 监听指定端口
	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	r.Run(addr)
}
