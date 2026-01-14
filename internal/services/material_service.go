package services

import (
	"stock-flow/internal/dao"
	"stock-flow/internal/models"
)

// MaterialService 耗材业务服务
// 处理耗材基础信息的创建与查询
type MaterialService struct {
	materialDao dao.MaterialDao
}

// CreateMaterial 创建新耗材
//
// 参数:
//
//	m: 耗材信息
//
// 返回值:
//
//	error: 失败返回错误
func (s *MaterialService) CreateMaterial(m *models.Material) error {
	// Check if exists
	_, err := s.materialDao.GetByCode(m.Code)
	if err == nil {
		// Exists
		return nil // Or return error? Requirement says "Batch Import" might reuse.
		// Ideally if it exists, we just skip or update. For strict creation, maybe return error.
		// Let's assume this is for manual creation or initial setup.
	}
	return s.materialDao.Create(m)
}

// DeleteMaterial 删除耗材
//
// 参数:
//
//	id: 耗材ID
//
// 返回值:
//
//	error: 删除错误
func (s *MaterialService) DeleteMaterial(id uint) error {
	return s.materialDao.Delete(id)
}

// GetMaterialByCode 根据编码获取耗材
//
// 参数:
//
//	code: 物料编码
//
// 返回值:
//
//	*models.Material: 耗材信息
//	error: 查询错误
func (s *MaterialService) GetMaterialByCode(code string) (*models.Material, error) {
	return s.materialDao.GetByCode(code)
}

// GetMaterialList 获取耗材列表
//
// 参数:
//
//	page, pageSize: 分页参数
//	name: 物料名称
//
// 返回值:
//
//	[]models.Material: 列表
//	int64: 总数
//	error: 错误
func (s *MaterialService) GetMaterialList(page, pageSize int, name string) ([]models.Material, int64, error) {
	return s.materialDao.List(page, pageSize, name)
}
