package models

import "time"

// Material 耗材基础信息模型
// 对应数据库表 wms_materials，存储耗材的静态属性
type Material struct {
	ID               uint      `gorm:"primaryKey" json:"id"`                              // 主键ID
	Code             string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"` // 物料编码(唯一标识)
	Name             string    `gorm:"type:varchar(100);not null" json:"name"`            // 物料名称
	Category         string    `gorm:"type:varchar(50)" json:"category"`                  // 物料类型(如: 试剂、耗材)
	Spec             string    `gorm:"type:varchar(50)" json:"spec"`                      // 规格型号
	Unit             string    `gorm:"type:varchar(20)" json:"unit"`                      // 计量单位
	Brand            string    `gorm:"type:varchar(50)" json:"brand"`                     // 厂家/品牌
	SafetyStock      int64     `gorm:"type:bigint;default:10" json:"safety_stock"`        // 安全库存
	OpenedExpiryDays int       `gorm:"type:int;default:180" json:"opened_expiry_days"`    // 开封后有效期(天)
	ExpiryAlertDays  int       `gorm:"type:int;default:60" json:"expiry_alert_days"`      // 有效期预警天数
	IsDeleted        bool      `gorm:"default:false;index" json:"is_deleted"`             // 软删除标记
	DeletedAt        *time.Time `json:"deleted_at"`                                       // 删除时间
	CreatedAt        time.Time `json:"created_at"`                                        // 创建时间
	UpdatedAt        time.Time `json:"updated_at"`                                        // 更新时间
}

// TableName 指定表名
// 返回值:
//
//	string: 数据库表名 "wms_materials"
func (Material) TableName() string {
	return "wms_materials"
}
