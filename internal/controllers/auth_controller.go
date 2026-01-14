package controllers

import (
	"stock-flow/internal/pkg/response"
	"stock-flow/internal/services"

	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
// 处理登录、注册等 HTTP 请求
type AuthController struct {
	authService services.AuthService
}

// LoginReq 登录请求参数
type LoginReq struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

// LoginResp 登录响应数据
type LoginResp struct {
	Token    string `json:"token"`    // JWT Token
	Username string `json:"username"` // 用户名
	Role     string `json:"role"`     // 角色
}

// RegisterReq 注册请求参数
type RegisterReq struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
	RealName string `json:"real_name"`                   // 真实姓名
	Role     string `json:"role"`                        // 角色: 管理员, 操作员, 普通用户
}

// Login
// @Summary 用户登录
// @Description 用户通过账号密码登录，获取 JWT Token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginReq true "登录参数"
// @Success 200 {object} response.Response "登录成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "认证失败"
// @Router /auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	token, user, err := ctrl.authService.Login(req.Username, req.Password)
	if err != nil {
		response.Error(c, response.CodeUnauthorized, err.Error())
		return
	}

	response.Success(c, LoginResp{
		Token:    token,
		Username: user.Username,
		Role:     user.Role,
	})
}

// Register
// @Summary 用户注册
// @Description 开放注册接口，默认角色为 普通用户
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterReq true "注册参数"
// @Success 200 {object} response.Response "注册成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /auth/register [post]
func (ctrl *AuthController) Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	// 角色映射: 中文 -> 内部标识
	roleMap := map[string]string{
		"管理员":  "Admin",
		"操作员":  "Keeper",
		"普通用户": "User",
		"":     "User", // 默认值
	}

	internalRole, ok := roleMap[req.Role]
	if !ok {
		response.Error(c, response.CodeBadRequest, "角色参数无效，可选值: 管理员, 操作员, 普通用户")
		return
	}

	if err := ctrl.authService.Register(req.Username, req.Password, req.RealName, internalRole); err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	response.Success[any](c, nil)
}
