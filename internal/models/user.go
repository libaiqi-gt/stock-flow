package models

import (
	"time"
)

// User 用户模型
// 对应数据库表 sys_users，存储用户账号、角色及状态信息
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`                    // 用户ID
	Username     string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"` // 用户名(唯一)
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`     // 密码哈希值(不返回给前端)
	RealName     string    `gorm:"type:varchar(50)" json:"real_name"`       // 真实姓名
	Role         string    `gorm:"type:varchar(20);not null" json:"role"`   // 角色: Admin, Keeper, User
	Status       int       `gorm:"type:tinyint;default:1" json:"status"`    // 状态: 1正常, 0禁用
	CreatedAt    time.Time `json:"created_at"`                              // 创建时间
	UpdatedAt    time.Time `json:"updated_at"`                              // 更新时间
}

// TableName 指定表名
// 返回值:
//   string: 数据库表名 "sys_users"
func (User) TableName() string {
	return "sys_users"
}
