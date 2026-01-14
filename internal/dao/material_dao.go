package dao

import (
	"stock-flow/internal/models"
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

// Delete 删除耗材
//
// 参数:
//
//	id: 耗材ID
//
// 返回值:
//
//	error: 错误信息
func (d *MaterialDao) Delete(id uint) error {
	return DB.Delete(&models.Material{}, id).Error
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
	err := DB.Where("code = ?", code).First(&m).Error
	return &m, err
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

	db := DB.Model(&models.Material{})
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
