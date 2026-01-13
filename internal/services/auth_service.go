package services

import (
	"errors"
	"stock-flow/internal/dao"
	"stock-flow/internal/models"
	"stock-flow/internal/pkg/utils"
)

// AuthService 认证服务
// 处理用户注册、登录等认证逻辑
type AuthService struct {
	userDao dao.UserDao
}

// Login 用户登录
//
// 参数:
//   username: 用户名
//   password: 密码(明文)
// 返回值:
//   string: JWT Token
//   *models.User: 用户信息
//   error: 登录失败返回错误
func (s *AuthService) Login(username, password string) (string, *models.User, error) {
	// 1. 查询用户
	user, err := s.userDao.GetByUsername(username)
	if err != nil {
		return "", nil, errors.New("用户不存在")
	}

	// 2. 校验密码
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", nil, errors.New("密码错误")
	}

	// 3. 校验状态
	if user.Status == 0 {
		return "", nil, errors.New("账号已禁用")
	}

	// 4. 生成 Token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// Register 用户注册(仅供内部调用或管理员使用)
//
// 参数:
//   username: 用户名
//   password: 密码
//   realName: 真实姓名
//   role: 角色
// 返回值:
//   error: 注册失败返回错误
func (s *AuthService) Register(username, password, realName, role string) error {
	// Check if user exists
	_, err := s.userDao.GetByUsername(username)
	if err == nil {
		return errors.New("用户名已存在")
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := &models.User{
		Username:     username,
		PasswordHash: hash,
		RealName:     realName,
		Role:         role,
		Status:       1,
	}

	return s.userDao.Create(user)
}
