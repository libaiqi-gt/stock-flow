package dao

import (
	"stock-flow/internal/models"
	"time"

	"gorm.io/gorm"
)

// InventoryDao 库存数据访问对象
// 封装对 wms_inventory 表的数据库操作
type InventoryDao struct{}

// Create 创建库存批次
//
// 参数:
//
//	inv: 库存模型
//
// 返回值:
//
//	error: 错误信息
func (d *InventoryDao) Create(inv *models.Inventory) error {
	return DB.Create(inv).Error
}

// Delete 删除库存 (软删除)
//
// 参数:
//
//	id: 库存ID
//
// 返回值:
//
//	error: 错误信息
func (d *InventoryDao) Delete(id uint) error {
	// 开启事务
	return DB.Transaction(func(tx *gorm.DB) error {
		// 1. 查询关联的待审批申请
		var pendingOutbounds []models.Outbound
		if err := tx.Where("inventory_id = ? AND approval_status = ? AND is_deleted = ?", id, "PENDING", false).Find(&pendingOutbounds).Error; err != nil {
			return err
		}

		// 2. 如果存在待审批申请，先删除它们
		if len(pendingOutbounds) > 0 {
			if err := tx.Model(&models.Outbound{}).
				Where("inventory_id = ? AND approval_status = ? AND is_deleted = ?", id, "PENDING", false).
				Updates(map[string]interface{}{
					"is_deleted": true,
					"deleted_at": time.Now(),
				}).Error; err != nil {
				return err
			}
		}

		// 3. 删除库存 (软删除)
		if err := tx.Model(&models.Inventory{}).
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"is_deleted": true,
				"deleted_at": time.Now(),
			}).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetByMaterialAndBatch 根据物料ID和批号查询库存
//
// 参数:
//
//	materialID: 物料ID
//	batchNo: 批号
//
// 返回值:
//
//	*models.Inventory: 库存模型
//	error: 错误信息
func (d *InventoryDao) GetByMaterialAndBatch(materialID uint, batchNo string) (*models.Inventory, error) {
	var inv models.Inventory
	err := DB.Where("is_deleted = ? AND material_id = ? AND batch_no = ?", false, materialID, batchNo).First(&inv).Error
	return &inv, err
}

// GetByInboundNo 根据入库单号查询库存
//
// 参数:
//
//	inboundNo: 入库单号
//
// 返回值:
//
//	*models.Inventory: 库存模型
//	error: 错误信息
func (d *InventoryDao) GetByInboundNo(inboundNo string) (*models.Inventory, error) {
	var inv models.Inventory
	err := DB.Where("is_deleted = ? AND inbound_no = ?", false, inboundNo).First(&inv).Error
	return &inv, err
}

// Update 更新库存信息
//
// 参数:
//
//	inv: 包含更新后信息的库存模型
//
// 返回值:
//
//	error: 错误信息
func (d *InventoryDao) Update(inv *models.Inventory) error {
	return DB.Save(inv).Error
}

// List 综合查询库存列表
//
// 参数:
//
//	page, pageSize: 分页参数
//	materialName: 物料名称(模糊)
//	code: 物料编码(精确)
//	batchNo: 批号(精确)
//	status: 状态(0全部, 1正常, 2临期, 3过期)
//
// 返回值:
//
//	[]models.Inventory: 库存列表
//	int64: 总数
//	error: 错误信息
func (d *InventoryDao) List(page, pageSize int, materialName, code, batchNo string, status int) ([]models.Inventory, int64, error) {
	// status: 0 all, 1 normal, 2 warning, 3 expired
	var list []models.Inventory
	var total int64

	// Explicitly specify table alias for is_deleted to avoid ambiguity when joining
	db := DB.Model(&models.Inventory{}).Where("wms_inventory.is_deleted = ?", false).Preload("Material")

	if materialName != "" {
		db = db.Joins("JOIN wms_materials ON wms_materials.id = wms_inventory.material_id").
			Where("wms_materials.name LIKE ?", "%"+materialName+"%")
	}
	if code != "" {
		if materialName == "" { // Avoid joining twice if possible, or just use association query
			db = db.Joins("JOIN wms_materials ON wms_materials.id = wms_inventory.material_id")
		}
		db = db.Where("wms_materials.code = ?", code)
	}
	if batchNo != "" {
		db = db.Where("batch_no = ?", batchNo)
	}

	now := time.Now()
	if status > 0 {
		switch status {
		case 2: // Warning: < 60 days
			warningDate := now.AddDate(0, 0, 60)
			db = db.Where("expiry_date > ? AND expiry_date <= ?", now, warningDate)
		case 3: // Expired
			db = db.Where("expiry_date <= ?", now)
		case 1: // Normal
			warningDate := now.AddDate(0, 0, 60)
			db = db.Where("expiry_date > ?", warningDate)
		}
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

// GetAvailableBatches 获取可用库存批次(FEFO策略)
//
// 参数:
//
//	materialID: 物料ID
//
// 返回值:
//
//	[]models.Inventory: 按有效期升序排列的可用库存列表
//	error: 错误信息
func (d *InventoryDao) GetAvailableBatches(materialID uint) ([]models.Inventory, error) {
	var list []models.Inventory
	// FEFO: Order by ExpiryDate ASC
	err := DB.Where("is_deleted = ? AND material_id = ? AND current_qty > 0", false, materialID).
		Order("expiry_date ASC").
		Find(&list).Error
	return list, err
}

// GetByID 根据ID获取库存详情
//
// 参数:
//
//	id: 库存ID
//
// 返回值:
//
//	*models.Inventory: 库存模型
//	error: 错误信息
func (d *InventoryDao) GetByID(id uint) (*models.Inventory, error) {
	var inv models.Inventory
	err := DB.Preload("Material").Where("is_deleted = ?", false).First(&inv, id).Error
	return &inv, err
}
