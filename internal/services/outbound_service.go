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
	InventoryID uint      // 库存ID
	UserID      uint      // 领用人ID
	Quantity    int64     // 数量
	Purpose     string    // 用途
	OpeningDate time.Time // 开封日期
	Remarks     string    // 备注
}

// ApplyOutbound 申请领用
// 创建领出记录，状态设为 PENDING，不扣减库存
//
// 参数:
//   dto: 申请信息
// 返回值:
//   error: 失败返回错误
func (s *OutboundService) ApplyOutbound(dto OutboundApplyDTO) error {
	// 1. 简单校验库存充足性 (非严格，实际扣减在审批时)
	inv, err := s.inventoryDao.GetByID(dto.InventoryID)
	if err != nil {
		return err
	}
	if inv.CurrentQty < dto.Quantity {
		return fmt.Errorf("库存不足，当前剩余: %d", inv.CurrentQty)
	}

	// 2. 创建领出记录 (待审批)
	outboundNo := fmt.Sprintf("LC%s%d", time.Now().Format("20060102"), time.Now().UnixNano()%10000)
	outbound := models.Outbound{
		OutboundNo:     outboundNo,
		InventoryID:    dto.InventoryID,
		UserID:         dto.UserID,
		Quantity:       dto.Quantity,
		Purpose:        dto.Purpose,
		Status:         "USING", // 审批通过后才真正开始使用，但此字段暂保留为USING或可设为WAITING，根据原逻辑保留USING不冲突，主要看ApprovalStatus
		ApprovalStatus: "PENDING",
		OpeningDate:    dto.OpeningDate,
		Remarks:        dto.Remarks,
		SnapExpiryDate: inv.ExpiryDate,
		ApplyDate:      time.Now(),
	}

	return s.outboundDao.Create(&outbound)
}

// AuditOutbound 审批领用
// 管理员审批通过后扣减库存，或驳回申请
//
// 参数:
//   id: 领出记录ID
//   approved: 是否通过
//   approverID: 审批人ID
//   opinion: 审批意见
// 返回值:
//   error: 错误信息
func (s *OutboundService) AuditOutbound(id uint, approved bool, approverID uint, opinion string) error {
	return dao.DB.Transaction(func(tx *gorm.DB) error {
		var out models.Outbound
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&out, id).Error; err != nil {
			return err
		}

		if out.ApprovalStatus != "PENDING" {
			return fmt.Errorf("该申请已被处理，当前状态: %s", out.ApprovalStatus)
		}

		now := time.Now()
		out.ApproverID = &approverID
		out.ApprovalTime = &now
		out.ApprovalOpinion = opinion

		if approved {
			// 1. 审批通过 -> 扣减库存
			var inv models.Inventory
			if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&inv, out.InventoryID).Error; err != nil {
				return err
			}

			if inv.CurrentQty < out.Quantity {
				return fmt.Errorf("库存不足，无法通过审批。当前剩余: %d", inv.CurrentQty)
			}

			inv.CurrentQty -= out.Quantity
			if err := tx.Save(&inv).Error; err != nil {
				return err
			}

			out.ApprovalStatus = "APPROVED"
			// 模拟通知操作员
			fmt.Printf("[Notification] User %d: Your application %s is APPROVED.\n", out.UserID, out.OutboundNo)
		} else {
			// 2. 审批驳回 -> 仅更新状态
			out.ApprovalStatus = "REJECTED"
			// 模拟通知操作员
			fmt.Printf("[Notification] User %d: Your application %s is REJECTED.\n", out.UserID, out.OutboundNo)
		}

		return tx.Save(&out).Error
	})
}

// GetOutboundList 获取领出记录列表
//
// 参数:
//   page, pageSize: 分页
//   userID: 用户ID (可选, 0表示所有用户)
// 返回值:
//   []models.Outbound: 列表
//   int64: 总数
//   error: 错误
func (s *OutboundService) GetOutboundList(page, pageSize int, userID uint) ([]models.Outbound, int64, error) {
	// 复用List，approvalStatus传空表示查询所有状态
	return s.outboundDao.List(page, pageSize, userID, "")
}

// GetAuditList 获取审批列表
//
// 参数:
//   page, pageSize: 分页
//   approvalStatus: 审批状态 (PENDING/APPROVED/REJECTED，空表示所有)
// 返回值:
//   []models.Outbound: 列表
//   int64: 总数
//   error: 错误
func (s *OutboundService) GetAuditList(page, pageSize int, approvalStatus string) ([]models.Outbound, int64, error) {
	// userID=0 表示管理员查询所有人的申请
	return s.outboundDao.List(page, pageSize, 0, approvalStatus)
}

// GetAllOutboundList 获取所有已审批通过的领用记录列表(无权限过滤)
//
// 参数:
//   page, pageSize: 分页
// 返回值:
//   []models.Outbound: 列表
//   int64: 总数
//   error: 错误
func (s *OutboundService) GetAllOutboundList(page, pageSize int) ([]models.Outbound, int64, error) {
	// userID=0, approvalStatus="APPROVED" 表示查询所有已通过审批的记录
	return s.outboundDao.List(page, pageSize, 0, "APPROVED")
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
