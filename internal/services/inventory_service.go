package services

import (
	"stock-flow/internal/dao"
	"stock-flow/internal/models"
	"time"
)

// InventoryService 库存业务服务
// 处理入库、库存查询及效期预警
type InventoryService struct {
	inventoryDao dao.InventoryDao
	materialDao  dao.MaterialDao
}

// InboundDTO 入库请求数据传输对象
type InboundDTO struct {
	MaterialCode string // 物料编码
	MaterialName string // 物料名称
	Category     string // 分类
	Spec         string // 规格
	Unit         string // 单位
	Brand        string // 品牌
	BatchNo      string // 内部批号
	ExpiryDate   string // 有效期 (YYYY-MM-DD)
	Quantity     int64  // 数量
	InboundNo    string // 入库单号
	Mode         string // 模式: "append" 追加, "overwrite" 覆盖 (默认追加)
}

// Inbound 耗材入库
// 包含物料自动创建、批次去重或追加逻辑
//
// 参数:
//   dto: 入库数据
// 返回值:
//   error: 入库失败返回错误
func (s *InventoryService) Inbound(dto InboundDTO) error {
	// 1. 查找或创建物料基础信息
	mat, err := s.materialDao.GetByCode(dto.MaterialCode)
	if err != nil {
		// Create new material
		mat = &models.Material{
			Code:     dto.MaterialCode,
			Name:     dto.MaterialName,
			Category: dto.Category,
			Spec:     dto.Spec,
			Unit:     dto.Unit,
			Brand:    dto.Brand,
		}
		if err := s.materialDao.Create(mat); err != nil {
			return err
		}
	}

	// 2. 检查库存批次是否存在
	inv, err := s.inventoryDao.GetByMaterialAndBatch(mat.ID, dto.BatchNo)
	
	expiry, _ := time.Parse("2006-01-02", dto.ExpiryDate)

	if err == nil {
		// 批次存在
		if dto.Mode == "overwrite" {
			inv.InitialQty = dto.Quantity
			inv.CurrentQty = dto.Quantity
			inv.ExpiryDate = expiry
			inv.InboundNo = dto.InboundNo
		} else {
			// 默认追加模式
			inv.InitialQty += dto.Quantity
			inv.CurrentQty += dto.Quantity
			// Expiry date usually shouldn't change for same batch, but if it does, update it?
			// Let's assume batch implies same expiry.
		}
		return s.inventoryDao.Update(inv)
	}

	// 3. 创建新批次
	newInv := &models.Inventory{
		MaterialID: mat.ID,
		BatchNo:    dto.BatchNo,
		InboundNo:  dto.InboundNo,
		InitialQty: dto.Quantity,
		CurrentQty: dto.Quantity,
		ExpiryDate: expiry,
	}
	return s.inventoryDao.Create(newInv)
}

// GetInventoryList 综合查询库存
//
// 参数:
//   page, pageSize: 分页
//   materialName: 物料名
//   code: 编码
//   batchNo: 批号
//   status: 状态(0:全部, 1:正常, 2:临期, 3:过期)
// 返回值:
//   []models.Inventory: 库存列表
//   int64: 总数
//   error: 错误
func (s *InventoryService) GetInventoryList(page, pageSize int, materialName, code, batchNo string, status int) ([]models.Inventory, int64, error) {
	return s.inventoryDao.List(page, pageSize, materialName, code, batchNo, status)
}

// GetRecommendedBatches 获取推荐批次 (FEFO)
//
// 参数:
//   materialID: 物料ID
// 返回值:
//   []models.Inventory: 推荐批次列表
//   error: 错误
func (s *InventoryService) GetRecommendedBatches(materialID uint) ([]models.Inventory, error) {
	// FEFO strategy: First Expired First Out
	return s.inventoryDao.GetAvailableBatches(materialID)
}
