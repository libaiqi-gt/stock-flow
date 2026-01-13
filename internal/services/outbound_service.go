package services

import (
	"fmt"
	"stock-flow/internal/dao"
	"stock-flow/internal/models"
	"time"

	"gorm.io/gorm"
)

// OutboundService 领出业务服务
// 处理领用申请、记录查询及状态更新
type OutboundService struct {
	outboundDao  dao.OutboundDao
	inventoryDao dao.InventoryDao
}

// OutboundApplyDTO 领用申请数据传输对象
type OutboundApplyDTO struct {
	InventoryID uint   // 库存ID
	UserID      uint   // 领用人ID
	Quantity    int64  // 数量
	Purpose     string // 用途
}

// ApplyOutbound 申请领用
// 使用数据库事务确保库存扣减的原子性
//
// 参数:
//   dto: 申请信息
// 返回值:
//   error: 失败返回错误(如库存不足)
func (s *OutboundService) ApplyOutbound(dto OutboundApplyDTO) error {
	return dao.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 锁定库存记录 (悲观锁)
		// 防止并发场景下的超卖
		var inv models.Inventory
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&inv, dto.InventoryID).Error; err != nil {
			return err
		}

		// 2. 校验库存充足
		if inv.CurrentQty < dto.Quantity {
			return fmt.Errorf("库存不足，当前剩余: %d", inv.CurrentQty)
		}

		// 3. 扣减库存
		inv.CurrentQty -= dto.Quantity
		if err := tx.Save(&inv).Error; err != nil {
			return err
		}

		// 4. 创建领出记录
		outboundNo := fmt.Sprintf("LC%s%d", time.Now().Format("20060102"), time.Now().UnixNano()%10000)
		outbound := models.Outbound{
			OutboundNo:     outboundNo,
			InventoryID:    dto.InventoryID,
			UserID:         dto.UserID,
			Quantity:       dto.Quantity,
			Purpose:        dto.Purpose,
			Status:         "USING",
			SnapExpiryDate: inv.ExpiryDate,
			ApplyDate:      time.Now(),
		}

		if err := tx.Create(&outbound).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetOutboundList 获取领出记录列表
//
// 参数:
//   page, pageSize: 分页
//   userID: 用户ID (可选)
// 返回值:
//   []models.Outbound: 列表
//   int64: 总数
//   error: 错误
func (s *OutboundService) GetOutboundList(page, pageSize int, userID uint) ([]models.Outbound, int64, error) {
	return s.outboundDao.List(page, pageSize, userID)
}

// UpdateStatus 更新领用状态
//
// 参数:
//   id: 记录ID
//   status: 新状态
// 返回值:
//   error: 错误
func (s *OutboundService) UpdateStatus(id uint, status string) error {
	return s.outboundDao.UpdateStatus(id, status)
}
