package dao

import (
	"fmt"
	"stock-flow/internal/models"
	"time"
)

// MaterialDao 耗材数据访问对象
// 封装对 wms_materials 表的数据库操作
type MaterialDao struct{}

// Create 创建耗材
//
// 参数:
//
//	m: 耗材模型指针
//
// 返回值:
//
//	error: 错误信息
func (d *MaterialDao) Create(m *models.Material) error {
	return DB.Create(m).Error
}

// Delete 删除耗材 (软删除)
//
// 参数:
//
//	id: 耗材ID
//
// 返回值:
//
//	error: 错误信息
func (d *MaterialDao) Delete(id uint) error {
	// 软删除: 更新 is_deleted = true, deleted_at = now
	return DB.Model(&models.Material{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}

// GetByCode 根据编码查询耗材
//
// 参数:
//
//	code: 物料编码
//
// 返回值:
//
//	*models.Material: 耗材模型
//	error: 错误信息
func (d *MaterialDao) GetByCode(code string) (*models.Material, error) {
	var m models.Material
	err := DB.Where("is_deleted = ? AND code = ?", false, code).First(&m).Error
	return &m, err
}

func (d *MaterialDao) GetByID(id uint) (*models.Material, error) {
	var m models.Material
	err := DB.Where("is_deleted = ? AND id = ?", false, id).First(&m).Error
	return &m, err
}

func (d *MaterialDao) UpdateByID(id uint, updates map[string]interface{}) error {
	tx := DB.Model(&models.Material{}).
		Where("is_deleted = ? AND id = ?", false, id).
		Updates(updates)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return fmt.Errorf("耗材不存在")
	}
	return nil
}

// List 分页查询耗材列表
//
// 参数:
//
//	page: 页码
//	pageSize: 每页数量
//	name: 物料名称(模糊查询)
//
// 返回值:
//
//	[]models.Material: 耗材列表
//	int64: 总数量
//	error: 错误信息
func (d *MaterialDao) List(page, pageSize int, name string) ([]models.Material, int64, error) {
	var materials []models.Material
	var total int64

	db := DB.Model(&models.Material{}).Where("is_deleted = ?", false)
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&materials).Error
	return materials, total, err
}
