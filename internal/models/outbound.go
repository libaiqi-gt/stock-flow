package models

import "time"

// Outbound 领出记录模型
// 对应数据库表 wms_outbound，记录每一次库存扣减操作
type Outbound struct {
	ID             uint      `gorm:"primaryKey" json:"id"`                              // 主键ID
	OutboundNo     string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"outbound_no"` // 领出单号(系统生成)
	InventoryID    uint      `gorm:"index;not null" json:"inventory_id"`                // 关联库存ID
	Inventory      Inventory `gorm:"foreignKey:InventoryID" json:"inventory"`           // 库存详情
	UserID         uint      `gorm:"index;not null" json:"user_id"`                     // 领用人ID
	User           User      `gorm:"foreignKey:UserID" json:"user"`                     // 领用人详情
	Quantity       int64     `gorm:"not null" json:"quantity"`                          // 领出数量
	Purpose         string    `gorm:"type:varchar(255)" json:"purpose"`                  // 领用用途
	Status          string    `gorm:"type:varchar(20);default:'USING'" json:"status"`    // 状态: USING(使用中), FINISHED(已用完)
	ApprovalStatus  string    `gorm:"type:varchar(20);default:'PENDING'" json:"approval_status"` // 审批状态: PENDING, APPROVED, REJECTED
	ApprovalOpinion string    `gorm:"type:varchar(255)" json:"approval_opinion"`         // 审批意见
	ApproverID      *uint     `gorm:"index" json:"approver_id"`                          // 审批人ID
	Approver        *User     `gorm:"foreignKey:ApproverID" json:"approver"`             // 审批人详情
	ApprovalTime    *time.Time `json:"approval_time"`                                    // 审批时间
	OpeningDate     time.Time `gorm:"type:date" json:"opening_date"`                     // 开封日期
	Remarks         string    `gorm:"type:varchar(500)" json:"remarks"`                  // 备注说明
	SnapExpiryDate  time.Time `gorm:"type:date" json:"snap_expiry_date"`                 // 快照有效期(冗余存储，防源数据变更)
	ApplyDate      time.Time `json:"apply_date"`                                        // 申请时间
	IsDeleted      bool      `gorm:"default:false;index" json:"is_deleted"`             // 软删除标记
	DeletedAt      *time.Time `json:"deleted_at"`                                       // 删除时间
	CreatedAt      time.Time `json:"created_at"`                                        // 创建时间
	UpdatedAt      time.Time `json:"updated_at"`                                        // 更新时间
}

// TableName 指定表名
// 返回值:
//   string: 数据库表名 "wms_outbound"
func (Outbound) TableName() string {
	return "wms_outbound"
}
