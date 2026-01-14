package dao

import (
	"time"
)

type StatisticsDao struct{}

type WarningBatch struct {
	ID              uint      `json:"id"`
	BatchNo         string    `json:"batch_no"`
	MaterialName    string    `json:"material_name"`
	ExpiryAlertDays int       `json:"expiry_alert_days"`
	ExpiryDate      time.Time `json:"expiry_date"`
}

type MonthlyOutbound struct {
	Month    string `json:"month"`
	TotalQty int64  `json:"total_qty"`
}

// CountTotalBatches 统计当前库存总批次数量 (current_qty > 0)
func (d *StatisticsDao) CountTotalBatches() (int64, error) {
	var count int64
	err := DB.Table("wms_inventory").Where("current_qty > 0").Count(&count).Error
	return count, err
}

// GetWarningBatches 获取临期预警库存批次
// 逻辑: expiry_date <= NOW + alert_days AND expiry_date > NOW
func (d *StatisticsDao) GetWarningBatches() ([]WarningBatch, error) {
	var results []WarningBatch
	// 使用 GORM 的 Join 和 Where 进行复杂查询
	// 注意: 这里的 SQL 语法针对 MySQL 优化
	err := DB.Table("wms_inventory").
		Select("wms_inventory.id, wms_inventory.batch_no, wms_materials.name as material_name, wms_materials.expiry_alert_days, wms_inventory.expiry_date").
		Joins("JOIN wms_materials ON wms_inventory.material_id = wms_materials.id").
		Where("wms_inventory.current_qty > 0").
		Where("wms_inventory.expiry_date <= DATE_ADD(NOW(), INTERVAL wms_materials.expiry_alert_days DAY)").
		Where("wms_inventory.expiry_date > NOW()").
		Scan(&results).Error

	return results, err
}

// CountExpiredBatches 统计已过期库存批次数量
// 逻辑: expiry_date <= NOW
func (d *StatisticsDao) CountExpiredBatches() (int64, error) {
	var count int64
	err := DB.Table("wms_inventory").
		Where("current_qty > 0").
		Where("expiry_date <= NOW()").
		Count(&count).Error
	return count, err
}

// GetOutboundTrend 近半年耗材出库数量统计 (按月分组)
// 逻辑: 过去6个月，approval_status = 'APPROVED'
func (d *StatisticsDao) GetOutboundTrend() ([]MonthlyOutbound, error) {
	var results []MonthlyOutbound
	// 获取6个月前的第一天
	sixMonthsAgo := time.Now().AddDate(0, -6, 0).Format("2006-01-02")

	err := DB.Table("wms_outbound").
		Select("DATE_FORMAT(created_at, '%Y-%m') as month, SUM(quantity) as total_qty").
		Where("created_at >= ?", sixMonthsAgo).
		Where("approval_status = ?", "APPROVED").
		Group("month").
		Order("month ASC").
		Scan(&results).Error

	return results, err
}
