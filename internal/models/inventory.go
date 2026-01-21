package models

import "time"

// Inventory 库存批次模型
// 对应数据库表 wms_inventory，存储每个入库批次的详细信息
type Inventory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`                         // 主键ID
	MaterialID uint      `gorm:"index;not null" json:"material_id"`            // 关联耗材ID
	Material   Material  `gorm:"foreignKey:MaterialID" json:"material"`        // 耗材详情(关联查询用)
	BatchNo    string    `gorm:"type:varchar(50);not null" json:"batch_no"`    // 内部批号(管控核心)
	InboundNo  string    `gorm:"type:varchar(50);unique;not null" json:"inbound_no"` // 入库单号(唯一标识)
	InitialQty int64     `gorm:"not null" json:"initial_qty"`                  // 初始入库数量
	CurrentQty int64     `gorm:"not null" json:"current_qty"`                  // 当前剩余数量(动态变化)
	ExpiryDate time.Time `gorm:"type:date;index" json:"expiry_date"`           // 有效期(用于效期预警)
	IsDeleted  bool      `gorm:"default:false;index" json:"is_deleted"`        // 软删除标记
	DeletedAt  *time.Time `json:"deleted_at"`                                  // 删除时间
	CreatedAt  time.Time `json:"created_at"`                                   // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`                                   // 更新时间
}

// TableName 指定表名
// 返回值:
//   string: 数据库表名 "wms_inventory"
func (Inventory) TableName() string {
	return "wms_inventory"
}
