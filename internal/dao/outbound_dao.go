package dao

import (
	"stock-flow/internal/models"
)

// OutboundDao 领出记录数据访问对象
// 封装对 wms_outbound 表的数据库操作
type OutboundDao struct{}

// Create 创建领出记录
//
// 参数:
//   out: 领出记录模型
// 返回值:
//   error: 错误信息
func (d *OutboundDao) Create(out *models.Outbound) error {
	return DB.Create(out).Error
}

// List 分页查询领出记录
//
// 参数:
//   page: 页码
//   pageSize: 每页数量
//   userID: 用户ID (0表示查询所有)
//   approvalStatus: 审批状态 (空字符串表示查询所有)
// 返回值:
//   []models.Outbound: 记录列表
//   int64: 总数
//   error: 错误信息
func (d *OutboundDao) List(page, pageSize int, userID uint, approvalStatus string) ([]models.Outbound, int64, error) {
	var list []models.Outbound
	var total int64

	db := DB.Model(&models.Outbound{}).Preload("Inventory.Material").Preload("User").Preload("Approver")
	
	if userID > 0 {
		db = db.Where("user_id = ?", userID)
	}

	if approvalStatus != "" {
		db = db.Where("approval_status = ?", approvalStatus)
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&list).Error
	return list, total, err
}

// UpdateStatus 更新领出记录状态
//
// 参数:
//   id: 记录ID
//   status: 新状态 (使用中/已用完)
// 返回值:
//   error: 错误信息
func (d *OutboundDao) UpdateStatus(id uint, status string) error {
	return DB.Model(&models.Outbound{}).Where("id = ?", id).Update("status", status).Error
}
