# API 统一响应结构体设计与实现计划

本计划旨在为项目引入统一的 API 响应格式，规范前后端交互协议。

## 1. 核心设计 (Design)

我们将创建 `internal/pkg/response` 包，核心组件如下：

### 1.1 泛型响应结构体 (`Response[T]`)
满足 JSON 输出规范，包含状态码、友好提示、错误详情、业务数据和时间戳。

```go
type Response[T any] struct {
    Code      int    `json:"code"`             // 业务状态码 (200成功)
    Msg       string `json:"msg"`              // 提示信息
    ErrMsg    string `json:"errMsg,omitempty"` // 错误详情 (仅调试返回)
    Data      T      `json:"data"`             // 业务数据
    Timestamp string `json:"timestamp"`        // ISO8601 时间戳
}
```

### 1.2 构造函数与辅助方法
- `Success[T](c *gin.Context, data T)`: 快速返回成功响应。
- `Error(c *gin.Context, code int, msg string, opts ...Option)`: 返回错误响应，支持自定义错误详情。
- `Result[T](...)`: 底层通用构造器。

### 1.3 国际化支持 (I18n)
- 建立简单的错误码与消息映射表 (`ErrorCodeMap`)。
- 提供 `GetMsg(code int, lang string)` 方法，支持根据 Accept-Language 头自动匹配消息。

## 2. 实现步骤 (Implementation)

1.  **创建包目录**: `internal/pkg/response`。
2.  **实现核心代码 (`response.go`)**:
    - 定义结构体与常量。
    - 实现 `Success`, `Error` 等 Gin 辅助函数。
    - 实现时间戳生成与 JSON 序列化逻辑。
3.  **实现国际化简易版 (`i18n.go`)**:
    - 定义默认的中英文错误消息字典。
    - 实现消息查找逻辑。
4.  **编写单元测试 (`response_test.go`)**:
    - 覆盖成功、失败、带详情、泛型数据等各种场景。
    - 验证 JSON 格式和时间戳。

## 3. 交付物
- `internal/pkg/response/response.go`: 核心实现。
- `internal/pkg/response/response_test.go`: 单元测试。
- `internal/pkg/response/codes.go`: 常用错误码定义 (可选，放入 response.go 或单独文件)。

注意：本次任务主要关注**工具库的实现与测试**，暂不全量替换现有 Controller 中的代码，以免影响现有业务逻辑，后续可逐步替换。