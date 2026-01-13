package response

// 常用错误码定义
const (
	CodeSuccess      = 200
	CodeBadRequest   = 400
	CodeUnauthorized = 401
	CodeForbidden    = 403
	CodeNotFound     = 404
	CodeServerError  = 500
)

// 简单的内存消息映射表
// Code -> Lang -> Message
var msgMap = map[int]map[string]string{
	CodeSuccess: {
		"zh-CN": "操作成功",
		"en-US": "Operation successful",
	},
	CodeBadRequest: {
		"zh-CN": "请求参数错误",
		"en-US": "Invalid request parameters",
	},
	CodeUnauthorized: {
		"zh-CN": "未授权，请登录",
		"en-US": "Unauthorized, please login",
	},
	CodeForbidden: {
		"zh-CN": "权限不足",
		"en-US": "Permission denied",
	},
	CodeNotFound: {
		"zh-CN": "资源未找到",
		"en-US": "Resource not found",
	},
	CodeServerError: {
		"zh-CN": "服务器内部错误",
		"en-US": "Internal server error",
	},
}

// GetMsg 获取国际化消息
//
// 参数:
//   code: 错误码
//   lang: 语言标识 (如 zh-CN, en-US)，支持简单前缀匹配
// 返回值:
//   对应语言的消息，默认为中文
func GetMsg(code int, lang string) string {
	msgs, ok := msgMap[code]
	if !ok {
		return "Unknown Error"
	}

	// 默认中文
	if lang == "" {
		return msgs["zh-CN"]
	}

	// 精确匹配
	if msg, exists := msgs[lang]; exists {
		return msg
	}

	// 简单前缀匹配 (例如 zh-HK -> zh-CN, en-GB -> en-US)
	if len(lang) >= 2 {
		prefix := lang[:2]
		if prefix == "zh" {
			return msgs["zh-CN"]
		}
		if prefix == "en" {
			return msgs["en-US"]
		}
	}

	return msgs["zh-CN"]
}
