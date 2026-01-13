# LIMS-Consumable (耗材管理系统) 后端

基于 Go + Gin + GORM + MySQL 构建的耗材全生命周期管理系统后端。

## 1. 项目简介

本系统旨在解决医疗/实验室耗材管理中的痛点，如库存不准、效期监管困难、追溯成本高等。核心功能包括：
- **精准库存**：支持批量入库，实时扣减。
- **效期强控**：红绿灯预警机制，FEFO 推荐策略。
- **全流程闭环**：入库 -> 领用 -> 消耗/归还全链路记录。

## 2. 技术栈

- **语言**: Go 1.20+
- **Web框架**: Gin
- **ORM**: GORM v2
- **数据库**: MySQL 8.0
- **配置管理**: Viper
- **文档**: Swagger (Swaggo)
- **鉴权**: JWT

## 3. 快速开始

### 3.1 前置要求

- Go 1.20+
- MySQL 8.0+

### 3.2 数据库初始化

1. 创建数据库 `stock_flow`。
2. 执行 `docs/schema.sql` 脚本建表。
3. 修改 `configs/config.yaml` 中的数据库连接信息。

### 3.3 运行项目

```bash
# 下载依赖
go mod tidy

# 运行服务 (入口文件已移动至根目录)
go run main.go
```

服务默认启动在 `http://localhost:8080`。

### 3.4 API 文档

启动服务后，访问 Swagger 文档：
http://localhost:8080/swagger/index.html

## 4. 目录结构

```
├── configs             # 配置文件
├── docs                # 文档 (Swagger/SQL)
├── internal
│   ├── config          # 配置加载
│   ├── controllers     # HTTP 控制器
│   ├── dao             # 数据访问层
│   ├── middleware      # Gin 中间件
│   ├── models          # 数据库模型
│   ├── pkg             # 公共工具
│   ├── routers         # 路由配置
│   └── services        # 业务逻辑
├── main.go             # 程序入口
└── go.mod
```

## 5. 核心业务逻辑说明

### 5.1 库存预警
系统根据 `wms_inventory.expiry_date` 计算状态：
- 🔴 **已过期**: `Now > ExpiryDate`
- 🟡 **临期**: `Now + 60d > ExpiryDate`
- 🟢 **正常**: 其他

### 5.2 FEFO 推荐
领用时，系统优先推荐 `ExpiryDate` 最早且 `CurrentQty > 0` 的批次。

### 5.3 事务控制
领用申请 (`/api/v1/outbound/apply`) 采用数据库事务：
1. `SELECT ... FOR UPDATE` 锁定库存记录。
2. 校验库存充足。
3. 扣减库存。
4. 生成领出记录。
5. 提交事务。
