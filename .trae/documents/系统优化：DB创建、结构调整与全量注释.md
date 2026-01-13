# 系统优化与完善计划

本计划旨在响应用户的新需求，包括数据库创建、项目结构调整以及代码注释规范化。

## 1. 数据库创建

* **执行方式**: 将尝试使用 `mysql` 命令行工具连接数据库并执行 `docs/schema.sql` 脚本。

* **命令**: `mysql -u root -proot < docs/schema.sql` (假设 root 密码为 root，如 config.yaml 所示)。

* **备选方案**: 如果连接失败，将在回答中明确提示用户执行该 SQL 脚本。

* **验证**: `schema.sql` 已包含 `CREATE DATABASE IF NOT EXISTS stock_flow DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;`，符合用户要求。

## 2. 调整项目结构

* **移动文件**: 将 `cmd/server/main.go` 移动到项目根目录 `D:\GO\stock-flow\main.go`。

* **清理**: 删除空的 `cmd/server` 目录。

* **代码修正**:

  * 更新 `main.go` 中的 Swagger `@title` 等元数据。

  * 检查 `main.go` 的 import 路径（通常不受影响，因为 module path 未变）。

* **文档更新**:

  * 更新 `README.md` 中的运行命令为 `go run main.go`。

  * 更新 Swagger 生成命令为 `swag init` (默认查找当前目录 main.go)。

## 3. 代码注释规范化

将遍历项目核心目录 (`internal/`, `main.go`)，为所有导出的结构体、函数、接口添加中文注释。

* **格式标准**:

  * **函数/方法**:

    ```go
    // FunctionName 功能简述
    //
    // 参数:
    //   param1: 参数说明
    // 返回值:
    //   return1: 返回值说明
    func FunctionName(...) ...
    ```

  * **结构体**:

    ```go
    // StructName 结构体说明
    type StructName struct {
        Field string // 字段说明
    }
    ```

  * **行内注释**: 在复杂逻辑（如库存扣减事务、Token 解析）处添加 `//` 注释。

## 4. 执行顺序

1. **移动文件**: 先调整结构，确保后续操作基于新结构。
2. **数据库初始化**: 尝试执行 SQL 脚本。
3. **添加注释**: 批量修改代码文件。
4. **生成文档**: 重新运行 `swag init` 更新 Swagger 文档。
5. **验证**: 运行测试用例确保未破坏功能。

