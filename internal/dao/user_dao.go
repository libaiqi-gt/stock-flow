package dao

import (
	"stock-flow/internal/models"
)

// UserDao 用户数据访问对象
// 封装对 sys_users 表的数据库操作
type UserDao struct{}

// Create 创建用户
//
// 参数:
//   user: 包含用户信息的模型指针
// 返回值:
//   error: 成功返回 nil，失败返回错误信息
func (d *UserDao) Create(user *models.User) error {
	return DB.Create(user).Error
}

// GetByUsername 根据用户名查询用户
//
// 参数:
//   username: 用户名
// 返回值:
//   *models.User: 用户模型指针
//   error: 查询失败返回错误
func (d *UserDao) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := DB.Where("is_deleted = ? AND username = ?", false, username).First(&user).Error
	return &user, err
}

// GetByID 根据ID查询用户
//
// 参数:
//   id: 用户ID
// 返回值:
//   *models.User: 用户模型指针
//   error: 查询失败返回错误
func (d *UserDao) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := DB.Where("is_deleted = ?", false).First(&user, id).Error
	return &user, err
}
