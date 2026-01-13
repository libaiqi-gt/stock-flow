# 耗材管理系统 (LIMS-Consumable) 后端实施计划

本计划旨在基于 `product.md` 构建高性能、规范化的 Go 后端服务。

## 1. 项目初始化与架构搭建
- **目录结构**: 采用标准的 Go Clean Architecture 分层结构。
  - `cmd/server`: 程序入口
  - `configs`: 配置文件 (config.yaml)
  - `internal`:
    - `models`: 数据库结构体 (GORM)
    - `dao`: 数据访问层
    - `services`: 业务逻辑层
    - `controllers`: HTTP 接口层
    - `routers`: 路由配置
    - `middleware`: 中间件 (JWT, CORS, Logger)
    - `pkg`: 公共工具 (Utils, Error codes)
  - `docs`: Swagger 文档
- **依赖管理**: 初始化 `go.mod`，引入 `gin`, `gorm`, `mysql`, `jwt-go`, `viper`, `swag` 等核心库。

## 2. 数据库设计 (Schema Design)
基于 3NF 设计 MySQL 表结构，确保数据一致性与扩展性：
- **sys_users**: 用户表 (RBAC 基础)
- **wms_materials**: 耗材基础信息表 (物料编码、名称、规格、单位、品牌) - *从 Inventory 拆分以满足 3NF*
- **wms_inventory**: 库存批次表 (关联 Material，包含批号、有效期、入库/当前数量)
- **wms_outbound**: 领出记录表 (关联 Inventory 和 User，包含用途、状态、快照有效期)

## 3. 核心功能开发
### 3.1 基础模块
- **配置管理**: 使用 Viper 读取 `config.yaml`。
- **数据库连接**: GORM 初始化及连接池配置。
- **全局日志**: 统一日志格式。
- **Swagger**: 集成 Swagger 自动生成 API 文档。

### 3.2 业务模块
- **用户认证 (Auth)**:
  - 登录接口 (JWT 生成)
  - 权限中间件 (Admin/Keeper/User 角色控制)
- **耗材管理 (Material)**:
  - 基础信息的增删改查。
- **库存管理 (Inventory)**:
  - **批量入库**: Excel 数据解析与批次去重/累加逻辑。
  - **库存查询**: 支持模糊搜索、效期筛选。
  - **效期预警**: 实现红绿灯逻辑 (Expired/Warning/Normal)。
- **领用管理 (Outbound)**:
  - **智能推荐**: 实现 FEFO (先失效先出) 算法，优先推荐近效期批次。
  - **领用申请**: 事务控制，确保高并发下库存不超卖。
  - **状态反馈**: 领用后状态流转 (使用中 -> 已用完)。

## 4. 测试与文档
- **单元测试**: 针对 Service 层编写测试用例 (Go Test)。
- **接口文档**: 使用 Swagger 注解生成在线文档。
- **部署文档**: 编写 `README.md` 和数据库 SQL 脚本。

## 5. 质量保证
- **代码规范**: 遵循 Go Code Review Comments 规范。
- **错误处理**: 统一错误码与 HTTP 响应结构。

请确认以上实施计划，确认后将立即开始代码编写。