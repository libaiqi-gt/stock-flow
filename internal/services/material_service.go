package services

import (
	"fmt"
	"io"
	"stock-flow/internal/dao"
	"stock-flow/internal/models"
	"strconv"

	"github.com/xuri/excelize/v2"
)

// MaterialService 耗材业务服务
// 处理耗材基础信息的创建与查询
type MaterialService struct {
	materialDao dao.MaterialDao
}

// BatchImport 批量导入耗材
//
// 参数:
//
//	r: 文件读取器
//	ext: 文件扩展名
//
// 返回值:
//
//	*BatchImportResult: 导入结果
//	error: 严重错误
func (s *MaterialService) BatchImport(r io.ReadSeeker, ext string) (*BatchImportResult, error) {
	// 检查支持的扩展名
	supportedExts := map[string]bool{
		".xlsx": true,
		".xlsm": true,
		".xltx": true,
		".xltm": true,
	}

	if !supportedExts[ext] {
		return nil, fmt.Errorf("不支持的文件格式: %s", ext)
	}

	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("打开Excel文件失败: %v", err)
	}
	defer f.Close()

	// 获取第一个工作表名称
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("工作簿为空")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取工作表内容失败: %v", err)
	}

	result := &BatchImportResult{
		Errors: []string{},
		Msg:    "统计数据已排除表头行",
	}

	if len(rows) < 2 {
		return result, nil // Empty or header only
	}

	// 映射表头
	headerMap := make(map[string]int)
	for i, cell := range rows[0] {
		headerMap[cell] = i
	}

	// 必填字段
	required := []string{"物料编号", "物料名称", "物料类型", "规格", "单位", "厂家/品牌", "安全库存", "有效期报警时限/天"}
	for _, field := range required {
		if _, ok := headerMap[field]; !ok {
			return nil, fmt.Errorf("缺少必填列: %s", field)
		}
	}

	for i := 1; i < len(rows); i++ {
		result.Total++
		rowIdx := i + 1

		// Helper to get cell value
		getVal := func(colName string) string {
			idx, ok := headerMap[colName]
			if !ok || idx >= len(rows[i]) {
				return ""
			}
			return rows[i][idx]
		}

		// 1. Validate Required Fields
		code := getVal("物料编号")
		name := getVal("物料名称")
		category := getVal("物料类型")
		spec := getVal("规格")
		unit := getVal("单位")
		brand := getVal("厂家/品牌")
		safetyStockStr := getVal("安全库存")
		expiryAlertStr := getVal("有效期报警时限/天")

		if code == "" || name == "" || category == "" || spec == "" || unit == "" || brand == "" || safetyStockStr == "" || expiryAlertStr == "" {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 缺少必填字段", rowIdx))
			continue
		}

		// 2. Parse Numbers
		safetyStock, err := strconv.ParseInt(safetyStockStr, 10, 64)
		if err != nil || safetyStock < 0 {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 安全库存格式错误", rowIdx))
			continue
		}

		expiryAlert, err := strconv.Atoi(expiryAlertStr)
		if err != nil || expiryAlert < 0 {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 有效期报警时限格式错误", rowIdx))
			continue
		}

		// Optional Field
		openedExpiryStr := getVal("开封效期/天")
		openedExpiry := 180 // Default
		if openedExpiryStr != "" {
			val, err := strconv.Atoi(openedExpiryStr)
			if err != nil || val < 0 {
				result.Failed++
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 开封效期格式错误", rowIdx))
				continue
			}
			openedExpiry = val
		}

		// 3. Check Existence
		_, err = s.materialDao.GetByCode(code)
		if err == nil {
			// Exists: Skip or Update? Requirement usually implies skip for batch import to avoid overwriting.
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 物料编号 %s 已存在", rowIdx, code))
			continue
		}

		// 4. Create
		newMat := &models.Material{
			Code:             code,
			Name:             name,
			Category:         category,
			Spec:             spec,
			Unit:             unit,
			Brand:            brand,
			SafetyStock:      safetyStock,
			ExpiryAlertDays:  expiryAlert,
			OpenedExpiryDays: openedExpiry,
		}

		if err := s.materialDao.Create(newMat); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 数据库写入失败 %v", rowIdx, err))
		} else {
			result.Success++
		}
	}

	return result, nil
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
